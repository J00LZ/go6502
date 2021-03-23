  .org $8000

reset:
  lda #$01
  tax
loop:
  jmp loop

  .org $fffc
  .word reset
  .word $0000
