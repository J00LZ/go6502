
#ifndef GFX_H
#define GFX_H

#include <stdint.h>

extern char * video_buffer;
extern uint8_t width;
extern uint8_t height;

typedef uint8_t color_t;

void gfx_init(uintptr_t video_buffer_address, uint8_t w, uint8_t h);
void clear(color_t color);

typedef enum {
    Palette0 = 0,
    Endesga = 1,
    Clear = 2,

    GrayPalette = 7,
} Palette;


color_t get_color(Palette p, uint8_t index);

#endif