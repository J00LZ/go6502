.autoimport +

___reset:
    jsr _reset

    ___loop:
        jmp ___loop



.segment "RESET"

.word $0000
.word ___reset
.word $0000
