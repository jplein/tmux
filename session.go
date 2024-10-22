package tmux

import (
	"fmt"
	"strings"
)

// Get the name of the active session. Returns an error if there is no active
// session.
func GetActiveSession(c Config) (string, error) {
	var activeSession []byte
	var err error

	if activeSession, err = Command(c, "display-message", "-p", "-F", "#{session_name}"); err != nil {
		fmt.Printf("err after calling display-message: '%s'\n", err)
		return "", err
	}

	s := Trim(string(activeSession))
	return s, nil
}

// Attach to the session with the provided name
func (r *Runner) AttachSession(sessionName string) error {
	_, err := r.Run(fmt.Sprintf("attach -t '%s'", sessionName))
	return err
}

// Returns a list of the names of the running sessions
func (r *Runner) ListSessions() ([]string, error) {
	result, err := r.Run("list-sessions -F '#{session_name}'")
	if err != nil {
		return nil, err
	}

	sessions := strings.Split(Trim(result), "\n")
	return sessions, nil

}

// Start a new session
func (r *Runner) StartSession(name string) error {
	sessions, err := r.ListSessions()
	if err != nil {
		return err
	}

	sessionRunning := false
	for _, s := range sessions {
		if s == name {
			sessionRunning = true
			break
		}
	}

	if !sessionRunning {
		_, err := r.Run(fmt.Sprintf("new-session -d -s %s", name))
		if err != nil {
			return err
		}
	}

	return nil

}
