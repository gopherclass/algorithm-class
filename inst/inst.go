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

func NewCounter() *Counter {
	return new(Counter)
}

func (c *Counter) Do(kind Kind) bool {
	if c == nil {
		return true
	}
	c.state[kind]++
	return true
}

func (c *Counter) Use(kind Kind, n uint) bool {
	if c == nil {
		return true
	}
	c.state[kind] += n
	return true
}

func (c *Counter) State() State {
	if c == nil {
		return State{}
	}
	return c.state
}

type State [NumKinds]uint

func (state State) Get(kind Kind) uint {
	return state[kind]
}

// TODO: State -> map[string]uint ?
