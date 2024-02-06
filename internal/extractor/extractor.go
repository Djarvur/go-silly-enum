// Package extractor enumerating the files and extracting the looks-like-enum lines
package extractor

import (
	"go/ast"
	"go/token"
	"go/types"
	"path"
	"strings"

	"golang.org/x/tools/go/packages"
)

const (
	testSuffix = "_test.go"
)

// StringMatcher is a regexp compatible matcher.
type StringMatcher interface {
	MatchString(s string) bool
}

// EnumDef describes the enum const record with all the details.
type EnumDef struct {
	Enum     string
	BaseType string
	Package  string
	Dir      string
	Test     bool
}

// Extract the enum constants relative records.
func Extract(
	pkgs []*packages.Package,
	fset *token.FileSet,
	enumName StringMatcher,
) map[EnumDef][]string {
	res := make(map[EnumDef][]string)

	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			var (
				fileName = fset.File(file.Pos()).Name()
				dirName  = path.Dir(fileName)
				isTest   = strings.HasSuffix(fileName, testSuffix)
			)

			for _, decl := range file.Decls {
				var typeLast ast.Expr

				for _, v := range extractEnumVals(decl) {
					if v.enumType == nil {
						v.enumType = typeLast
					}

					typeName, typeBase, ok := extractType(pkg.TypesInfo.Types, v.enumType)
					if !ok || !enumName.MatchString(typeName) {
						continue
					}

					enumDef := EnumDef{
						Enum:     typeName,
						BaseType: typeBase,
						Package:  pkg.Name,
						Dir:      dirName,
						Test:     isTest,
					}

					res[enumDef] = append(res[enumDef], v.name)
				}
			}
		}
	}

	return res
}

func extractType(info map[ast.Expr]types.TypeAndValue, expr ast.Expr) (string, string, bool) {
	decl, ok := info[expr]
	if !ok {
		return "", "", false
	}

	named, ok := decl.Type.(*types.Named)
	if !ok {
		return "", "", false
	}

	return named.Obj().Name(), decl.Type.Underlying().String(), true
}

type enumValue struct {
	name     string
	enumType ast.Expr
}

func extractEnumVals(raw ast.Decl) []enumValue {
	decl, ok := raw.(*ast.GenDecl)
	if !ok {
		return nil
	}

	if decl.Tok != token.CONST {
		return nil
	}

	res := make([]enumValue, 0, 8) //nolint:gomnd

	for _, rawSpec := range decl.Specs {
		if v, parsed := parseSpec(rawSpec); parsed {
			res = append(res, v)
		}
	}

	return res
}

func parseSpec(raw ast.Spec) (enumValue, bool) {
	spec, isValue := raw.(*ast.ValueSpec)
	if !isValue || len(spec.Names) < 1 {
		return enumValue{}, false //nolint:exhaustruct
	}

	if spec.Type == nil && len(spec.Values) > 0 {
		return enumValue{}, false //nolint:exhaustruct
	}

	return enumValue{name: spec.Names[0].Name, enumType: spec.Type}, true
}
