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
		url      string
		headers  [][]string
		err      error
		contents []byte
		request  *http.Request
	)

	if len(values) > 0 {
		url = values[0].Text()
	} else if len(values) > 1 {
		heads := values[1].Text()
		headers = ParseHeaders(heads)
	} else {
		err := errors.New("input a single url as the argument to http get or a url with headers")
		ctx.ResultError(err)
	}
	request, err = httpGet(url, headers)
	if err != nil {
		ctx.ResultError(err)
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		ctx.ResultError(err)
	}
	contents, err = ioutil.ReadAll(response.Body)
	if err != nil {
		ctx.ResultError(err)
	}
	ctx.ResultText(string(contents))
}

func httpGet(requestUrl string, headers [][]string) (*http.Request, error) {
	request, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return nil, err
	}
	for _, header := range headers {
		request.Header.Add(header[0], header[1])
	}
	return request, nil
}

// NewHTTPGet returns a sqlite function for reading the contents of a file
func NewHTTPGet() sqlite.Function {
	return &get{}
}
