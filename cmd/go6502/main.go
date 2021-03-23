package main

import (
	"github.com/J00LZZ/go6502/pkg/bus"
	"github.com/J00LZZ/go6502/pkg/cpu"
	"io/ioutil"
	"log"
	"time"
)

func main() {
	ram := &bus.ByteBus{
		StartVal: 0,
		Arr:      make([]byte, 0x1000),
		Name:     "RAM",
		Type:     bus.RW,
	}
	//romm := make([]byte, 0x8000)
	//romm[0] = 0xA9 // lda
	//romm[1] = 0x01 // #$01
	//romm[2] = 0xAA // tax
	//romm[3] = 0x4C // JMP
	//romm[4] = 0x03 // #$03
	//romm[5] = 0x80 // #$80
	//romm[cpu.RstVectorH-0x8000] = 0x80
	//romm[cpu.RstVectorL-0x8000] = 0x00
	// blink stolen from Ben Eater (he inspired this project)
	blink, err := ioutil.ReadFile("./code/blink")
	if err != nil {
		log.Panic(err)
	}

	rom := &bus.ByteBus{StartVal: 0x8000, Arr: blink, Name: "ROM", Type: bus.R}
	b := bus.Bus{Devices: []bus.Device{ram, rom}}
	c := cpu.CPU{
		Bus: &b,
	}
	c.Reset()
	tickSpeed := time.Second / 10
	ticker := time.NewTicker(tickSpeed)
	c.Run(ticker)

}
