package readfile

import (
	"io/ioutil"
	"os"

	"go.riyazali.net/sqlite"
)

type readFile struct{}

func (m *readFile) Args() int           { return -1 }
func (m *readFile) Deterministic() bool { return false }
func (m *readFile) Apply(ctx *sqlite.Context, values ...sqlite.Value) {
	var (
		filePath string
		err      error
		contents []byte
	)

	if len(values) > 0 {
		filePath = values[0].Text()
	}

	if filePath == "" {
		contents, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			ctx.ResultError(err)
		}

	} else {
		contents, err = os.ReadFile(filePath)
		if err != nil {
			ctx.ResultError(err)
		}
	}

	ctx.ResultText(string(contents))
}

func NewReadFile() sqlite.Function {
	return &readFile{}
}
