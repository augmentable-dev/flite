package main

import (
	"github.com/augmentable-dev/jqlite/internal/ndjson"
	"go.riyazali.net/sqlite"
)

func init() {
	sqlite.Register(func(api *sqlite.ExtensionApi) (sqlite.ErrorCode, error) {
		if err := api.CreateModule("ndjson", ndjson.NewVTab(),
			sqlite.EponymousOnly(true), sqlite.ReadOnly(true)); err != nil {
			return sqlite.SQLITE_ERROR, err
		}
		return sqlite.SQLITE_OK, nil
	})
}

func main() {}
