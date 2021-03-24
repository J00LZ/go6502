package bus

import (
	"fmt"
	"github.com/J00LZZ/go6502/pkg/deviceinfo"
	"io/ioutil"
)

type Rom struct {
	buf       []byte
	StartAddr uint16
	Filename  string
}

func NewRom(start uint16, filename string) *Rom {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return &Rom{f, start, filename}
}

func (r *Rom) Start() uint16 {
	return r.StartAddr
}

func (r *Rom) End() uint16 {
	return r.StartAddr + uint16(len(r.buf)) - 1
}

func (r *Rom) LoadAddress(address uint16) byte {
	return r.buf[address-r.StartAddr]
}

func (r *Rom) WriteAddress(address uint16, data byte) {
	panic("No writing to rom!")
}

func (r *Rom) GetName() string {
	return fmt.Sprintf("Rom: %v", r.Filename)
}

func (r *Rom) GetRWMode() deviceinfo.RWMode {
	return deviceinfo.R
}

func (r *Rom) GetType() deviceinfo.DeviceType {
	return deviceinfo.ROM
}
