package bus

import (
	"fmt"
)

type Bus struct {
	Devices []Device
}

func New(devices ...Device) (*Bus, error) {
	for _, device1 := range devices {
		for _, device2 := range devices {
			if device1 != device2 {
				if device1.Start() <= device2.End() && device2.Start() <= device1.End() {
					return nil, fmt.Errorf(
						"Memory ranges of %v (%x..%x) and %v  (%x..%x)  overlap",
						device1.GetName(),
						device1.Start(),
						device1.End(),

						device2.GetName(),
						device2.Start(),
						device2.End())
				}
			}
		}
	}


	return &Bus{
		Devices: devices,
	}, nil
}

func (b *Bus) ReadAddress(address uint16) byte {
	for _, d := range b.Devices {
		//n := d.GetName()
		if d.GetType().HasFlag(R) && address >= d.Start() && address < d.End() {
			//log.Printf("%v had address %X", n, address)
			z := d.LoadAddress(address)
			return z
		}
	}
	return 0
}

func (b *Bus) WriteAddress(address uint16, data byte) {
	for _, d := range b.Devices {
		if d.GetType().HasFlag(W) && address >= d.Start() && address < d.End() {
			d.WriteAddress(address, data)
		}
	}
}

type Device interface {
	Start() uint16
	End() uint16
	LoadAddress(address uint16) byte
	WriteAddress(address uint16, data byte)
	GetName() string
	GetType() Type
}
