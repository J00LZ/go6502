package graphics

type Color struct{ R, G, B byte }

func Gray(scale byte) Color {
	return Color{scale, scale, scale}
}

func scale() [32]Color {
	c := [32]Color{}
	for i := range c {
		c[i] = Gray(byte(i) * 8)
	}
	return c
}

var Palette0 = [32]Color{{0, 0, 0}, {0xFF, 0xFF, 0xFF}, {0xFF, 0, 0}, {0, 0xFF, 0},
	{0, 0, 0xFF}, {0xFF, 0xFF, 0}, {0xFF, 0, 0xFF}, {0, 0xFF, 0xFF}, {0xFF, 0xA5, 0}, {}, {}, {}, {}, {}, {}, {}, {}}

var GrayPalette = scale()

var Pallets = [8][32]Color{Palette0, Palette0, Palette0, Palette0, Palette0, Palette0, Palette0, GrayPalette}

func Yeet() {

}
