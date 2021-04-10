package ndjson

import (
	"github.com/augmentable-dev/vtab"
	"go.riyazali.net/sqlite"
)

//
var cols = []vtab.Column{
	{Name: "line", Type: sqlite.SQLITE_INTEGER, NotNull: false, Hidden: false, Filters: nil},
	{Name: "json", Type: sqlite.SQLITE_TEXT, NotNull: false, Hidden: false, Filters: nil},
	{Name: "file_path", Type: sqlite.SQLITE_TEXT, NotNull: false, Hidden: true, Filters: []sqlite.ConstraintOp{sqlite.INDEX_CONSTRAINT_EQ}},
}

// NewVTab returns an ndjson virtual table
func NewVTab() sqlite.Module {
	return vtab.NewTableFunc("ndjson", cols, func(constraints []vtab.Constraint) (vtab.Iterator, error) {
		var filePath string
		for _, constraint := range constraints {
			if constraint.Op == sqlite.INDEX_CONSTRAINT_EQ {
				switch constraint.ColIndex {
				case 2:
					filePath = constraint.Value.Text()
				}
			}
		}

		iter, err := newIter(filePath)
		if err != nil {
			return nil, err
		}

		return iter, nil
	})
}
