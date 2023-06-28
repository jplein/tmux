package tmux

import (
	"fmt"
	"strconv"
	"strings"
)

func (r *Runner) GetActiveWindow() (string, error) {
	return r.Run("list-windows -F '#{window_id}' -f '#{m:#{window_active},1}'")
}

func (r *Runner) GetWindowDimensions(windowName string) (int, int, error) {
	// ['list-windows', '-F', '#{window_width} #{window_height}', '-f', '#{m:#{window_active},1}']
	var err error

	var output string
	if output, err = r.Run("list-windows -F '#{window_width} #{window_height}' -f '#{m:#{window_active},1}'"); err != nil {
		return 0, 0, err
	}

	dimensions := strings.Split(Trim(output), " ")
	if len(dimensions) != 2 {
		return 0, 0, fmt.Errorf(fmt.Sprintf("expected command to return a string with two elements separated by a space but found '%s'", output))
	}

	var width, height int
	if width, err = strconv.Atoi(dimensions[0]); err != nil {
		return 0, 0, err
	}
	if height, err = strconv.Atoi(dimensions[1]); err != nil {
		return 0, 0, err
	}

	return width, height, nil
}
