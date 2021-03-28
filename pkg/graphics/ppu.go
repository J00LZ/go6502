package graphics

import (
	"fmt"
	"github.com/J00LZZ/go6502/pkg/cpu"
	"github.com/J00LZZ/go6502/pkg/deviceinfo"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
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

func (p *PPU) RunWindow(c *cpu.CPU, ticker *cpu.Ticker) {
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
	clearColor := color.Gray{Y: 0x20}
	p.Clear(clearColor)

	canvas := pixelgl.NewCanvas(pixel.R(0, 0, 64, 64))
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)

	for !p.Closed() {

		c.Lock()
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
		p.Clear(clearColor)
		sprite := pixel.NewSprite(canvas, canvas.Bounds())
		canvas.SetPixels(pixels2)
		sprite.Draw(p, pixel.IM.Scaled(pixel.ZV, p.Bounds().H()/64).Moved(p.Bounds().Center().Add(pixel.V(120, 0))))

		txt := text.New(pixel.V(0, p.Bounds().H()-atlas.LineHeight()), atlas)
		txt.LineHeight = atlas.LineHeight() * 1.5

		s := c.DisassembleMore()
		idx := findCurrentIndex(s, c.DisassembleCurrent())
		startIdx := idx - 5
		endIdx := idx + 5

		txt.Color = colornames.White
		for i := startIdx; i < idx; i++ {
			if i < 0 {
				_, _ = txt.WriteString("\n")
			} else {
				_, _ = fmt.Fprintf(txt, "%v\n", s[i].Text)
			}
		}
		txt.Color = colornames.Green
		_, _ = fmt.Fprintf(txt, "%v\n", c.DisassembleCurrent().Text)
		txt.Color = colornames.White
		for i := idx + 1; i < endIdx; i++ {
			if i >= len(s) {
				_, _ = txt.WriteString("\n")
			} else {
				_, _ = fmt.Fprintf(txt, "%v\n", s[i].Text)
			}
		}

		if c.Manual {
			_, _ = fmt.Fprintf(txt, "Manual (press X)\n")

		} else {
			_, _ = fmt.Fprintf(txt, "speed: %s\n", c.Speed.String())
		}

		_, _ = fmt.Fprintf(txt, "A: %02X, X: %02X, Y: %02X\n", c.AC, c.X, c.Y)
		_, _ = fmt.Fprintf(txt, "PC: %04X\n", c.PC)
		_, _ = fmt.Fprintf(txt, "CZIDB VN\n")
		_, _ = fmt.Fprintf(txt, "%b%b%b%b%b%b%b%b\n", c.Val(cpu.C), c.Val(cpu.Z), c.Val(cpu.I), c.Val(cpu.D), c.Val(cpu.B), c.Val(cpu.Unused), c.Val(cpu.V), c.Val(cpu.N))

		sidx := int(c.SP)

		txt.Color = colornames.Cyan
		_, _ = fmt.Fprintf(txt, "Stack\n")
		_, _ = fmt.Fprintf(txt, "%04X: %02X\n", uint16(sidx)+0x0100, c.ReadAddress(uint16(sidx)+0x0100))
		txt.Color = colornames.White
		for i := 1; i < 10; i++ {
			target := 0x0100 + uint16(i+sidx)
			if target <= 0x01FF {
				_, _ = fmt.Fprintf(txt, "%04X: %02X\n", target, c.ReadAddress(target))
			}

		}


		txt.Draw(p, pixel.IM)

		if p.JustPressed(pixelgl.KeyLeftBracket) {
			c.Speed *= 10
			ticker.TTicker.Reset(c.Speed)
		}
		if p.JustPressed(pixelgl.KeyRightBracket) {
			c.Speed /= 10
			ticker.TTicker.Reset(c.Speed)
		}
		if p.JustPressed(pixelgl.KeyEscape) {
			p.SetClosed(true)
		}
		if p.JustPressed(pixelgl.KeyC) {
			c.Manual = !c.Manual
			if c.Manual {
				ticker.TTicker.Stop()
			} else {
				ticker.TTicker.Reset(c.Speed)
			}
		}
		if p.JustPressed(pixelgl.KeyX) && c.Manual {
			ticker.T <- time.Now()
		}
		//if p.JustPressed(pixelgl.KeyC) {
		//	c.Manual = !c.Manual
		//}
		//if p.JustPressed(pixelgl.KeyX) && c.Manual {
		//
		//}
		p.Update()
		c.Unlock()
	}
}

func findCurrentIndex(d []cpu.Dis, current cpu.Dis) int {
	for i, s := range d {
		if s == current {
			return i
		}
	}
	return 0
}
