package deviceinfo


type DeviceType int

//go:generate enumer -type=DeviceType -json -yaml
const (
	ROM DeviceType = iota
	RAM
	PPU
)
