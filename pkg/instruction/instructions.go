package instruction

type Instruction struct {
	Opcode
	Mode
	Cycles int
}

func ins(instructionType Opcode, mode Mode, cycles int) *Instruction {
	return &Instruction{instructionType, mode, cycles}
}

func (i *Instruction) Size() uint16 {
	x := 1
	switch i.Mode {
	case Acc, Impl:
		x += 0
	case Immediate, Rel, Zpg, ZpgX, ZpgY, Xind, IndY:
		x += 1
	case Abs, AbsX, AbsY, Ind:
		x += 2
	}
	return uint16(x)
}

type Opcode int
type Mode int

// Instructions
//go:generate stringer -type=Opcode
const (
	ADC Opcode = iota
	AND
	ASL
	ASL_ACC
	BCC
	BCS
	BEQ
	BIT
	BMI
	BNE
	BPL
	BRK
	BVC
	BVS
	CLC
	CLD
	CLI
	CLV
	CMP
	CPX
	CPY
	DEC
	DEX
	DEY
	EOR
	INC
	INX
	INY
	JMP
	JSR
	LDA
	LDX
	LDY
	LSR
	LSR_ACC
	NOP
	ORA
	PHA
	PHP
	PLA
	PLP
	ROL
	ROL_ACC
	ROR
	ROR_ACC
	RTI
	RTS
	SBC
	SEC
	SED
	SEI
	STA
	STX
	STY
	TAX
	TAY
	TSX
	TXA
	TXS
	TYA
)

// Addressing modes
//go:generate stringer -type=Mode
const (
	Acc Mode = iota
	Abs
	AbsX
	AbsY
	Immediate
	Impl
	Ind
	Xind
	IndY
	Rel
	Zpg
	ZpgX
	ZpgY
)

func lowNibble(b byte) byte {
	return (b >> 4) & 0x0F
}
func highNibble(b byte) byte {
	return b & 0x0F
}

