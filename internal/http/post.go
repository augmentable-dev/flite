package http

import (
	"errors"
	"io/ioutil"
	"net/http"

	"go.riyazali.net/sqlite"
)

type post struct {
	client *http.Client
}

// TODO add PUT and POST stuff

func (f *post) Args() int           { return -1 }
func (f *post) Deterministic() bool { return false }
func (f *post) Apply(ctx *sqlite.Context, values ...sqlite.Value) {
	var (
		url      string
		headers  [][]string
		err      error
		contents []byte
		request  *http.Request
	)
	if len(values) > 1 {
		url = values[0].Text()
		heads := values[1].Text()
		headers = ParseHeaders(heads)
	} else {
		err := errors.New("input a url with headers as the argument to post")
		ctx.ResultError(err)
	}
	request, err = HttpRequest(url, headers, "POST")
	if err != nil {
		ctx.ResultError(err)
	}
	response, err := f.client.Do(request)
	if err != nil {
		ctx.ResultError(err)
	}
	contents, err = ioutil.ReadAll(response.Body)
	if err != nil {
		ctx.ResultError(err)
	}
	ctx.ResultText(string(contents))
}

// NewHTTPpost returns a sqlite function for reading the contents of a file
func NewHTTPpost() sqlite.Function {
	return &post{http.DefaultClient}
}
