package tmux

import (
	"fmt"
	"strconv"
	"strings"
)

// Tmux doesn't have a built-in notion of a 'column'. A column for the purpose
// of these functions is one or more panes stacked on top of each other. For
// example in a layout like this:
//
//   +---+---+---+
//   | 0 |   |   |
//   +---+ 2 | 3 |
//   | 1 |   |   |
//   +---+---+---+
//
// the first column is panes 0 and 1, the second is pane 2, and the third is
// pane 3.
//
// There are edge cases, like this:
//
//   +---+---+---+
//   | 0 | 1 |   |
//   +---+---+ 3 |
//   |   2   |   |
//   +-------+---+
//
// Is this two columns, or three? These methods count it as 3: each pane at the
// top of the window is the top of a column.

type Column struct {
	// The pane ID of the pane at the top of this column
	Pane string

	// The width of this column
	Width int
}

func (r *Runner) ListColumns() ([]Column, error) {
	var err error

	var output string
	var cmd string = "list-panes -F '#{pane_id} #{pane_width}' -f '#{m:#{pane_at_top},1}'"
	if output, err = r.Run(cmd); err != nil {
		return nil, err
	}

	columns := make([]Column, 0)

	lines := strings.Split(Trim(output), "\n")
	for _, line := range lines {
		tokens := strings.Split(line, " ")
		if len(tokens) != 2 {
			return nil, fmt.Errorf("expected line to be a string with two elements separated by spaces but found '%s'", line)
		}

		pane := tokens[0]

		var width int
		if width, err = strconv.Atoi(tokens[1]); err != nil {
			return nil, fmt.Errorf("Error parsing second element of line '%s': '%s'", line, err.Error())
		}

		columns = append(columns, Column{Pane: pane, Width: width})
	}

	return columns, nil
}
