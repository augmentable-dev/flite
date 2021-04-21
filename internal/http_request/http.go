package http_request

import (
	"io/ioutil"
	"net/http"
	"os"

	"go.riyazali.net/sqlite"
)

type http_request struct{}

func (m *http_request) Args() int           { return -1 }
func (m *http_request) Deterministic() bool { return false }
func (m *http_request) Apply(ctx *sqlite.Context, values ...sqlite.Value) {
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

// Newhttp_request returns a sqlite function for reading the contents of a file
func NewHttp_request() sqlite.Function {
	return &http_request{}
}