var z = [][]*Instruction{
	{
		ins(BRK, Impl, 7), ins(BPL, Rel, 2), ins(JSR, Abs, 6), ins(BMI, Rel, 2), ins(RTI, Impl, 6),
		ins(BVC, Rel, 2), ins(RTS, Impl, 6), ins(BVS, Rel, 2), nil, ins(BCC, Rel, 2), ins(LDY, Immediate, 2),
		ins(BCS, Rel, 2), ins(CPY, Immediate, 2), ins(BNE, Rel, 2), ins(CPX, Immediate, 2), ins(BEQ, Rel, 2),
	},
	{
		ins(ORA, Xind, 6), ins(ORA, IndY, 5), ins(AND, Xind, 6), ins(AND, IndY, 5),
		ins(EOR, Xind, 6), ins(EOR, IndY, 5), ins(ADC, Xind, 6), ins(ADC, IndY, 5),
		ins(STA, Xind, 6), ins(STA, IndY, 6), ins(LDA, Xind, 6), ins(LDA, IndY, 5),
		ins(CMP, Xind, 6), ins(CMP, IndY, 5), ins(SBC, Xind, 6), ins(SBC, IndY, 5),
	},
	{ //pain
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, ins(LDX, Immediate, 2), nil, nil, nil, nil, nil,
	},
	nil,
	{ //less pain but still pain
		nil, nil, ins(BIT, Zpg, 3), nil, nil, nil, nil, nil,
		ins(STY, Zpg, 3), ins(STY, ZpgX, 4), ins(LDY, Zpg, 3), ins(LDY, AbsX, 4),
		ins(CPY, Zpg, 3), nil, ins(CPX, Zpg, 3), nil,
	},
	{
		ins(ORA, Zpg, 3), ins(ORA, ZpgX, 4), ins(AND, Zpg, 3), ins(AND, ZpgX, 4),
		ins(EOR, Zpg, 3), ins(EOR, ZpgX, 4), ins(ADC, Zpg, 3), ins(ADC, ZpgX, 4),
		ins(STA, Zpg, 3), ins(STA, ZpgX, 4), ins(LDA, Zpg, 3), ins(LDA, ZpgX, 4),
		ins(CMP, Zpg, 3), ins(CMP, ZpgX, 4), ins(SBC, Zpg, 3), ins(SBC, ZpgX, 4),
	},
	{
		ins(ASL, Zpg, 5), ins(ASL, ZpgX, 6), ins(ROL, Zpg, 5), ins(ROL, ZpgX, 6),
		ins(LSR, Zpg, 5), ins(LSR, ZpgX, 6), ins(ROR, Zpg, 5), ins(ROR, ZpgX, 6),
		ins(STX, Zpg, 3), ins(STX, ZpgY, 4), ins(LDX, Zpg, 3), ins(LDX, ZpgY, 4),
		ins(DEC, Zpg, 5), ins(DEC, ZpgX, 6), ins(INC, Zpg, 5), ins(INC, ZpgX, 6),
	},
	nil,
	{
		ins(PHP, Impl, 3), ins(CLC, Impl, 2), ins(PLP, Impl, 4), ins(SEC, Impl, 2),
		ins(PHA, Impl, 3), ins(CLI, Impl, 2), ins(PLA, Impl, 4), ins(SEI, Impl, 2),
		ins(DEY, Impl, 2), ins(TYA, Impl, 2), ins(TAY, Impl, 2), ins(CLV, Impl, 2),
		ins(INY, Impl, 2), ins(CLD, Impl, 2), ins(INX, Impl, 2), ins(SED, Impl, 2),
	},
	{
		ins(ORA, Immediate, 2), ins(ORA, AbsY, 4), ins(AND, Immediate, 2), ins(AND, AbsY, 4),
		ins(EOR, Immediate, 2), ins(EOR, AbsY, 4), ins(ADC, Immediate, 2), ins(ADC, AbsY, 4),
		nil, ins(STA, AbsY, 5), ins(LDA, Immediate, 2), ins(LDA, AbsY, 4),
		ins(CMP, Immediate, 2), ins(CMP, AbsY, 4), ins(SBC, Immediate, 2), ins(SBC, AbsY, 4),
	},
	{
		ins(ASL_ACC, Acc, 2), nil, ins(ROL_ACC, Acc, 2), nil, ins(LSR_ACC, Acc, 2), nil, ins(ROR_ACC, Acc, 2), nil,
		ins(TXA, Impl, 2), ins(TXS, Impl, 2), ins(TAX, Impl, 2), ins(TSX, Impl, 2),
		ins(DEX, Impl, 2), nil, ins(NOP, Impl, 2), nil,
	},
	nil,
	{
		nil, nil, ins(BIT, Abs, 3), nil, ins(JMP, Abs, 3), nil, ins(JMP, Ind, 5), nil,
		ins(STY, Abs, 4), nil, ins(LDY, Abs, 4), ins(LDY, AbsX, 4), ins(CPY, Abs, 4),
		nil, ins(CPX, Abs, 4), nil,
	},
	{
		ins(ORA, Abs, 4), ins(ORA, AbsX, 4), ins(AND, Abs, 4), ins(AND, AbsX, 4),
		ins(EOR, Abs, 4), ins(EOR, AbsX, 4), ins(ADC, Abs, 4), ins(ADC, AbsX, 4),
		ins(STA, Abs, 4), ins(STA, AbsX, 5), ins(LDA, Abs, 4), ins(LDA, AbsX, 4),
		ins(CMP, Abs, 4), ins(CMP, AbsX, 4), ins(SBC, Abs, 4), ins(SBC, AbsX, 4),
	},
	{
		ins(ASL, Abs, 6), ins(ASL, AbsX, 7), ins(ROL, Abs, 6), ins(ROL, AbsX, 7),
		ins(LSR, Abs, 6), ins(LSR, AbsX, 7), ins(ROR, Abs, 6), ins(ROR, AbsX, 7),
		ins(STX, Abs, 4), nil, ins(LDX, Abs, 4), ins(LDX, AbsY, 4),
		ins(DEC, Abs, 6), ins(DEC, AbsX, 7), ins(INC, Abs, 6), ins(INC, AbsX, 7),
	},
	nil,
}

func FetchInstruction(i byte) *Instruction {
	low := lowNibble(i)
	high := highNibble(i)

	l := z[high]
	if l == nil {
		return nil
	} else {
		return l[low]
	}
}
