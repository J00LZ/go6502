package cpu

import (
	"github.com/J00LZZ/go6502/pkg/bus"
	"github.com/J00LZZ/go6502/pkg/deviceinfo"
	"testing"
)

var vectorRam = bus.NewRam(0xFFFA, 6)

func execProcessor2(t *testing.T, acc, x, y byte, time int, ops ...byte) CPU {
	data := bus.ByteBus{
		StartVal: 0,
		Arr:      ops,
		Name:     "Test",
		RWMode:   deviceinfo.RW,
	}
	actualRam := bus.ByteBus{
		StartVal: 0x80,
		Arr:      make([]byte, 0x200),
		Name:     "ActualRam",
		RWMode:   deviceinfo.RW,
	}
	b := &bus.Bus{Devices: []bus.Device{vectorRam, &data, &actualRam}}
	c := CPU{}
	c.Bus = b
	c.Reset()
	c.AC = acc
	c.X = x
	c.Y = y
	time2 := c.RunInstruction()
	if time2 != time {
		t.Errorf("%d != %d", time, 7)
	}
	return c
}

func execProcessor(t *testing.T, time int, ops ...byte) CPU {
	return execProcessor2(t, 0, 0, 0, time, ops...)
}

func TestBRK(t *testing.T) {
	c := execProcessor(t, 7, 0x00)
	if c.PC != 0 {
		t.Errorf("Processor did not set PC after break!")
	}
	if !c.Has(I) {
		t.Errorf("Processor did not set Interrupt signal")
	}
}

func TestORAXInd(t *testing.T) {
	c := execProcessor(t, 6, 0x01, 0x2, 0, 0)
	if c.AC != 0x01 {
		t.Errorf("Acc was %02X, should be 01", c.AC)
	}
}

func TestORAZPG(t *testing.T) {
	c := execProcessor(t, 3, 0x05, 0x1, 0, 0)
	if c.AC != 0x01 {
		t.Errorf("Acc was %02X, should be 01", c.AC)
	}
}

func TestASLZPG(t *testing.T) {
	c := execProcessor(t, 5, 0x06, 0x1, 0, 0)
	if c.ReadAddress(0x01) != 0b10 {
		t.Errorf("Address was not shifted 1 index!")
	}
	c = execProcessor(t, 5, 0x06, 0x2, 0xFF, 0)
	if c.ReadAddress(0x02) != 0xFE {
		t.Errorf("Address was not shifted 1 index!")
	}
	if !c.Has(C) {
		t.Errorf("Carry was not set on overflow!")
	}
}

func TestPHP(t *testing.T) {
	c := execProcessor(t, 3, 0x08)
	stack := c.popStack()
	if stack != byte(c.SR) {
		t.Errorf("Did not push status correctly, %b != %b", stack, byte(c.SR))
	}
}

func TestORAIMM(t *testing.T) {
	c := execProcessor(t, 2, 0x09, 0x01)
	if c.AC != 0x01 {
		t.Errorf("Acc was %02X, should be 01", c.AC)
	}
}

func TestASLAcc(t *testing.T) {
	c := execProcessor2(t, 1, 0, 0, 2, 0x0A)
	if c.AC != 0b10 {
		t.Errorf("Acc was %02X, should be 0x10", c.AC)
	}
}

func TestORAABS(t *testing.T) {
	c := execProcessor(t, 4, 0x0D, 0, 0, 0)
	if c.AC != 0x0D {
		t.Errorf("Acc was %02X, should be 01", c.AC)
	}
}

func TestASLABS(t *testing.T) {
	c := execProcessor(t, 6, 0x0E, 0x1, 0, 0)
	if c.ReadAddress(0x01) != 0b10 {
		t.Errorf("Address was not shifted 1 index!")
	}
	c = execProcessor(t, 5, 0x06, 0x2, 0xFF, 0)
	if c.ReadAddress(0x02) != 0xFE {
		t.Errorf("Address was not shifted 1 index!")
	}
	if !c.Has(C) {
		t.Errorf("Carry was not set on overflow!")
	}
}

func TestORAIndY(t *testing.T) {
	c := execProcessor(t, 5, 0x11, 0x02, 0, 0)
	if c.AC != 0x11 {
		t.Errorf("Acc was %02X, should be 01", c.AC)
	}
}
func TestORAZPGX(t *testing.T) {
	c := execProcessor(t, 4, 0x15, 0x1, 0, 0)
	if c.AC != 0x01 {
		t.Errorf("Acc was %02X, should be 01", c.AC)
	}
}
func TestASLZPGX(t *testing.T) {
	c := execProcessor(t, 6, 0x16, 0x1, 0, 0)
	if c.ReadAddress(0x01) != 0b10 {
		t.Errorf("Address was not shifted 1 index!")
	}
	c = execProcessor(t, 6, 0x16, 0x2, 0xFF, 0)
	if c.ReadAddress(0x02) != 0xFE {
		t.Errorf("Address was not shifted 1 index!")
	}
	if !c.Has(C) {
		t.Errorf("Carry was not set on overflow!")
	}
}

func TestCLC(t *testing.T) {
	c := execProcessor(t, 2, 0x18)
	if c.Has(C) {
		t.Errorf("Carry is set!")
	}
}

func TestORAABSY(t *testing.T) {
	c := execProcessor(t, 4, 0x19, 0x1, 0, 0)
	if c.AC != 0x01 {
		t.Errorf("Acc was %02X, should be 01", c.AC)
	}
}

func TestORAABSX(t *testing.T) {
	c := execProcessor(t, 4, 0x1D, 0x1, 0, 0)
	if c.AC != 0x01 {
		t.Errorf("Acc was %02X, should be 01", c.AC)
	}
}

func TestASLABSX(t *testing.T) {
	c := execProcessor(t, 7, 0x1E, 0x1, 0, 0)
	if c.ReadAddress(0x01) != 0b10 {
		t.Errorf("Address was not shifted 1 index!")
	}
	c = execProcessor(t, 7, 0x1E, 0x3, 0, 0xFF, 0)
	if c.ReadAddress(0x03) != 0xFE {
		t.Errorf("Address was not shifted 1 index!")
	}
	if !c.Has(C) {
		t.Errorf("Carry was not set on overflow!")
	}
}

func TestJSR(t *testing.T) {
	c := execProcessor(t, 6, 0x20, 0x20, 0x00)
	if c.PC != 0x20 {
		t.Errorf("JSR did not update PC, %04X!", c.PC)
	}
}


///////
func TestINX(t *testing.T) {
	c := execProcessor(t, 2, 0xE8)
	if c.X != 1 {
		t.Errorf("X is not 1, it is %c", c.X)
	}
}
func TestINY(t *testing.T) {
	c := execProcessor(t, 2, 0xC8)

	if c.Y != 1 {
		t.Errorf("Y is not 1, it is %c", c.Y)
	}
}
