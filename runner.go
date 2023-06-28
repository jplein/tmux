package tmux

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Runner struct {
	writePipe   io.WriteCloser
	readPipe    io.ReadCloser
	readScanner bufio.Scanner
	tmpSession  string
	tmuxCommand *exec.Cmd
}

func (r *Runner) readNextLine() (string, error) {
	b := r.readScanner.Scan()
	if !b {
		err := r.readScanner.Err()
		if err != nil {
			return "", err
		} else {
			return "", io.EOF
		}
	}

	return r.readScanner.Text(), nil
}

const tmuxBeginMarker = "%begin"
const tmuxEndMarker = "%end"
const tmuxErrorMarker = "%error"

func (r *Runner) isBeginLine(line string) bool {
	if len(line) < len(tmuxBeginMarker) {
		return false
	}

	beginMarkerLength := len(tmuxBeginMarker)
	return line[:beginMarkerLength] == tmuxBeginMarker
}

func (r *Runner) getExpectedEndLine(beginLine string) string {
	return fmt.Sprintf("%s %s", tmuxEndMarker, beginLine[len(tmuxBeginMarker)+1:])
}

func (r *Runner) getExpectedErrorLine(beginLine string) string {
	return fmt.Sprintf("%s %s", tmuxErrorMarker, beginLine[len(tmuxBeginMarker)+1:])
}

type readState int

const (
	stateBeforeOutput readState = 0
	stateOutput       readState = 1
	stateError        readState = 2
	stateEnd          readState = 3
)

func (r *Runner) readCommandOutput() (string, error) {
	done := false

	var expectedEndLine string
	var expectedErrorLine string

	outputLines := make([]string, 0)

	var state readState = stateBeforeOutput

	type returnval struct {
		output string
		err    error
	}
	var result returnval

	for !done {
		switch state {
		case stateBeforeOutput:
			line, err := r.readNextLine()
			if err != nil {
				return "", err
			}

			if r.isBeginLine(line) {
				state = stateOutput
			}

			expectedEndLine = r.getExpectedEndLine(line)
			expectedErrorLine = r.getExpectedErrorLine(line)
		case stateOutput:
			line, err := r.readNextLine()
			if err != nil {
				return "", err
			}

			if line == expectedEndLine {
				state = stateEnd
			} else if line == expectedErrorLine {
				state = stateError
			} else {
				outputLines = append(outputLines, line)
			}
		case stateEnd:
			result = returnval{
				output: strings.Join(outputLines, "\n"),
				err:    nil,
			}
			done = true
		case stateError:
			result = returnval{
				output: "",
				err: fmt.Errorf(
					fmt.Sprintf(
						"tmux error: %s",
						strings.Join(outputLines, "\n"),
					),
				),
			}
			done = true
		}
	}

	return result.output, result.err
}

// Before the tmux -C process used by the runner has started, use this to get
// the list of session names
func (r *Runner) getSessionNamesByCommand() ([]string, error) {
	var err error

	var output []byte
	if output, err = Command("list-sessions", "-F", "#{session_name}"); err != nil {
		return nil, err
	}

	sessionNames := strings.Split(string(output), "\n")
	return sessionNames, nil
}

// After the tmux -C process used by the runner has already started, use this to
// get the list of session names without spawning a new process
func (r *Runner) getSessionNames() ([]string, error) {
	var err error
	var output string

	if output, err = r.Run("list-sessions -F '#{session_name}'"); err != nil {
		return nil, err
	}

	sessionNames := strings.Split(Trim(output), "\n")
	return sessionNames, nil
}

func (r *Runner) Init() error {
	var err error

	var tmuxPath string
	if tmuxPath, err = Tmux(); err != nil {
		return err
	}

	var sessionsBeforeStart []string
	if sessionsBeforeStart, err = r.getSessionNamesByCommand(); err != nil {
		return err
	}

	r.tmuxCommand = exec.Command(tmuxPath, "-C")

	writePipe, err := r.tmuxCommand.StdinPipe()
	if err != nil {
		return err
	}
	r.writePipe = writePipe

	readPipe, err := r.tmuxCommand.StdoutPipe()
	if err != nil {
		return err
	}
	r.readPipe = readPipe

	r.readScanner = *bufio.NewScanner(readPipe)

	r.tmuxCommand.Start()

	// When tmux -C first runs, it prints a pair of %begin and %end lines with
	// nothing in between
	_, err = r.readCommandOutput()
	if err != nil {
		return err
	}

	var sessionsAfterStart []string
	if sessionsAfterStart, err = r.getSessionNames(); err != nil {
		return err
	}

	sessionMap := make(map[string]bool)
	for _, session := range sessionsBeforeStart {
		sessionMap[session] = true
	}

	newSessions := make([]string, 0)
	for _, session := range sessionsAfterStart {
		if !sessionMap[session] {
			newSessions = append(newSessions, session)
		}
	}

	if len(newSessions) != 1 {
		return fmt.Errorf("expected exactly 1 new session but found %d: %s", len(newSessions), strings.Join(newSessions, ","))
	}

	r.tmpSession = newSessions[0]

	return nil
}

func (r *Runner) Run(cmd string) (string, error) {
	cmdBuf := []byte(fmt.Sprintf("%s\n", cmd))
	bytesWritten, err := r.writePipe.Write(cmdBuf)
	if err != nil {
		return "", err
	}

	if bytesWritten != len(cmd)+1 {
		fmt.Printf("Expected to write %d bytes but wrote %d", len(cmd)+1, bytesWritten)
	}

	var output string
	if output, err = r.readCommandOutput(); err != nil {
		return "", fmt.Errorf(fmt.Sprintf("Error running command '%s': '%s", cmd, err.Error()))
	}

	return output, nil
}

func (r *Runner) Close() error {
	defer func() {
		e := r.tmuxCommand.Process.Kill()
		if e != nil {
			os.Stderr.Write([]byte(fmt.Sprintf("Error killing tmux -C process: '%s'", e.Error())))
		}
	}()

	var err error
	if _, err = r.Run(fmt.Sprintf("kill-session -t %s", r.tmpSession)); err != nil {
		return err
	}

	return err
}
