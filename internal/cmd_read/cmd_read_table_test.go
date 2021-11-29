package cmd_read

import (
	"database/sql"
	"strings"
	"testing"

	_ "github.com/augmentable-dev/flite/internal/sqlite"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"go.riyazali.net/sqlite"
)

func RowContent(rows *sql.Rows) (colCount int, contents [][]string, err error) {
	columns, err := rows.Columns()
	if err != nil {
		return colCount, nil, err
	}

	colCount = len(columns)

	pointers := make([]interface{}, len(columns))
	container := make([]sql.NullString, len(columns))

	for i := range pointers {
		pointers[i] = &container[i]
	}

	for rows.Next() {
		err = rows.Scan(pointers...)
		if err != nil {
			return colCount, nil, err
		}

		r := make([]string, len(columns))
		for i, c := range container {
			if c.Valid {
				r[i] = c.String
			} else {
				r[i] = "NULL"
			}
		}
		contents = append(contents, r)
	}
	return colCount, contents, rows.Err()

}

func TestReadTableCommand(t *testing.T) {
	commandToExec := "echo Hello World"
	checkOutput := strings.Split("Hello World", " ")
	sqlite.Register(func(api *sqlite.ExtensionApi) (sqlite.ErrorCode, error) {
		if err := api.CreateModule("cmd_table", NewCommandModule(),
			sqlite.EponymousOnly(true), sqlite.ReadOnly(true)); err != nil {
			return sqlite.SQLITE_ERROR, err
		}
		return sqlite.SQLITE_OK, nil
	})
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT * from cmd_table($1 ,' ')", commandToExec)
	if err != nil {
		t.Fatal(err)
	}
	_, content, err := RowContent(rows)
	for i := 0; i < len(content); i++ {
		if strings.TrimSpace(content[i][1]) != strings.TrimSpace(checkOutput[i]) {
			t.Fatalf("expected response: %s, got: %s", checkOutput[i], content[i][1])
		}
	}
	if rows.Err() != nil {
		t.Fatal(err)
	}
}
