package tmux

import "os/exec"

const TmuxExec = "tmux"

// Get the path to the "tmux" executable, if found in the current PATH. Returns
// an error if not found.
func Tmux() (string, error) {
	return exec.LookPath("tmux")
}
