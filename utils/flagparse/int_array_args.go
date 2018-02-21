package flagparse

// ArrayStringFlags command line flag for []string.
type ArrayStringFlags []string

func (self *ArrayStringFlags) String() string {
	return "my string representation"
}

func (self *ArrayStringFlags) Set(value string) error {
	*self = append(*self, value)
	return nil
}
