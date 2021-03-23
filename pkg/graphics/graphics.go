package graphics

type Color struct{ R, G, B byte }

func Gray(scale byte) Color {
	return Color{scale, scale, scale}
}

var Palette0 = [32]Color{{0, 0, 0}, {0xFF, 0xFF, 0xFF}, {0xFF, 0, 0}, {0, 0xFF, 0},
	{0, 0, 0xFF}, {0xFF, 0xFF, 0}, {0xFF, 0, 0xFF}, {0, 0xFF, 0xFF}, {0xFF, 0xA5, 0},
	Gray(0x10), Gray(0x20), Gray(0x30), Gray(0x40), Gray(0x50), Gray(0x60),
	Gray(0x70), Gray(0x80), Gray(0x90), Gray(0xA0), Gray(0xB0), Gray(0xC0),
	Gray(0xD0), Gray(0xE0), Gray(0xF0)}

var Pallets = [][32]Color{Palette0, Palette0, Palette0, Palette0, Palette0, Palette0, Palette0, Palette0}

func Yeet() {

}
