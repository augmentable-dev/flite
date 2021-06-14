package html

import (
	"fmt"
	"strings"

	"github.com/augmentable-dev/vtab"
	"golang.org/x/net/html"
)

type iter struct {
	html_body  string
	tokenizer  *html.Tokenizer
	node       *html.Node
	token      string
	token_type html.TokenType
	raw_token  string
	end        error
}

func newIter(html_body string) (*iter, error) {
	// var (
	// 	tokenizer *html.Tokenizer
	// )
	// println("new iter")
	r := strings.NewReader(html_body)
	tokenizer := html.NewTokenizer(r)

	// if html_body != "" {
	// 	tokenizer = html.NewTokenizer(r)
	// }
	// token := tokenizer.Token().Data
	// raw_token := tokenizer.Raw()
	return &iter{
		html_body: html_body,
		// node:      node,
		tokenizer: tokenizer,
		// raw_token: string(raw_token),
		// token:     token,
	}, nil
}

func (i *iter) Column(c int) (interface{}, error) {
	println("column", c)
	switch c {
	case 0:
		return strings.Trim(string(i.tokenizer.Raw()), "\n"), nil
	case 1:
		return strings.Trim(i.tokenizer.Token().Data, "\n"), nil
	case 2:
		return i.token_type.String(), nil
	case 3:
		return i.html_body, nil
	}

	return nil, fmt.Errorf("unknown column")
}

func (i *iter) Next() (vtab.Row, error) {
	println("next")
	i.token_type = i.tokenizer.Next()
	// println(i.token_type.String())

	// keepGoing := i.tokenizer.Next()
	// i.raw_token = string(i.tokenizer.Raw())
	// i.token = i.tokenizer.Token().Data
	//println(i.tokenizer.Token().Type.String())
	for strings.TrimSpace(string(i.tokenizer.Raw())) == "" && i.token_type != html.ErrorToken {
		i.token_type = i.tokenizer.Next()
	}
	if i.token_type == html.ErrorToken {
		println(i.tokenizer.Err().Error())
		return nil, i.tokenizer.Err()
	}
	return i, nil
}
