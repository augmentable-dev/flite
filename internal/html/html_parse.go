package html

import (
	"github.com/augmentable-dev/vtab"

	"go.riyazali.net/sqlite"
)

var cols = []vtab.Column{
	{Name: "content", Type: sqlite.SQLITE_TEXT, NotNull: false, Hidden: false, Filters: nil},
	{Name: "html_body", Type: sqlite.SQLITE_TEXT, NotNull: false, Hidden: true, Filters: []sqlite.ConstraintOp{sqlite.INDEX_CONSTRAINT_EQ}},
	{Name: "type", Type: sqlite.SQLITE_TEXT, NotNull: false, Hidden: false, Filters: nil},
	{Name: "data", Type: sqlite.SQLITE_TEXT, NotNull: false, Hidden: false, Filters: nil},
	{Name: "attributes", Type: sqlite.SQLITE_TEXT, NotNull: false, Hidden: false, Filters: nil},
}

// NewVTab returns a line reader virtual table
func NewVTab() sqlite.Module {
	return vtab.NewTableFunc("html_parse", cols, func(constraints []vtab.Constraint) (vtab.Iterator, error) {
		var body string
		for _, constraint := range constraints {
			if constraint.Op == sqlite.INDEX_CONSTRAINT_EQ {
				switch constraint.ColIndex {
				case 1:
					body = constraint.Value.Text()
				}
			}
		}
		//println("body", body)
		iter, err := newIter(body)
		if err != nil {
			return nil, err
		}

		return iter, nil
	})
}

// type html_parse struct{}

// func (f *html_parse) Args() int           { return -1 }
// func (f *html_parse) Deterministic() bool { return false }
// func (f *html_parse) Apply(ctx *sqlite.Context, values ...sqlite.Value) {
// 	var (
// 		body string
// 		//err      error
// 		//node     *html.Node
// 		contents string
// 	)

// 	if len(values) > 0 {
// 		body = values[0].Text()
// 	}

// 	bReader := strings.NewReader(body)
// 	tokenizer := html.NewTokenizer(bReader)
// 	tokenizer.Token().

// 	for token := tokenizer.Next(); token.String() != "nil"; token = tokenizer.Next() {
// 		print(tokenizer.Token().Data)
// 		contents += tokenizer.Token().String() + "\n" + tokenizer.Token().Data
// 	}

// 	ctx.ResultText(string(contents))
// }

// // New returns a sqlite function for reading the contents of a file
// func New() sqlite.Function {
// 	return &html_parse{}
// }
