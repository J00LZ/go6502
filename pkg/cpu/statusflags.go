package cpu

type StatusFlags uint8

//go:generate stringer -type=StatusFlags
const (
	C StatusFlags = 1 << iota
	Z
	I
	D
	B
	Unused
	V
	N
)

func (c *CPU) Set(flag StatusFlags)      { c.SR |= flag }
func (c *CPU) Clear(flag StatusFlags)    { c.SR = c.SR &^ flag }
func (c *CPU) Toggle(flag StatusFlags)   { c.SR = ^ flag }
func (c *CPU) Has(flag StatusFlags) bool { return c.SR&flag != 0 }

func (c *CPU) IfCarry() byte {
	if c.Has(C) {
		return 1
	} else {
		return 0
	}
}
