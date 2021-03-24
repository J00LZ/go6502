package graphics

type Color struct{ R, G, B byte }

func Gray(scale byte) Color {
	return Color{scale, scale, scale}
}

func Html(color int) Color {
	return Color{
		R: byte((color >> 16) & 0xFF),
		G: byte((color >> 8) & 0xFF),
		B: byte((color >> 0) & 0xFF),
	}
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

var Endesga = [32]Color{
	Html(0xbe4a2f), Html(0xd77643), Html(0xead4aa), Html(0xe4a672), Html(0xb86f50),
	Html(0x733e39), Html(0x3e2731), Html(0xa22633), Html(0xe43b44), Html(0xf77622),
	Html(0xfeae34), Html(0xfee761), Html(0x63c74d), Html(0x3e8948), Html(0x265c42),
	Html(0x193c3e), Html(0x124e89), Html(0x0099db), Html(0x2ce8f5), Html(0xffffff),
	Html(0xc0cbdc), Html(0x8b9bb4), Html(0x5a6988), Html(0x3a4466), Html(0x262b44),
	Html(0x181425), Html(0xff0044), Html(0x68386c), Html(0xb55088), Html(0xf6757a),
	Html(0xe8b796), Html(0xc28569),
}

var Clear = [32]Color{
	Html(0x1b0813),
	Html(0x1c102e),
	Html(0x32303c),
	Html(0x595a61),
	Html(0x919191),
	Html(0xbdbdbd),
	Html(0xe9e9e9),
	Html(0x6fd3e9),
	Html(0x3e95cb),
	Html(0x1c47ab),
	Html(0x12197e),
	Html(0x151e4a),
	Html(0x124a43),
	Html(0x0f6c3f),
	Html(0x259d15),
	Html(0x9fd228),
	Html(0xfbf135),
	Html(0xfb981e),
	Html(0xf55718),
	Html(0xb51637),
	Html(0x7e0b36),
	Html(0x500839),
	Html(0x360736),
	Html(0x301911),
	Html(0x573824),
	Html(0x7d5c38),
	Html(0xbfa363),
	Html(0xe0d18c),
	Html(0xfbc2ff),
	Html(0xf780ff),
	Html(0xee31eb),
	Html(0xc5179c),
}

var Pallets = [8][32]Color{Palette0, Endesga, Clear, Palette0, Palette0, Palette0, Palette0, GrayPalette}
