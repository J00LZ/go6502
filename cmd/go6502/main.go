package main

import (
	"github.com/J00LZZ/go6502/pkg/bus"
	"github.com/J00LZZ/go6502/pkg/cpu"
	"github.com/J00LZZ/go6502/pkg/devices"
	"github.com/J00LZZ/go6502/pkg/graphics"
	"github.com/faiface/pixel/pixelgl"
	"io/ioutil"
	"log"
	"time"
)

func run() {
	ram := &bus.ByteBus{
		StartVal: 0,
		Arr:      make([]byte, 0x1000),
		Name:     "RAM",
		Type:     bus.RW,
	}

	ppu := graphics.CreatePPU(0x1000)
	imu := devices.NewInterruptManager(0x6000, 16)

	// blink stolen from Ben Eater (he inspired this project)
	blink, err := ioutil.ReadFile("./code/graphics")
	if err != nil {
		log.Panic(err)
	}

	rom := &bus.ByteBus{StartVal: 0x8000, Arr: blink, Name: "ROM", Type: bus.R}
	b, err := bus.New(ram, rom, ppu, imu)
	if err != nil {
		log.Fatal(err)
	}
	c := cpu.New(b)

	imu.SetNMIFunc(c.NMI)
	imu.SetIRQFunc(c.IRQ)

	c.Reset()
	tickSpeed := time.Second / 10000
	ticker := time.NewTicker(tickSpeed)
	go func() { c.Run(ticker.C) }()
	ppu.RunWindow(c, ticker)
}

func main() {
	// handoff needed if graphical output is enabled,
	// if not enabled it does not add extra overhead.
	pixelgl.Run(run)
}
