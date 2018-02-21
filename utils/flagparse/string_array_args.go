package flagparse

// ArrayIntFlags command line flag for []int.
type ArrayIntFlags []int

func (self *ArrayIntFlags) String() string {
	return "my string representation"
}

func (self *ArrayIntFlags) Set(value int) error {
	*self = append(*self, value)
	return nil
}
