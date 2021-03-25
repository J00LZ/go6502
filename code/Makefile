
GFX_FILES = gfx/gfx.c
GFX = $(patsubst %.c,%.o,$(GFX_FILES))

%.o: %.s
	ca65 -o $@ $^ -t none

%.s: %.c
	cc65 -o $@ $^ -t none

%.o: %.c
	cl65 -o $@ $^ -t none -c


graphics: graphics.o $(GFX) reset.o
	ld65 -o $@ $^ none.lib -C link.ld

blinkc: blinkc.o reset.o
	ld65 -o $@ $^ -C link.ld
