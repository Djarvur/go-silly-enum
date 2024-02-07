// Package extractor enumerating the files and extracting the looks-like-enum lines
package extractor

import (
	"go/ast"
	"go/token"
	"go/types"
	"path"
	"sort"
	"strings"

	"golang.org/x/exp/slices"
	"golang.org/x/tools/go/packages"
)

const (
	testSuffix = "_test.go"
)

// StringMatcher is a regexp compatible matcher.
type StringMatcher interface {
	MatchString(s string) bool
}

// Package describes the package enum was found.
type Package struct {
	Name   string
	IsTest bool
	Dir    string
}

func newPackage(pkgName, fileName string) Package {
	return Package{
		Name:   pkgName,
		IsTest: strings.HasSuffix(fileName, testSuffix),
		Dir:    path.Dir(fileName),
	}
}

// Enum describes the enum type.
type Enum struct {
	Name     string
	BaseType string
	Values   []string
}

// ForEach is going through the packages calling extract function for each declaration found
// and combining the records extracted per package using the merge function.
func ForEach[R any](
	pkgs []*packages.Package,
	fset *token.FileSet,
	extract func(ast.Decl, map[ast.Expr]types.TypeAndValue) []R,
	merge func([]R, []R) []R,
) map[Package][]R {
	res := make(map[Package][]R)

	uniq := make(map[string]int)

	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			var (
				fileName = fset.File(file.Pos()).Name()
				pkgDef   = newPackage(pkg.Name, fileName)
			)

			if uniq[fileName]++; uniq[fileName] > 1 {
				continue
			}

			for _, decl := range file.Decls {
				res[pkgDef] = merge(res[pkgDef], extract(decl, pkg.TypesInfo.Types))
			}
		}
	}

	return res
}

// Extract the enum constants relative records.
func Extract(
	pkgs []*packages.Package,
	fset *token.FileSet,
	enumName StringMatcher,
) map[Package][]Enum {
	return ForEach(
		pkgs,
		fset,
		(&enumExtractor{StringMatcher: enumName}).extract,
		mergeEnums,
	)
}

func mergeEnums(a, b []Enum) []Enum {
	res := append(a, b...)

	if len(res) < 2 || len(b) == 0 {
		return res
	}

	sort.Slice(res, func(i, j int) bool { return res[i].Name < res[j].Name })

	j := 0
	for i := 1; i < len(res); i++ {
		if res[j].Name == res[i].Name {
			res[j].Values = append(res[j].Values, res[i].Values...)
			continue
		}
		j++
		res[j] = res[i]
	}

	for i, v := range res {
		sort.Strings(v.Values)
		v.Values = slices.Compact(v.Values)
		res[i] = v
	}

	return res[:j+1]
}

type enumExtractor struct {
	lastType ast.Expr
	StringMatcher
}

func (e *enumExtractor) extract(decl ast.Decl, typesInfo map[ast.Expr]types.TypeAndValue) []Enum {
	var (
		found = extractConstants(decl)
		res   = make([]Enum, 0, len(found))
	)

	for _, v := range found {
		if v.Type == nil {
			v.Type = e.lastType
		} else {
			e.lastType = v.Type
		}

		name, base, ok := findType(typesInfo, v.Type)
		if !ok || !e.MatchString(name) {
			continue
		}

		res = append(res, Enum{Name: name, BaseType: base, Values: []string{v.Name}})
	}

	return res
}

func findType(info map[ast.Expr]types.TypeAndValue, expr ast.Expr) (string, string, bool) {
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

type constValue struct {
	Name string
	Type ast.Expr
}

func extractConstants(raw ast.Decl) []constValue {
	decl, ok := raw.(*ast.GenDecl)
	if !ok {
		return nil
	}

	if decl.Tok != token.CONST {
		return nil
	}

	res := make([]constValue, 0, len(decl.Specs))

	for _, spec := range decl.Specs {
		if v, parsed := parseSpec(spec); parsed {
			res = append(res, v)
		}
	}

	return res
}

func parseSpec(raw ast.Spec) (constValue, bool) {
	spec, isValue := raw.(*ast.ValueSpec)
	if !isValue || len(spec.Names) != 1 {
		return constValue{}, false //nolint:exhaustruct
	}

	if spec.Type == nil && len(spec.Values) > 0 {
		return constValue{}, false //nolint:exhaustruct
	}

	return constValue{Name: spec.Names[0].Name, Type: spec.Type}, true
}
