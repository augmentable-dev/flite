package cmd_read

import (
	"os/exec"

	"go.riyazali.net/sqlite"
)

type cmdRead struct{}

func (f *cmdRead) Args() int           { return 1 }
func (f *cmdRead) Deterministic() bool { return false }
func (f *cmdRead) Apply(ctx *sqlite.Context, values ...sqlite.Value) {
	var (
		command string
		err     error
	)

	if len(values) > 0 {
		command = values[0].Text()
	}
	cmd := exec.Command("bash", "-c", command)
	stdout, err := cmd.Output()

	if err != nil {
		ctx.ResultError(err)
	}

	ctx.ResultText(string(stdout))
}

// New returns a sqlite function for reading the contents of a cmd
func New() sqlite.Function {
	return &cmdRead{}
}
