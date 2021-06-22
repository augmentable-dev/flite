package cmd

import (
	"fmt"
	"strings"
	"testing"

	_ "github.com/augmentable-dev/flite/internal/sqlite"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"go.riyazali.net/sqlite"
)

func TestExec(t *testing.T) {
	exec := New()
	text := "Hello, World!"
	sqlite.Register(func(api *sqlite.ExtensionApi) (sqlite.ErrorCode, error) {
		if err := api.CreateFunction("exec", exec); err != nil {
			return sqlite.SQLITE_ERROR, err
		}
		return sqlite.SQLITE_OK, nil
	})
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	row := db.QueryRow(fmt.Sprintf("SELECT exec('bash','echo %s')", text))
	err = row.Err()
	if err != nil {
		t.Fatal(err)
	}
	var res string
	err = row.Scan(&res)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Compare(text, res) == 0 {
		t.Fatalf("expected response: %s, got: %s", text, res)
	}
}
