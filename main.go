package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"

	_ "github.com/augmentable-dev/flite/internal/sqlite"
	_ "github.com/augmentable-dev/flite/pkg/ext"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var query string
	if len(os.Args) > 0 {
		query = os.Args[1]
	}

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = new(interface{})
	}

	enc := json.NewEncoder(os.Stdout)

	for rows.Next() {
		err = rows.Scan(values...)
		if err != nil {
			log.Fatal(err)
		}

		dest := make(map[string]interface{})

		for i, column := range columns {
			dest[column] = *(values[i].(*interface{}))
		}

		err := enc.Encode(dest)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
