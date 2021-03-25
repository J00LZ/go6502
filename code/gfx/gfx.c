#include "gfx.h"
#include <stdint.h>

char * video_buffer;
uint8_t width;
uint8_t height;

void gfx_init(uintptr_t video_buffer_address, uint8_t w, uint8_t h) {
    video_buffer = (char *)video_buffer_address;
    width = w;
    height = h;
}

void clear(color_t color) {
    uint8_t i;
    uint8_t j;

    for (i = 0; i < width; ++i) {
        for (j = 0; j < height; ++j) {
            video_buffer[j * width + i] = color;
        }
    }
}

color_t get_color(Palette p, uint8_t index) {
    return (uint8_t) p << 5 | index && 0b11111;
}