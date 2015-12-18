package cluster

import (
	"fmt"
)

type MultiError []error

func (e MultiError) Error() string {
	return fmt.Sprint(e)
}
