package bus

type Type byte

const (
	R Type = 1 << iota
	W
)

const RW = R | W

func (t Type) HasFlag(flag Type) bool {
	return t&flag != 0
}
