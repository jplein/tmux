package tmux

import "os/exec"

// Run a tmux shell command with the provided arguments, and return its output.
func Command(c Config, args ...string) ([]byte, error) {
	var tmuxPath string
	var err error

	if tmuxPath, err = Tmux(); err != nil {
		return []byte(""), err
	}

	if c.Socket != "" {
		args = append([]string{"-L", c.Socket}, args...)
	}
	return exec.Command(tmuxPath, args...).Output()
}
