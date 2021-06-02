package file_read

import (
	"io/ioutil"
	"os"

	"go.riyazali.net/sqlite"
)

type fileRead struct{}

func (f *fileRead) Args() int           { return -1 }
func (f *fileRead) Deterministic() bool { return false }
func (f *fileRead) Apply(ctx *sqlite.Context, values ...sqlite.Value) {
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
			return
		}

	} else {
		contents, err = os.ReadFile(filePath)
		if err != nil {
			ctx.ResultError(err)
			return
		}
	}

	ctx.ResultText(string(contents))
}

// New returns a sqlite function for reading the contents of a file
func New() sqlite.Function {
	return &fileRead{}
}
