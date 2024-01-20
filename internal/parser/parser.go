package parser

import (
	"go/build"

	"golang.org/x/tools/go/loader"
)

func Parse(pkgs, tags []string, includeTests bool) (*loader.Program, []string, error) {
	conf := loader.Config{
		Build: &build.Default,
	}
	conf.Build.BuildTags = append(conf.Build.BuildTags, tags...)

	rest, err := conf.FromArgs(pkgs, includeTests)
	if err != nil {
		return nil, rest, err
	}

	prog, err := conf.Load()

	return prog, rest, err
}
