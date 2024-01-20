package generator

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"text/template"

	"golang.org/x/exp/slog"

	"github.com/Djarvur/go-silly-enum/internal/extractor"
	"github.com/Djarvur/go-silly-enum/internal/parser"
)

//go:embed codegen.go.tmpl
var fileContentTmplSrc string

const fileNameTmplSrc = `{{.Dir}}/enum_silly_codegen_{{.Enum}}{{if .Test}}_test{{end}}.go`

var (
	fileContentTmpl = template.Must(template.New("fileContent").Parse(fileContentTmplSrc))
	fileNameTmpl    = template.Must(template.New("fileName").Parse(fileNameTmplSrc))
)

func Generate(pkgs, tags []string, includeTests bool, log *slog.Logger) ([]string, error) {
	prog, rest, err := parser.Parse(pkgs, tags, includeTests)
	if err != nil {
		return rest, fmt.Errorf("parsing sources: %w", err)
	}

	for enumDef, values := range extractor.Extract(prog) {
		if err = writeFile(enumDef, values); err != nil {
			return rest, fmt.Errorf("generating: %w", err)
		}

		log.Debug("Generate", "enum", enumDef, "values", values)
	}

	return rest, nil
}

func buildFileName(data extractor.EnumDef) (string, error) {
	var b bytes.Buffer

	if err := fileNameTmpl.Execute(&b, data); err != nil {
		return "", fmt.Errorf("%+v: %w", data, err)
	}

	return b.String(), nil
}

func writeFile(enumDef extractor.EnumDef, values []string) error {
	fileName, err := buildFileName(enumDef)
	if err != nil {
		return fmt.Errorf("building file name: %w", err)
	}

	fileNameTmp := fileName + ".tmp"

	file, err := os.Create(fileName + ".tmp")
	if err != nil {
		return fmt.Errorf("opening file %q: %w", fileNameTmp, err)
	}

	defer file.Close()

	type renderData struct {
		extractor.EnumDef
		Values []string
	}

	if err = fileContentTmpl.Execute(file, renderData{EnumDef: enumDef, Values: values}); err != nil {
		return fmt.Errorf("writing file %q: %w", fileNameTmp, err)
	}

	if err = os.Rename(fileNameTmp, fileName); err != nil {
		return fmt.Errorf("renaming file %q to %q: %w", fileNameTmp, fileName, err)
	}

	return nil
}
