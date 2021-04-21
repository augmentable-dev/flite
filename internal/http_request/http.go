package readHTTP

import (
	"io/ioutil"
	"net/http"
	"os"

	"go.riyazali.net/sqlite"
)

type readHTTP struct{}

func (m *readHTTP) Args() int           { return -1 }
func (m *readHTTP) Deterministic() bool { return false }
func (m *readHTTP) Apply(ctx *sqlite.Context, values ...sqlite.Value) {
	var (
		request  string
		err      error
		contents []byte
		response *http.Response
	)

	if len(values) > 0 {
		request = values[0].Text()
	}

	if request == "" {
		contents, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			ctx.ResultError(err)
		}
		response, err = http.Get(string(contents))

	} else {
		response, err = http.Get(request)
		if err != nil {
			ctx.ResultError(err)
		}
	}
	ret, err := ioutil.ReadAll(response.Body)
	ctx.ResultText(string(ret))
}

// NewreadHTTP returns a sqlite function for reading the contents of a file
func NewReadHTTP() sqlite.Function {
	return &readHTTP{}
}
