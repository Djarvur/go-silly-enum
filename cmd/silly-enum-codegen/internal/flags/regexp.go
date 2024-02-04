package flags

import (
	"fmt"
	"regexp"

	"github.com/spf13/pflag"
)

var _ pflag.Value = (*Regexp)(nil)

func NewRegexp(s string) (*Regexp, error) {
	r := &Regexp{}

	if err := r.Set(s); err != nil {
		return nil, err
	}

	return r, nil
}

type Regexp struct {
	*regexp.Regexp
}

func (r *Regexp) Set(s string) error {
	compiled, err := regexp.Compile(s)
	if err != nil {
		return fmt.Errorf("compiling regexp %q: %w", s, err)
	}

	r.Regexp = compiled

	return nil
}

func (*Regexp) Type() string {
	return "regexp"
}
