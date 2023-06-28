package tmux

import "os/exec"

const TmuxExec = "tmux"

func Tmux() (string, error) {
	return exec.LookPath("tmux")
}
