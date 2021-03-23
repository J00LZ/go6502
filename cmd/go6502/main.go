package main

import (
	"github.com/J00LZZ/go6502/pkg/bus"
	"github.com/J00LZZ/go6502/pkg/cpu"
	"github.com/J00LZZ/go6502/pkg/graphics"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"image/color"
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

	vram := &bus.ByteBus{
		StartVal: 0x1000,
		Arr:      make([]byte, 0x1000),
		Name:     "VRAM",
		Type:     bus.RW,
	}

	// blink stolen from Ben Eater (he inspired this project)
	blink, err := ioutil.ReadFile("./code/graphics")
	if err != nil {
		log.Panic(err)
	}

	rom := &bus.ByteBus{StartVal: 0x8000, Arr: blink, Name: "ROM", Type: bus.R}
	b := bus.Bus{Devices: []bus.Device{ram, vram, rom}}
	c := cpu.CPU{
		Bus: &b,
	}
	c.Reset()
	tickSpeed := time.Second / 100000
	ticker := time.NewTicker(tickSpeed)
	go func() { c.Run(ticker.C) }()
	window(c, ticker, vram)
}
func window(c cpu.CPU, ticker *time.Ticker, vram *bus.ByteBus) {
	cfg := pixelgl.WindowConfig{
		Title:     "go6502",
		Bounds:    pixel.R(0, 0, 1024, 768),
		VSync:     true,
		Resizable: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.Clear(color.Gray{Y: 0x20})

	canvas := pixelgl.NewCanvas(pixel.R(0, 0, 64, 64))
	for !win.Closed() {

		c.Pause = true
		<-ticker.C
		pixels2 := make([]byte, 0x4000)
		for i, v := range vram.Arr {
			//R = (byte) ((rgb8 & 0xE0) >> 5);     // rgb8 & 1110 0000  >> 5
			//G = (byte ((rgb8 & 0x1C) >> 2);     // rgb8 & 0001 1100  >> 2
			//B = (byte (rgb8 & 0x03);            // rgb8 & 0000 0011

			pallet := v >> 5
			index := v & 0b11111
			c := graphics.Pallets[pallet][index]
			pixels2[i*4] = c.R
			pixels2[i*4+1] = c.G
			pixels2[i*4+2] = c.B
			pixels2[i*4+3] = 255
		}
		sprite := pixel.NewSprite(canvas, canvas.Bounds())
		canvas.SetPixels(pixels2)
		sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()).Scaled(win.Bounds().Center(), win.Bounds().H()/64))
		win.Update()
		c.Pause = false
	}
}
func main() {
	log.Printf("%d", 0x1000)
	pixelgl.Run(run)
}
