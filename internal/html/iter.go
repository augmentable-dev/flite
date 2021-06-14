package html

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/augmentable-dev/vtab"

	"golang.org/x/net/html"
)

type iter struct {
	html_body string
	tokenizer *html.Tokenizer
	// node       *html.Node
	// token      string
	token_type html.TokenType
	data       string
	raw        string
	// end        error
}

func newIter(html_body string) (*iter, error) {
	// var (
	// 	tokenizer *html.Tokenizer
	// )
	// println("new iter")

	x, err := http.Get(html_body)
	if err != nil {
		return nil, err
	}

	tokenizer := html.NewTokenizer(x.Body)
	token_type := tokenizer.Next()

	// if html_body != "" {
	// 	tokenizer = html.NewTokenizer(r)
	// }
	// token := tokenizer.Token().Data
	// raw_token := tokenizer.Raw()
	return &iter{
		html_body: html_body,
		data:      tokenizer.Token().Data,
		raw:       string(tokenizer.Raw()),
		// node:      node,
		tokenizer: tokenizer,
		// raw_token: string(raw_token),
		token_type: token_type,
	}, nil
}

func (i *iter) Column(c int) (interface{}, error) {
	//println(c)
	switch c {
	case 0:
		return i.raw, nil
	case 2:
		return i.token_type.String(), nil
	case 1:
		return i.html_body, nil
	}

	return nil, fmt.Errorf("unknown column")
}

func (i *iter) Next() (vtab.Row, error) {

	for strings.TrimSpace(string(i.tokenizer.Raw())) == "" && i.token_type != html.ErrorToken {
		i.token_type = i.tokenizer.Next()
	}

	//println("next")
	if i.tokenizer.Next() == html.ErrorToken {
		println(i.tokenizer.Err().Error())
		return nil, i.tokenizer.Err()
	}
	i.token_type = i.tokenizer.Token().Type
	i.raw = string(i.tokenizer.Raw())

	//if strings.TrimSpace(string)
	return i, nil
}
