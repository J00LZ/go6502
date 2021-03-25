package cpu

import (
	"github.com/J00LZZ/go6502/pkg/bus"
	"github.com/J00LZZ/go6502/pkg/instruction"
	"log"
	"sync"
	"time"
)

const IrqVectorH uint16 = 0xFFFF

const IrqVectorL uint16 = 0xFFFE

const RstVectorH uint16 = 0xFFFD

const RstVectorL uint16 = 0xFFFC

const NmiVectorH uint16 = 0xFFFB

const NmiVectorL uint16 = 0xFFFA

type CPU struct {
	PC    uint16
	AC    byte
	X     byte
	Y     byte
	SR    StatusFlags
	SP    byte
	Pause bool
	*bus.Bus

	// TODO: locks may be slow, could be solved by having interrupts set an atomic variable
	// which is checked before executing each instruction
	sync.Mutex
}

func New(bus *bus.Bus) *CPU {
	return &CPU{
		PC:    0,
		AC:    0,
		X:     0,
		Y:     0,
		SR:    0,
		SP:    0,
		Pause: false,
		Bus:   bus,
		Mutex: sync.Mutex{},
	}
}

func (c *CPU) pushStack(d byte) {
	c.Bus.WriteAddress(0x0100+uint16(c.SP), d)

	if c.SP == 0x00 {
		c.SP = 0xFF
	} else {
		c.SP--
	}
}

func (c *CPU) pushPC() {
	c.pushStack(byte((c.PC >> 8) & 0xFF))
	c.pushStack(byte(c.PC & 0xFF))
}

func (c *CPU) popStack() byte {
	if c.SP == 0xFF {
		c.SP = 0x00
	} else {
		c.SP++
	}
	return c.Bus.ReadAddress(0x0100 + uint16(c.SP))
}

func (c *CPU) IRQ() {
	if c.Has(I) {
		c.Clear(B)
		c.pushPC()
		c.pushStack(byte(c.SR))
		c.Set(I)
		c.PC = uint16(c.ReadAddress(IrqVectorH))<<8 + uint16(c.ReadAddress(IrqVectorL))
	}
}

func (c *CPU) NMI() {
	c.Lock()
	defer c.Unlock()

	c.Clear(B)
	c.pushPC()
	c.pushStack(byte(c.SR))
	c.Set(I)
	c.PC = uint16(c.ReadAddress(NmiVectorH))<<8 + uint16(c.ReadAddress(NmiVectorL))

}

func (c *CPU) Reset() {
	c.AC = 0
	c.Y = 0
	c.X = 0

	c.PC = uint16(c.ReadAddress(RstVectorH))<<8 + uint16(c.ReadAddress(RstVectorL))
	c.SP = 0xFD
	c.SR |= Unused
}

func (c *CPU) Run(clock <-chan time.Time) {
	var op byte
	var instr *instruction.Instruction
	for {
		for c.Pause {
		}
		op = c.readAddr()
		c.PC++
		instr = instruction.FetchInstruction(op)
		if instr == nil {
			log.Panicf("Instruction does not exist, %02X", op)
		}
		log.Printf("OP:%v Mode:%v", instr.Opcode, instr.Mode)
		extra := c.LoadInstruction(instr)
		tts := instr.Cycles + int(extra)
		log.Printf("A: %02X, X: %02X, Y: %02X", c.AC, c.X, c.Y)
		for tts > 0 {
			tts--
			<-clock
		}
	}
}

func (c *CPU) readAddr() byte {
	return c.ReadAddress(c.PC)
}

func (c *CPU) LoadInstruction(instr *instruction.Instruction) byte {
	addr, extra := c.dataAddr(instr.Mode)
	extra += c.execute(instr.Opcode, addr)
	return extra
}

