package cpu

func (c *CPU) SetNegative16(b uint16) {
	if b != 0 {
		c.Set(N)
	} else {
		c.Clear(N)
	}
}

func (c *CPU) SetNegative8(b uint8) {
	if b&0x80 != 0 {
		c.Set(N)
	} else {
		c.Clear(N)
	}
}

func (c *CPU) SetZero16(b uint16) {
	if b&0xFF == 0 {
		c.Set(Z)
	} else {
		c.Clear(Z)
	}
}

func (c *CPU) SetZero8(b uint8) {
	if b&0xFF == 0 {
		c.Set(Z)
	} else {
		c.Clear(Z)
	}
}

func (c *CPU) SetCarry(b bool) {
	if b {
		c.Set(C)
	} else {
		c.Clear(C)
	}
}
