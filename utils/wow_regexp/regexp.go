package wow_regexp

import "regexp"

type RegExp struct {
	*regexp.Regexp
}

// CreateRegExp create
func CreateRegExp(expr string) (*RegExp, error) {
	reg, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	re := RegExp{reg}
	return &re, nil
}
