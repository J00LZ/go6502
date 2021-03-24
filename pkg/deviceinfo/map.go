package deviceinfo

type DeviceMap map[uint16]DeviceEntry

type DeviceEntry struct {
	Type DeviceType `yaml:"type"`
	File string     `yaml:"file,omitempty"`
	Size uint16     `yaml:"size,omitempty"`
}
