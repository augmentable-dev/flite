package http

import (
	"errors"
	"io/ioutil"
	"net/http"

	"go.riyazali.net/sqlite"
)

type get struct{}

// TODO add PUT and POST stuff

func (m *get) Args() int           { return -1 }
func (m *get) Deterministic() bool { return false }
func (m *get) Apply(ctx *sqlite.Context, values ...sqlite.Value) {
	var (
		request  string
		err      error
		contents []byte
		response *http.Response
	)

	if len(values) > 0 {
		request = values[0].Text()
	} else {
		err := errors.New("input a single url as the argument to http get")
		ctx.ResultError(err)
	}

	response, err = http.Get(request)
	if err != nil {
		ctx.ResultError(err)
	}
	contents, err = ioutil.ReadAll(response.Body)
	if err != nil {
		ctx.ResultError(err)
	}

	ctx.ResultText(string(contents))
}

// NewHTTPGet returns a sqlite function for reading the contents of a file
func NewHTTPGet() sqlite.Function {
	return &get{}
}
