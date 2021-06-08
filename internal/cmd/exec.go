package cmd

import (
	"fmt"
	exe "os/exec"

	"go.riyazali.net/sqlite"
)

type exec struct{}

func (f *exec) Args() int           { return -1 }
func (f *exec) Deterministic() bool { return false }
func (f *exec) Apply(ctx *sqlite.Context, values ...sqlite.Value) {
	var (
		args  string
		shell string
		err   error
		out   []byte
		cmd   *exe.Cmd
	)
	if len(values) == 2 {
		shell = values[0].Text()
		args = values[1].Text()
	} else {
		err = fmt.Errorf("must input args in format (<shell_to_execute>, <args>)")
		ctx.ResultError(err)
		return
	}
	cmd = exe.Command(shell, "-c", args)
	out, err = cmd.CombinedOutput()
	if err != nil {
		ctx.ResultError(err)
		return
	}

	ctx.ResultText(string(out))

}

// New returns a sqlite function for reading the contents of a file
func New() sqlite.Function {
	return &exec{}
}
