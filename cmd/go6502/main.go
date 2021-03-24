package main

import (
	"github.com/J00LZZ/go6502/pkg/bus"
	"github.com/J00LZZ/go6502/pkg/cpu"
	"github.com/J00LZZ/go6502/pkg/graphics"
	"github.com/faiface/pixel/pixelgl"
	"time"
)

var deviceMap bus.Bus

func run() {
	deviceMap = bus.Bus{Devices: []bus.Device{bus.NewRam(0x0, 0x1000), bus.NewRom(0x8000, "./code/graphics"), graphics.CreatePPU(0x1000)}}
	loadMapFile()

	c := cpu.CPU{
		Bus: &deviceMap,
	}
	c.Reset()
	tickSpeed := time.Second / 10000
	ticker := time.NewTicker(tickSpeed)

	if ppu := findPPU(c.Bus); ppu != nil {
		go func() {
			c.Run(ticker.C)
		}()
		ppu.RunWindow(c, ticker)
	} else {
		c.Run(ticker.C)
	}

}

func main() {
	// handoff needed if graphical output is enabled,
	// if not enabled it does not add extra overhead.
	pixelgl.Run(run)

}
