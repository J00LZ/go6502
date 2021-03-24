package bus

import "github.com/J00LZZ/go6502/pkg/deviceinfo"

type ByteBus struct {
	StartVal uint16
	Arr      []byte
	Name     string
	RWMode   deviceinfo.RWMode
}

func NewRam(location, size uint16) *ByteBus {
	return &ByteBus{
		StartVal: location,
		Arr:      make([]byte, size),
		Name:     "RAM",
		RWMode:   deviceinfo.RW,
	}
}

func (b *ByteBus) Start() uint16 {
	return b.StartVal
}

func (b *ByteBus) End() uint16 {
	return b.StartVal + uint16(len(b.Arr)-1)
}

func (b *ByteBus) LoadAddress(address uint16) byte {
	address = address - b.StartVal
	return b.Arr[address]
}

func (b *ByteBus) WriteAddress(address uint16, data byte) {
	address = address - b.StartVal
	b.Arr[address] = data
}

func (b *ByteBus) GetName() string {
	return b.Name
}

func (b *ByteBus) GetRWMode() deviceinfo.RWMode {
	return b.RWMode
}

func (b *ByteBus) GetType() deviceinfo.DeviceType {
	return deviceinfo.RAM
}
