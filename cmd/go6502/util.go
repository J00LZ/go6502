package main

import (
	"errors"
	"github.com/J00LZZ/go6502/pkg/bus"
	"github.com/J00LZZ/go6502/pkg/deviceinfo"
	"github.com/J00LZZ/go6502/pkg/graphics"
	"github.com/J00LZZ/go6502/pkg/interrupt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func loadMapFile() (*bus.Bus, error) {
	f, err := ioutil.ReadFile("./map.yml")
	if err == nil {
		var devMap deviceinfo.DeviceMap
		err = yaml.Unmarshal(f, &devMap)
		if err == nil {
			devices := make([]bus.Device, 0, len(devMap))
			for addr, def := range devMap {
				switch def.Type {
				case deviceinfo.ROM:
					f := def.File
					if len(f) == 0 {
						if len(os.Args) >= 2 {
							f = os.Args[1]
						} else {
							return nil, errors.New("no file provided, can't create ROM")
						}
					}
					devices = append(devices, bus.NewRom(addr, f))
				case deviceinfo.RAM:
					devices = append(devices, bus.NewRam(addr, def.Size))
				case deviceinfo.PPU:
					devices = append(devices, graphics.CreatePPU(addr))
				case deviceinfo.IMU:
					devices = append(devices, interrupt.NewInterruptManager(addr, uint8(def.Size)))
				}
			}
			deviceMap, err := bus.New(devices...)

			if err != nil {
				return nil, err
			}
			return deviceMap, nil

		}
	}

	return nil, err
}

func findPPU(b *bus.Bus) *graphics.PPU {
	for _, d := range b.Devices {
		if d.GetType() == deviceinfo.PPU {
			return d.(*graphics.PPU)
		}
	}
	return nil
}

func findIMU(b *bus.Bus) *interrupt.InterruptManager {
	for _, d := range b.Devices {
		if d.GetType() == deviceinfo.IMU {
			return d.(*interrupt.InterruptManager)
		}
	}
	return nil
}
