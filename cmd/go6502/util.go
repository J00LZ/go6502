package main

import (
	"github.com/J00LZZ/go6502/pkg/bus"
	"github.com/J00LZZ/go6502/pkg/deviceinfo"
	"github.com/J00LZZ/go6502/pkg/graphics"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func loadMapFile() {
	f, err := ioutil.ReadFile("./map.yml")
	if err == nil {
		var devMap deviceinfo.DeviceMap
		err = yaml.Unmarshal(f, &devMap)
		if err == nil {
			devices := make([]bus.Device, 0, len(devMap))
			for addr, def := range devMap {
				switch def.Type {
				case deviceinfo.ROM:
					devices = append(devices, bus.NewRom(addr, def.File))
				case deviceinfo.RAM:
					devices = append(devices, bus.NewRam(addr, def.Size))
				case deviceinfo.PPU:
					devices = append(devices, graphics.CreatePPU(addr))
				}
			}
			deviceMap = bus.Bus{Devices: devices}
		}
	}
}


func findPPU(b *bus.Bus) *graphics.PPU {
	for _, d := range b.Devices {
		if d.GetType() == deviceinfo.PPU {
			return d.(*graphics.PPU)
		}
	}
	return nil
}