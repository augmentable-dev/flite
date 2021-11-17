package cmd_read

import (
	"os/exec"

	"go.riyazali.net/sqlite"
)

type cmdRead struct{}

func (f *cmdRead) Args() int           { return -1 }
func (f *cmdRead) Deterministic() bool { return false }
func (f *cmdRead) Apply(ctx *sqlite.Context, values ...sqlite.Value) {
	var (
		cmdName string
		args    []string
		err     error
	)

	if len(values) > 0 {
		cmdName = values[0].Text()
	}
	for i := 1; i < len(values); i++ {
		args = append(args, values[i].Text())
	}
	command := exec.Command(cmdName, args...)
	stdout, err := command.Output()

	if err != nil {
		ctx.ResultError(err)
	}

	ctx.ResultText(string(stdout))
}

// New returns a sqlite function for reading the contents of a cmd
func New() sqlite.Function {
	return &cmdRead{}
}
