package extractor

import (
	"go/ast"
	"go/token"
	"path"
	"strings"

	"golang.org/x/exp/slog"
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
						Enum:     v.enum,
						BaseType: v.base,
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

type enumValue struct {
	name string
	enum string
	base string
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

	for _, rawSpec := range decl.Specs {
		if v, parsed := parseSpec(rawSpec); parsed {
			res = append(res, v)
		}
	}

	return res
}

func parseSpec(raw ast.Spec) (enumValue, bool) {
	spec, ok := raw.(*ast.ValueSpec)
	if !ok || spec.Type == nil || len(spec.Names) < 1 {
		return enumValue{}, false
	}

	specType, ok := spec.Type.(*ast.Ident)
	if !ok || !strings.HasSuffix(specType.Name, EnumSuffix) {
		return enumValue{}, false
	}

	return enumValue{
		name: spec.Names[0].Name,
		enum: specType.Name,
		base: digTypeName(spec.Type.(*ast.Ident)),
	}, true
}

func digTypeName(decl *ast.Ident) string {
	for decl.Obj != nil {
		slog.Info("digTypeName", "decl", decl)
		decl = decl.Obj.Decl.(*ast.TypeSpec).Type.(*ast.Ident)
	}

	slog.Info("digTypeName", "decl", decl)

	return decl.Name
}
