package inst

//go:generate stringer -type=Kind

type Kind uint

const (
	Trivial  Kind = iota // Trivial instructions (e.g 1 + 2)
	Indirect             // Pointer indirection
	Swap                 // Swap
	Compare              // Compare

	NumKinds
)

type Counter struct {
	state State
}

func (c *Counter) Do(kind Kind) {
	c.state[kind]++
}

func (c *Counter) State() State {
	return c.state
}

type State [NumKinds]uint

// TODO: State -> map[string]uint ?
