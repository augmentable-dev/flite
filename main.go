package main

import (
	"fmt"
	"os"

	"github.com/augmentable-dev/flite/internal/lines"
	"github.com/augmentable-dev/flite/internal/readfile"
	"github.com/augmentable-dev/flite/internal/yaml"
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

		if err := api.CreateFunction("yaml_to_json", yaml.NewYAMLToJSON()); err != nil {
			return sqlite.SQLITE_ERROR, err
		}

		if err := api.CreateFunction("json_to_yaml", yaml.NewJSONToYAML()); err != nil {
			return sqlite.SQLITE_ERROR, err
		}

		return sqlite.SQLITE_OK, nil
	})
}

func main() {
	var query string
	if len(os.Args) > 0 {
		query = os.Args[0]
	}

	fmt.Println(query)
}
