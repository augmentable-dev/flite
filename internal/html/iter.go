package html

import (
	"fmt"
	"io"
	"strings"

	"github.com/augmentable-dev/vtab"
	"golang.org/x/net/html"
)

type iter struct {
	html_body string
	tokenizer *html.Tokenizer
	node      *html.Node
	token     string
	raw_token string
	end       error
}

func newIter(html_body string) (*iter, error) {
	// var (
	// 	tokenizer *html.Tokenizer
	// )
	println("new iter")
	r := strings.NewReader(html_body)
	node, err := html.Parse(r)
	if err != nil {
		println(err.Error())
		return nil, err
	}
	node = node.FirstChild
	println(node.Namespace)

	// if html_body != "" {
	// 	tokenizer = html.NewTokenizer(r)
	// }
	// token := tokenizer.Token().Data
	// raw_token := tokenizer.Raw()
	return &iter{
		html_body: html_body,
		node:      node,
		// tokenizer: tokenizer,
		// raw_token: string(raw_token),
		// token:     token,
	}, nil
}

func (i *iter) Column(c int) (interface{}, error) {
	println("column", c)
	switch c {
	case 0:
		return i.node.Data, nil
	case 1:
		return i.node.Type, nil
	case 2:
		return "jubaloo", nil
	}

	return nil, fmt.Errorf("unknown column")
}

func (i *iter) Next() (vtab.Row, error) {
	println("next")
	// keepGoing := i.tokenizer.Next()
	// i.raw_token = string(i.tokenizer.Raw())
	// i.token = i.tokenizer.Token().Data
	//println(i.tokenizer.Token().Type.String())
	println(i.node == nil)
	for _, j := range i.node.Attr {
		println(j.Namespace)
		println(j.Key)
		println(j.Val)

	}

	println(i.node == nil)
	if i.node != nil {
		println(i.node.Data)
		println(i.node.LastChild.Data)
		println(i.node.Type)
	}
	if i.node == nil {
		return nil, io.EOF
	}
	return i, nil
}
