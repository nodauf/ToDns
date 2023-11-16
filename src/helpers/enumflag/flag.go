package enumflag

import "fmt"

type enumFlag struct {
	target  *string
	options []string
}

// New returns a flag.Value implementation for parsing flags with a one-of-a-set value
func New(target *string, defaultValue string, options ...string) *enumFlag {
	var flag = &enumFlag{target: target, options: options}
	flag.Set(defaultValue)
	return flag
}

func (f *enumFlag) String() string {
	return *f.target
}

func (f *enumFlag) Type() string {
	return "string"
}

func (f *enumFlag) Set(value string) error {
	for _, v := range f.options {
		if v == value {
			*f.target = value
			return nil
		}
	}

	return fmt.Errorf("expected one of %q", f.options)
}
