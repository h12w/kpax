package consumer

type Config struct {
}

type C struct{}

func New() (*C, error) {
	return &C{}, nil
}

func (c *C)
