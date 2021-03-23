package main

import (
	"github.com/J00LZZ/go6502/pkg/bus"
	"github.com/J00LZZ/go6502/pkg/cpu"
	"log"
	"time"
)

func main() {
	ram := &bus.ByteBus{
		StartVal: 0,
		Arr:      make([]byte, 0x1000),
		Name:     "RAM",
	}
	romm := make([]byte, 0x8000)
	romm[0] = 0xA9
	romm[1] = 0x01
	romm[2] = 0xAA
	romm[3] = 0x4C
	romm[4] = 0x03
	romm[5] = 0x80
	romm[cpu.RstVectorH-0x8000] = 0x80
	romm[cpu.RstVectorL-0x8000] = 0x00
	rom := &bus.ByteBus{StartVal: 0x8000, Arr: romm, Name: "ROM"}
	log.Printf("Start %X end %X", rom.Start(), rom.End())
	b := bus.Bus{Devices: []bus.Device{ram, rom}}
	c := cpu.CPU{
		PC:  0x8000,
		X:   0,
		Bus: &b,
	}
	tickSpeed := time.Second
	ticker := time.NewTicker(tickSpeed)
	c.Run(ticker)

}
