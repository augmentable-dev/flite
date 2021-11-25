package cmd_read

import (
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/augmentable-dev/vtab"
	"go.riyazali.net/sqlite"
)

type iterCommands struct {
	command   string
	delimiter string
	contents  []string
	index     int
}

func (i *iterCommands) Column(c int) (interface{}, error) {
	switch c {
	case 0:
		return i.command, nil
	case 1:
		return i.delimiter, nil
	case 2:
		return i.index, nil
	case 3:
		return i.contents[i.index], nil
	}
	return nil, fmt.Errorf("unknown column")
}
func (i *iterCommands) Next() (vtab.Row, error) {
	i.index += 1
	if i.index >= len(i.contents) {
		return nil, io.EOF
	}

	return i, nil
}

var commandCols = []vtab.Column{
	{Name: "command", Type: sqlite.SQLITE_TEXT, NotNull: true, Hidden: true, Filters: []sqlite.ConstraintOp{sqlite.INDEX_CONSTRAINT_EQ}},
	{Name: "delimiter", Type: sqlite.SQLITE_TEXT, NotNull: true, Hidden: true, Filters: []sqlite.ConstraintOp{sqlite.INDEX_CONSTRAINT_EQ}},
	{Name: "line_no", Type: sqlite.SQLITE_INTEGER},
	{Name: "results", Type: sqlite.SQLITE_TEXT},
}

func NewCommandModule() sqlite.Module {
	return vtab.NewTableFunc("command", commandCols, func(constraints []vtab.Constraint) (vtab.Iterator, error) {
		var command, delimiter string
		for _, constraint := range constraints {
			if constraint.Op == sqlite.INDEX_CONSTRAINT_EQ {
				switch constraint.ColIndex {
				case 0:
					command = constraint.Value.Text()
				case 1:
					delimiter = constraint.Value.Text()
				}
			}
		}

		if delimiter == "" {
			delimiter = "\n"
		}
		cmd := exec.Command("bash", "-c", command)
		stdout, err := cmd.Output()
		if err != nil {
			return nil, err
		}
		contents := strings.Split(string(stdout), delimiter)
		iter := &iterCommands{command, delimiter, contents, -1}
		return iter, nil
	})
}
