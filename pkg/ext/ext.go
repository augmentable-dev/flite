package ext

import (
	"github.com/augmentable-dev/flite/internal/http"
	"github.com/augmentable-dev/flite/internal/lines"
	"github.com/augmentable-dev/flite/internal/readfile"
	"github.com/augmentable-dev/flite/internal/yaml"
	_ "github.com/mattn/go-sqlite3"
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

		if err := api.CreateFunction("http_get", http.NewHTTPGet()); err != nil {
			return sqlite.SQLITE_ERROR, err
		}

		return sqlite.SQLITE_OK, nil
	})
}
