package extractor

import (
	"go/ast"
	"go/token"
	"path"
	"strings"

	"golang.org/x/tools/go/loader"
)

const (
	EnumSuffix = "Enum"
	TestSuffix = "_test.go"
)

type EnumDef struct {
	Enum     string
	BaseType string
	Package  string
	Dir      string
	Test     bool
}

func Extract(prog *loader.Program) map[EnumDef][]string {
	res := make(map[EnumDef][]string)

	for _, pkgInfo := range prog.InitialPackages() {
		for _, file := range pkgInfo.Files {
			var (
				fileName = prog.Fset.File(file.Pos()).Name()
				dirName  = path.Dir(fileName)
				isTest   = strings.HasSuffix(fileName, TestSuffix)
			)

			for _, decl := range file.Decls {
				for _, v := range extractEnumVals(decl) {
					enumDef := EnumDef{
						Enum:     v.typeName,
						BaseType: v.baseType,
						Package:  pkgInfo.Pkg.Name(),
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

	res := make([]enumValue, 0, 8)

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
		return enumValue{}, false
	}

	specType := lastType

	switch {
	case spec.Type != nil:
		specType = extractEnumType(spec.Type)
	case len(spec.Values) != 0:
		return enumValue{}, false
	}

	if specType.typeName == "" {
		return enumValue{}, false
	}

	return enumValue{
		name:     spec.Names[0].Name,
		enumType: specType,
	}, true
}

func extractEnumType(expr ast.Expr) enumType {
	typeIdent, ok := expr.(*ast.Ident)
	if !ok || !strings.HasSuffix(typeIdent.Name, EnumSuffix) {
		return enumType{}
	}

	return enumType{
		typeName: typeIdent.Name,
		baseType: digTypeName(typeIdent),
	}
}

func digTypeName(decl *ast.Ident) string {
	for decl.Obj != nil {
		decl = decl.Obj.Decl.(*ast.TypeSpec).Type.(*ast.Ident)
	}

	return decl.Name
}
