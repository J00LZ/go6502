package main

import (
	"github.com/J00LZZ/go6502/pkg/bus"
	"github.com/J00LZZ/go6502/pkg/cpu"
	"github.com/J00LZZ/go6502/pkg/graphics"
	"github.com/J00LZZ/go6502/pkg/interrupt"
	"github.com/faiface/pixel/pixelgl"
	"log"
	"time"
)

func run() {
	imu := interrupt.NewInterruptManager(0x6000, 16)

	deviceMap, err := bus.New(
		bus.NewRam(0x0, 0x1000),
		bus.NewRom(0x8000, "./code/blinkc"),
		graphics.CreatePPU(0x1000),
		imu,
	)
	if err != nil {
		log.Fatal(err)
	}
	m, err := loadMapFile()
	if err != nil {
		log.Fatal(err)
	} else {
		deviceMap = m
	}

	c := cpu.New(deviceMap)
	imu.SetNMIFunc(c.NMI)
	imu.SetIRQFunc(c.IRQ)

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
