package cmd_read

import (
	"strings"
	"testing"

	_ "github.com/augmentable-dev/flite/internal/sqlite"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"go.riyazali.net/sqlite"
)

func TestReadCommand(t *testing.T) {
	commandToExec := "echo Hello World"
	sqlite.Register(func(api *sqlite.ExtensionApi) (sqlite.ErrorCode, error) {
		if err := api.CreateFunction("cmd", New()); err != nil {
			return sqlite.SQLITE_ERROR, err
		}
		return sqlite.SQLITE_OK, nil
	})
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	row := db.QueryRow("SELECT cmd($1)", commandToExec)
	err = row.Err()
	if err != nil {
		t.Fatal(err)
	}
	var res string
	err = row.Scan(&res)
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(res) != "Hello World" {
		t.Fatalf("expected response: %s, got: %s", "Hello World", res)
	}
}
