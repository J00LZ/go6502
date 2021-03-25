# go6502
> emulating a cpu, but in Go!

I was bored, so I made a 6502 emulator. Partially due to [Ben Eater's](https://eater.net/) video's, 
partially due to my father talking about it a lot since he used a lot of 6502 based
computers when he was younger. 

## Features (and todo)
 * [x] Supports the full instruction set of the MOS 6502
 * [x] Allows you to create a YAML file to specify the memory locations of devices
 * [ ] Add more devices
    * [x] Add a basic graphics device (pixels+pallettes)
    * [ ] Add a realistic graphics device (vga-like output)
    * [ ] Add a device that captures "input"
    * [ ] Add a MMU (memory mapper)
    * [ ] Add a FPU (floating point support)
    * [ ] Add an APU (audio support)
    * [ ] Add an IMU (interrupt management unit)
    * [ ] Add a DTB (Device Tree blob, tihi)
 * [ ] maybe more?

## Creating a `map.yml`
This is the default map.yml file, it loads 4kb of ram at address 0, the PPU at address 4096 (also 4kb, 64x64 display)
and the rom file at address 0x8000, or 32k. 
```yaml
---
0x0:
  type: RAM
  size: 0x1000
0x1000:
  type: PPU
0x8000:
  type: ROM
  file: ./code/graphics
```
Currently, only 1 PPU is supported, and the cpu will only ever use the one with the lowest address, 
the higher ones will simply provide 4kb of vram. 

The types available are currently
* RAM: Random access memory, Readable and Writable
* ROM: Read Only Memory, only Readable
* PPU: Attaches a display, and provides 4kb of vram.

## Writing C for the 6502
You can use the `code/Makefile` as an example
* Always link to reset.o (or reset.s)
* Call your main function `reset`, in your main.c file
* have fun!
