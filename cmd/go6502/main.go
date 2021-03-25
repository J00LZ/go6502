package main

import (
	"github.com/J00LZZ/go6502/pkg/cpu"
	"github.com/faiface/pixel/pixelgl"
	"log"
	"time"
)

func run() {
	deviceMap, err := loadMapFile()
	if err != nil {
		log.Fatal(err)
	}

	c := cpu.New(deviceMap)
	if imu := findIMU(deviceMap); imu != nil {
		imu.SetNMIFunc(c.NMI)
		imu.SetIRQFunc(c.IRQ)
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
