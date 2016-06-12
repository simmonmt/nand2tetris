// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Fill.asm

// Runs an infinite loop that listens to the keyboard input.
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel. When no key is pressed, the
// program clears the screen, i.e. writes "white" in every pixel.

// last = 0
@last
M=0

// R14 = LOOP
@LOOP
D=A
@R14
M=D

(LOOP)
// kbd -> R15
@KBD
D=M
@R15
M=D

// if R15 == 0 goto KBDNORMALIZED
@KBDNORMALIZED
D; JEQ

// R15 = -1
@R15
M=-1

(KBDNORMALIZED)
// if R15 == last goto LOOP
@R15
D=M
@last
D=D-M
@LOOP
D; JEQ

// last = R15
@R15
D=M
@last
M=D

// goto FILL
@FILL
0; JMP

@R15
M=-1
@END
D=A
@R14
M=D

@FILL
0; JMP

(END)
@END
0; JMP

(FILL)
// IN: value in R15; ret in R14
// CLOBBERS: R0, R1, R2
//
// R0=SCREEN
@SCREEN
D=A
@R0
M=D

// R1=8192
@8192
D=A
@R1
M=D

(FILLLOOP)
// if [R1] == 0 goto FILLRET
@R1
D=M
@FILLRET
D; JEQ

// [R0] = value
// R0++
@R15
D=M
@R0
A=M
M=D
D=A+1
@R0
M=D

// R1--
@R1
M=M-1

// goto FILLLOOP
@FILLLOOP
0; JMP

(FILLRET)
@R14
A=M
0; JMP
