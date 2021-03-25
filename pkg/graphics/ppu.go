package graphics

import (
	"github.com/J00LZZ/go6502/pkg/cpu"
	"github.com/J00LZZ/go6502/pkg/deviceinfo"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"image/color"
	"time"
)

type PPU struct {
	*pixelgl.Window
	vram       []byte
	RangeStart uint16
}

func (p *PPU) Start() uint16 {
	return p.RangeStart
}

func (p *PPU) End() uint16 {
	return p.RangeStart + 0x1000 - 1
}

func (p *PPU) LoadAddress(address uint16) byte {
	return 0
}

func (p *PPU) WriteAddress(address uint16, data byte) {
	p.vram[address-p.RangeStart] = data
}

func (p *PPU) GetName() string {
	return "PPU"
}

func (p *PPU) GetRWMode() deviceinfo.RWMode {
	return deviceinfo.W
}

func (p *PPU) GetType() deviceinfo.DeviceType {
	return deviceinfo.PPU
}

func CreatePPU(start uint16) *PPU {

	return &PPU{nil, make([]byte, 0x1000), start}
}

func (p *PPU) RunWindow(c *cpu.CPU, ticker *time.Ticker) {
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
	p.Window = win
	p.Clear(color.Gray{Y: 0x20})

	canvas := pixelgl.NewCanvas(pixel.R(0, 0, 64, 64))
	for !p.Closed() {

		c.Pause = true
		<-ticker.C
		pixels2 := make([]byte, 0x4000)
		for i, v := range p.vram {
			pallet := v >> 5
			index := v & 0b11111
			c := Pallets[pallet][index]
			pixels2[i*4] = c.R
			pixels2[i*4+1] = c.G
			pixels2[i*4+2] = c.B
			pixels2[i*4+3] = 255
		}
		sprite := pixel.NewSprite(canvas, canvas.Bounds())
		canvas.SetPixels(pixels2)
		sprite.Draw(p, pixel.IM.Moved(p.Bounds().Center()).Scaled(p.Bounds().Center(), p.Bounds().H()/64))
		p.Update()
		c.Pause = false
	}
}
