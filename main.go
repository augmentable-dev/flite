package main

import (
	"github.com/augmentable-dev/jqlite/internal/lines"
	"github.com/augmentable-dev/jqlite/internal/readfile"
	"go.riyazali.net/sqlite"
)

func init() {
	sqlite.Register(func(api *sqlite.ExtensionApi) (sqlite.ErrorCode, error) {
		if err := api.CreateModule("lines", lines.NewVTab(),
			sqlite.EponymousOnly(true), sqlite.ReadOnly(true)); err != nil {
			return sqlite.SQLITE_ERROR, err
		}

		if err := api.CreateFunction("readfile", readfile.NewReadFile()); err != nil {
			return sqlite.SQLITE_ERROR, err
		}

		return sqlite.SQLITE_OK, nil
	})
}

func main() {}
