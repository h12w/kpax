package cluster

import (
	"strings"
)

type MultiError []error

func (es MultiError) Error() string {
	ss := make([]string, len(es))
	for i := range ss {
		ss[i] = es[i].Error()
	}
	return strings.Join(ss, ", ")
}
