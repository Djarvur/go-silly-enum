// Package extractor enumerating the files and extracting the looks-like-enum lines
package extractor

import (
	"go/ast"
	"go/token"
	"path"
	"strings"

	"golang.org/x/tools/go/packages"
)

const (
	enumSuffix = "Enum"
	testSuffix = "_test.go"
)

// EnumDef describes the enum const record with all the details.
type EnumDef struct {
	Enum     string
	BaseType string
	Package  string
	Dir      string
	Test     bool
}

// Extract the enum constants relative records.
func Extract(pkgs []*packages.Package, fset *token.FileSet) map[EnumDef][]string {
	res := make(map[EnumDef][]string)

	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			var (
				fileName = fset.File(file.Pos()).Name()
				dirName  = path.Dir(fileName)
				isTest   = strings.HasSuffix(fileName, testSuffix)
			)

			for _, decl := range file.Decls {
				for _, v := range extractEnumVals(decl) {
					enumDef := EnumDef{
						Enum:     v.typeName,
						BaseType: v.baseType,
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

type enumType struct {
	typeName string
	baseType string
}

type enumValue struct {
	name string
	enumType
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

	var lastType enumType

	for _, rawSpec := range decl.Specs {
		if v, parsed := parseSpec(rawSpec, lastType); parsed {
			lastType = v.enumType
			res = append(res, v)
		}
	}

	return res
}

func parseSpec(raw ast.Spec, lastType enumType) (enumValue, bool) {
	spec, isValue := raw.(*ast.ValueSpec)
	if !isValue || len(spec.Names) < 1 {
		return enumValue{}, false //nolint:exhaustruct
	}

	specType := lastType

	switch {
	case spec.Type != nil:
		specType = extractEnumType(spec.Type)
	case len(spec.Values) != 0:
		return enumValue{}, false //nolint:exhaustruct
	}

	if specType.typeName == "" {
		return enumValue{}, false //nolint:exhaustruct
	}

	return enumValue{
		name:     spec.Names[0].Name,
		enumType: specType,
	}, true
}

func extractEnumType(expr ast.Expr) enumType {
	typeIdent, ok := expr.(*ast.Ident)
	if !ok || !strings.HasSuffix(typeIdent.Name, enumSuffix) {
		return enumType{} //nolint:exhaustruct
	}

	return enumType{
		typeName: typeIdent.Name,
		baseType: digTypeName(typeIdent),
	}
}

func digTypeName(decl *ast.Ident) string {
	for decl.Obj != nil {
		decl = decl.Obj.Decl.(*ast.TypeSpec).Type.(*ast.Ident) //nolint:forcetypeassert
	}

	return decl.Name
}
