package cpu

import (
	"fmt"
	"github.com/J00LZZ/go6502/pkg/instruction"
)

type Dis struct {
	Text     string
	Location uint16
}

func (c *CPU) DisassembleMore() []Dis {
	pc := uint16(c.ReadAddress(0xFFFC))
	pc = pc + uint16(c.ReadAddress(0xFFFd))<<8
	text := make([]Dis, 0)
	for i := 0; i < 0x8000; i++ {
		t, offset := c.DisassembleAt(pc)
		text = append(text, Dis{t, pc})
		pc += offset
	}
	return text
}

func (c *CPU) DisassembleAt(addr uint16) (string, uint16) {
	current := c.ReadAddress(addr)

	instr := instruction.FetchInstruction(current)
	if instr == nil {
		return "ERR", 1
	}
	s := instr.Size()
	if s == 1 {
		return fmt.Sprintf("%s", instr.Opcode), 1
	}
	param := make([]byte, 0, 2)
	for i := uint16(1); i < s; i++ {
		param = append(param, c.ReadAddress(addr+i))
	}
	switch instr.Mode {
	case instruction.Abs:
		return fmt.Sprintf("%s $%02X%02X", instr.Opcode, param[1], param[0]), s
	case instruction.AbsX:
		return fmt.Sprintf("%s $%02X%02X,X", instr.Opcode, param[1], param[0]), s
	case instruction.AbsY:
		return fmt.Sprintf("%s $%02X%02X,Y", instr.Opcode, param[1], param[0]), s
	case instruction.Immediate:
		return fmt.Sprintf("%s #$%02X", instr.Opcode, param[0]), s
	case instruction.Ind:
		return fmt.Sprintf("%s ($%02X%02X)", instr.Opcode, param[1], param[0]), s
	case instruction.Xind:
		return fmt.Sprintf("%s ($%02X,X)", instr.Opcode, param[0]), s
	case instruction.IndY:
		return fmt.Sprintf("%s ($%02X),Y", instr.Opcode, param[0]), s
	case instruction.Rel:
		offset := param[0]
		off2 := int8(offset)
		return fmt.Sprintf("%s $%02X", instr.Opcode, off2), s
	case instruction.Zpg:
		return fmt.Sprintf("%s $%02X", instr.Opcode, param[0]), s
	case instruction.ZpgX:
		return fmt.Sprintf("%s $%02X,X", instr.Opcode, param[0]), s
	case instruction.ZpgY:
		return fmt.Sprintf("%s $%02X,Y", instr.Opcode, param[0]), s
	}
	return fmt.Sprintf("%v", instr), s
}

func (c *CPU) DisassembleCurrent() Dis {
	t, _ := c.DisassembleAt(c.PC)
	return Dis{t, c.PC}
}
