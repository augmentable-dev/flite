package html

import (
	"fmt"
	"math"
	"net/http"
	"strings"

	"github.com/augmentable-dev/vtab"
	"golang.org/x/net/html"
)

type iter struct {
	html_body  string
	tokenizer  *html.Tokenizer
	token_type html.TokenType
	data       string
	raw        string
	path       []string
}

func newIter(html_body string) (*iter, error) {
	var path []string
	x, err := http.Get(html_body)
	if err != nil {
		return nil, err
	}
	tokenizer := html.NewTokenizer(x.Body)
	tokenizer.AllowCDATA(true)
	token_type := tokenizer.Next()

	if token_type == html.StartTagToken {
		path = append(path, tokenizer.Token().Data)

	}
	return &iter{
		html_body:  html_body,
		data:       tokenizer.Token().DataAtom.String(),
		raw:        string(tokenizer.Raw()),
		path:       path,
		tokenizer:  tokenizer,
		token_type: token_type,
	}, nil
}

func (i *iter) Column(c int) (interface{}, error) {
	switch c {
	case 0:
		return i.raw, nil
	case 1:
		return i.html_body, nil
	case 2:
		return i.token_type.String(), nil
	case 3:
		return strings.Join(i.path, "/"), nil
	case 4:
		s := ""
		for k, v, n := i.tokenizer.TagAttr(); n; k, v, n = i.tokenizer.TagAttr() {
			fmt.Println(string(k), " : ", string(v))
			s += string(k) + " : " + string(v)
		}
		s = strings.ReplaceAll(s[:int(math.Max(0, float64(len(s)-1)))], "\n", " ")
		return s, nil
	}

	return nil, fmt.Errorf("unknown column")
}

func (i *iter) Next() (vtab.Row, error) {
	i.token_type = i.tokenizer.Next()
	for strings.TrimSpace(string(i.tokenizer.Raw())) == "" && i.token_type != html.ErrorToken {
		i.token_type = i.tokenizer.Next()
	}
	switch i.token_type {

	case html.ErrorToken:
		return nil, i.tokenizer.Err()

	case html.TextToken:
		i.data = i.token_type.String()

	case html.StartTagToken, html.EndTagToken:
		x := i.tokenizer.Token()
		i.data = x.Data
		i.raw = string(i.tokenizer.Raw())
		if i.token_type == html.StartTagToken {
			i.path = append(i.path, x.Data)
		} else {
			i.path = i.path[:len(i.path)-1]
		}

	}
	return i, nil
}
