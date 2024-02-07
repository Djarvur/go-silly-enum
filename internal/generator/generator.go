// Package generator generates the files
package generator

import (
	_ "embed"
	"fmt"
	"os"
	"path"
	"text/template"

	"golang.org/x/exp/slog"

	"github.com/Djarvur/go-silly-enum/internal/extractor"
	"github.com/Djarvur/go-silly-enum/internal/parser"
)

const fileNameBase = "enum_silly_codegen"

//nolint:gochecknoglobals
var (
	//go:embed enum_silly_codegen.go.tmpl
	fileContentTmplSrc string

	fileContentTmpl = template.Must(template.New("fileContent").Parse(fileContentTmplSrc))
)

// Generate generates the files, calling parser and extractor.
func Generate(
	dirs []string,
	tags []string,
	env []string,
	includeTests bool,
	enumName extractor.StringMatcher,
	log *slog.Logger,
) error {
	pkgs, fset, err := parser.Parse(dirs, tags, env, includeTests, log)
	if err != nil {
		return fmt.Errorf("parsing sources: %w", err)
	}

	for pkg, enums := range extractor.Extract(pkgs, fset, enumName) {
		if err = writeFile(pkg, enums); err != nil {
			return fmt.Errorf("generating: %w", err)
		}

		log.Debug("Generate", "enum", pkg, "enums", enums)
	}

	return nil
}

func buildFileName(pkg extractor.Package) string {
	if pkg.IsTest {
		return path.Join(pkg.Dir, fileNameBase+"_test.go")
	}

	return path.Join(pkg.Dir, fileNameBase+".go")
}

func writeFile(pkg extractor.Package, enums []extractor.Enum) error {
	var (
		fileName    = buildFileName(pkg)
		fileNameTmp = fileName + ".tmp"
	)

	file, err := os.Create(fileNameTmp) //nolint:gosec
	if err != nil {
		return fmt.Errorf("opening file %q: %w", fileNameTmp, err)
	}

	defer file.Close() //nolint:errcheck

	type renderData struct {
		Package extractor.Package
		Enums   []extractor.Enum
	}

	if err = fileContentTmpl.Execute(file, renderData{Package: pkg, Enums: enums}); err != nil {
		return fmt.Errorf("writing file %q: %w", fileNameTmp, err)
	}

	if err = os.Rename(fileNameTmp, fileName); err != nil {
		return fmt.Errorf("renaming file %q to %q: %w", fileNameTmp, fileName, err)
	}

	return nil
}
