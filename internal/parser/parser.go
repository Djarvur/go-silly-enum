// Package parser is a wrapper for standard packages tool
package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"

	"golang.org/x/exp/slog"
	"golang.org/x/tools/go/packages"
)

// Parse the packages.
func Parse(
	dirs []string,
	tags []string,
	env []string,
	includeTests bool,
	log *slog.Logger,
) ([]*packages.Package, *token.FileSet, error) {
	var (
		allPkgs = make([]*packages.Package, 0, len(dirs))
		fset    = token.NewFileSet()
	)

	for _, dir := range dirs {
		cfg := &packages.Config{ //nolint:exhaustruct
			Fset:       fset,
			Mode:       packages.NeedName | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo,
			Dir:        dir,
			BuildFlags: tags,
			Env:        append(os.Environ(), env...),
			Tests:      includeTests,
			ParseFile: func(fset *token.FileSet, filename string, src []byte) (*ast.File, error) {
				return parseFile(fset, filename, src, log)
			},
		}

		pkgs, err := packages.Load(cfg, "./...")
		if err != nil {
			return nil, nil, fmt.Errorf("loading sources: %w", err)
		}

		allPkgs = append(allPkgs, pkgs...)
	}

	return allPkgs, fset, nil
}

func parseFile(
	fset *token.FileSet,
	filename string,
	src []byte,
	log *slog.Logger,
) (*ast.File, error) {
	log.Debug("parsing", "file", filename)

	return parser.ParseFile(fset, filename, src, parser.AllErrors|parser.ParseComments) //nolint:wrapcheck
}
