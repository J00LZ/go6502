package bus

type Bus struct {
	Devices []Device
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
