package tmux

import (
	"fmt"
	"strings"
)

func GetActiveSession() (string, error) {
	var activeSession []byte
	var err error

	if activeSession, err = Command("display-message", "-p", "-F", "#{session_name}"); err != nil {
		fmt.Printf("err after calling display-message: '%s'\n", err)
		return "", err
	}

	s := Trim(string(activeSession))
	return s, nil
}

func (r *Runner) AttachSession(sessionName string) error {
	_, err := r.Run(fmt.Sprintf("attach -t %s", sessionName))
	return err
}

func (r *Runner) ListSessions() ([]string, error) {
	result, err := r.Run("list-sessions -F '#{session_name}'")
	if err != nil {
		return nil, err
	}

	sessions := strings.Split(Trim(result), "\n")
	return sessions, nil

}

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
