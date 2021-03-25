package cpu

import (
	"fmt"
	"github.com/J00LZZ/go6502/pkg/instruction"
)

func (c *CPU) DisassembleCurrent() (string, uint16) {
	current := c.ReadAddress(c.PC)
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
		param = append(param, c.ReadAddress(c.PC+i))
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
	case instruction.Rel, instruction.Zpg:
		return fmt.Sprintf("%s $%02X", instr.Opcode, param[0]), s
	case instruction.ZpgX:
		return fmt.Sprintf("%s $%02X,X", instr.Opcode, param[0]), s
	case instruction.ZpgY:
		return fmt.Sprintf("%s $%02X,Y", instr.Opcode, param[0]), s
	}
	return fmt.Sprintf("%v", instr), s
}
