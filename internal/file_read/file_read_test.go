package file_read

import (
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"

	_ "github.com/augmentable-dev/flite/internal/sqlite"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"go.riyazali.net/sqlite"
)

func TestReadFile(t *testing.T) {
	readFile := New()
	tmpFile, err := ioutil.TempFile(os.TempDir(), "prefix")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	text := ""
	rand.Seed(time.Now().Unix())
	x := rand.Intn(100)
	for j := x; j > 0; j-- {
		text += "a\n"
	}
	if _, err = tmpFile.Write([]byte(text)); err != nil {
		t.Fatal(err)
	}

	sqlite.Register(func(api *sqlite.ExtensionApi) (sqlite.ErrorCode, error) {
		if err := api.CreateFunction("file_read", readFile); err != nil {
			return sqlite.SQLITE_ERROR, err
		}
		return sqlite.SQLITE_OK, nil
	})
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	row := db.QueryRow("SELECT file_read($1)", tmpFile.Name())
	err = row.Err()
	if err != nil {
		t.Fatal(err)
	}
	var res string
	err = row.Scan(&res)
	if err != nil {
		t.Fatal(err)
	}
	if res != text {
		t.Fatalf("expected response: %s, got: %s", text, res)
	}
}
