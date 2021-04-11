package yaml

import (
	"go.riyazali.net/sqlite"
	"sigs.k8s.io/yaml"
)

type yamlToJSON struct{}

func (m *yamlToJSON) Args() int           { return 1 }
func (m *yamlToJSON) Deterministic() bool { return true }
func (m *yamlToJSON) Apply(ctx *sqlite.Context, values ...sqlite.Value) {
	input := values[0].Blob()

	asJSON, err := yaml.YAMLToJSON(input)
	if err != nil {
		ctx.ResultError(err)
	}

	ctx.ResultText(string(asJSON))
}

func NewYAMLToJSON() sqlite.Function {
	return &yamlToJSON{}
}

type jsonToYaml struct{}

func (m *jsonToYaml) Args() int           { return 1 }
func (m *jsonToYaml) Deterministic() bool { return true }
func (m *jsonToYaml) Apply(ctx *sqlite.Context, values ...sqlite.Value) {
	input := values[0].Blob()

	asYAML, err := yaml.JSONToYAML(input)
	if err != nil {
		ctx.ResultError(err)
	}

	ctx.ResultText(string(asYAML))
}

func NewJSONToYAML() sqlite.Function {
	return &jsonToYaml{}
}