func (c *CPU) dataAddr(mode instruction.Mode) (uint16, byte) {
	switch mode {
	case instruction.Acc:
		return 0, 0
	case instruction.Immediate:
		d := c.PC
		c.PC++
		return d, 0
	case instruction.Abs:
		low := uint16(c.readAddr())
		c.PC++
		high := uint16(c.readAddr())
		c.PC++
		addr := low | (high << 8)
		return addr, 0
	case instruction.Zpg:
		d := uint16(c.readAddr())
		c.PC++
		return d, 0
	case instruction.Impl:
		return 0, 0
	case instruction.Rel:
		offset := c.readAddr()
		c.PC++
		off2 := int8(offset)
		if off2 < 0 {
			off2 *= -1
			addr := c.PC - uint16(off2)
			if (c.PC >> 8) != (addr >> 8) {
				return addr, 1
			}
			return addr, 0
		} else {
			addr := c.PC + uint16(off2)
			if (c.PC >> 8) != (addr >> 8) {
				return addr, 1
			}
			return addr, 0
		}
	case instruction.Ind:
		low := uint16(c.readAddr())
		c.PC++
		high := uint16(c.readAddr())
		c.PC++
		abs := (high << 8) + low
		effL := uint16(c.Bus.ReadAddress(abs))
		effH := uint16(c.Bus.ReadAddress((abs & 0xFF00) + ((abs + 1) & 0x00FF)))
		return effL + (effH << 8), 0
	case instruction.ZpgX:
		z := uint16(c.readAddr())
		c.PC++
		z += uint16(c.X)
		return z & 0xFF, 0
	case instruction.ZpgY:
		z := uint16(c.readAddr())
		c.PC++
		z += uint16(c.Y)
		return z & 0xFF, 0
	case instruction.AbsX:
		low := uint16(c.readAddr())
		c.PC++
		high := uint16(c.readAddr())
		c.PC++
		abs := (high << 8) + low + uint16(c.X)
		if high != (abs >> 8) {
			return abs, 1
		}
		return abs, 0
	case instruction.AbsY:
		low := uint16(c.readAddr())
		c.PC++
		high := uint16(c.readAddr())
		c.PC++
		abs := (high << 8) + low + uint16(c.Y)
		if high != (abs >> 8) {
			return abs, 1
		}
		return abs, 0
	case instruction.Xind:
		low := (uint16(c.readAddr()) + uint16(c.X)) & 0xFF
		high := (low + 1) & 0xFF
		c.PC++
		addr := uint16(c.ReadAddress(low)) + (uint16(c.ReadAddress(high)) << 8)
		return addr, 0
	case instruction.IndY:
		low := uint16(c.readAddr())
		high := (low + 1) & 0xFF
		c.PC++
		addr := uint16(c.ReadAddress(low)) + (uint16(c.ReadAddress(high)) << 8) + uint16(c.Y)
		return addr, 0
	}

	return 0, 0
}

