// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Mult.asm

// Multiplies R0 and R1 and stores the result in R2.
// (R0, R1, R2 refer to RAM[0], RAM[1], and RAM[2], respectively.)

// i = R1
@R1
D=M
@i
M=D

// R2 = 0
@R2
M=0

(LOOP)
// if i == 0; goto STOP
@i
D=M
@STOP
D; JEQ

// R2 += R0
@R2
D=M
@R0
D=D+M
@R2
M=D

// i--
@i
M=M-1

// goto LOOP
@LOOP
0; JMP

(STOP)
@STOP
0; JMP

