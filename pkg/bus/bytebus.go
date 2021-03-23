package bus

type ByteBus struct {
	StartVal uint16
	Arr      []byte
	Name     string
	Type
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

func (b *ByteBus) GetType() Type {
	return b.Type
}
