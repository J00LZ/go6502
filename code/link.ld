
MEMORY {
    RAM1:
        start $0000
        size $1000;
    ROM1:
        start $8000
        size $8000;
}

SEGMENTS {
    CODE:
        load = ROM1
        type=ro;
    RODATA:
        load = ROM1
        type=ro;
    DATA:
        load = RAM1
        type=rw;
    BSS:
        load = RAM1
        type=bss
        define=yes;
    ZEROPAGE:
        load = RAM1
        type = zp
        define = yes;
    RESET:
        load = ROM1
        start = $fffa
        type = ro;
}