func (c *CPU) execute(opcode instruction.Opcode, data uint16) byte {
	// Lock the cpu lock to make sure only *one* thread can execute code on the cpu
	// at a time. This to make sure interrupts will never interrupt while in the middle
	// of executing an instruction.
	c.Lock()
	defer c.Unlock()

	switch opcode {
	case instruction.ADC:
		m := c.ReadAddress(data)
		res := uint16(m) + uint16(c.AC) + uint16(c.IfCarry())
		c.SetZero16(res)
		if c.Has(D) {
			if (c.AC&0xF)+(m&0xF)+c.IfCarry() > 9 {
				res += 6
			}
			c.SetNegative16(res & 0x80)
			if (c.AC^m)&0x80 == 0 && ((uint16(c.AC)^res)&0x80) != 0 {
				c.Set(V)
			} else {
				c.Clear(V)
			}
			if res > 0x99 {
				res += 96
			}
			c.SetCarry(res > 0x99)
		} else {
			c.SetNegative16(res & 0x80)
			if (c.AC^m)&0x80 == 0 && ((uint16(c.AC)^res)&0x80) != 0 {
				c.Set(V)
			} else {
				c.Clear(V)
			}
			c.SetCarry(res > 0xFF)
		}
		c.AC = byte(res & 0xFF)
	case instruction.AND:
		d := c.ReadAddress(data)
		res := d & c.AC
		c.SetNegative8(res)
		c.SetZero8(res)
		c.AC = res
	case instruction.ASL:
		m := c.ReadAddress(data)
		c.SetCarry(m&0x80 != 0)
		m <<= 1
		m &= 0xFF
		c.SetNegative8(m & 0x80)
		c.SetZero8(m)
		c.Bus.WriteAddress(data, m)
	case instruction.ASL_ACC:
		m := c.AC
		c.SetCarry(m&0x80 != 0)
		m <<= 1
		m &= 0xFF
		c.SetNegative8(m & 0x80)
		c.SetZero8(m)
		c.AC = m
	case instruction.BCC:
		if !c.Has(C) {
			x := byte(0)
			if (data >> 8) == c.PC>>8 {
				x = 1
			} else {
				x = 2
			}
			c.PC = data
			return x
		}
	case instruction.BCS:
		if c.Has(C) {
			x := byte(0)
			if (data >> 8) == c.PC>>8 {
				x = 1
			} else {
				x = 2
			}
			c.PC = data
			return x
		}
	case instruction.BEQ:
		if c.Has(Z) {
			x := byte(0)
			if (data >> 8) == c.PC>>8 {
				x = 1
			} else {
				x = 2
			}
			c.PC = data
			return x
		}
	case instruction.BIT:
		m := c.ReadAddress(data)
		res := m & c.AC
		c.SetNegative8(res & 0x80)
		c.SR = (c.SR & 0x3F) | StatusFlags(m&0xC0)
		c.SetZero8(res)
	case instruction.BMI:
		if c.Has(N) {
			x := byte(0)
			if (data >> 8) == c.PC>>8 {
				x = 1
			} else {
				x = 2
			}
			c.PC = data
			return x
		}
	case instruction.BNE:
		if !c.Has(Z) {
			x := byte(0)
			if (data >> 8) == c.PC>>8 {
				x = 1
			} else {
				x = 2
			}
			c.PC = data
			return x
		}
	case instruction.BPL:
		if !c.Has(N) {
			x := byte(0)
			if (data >> 8) == c.PC>>8 {
				x = 1
			} else {
				x = 2
			}
			c.PC = data
			return x
		}
	case instruction.BRK:
		c.PC++
		c.pushPC()
		c.pushStack(byte(c.SR | B))
		c.PC = uint16(c.ReadAddress(IrqVectorH))<<8 + uint16(c.ReadAddress(IrqVectorL))
	case instruction.BVC:
		if !c.Has(V) {
			x := byte(0)
			if (data >> 8) == c.PC>>8 {
				x = 1
			} else {
				x = 2
			}
			c.PC = data
			return x
		}
	case instruction.BVS:
		if c.Has(V) {
			x := byte(0)
			if (data >> 8) == c.PC>>8 {
				x = 1
			} else {
				x = 2
			}
			c.PC = data
			return x
		}
	case instruction.CLC:
		c.Clear(C)
	case instruction.CLD:
		c.Clear(D)
	case instruction.CLI:
		c.Clear(I)
	case instruction.CLV:
		c.Clear(V)
	case instruction.CMP:
		tmp := uint16(c.AC) - uint16(c.ReadAddress(data))
		c.SetCarry(tmp < 0x100)
		c.SetNegative16(tmp & 0x80)
		c.SetZero16(tmp)
	case instruction.CPX:
		tmp := uint16(c.X) - uint16(c.ReadAddress(data))
		c.SetCarry(tmp < 0x100)
		c.SetNegative16(tmp & 0x80)
		c.SetZero16(tmp)
	case instruction.CPY:
		tmp := uint16(c.Y) - uint16(c.ReadAddress(data))
		c.SetCarry(tmp < 0x100)
		c.SetNegative16(tmp & 0x80)
		c.SetZero16(tmp)
	case instruction.DEC:
		tmp := c.ReadAddress(data)
		tmp = (tmp - 1) & 0xFF
		c.SetNegative8(tmp)
		c.SetZero8(tmp)
		c.WriteAddress(data, tmp)
	case instruction.DEX:
		tmp := c.X
		tmp = (tmp - 1) & 0xFF
		c.SetNegative8(tmp)
		c.SetZero8(tmp)
		c.X = tmp
	case instruction.DEY:
		tmp := c.Y
		tmp = (tmp - 1) & 0xFF
		c.SetNegative8(tmp)
		c.SetZero8(tmp)
		c.Y = tmp
	case instruction.EOR:
		m := c.ReadAddress(data)
		m = c.AC ^ m
		c.SetNegative8(m)
		c.SetZero8(m)
		c.AC = m
	case instruction.INC:
		tmp := c.ReadAddress(data)
		tmp = (tmp + 1) & 0xFF
		c.SetNegative8(tmp)
		c.SetZero8(tmp)
		c.WriteAddress(data, tmp)
	case instruction.INX:
		tmp := c.X
		tmp = (tmp + 1) & 0xFF
		c.SetNegative8(tmp)
		c.SetZero8(tmp)
		c.X = tmp
	case instruction.INY:
		tmp := c.Y
		tmp = (tmp + 1) & 0xFF
		c.SetNegative8(tmp)
		c.SetZero8(tmp)
		c.Y = tmp
	case instruction.JMP:
		c.PC = data
	case instruction.JSR:
		c.PC--
		c.pushPC()
		c.PC = data
	case instruction.LDA:
		m := c.ReadAddress(data)
		c.SetNegative8(m)
		c.SetZero8(m)
		c.AC = m
	case instruction.LDX:
		m := c.ReadAddress(data)
		c.SetNegative8(m)
		c.SetZero8(m)
		c.X = m
	case instruction.LDY:
		m := c.ReadAddress(data)
		c.SetNegative8(m)
		c.SetZero8(m)
		c.Y = m
	case instruction.LSR:
		m := c.ReadAddress(data)
		c.SetCarry(m&0x01 != 0)
		m >>= 1
		c.SetNegative8(0)
		c.SetZero8(m)
		c.WriteAddress(data, m)
	case instruction.LSR_ACC:
		m := c.AC
		c.SetCarry(m&0x01 != 0)
		m >>= 1
		c.SetNegative8(0)
		c.SetZero8(m)
		c.AC = m
	case instruction.NOP:
		return 0
	case instruction.ORA:
		d := c.ReadAddress(data)
		res := d | c.AC
		c.SetNegative8(res)
		c.SetZero8(res)
		c.AC = res
	case instruction.PHA:
		c.pushStack(c.AC)
	case instruction.PHP:
		c.pushStack(byte(c.SR | B))
	case instruction.PLA:
		c.AC = c.popStack()
		c.SetNegative8(c.AC)
		c.SetZero8(c.AC)
	case instruction.PLP:
		c.SR = StatusFlags(c.popStack())
		c.Set(Unused)
	case instruction.ROL:
		m := uint16(c.ReadAddress(data))
		m <<= 1
		if c.Has(C) {
			m |= 0x01
		}
		c.SetCarry(m > 0xFF)
		m &= 0xFF
		c.SetNegative16(m & 0x80)
		c.SetZero16(m)
		c.WriteAddress(data, byte(m))
	case instruction.ROL_ACC:
		m := uint16(c.AC)
		m <<= 1
		if c.Has(C) {
			m |= 0x01
		}
		c.SetCarry(m > 0xFF)
		m &= 0xFF
		c.SetNegative16(m & 0x80)
		c.SetZero16(m)
		c.AC = byte(m)
	case instruction.ROR:
		m := uint16(c.ReadAddress(data))
		if c.Has(C) {
			m |= 0x0100
		}
		c.SetCarry(m&0x01 != 0)
		m >>= 1
		m &= 0xFF
		c.SetNegative16(m & 0x80)
		c.SetZero16(m)
		c.WriteAddress(data, byte(m))
	case instruction.ROR_ACC:
		m := uint16(c.AC)
		if c.Has(C) {
			m |= 0x0100
		}
		c.SetCarry(m&0x01 != 0)
		m >>= 1
		m &= 0xFF
		c.SetNegative16(m & 0x80)
		c.SetZero16(m)
		c.AC = byte(m)
	case instruction.RTI:
		c.SR = StatusFlags(c.popStack())
		low := uint16(c.popStack())
		high := uint16(c.popStack())
		c.PC = low | (high << 8)
	case instruction.RTS:
		low := uint16(c.popStack())
		high := uint16(c.popStack())
		c.PC = (low | (high << 8)) + 1
	case instruction.SBC:
		m := uint16(c.ReadAddress(data))
		tmp := uint16(c.AC) - m - uint16(1-c.IfCarry())
		c.SetNegative16(tmp & 0x80)
		c.SetZero16(tmp)
		if (uint16(c.AC)^tmp)&0x80 == 0 && ((uint16(c.AC)^m)&0x80) != 0 {
			c.Set(V)
		} else {
			c.Clear(V)
		}
		if c.Has(D) {
			if uint16((c.AC&0x0F)-(1-c.IfCarry())) < m&0x0F {
				tmp -= 6
			}
			if tmp > 0x99 {
				tmp -= 0x60
			}
		}
		c.SetCarry(tmp < 0x100)
		c.AC = byte(tmp & 0xFF)
	case instruction.SEC:
		c.Set(C)
	case instruction.SED:
		c.Set(D)
	case instruction.SEI:
		c.Set(I)
	case instruction.STA:
		c.WriteAddress(data, c.AC)
	case instruction.STX:
		c.WriteAddress(data, c.X)
	case instruction.STY:
		c.WriteAddress(data, c.Y)
	case instruction.TAX:
		m := c.AC
		c.SetNegative8(m)
		c.SetZero8(m)
		c.X = m
	case instruction.TAY:
		m := c.AC
		c.SetNegative8(m)
		c.SetZero8(m)
		c.Y = m
	case instruction.TSX:
		m := c.SP
		c.SetNegative8(m)
		c.SetZero8(m)
		c.X = m
	case instruction.TXA:
		m := c.X
		c.SetNegative8(m)
		c.SetZero8(m)
		c.AC = m
	case instruction.TXS:
		c.SP = c.X
	case instruction.TYA:
		m := c.Y
		c.SetNegative8(m)
		c.SetZero8(m)
		c.AC = m
	}

	return 0
}
