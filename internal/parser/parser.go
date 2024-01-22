// Package parser is a wrapper for standard packages tool
package parser

import (
	"fmt"
	"go/token"

	"golang.org/x/tools/go/packages"
)

// Parse the packages.
func Parse(dir string, tags []string, includeTests bool) ([]*packages.Package, *token.FileSet, error) {
	cfg := &packages.Config{ //nolint:exhaustruct
		Fset:       token.NewFileSet(),
		Mode:       packages.NeedName | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo,
		Dir:        dir,
		BuildFlags: tags,
		Tests:      includeTests,
	}

	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		return nil, nil, fmt.Errorf("loading sources: %w", err)
	}

	return pkgs, cfg.Fset, nil
}
