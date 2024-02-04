package flags

import (
	"fmt"
	"regexp"

	"github.com/spf13/pflag"
)

var _ pflag.Value = (*regexpValue)(nil)

func newRegexp(s string) (*regexpValue, error) {
	r := &regexpValue{Regexp: nil}

	if err := r.Set(s); err != nil {
		return nil, err
	}

	return r, nil
}

type regexpValue struct {
	*regexp.Regexp
}

// Set is a method to set the regexp value.
func (r *regexpValue) Set(s string) error {
	compiled, err := regexp.Compile(s)
	if err != nil {
		return fmt.Errorf("compiling regexp %q: %w", s, err)
	}

	r.Regexp = compiled

	return nil
}

// Type required to implement pflag.Value.
func (*regexpValue) Type() string {
	return "regexp"
}
