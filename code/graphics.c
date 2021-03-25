#include "gfx/gfx.h"

void reset() {
    color_t red;

    gfx_init(0x1000, 64, 64);

    red = get_color(Endesga, 3);
    clear(red);
}