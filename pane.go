package tmux

import "fmt"

// Set the width of the given pane
func (r *Runner) SetPaneWidth(pane string, width int) error {
	var cmd string = fmt.Sprintf("resize-pane -x '%d' -t '%s'", width, pane)

	_, err := r.Run(cmd)
	return err
}
