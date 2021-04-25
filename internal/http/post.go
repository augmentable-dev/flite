package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"go.riyazali.net/sqlite"
)

/*
* post takes in an input in the format of http_post(url, args...) and will return the response
* from the url
 */
type post struct{}

func (m *post) Args() int           { return -1 }
func (m *post) Deterministic() bool { return false }
func (m *post) Apply(ctx *sqlite.Context, values ...sqlite.Value) {
	var (
		postUrl  string
		err      error
		contents []byte
		response *http.Response
		vars     map[string]interface{}
	)

	err = json.Unmarshal([]byte(values[1].Text()), &vars)
	if err != nil {
		ctx.ResultError(err)
	}

	actualVars := url.Values{}
	for k, v := range vars {
		actualVars.Set(k, fmt.Sprint(v))
	}

	postUrl = values[0].Text()
	response, err = http.PostForm(postUrl, actualVars)
	if err != nil {
		contents, _ = ioutil.ReadAll(response.Body)
		println("Post request returned an error", err.Error(), "with code ", contents)
		ctx.ResultError(err)
	}
	contents, err = ioutil.ReadAll(response.Body)
	if err != nil {
		ctx.ResultError(err)
	}

	ctx.ResultText(string(contents))
}

// NewHTTPPost returns a sqlite function for reading the contents of a file
func NewHTTPPost() sqlite.Function {
	return &post{}
}
