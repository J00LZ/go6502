
%.o: %.s
	ca65 -o $@ $^ -t none

%.s: %.c
	cc65 -o $@ $^ -t none

%.o: %.c
	cl65 -o $@ $^ -t none -c


blinkc: blinkc.o reset.o
	ld65 -o $@ $^ -C link.ld
