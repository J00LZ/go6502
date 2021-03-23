# go6502
> emulating a cpu, but in Go!

I was bored, so I made a 6502 emulator. Partially due to [Ben Eater's](https://eater.net/) video's, 
partially due to my father talking about it a lot since he used a lot of 6502 based
computers when he was younger. 

## Features (and todo)
 * [x] Supports the full instruction set of the MOS 6502
 * [x] Supports loading a file as rom (currently hardcoded)
 * [x] Has actual ram it can load
 * [ ] Add more devices
    * [ ] Add a basic graphics device
    * [ ] Add a realistic graphics device (vga-like output)
    * [ ] Add a device that captures "input"
    * [ ] Add a MMU (memory mapper)
    * [ ] Add a FPU (floating point support)
    * [ ] Add an APU (audio support)
    * [ ] Add an IMU (interrupt management unit)
    * [ ] Add a DTB (Device Tree blob, tihi)
 * [ ] maybe more?

