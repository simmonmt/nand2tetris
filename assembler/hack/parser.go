package hack

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var (
	labelRegexp  = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_.$]*$")
	cInstrRegexp = regexp.MustCompile("^(?:(A?M?D?)=)?([^;]+)(?:;(J..))?$")

	compStrToBits = map[string]int{
		"0":   COMP_0,
		"1":   COMP_1,
		"-1":  COMP_NEG1,
		"D":   COMP_D,
		"A":   COMP_A,
		"!D":  COMP_NOTD,
		"!A":  COMP_NOTA,
		"-D":  COMP_NEGD,
		"-A":  COMP_NEGA,
		"D+1": COMP_DPLUS1,
		"A+1": COMP_APLUS1,
		"D-1": COMP_DMIN1,
		"A-1": COMP_AMIN1,
		"D+A": COMP_DPLUSA,
		"D-A": COMP_DMINA,
		"A-D": COMP_AMIND,
		"D&A": COMP_DANDA,
		"D|A": COMP_DORA,
		"M":   COMP_M,
		"!M":  COMP_NOTM,
		"-M":  COMP_NEGM,
		"M+1": COMP_MPLUS1,
		"M-1": COMP_MMIN1,
		"D+M": COMP_DPLUSM,
		"D-M": COMP_DMINM,
		"M-D": COMP_MMIND,
		"D&M": COMP_DANDM,
		"D|M": COMP_DORM,
	}

	jumpStrToBits = map[string]int{
		"":    0,
		"JGT": JUMP_GT,
		"JEQ": JUMP_EQ,
		"JGE": JUMP_GT | JUMP_EQ,
		"JLT": JUMP_LT,
		"JNE": JUMP_GT | JUMP_LT,
		"JLE": JUMP_EQ | JUMP_LT,
		"JMP": JUMP_GT | JUMP_EQ | JUMP_LT,
	}

	destStrToBits = map[string]int{
		"":    0,
		"M":   DEST_M,
		"D":   DEST_D,
		"MD":  DEST_M | DEST_D,
		"A":   DEST_A,
		"AM":  DEST_A | DEST_M,
		"AD":  DEST_A | DEST_D,
		"AMD": DEST_A | DEST_M | DEST_D,
	}
)

type LineError struct {
	msg       string
	inputLine InputLine
}

func (e *LineError) Error() string {
	return fmt.Sprintf("%d: %s", e.inputLine.lineNum, e.msg)
}

func newLineError(inputLine *InputLine, msg string) *LineError {
	return &LineError{
		msg:       msg,
		inputLine: *inputLine,
	}
}

type InputLine struct {
	str     string
	lineNum int
}

func newInputLine(str string, lineNum int) *InputLine {
	return &InputLine{
		str:     str,
		lineNum: lineNum,
	}
}

type Instruction interface {
	LineNum() int
}

type baseInstruction struct {
	lineNum int
}

func (i *baseInstruction) LineNum() int {
	return i.lineNum
}

type LabelInstruction struct {
	baseInstruction
	LabelName string
}

func parseLabelInstruction(line *InputLine) (Instruction, error) {
	if !strings.HasPrefix(line.str, "(") || !strings.HasSuffix(line.str, ")") {
		return nil, nil
	}

	label := strings.TrimSuffix(line.str, ")")
	label = strings.TrimPrefix(label, "(")

	if !labelRegexp.MatchString(label) {
		return nil, newLineError(line, "invalid label")
	}

	return &LabelInstruction{
		baseInstruction: baseInstruction{lineNum: line.lineNum},
		LabelName:       label,
	}, nil
}

type AInstruction struct {
	baseInstruction
	Address uint16
	Label   string
}

func parseAInstruction(line *InputLine) (Instruction, error) {
	if !strings.HasPrefix(line.str, "@") {
		return nil, nil
	}

	val := strings.TrimPrefix(line.str, "@")
	if labelRegexp.MatchString(val) {
		return &AInstruction{
			Label: val,
		}, nil
	}

	addr, err := strconv.ParseUint(val, 10, 16)
	if err != nil {
		return nil, newLineError(line, "invalid address")
	}

	return &AInstruction{
		baseInstruction: baseInstruction{lineNum: line.lineNum},
		Address:         uint16(addr),
	}, nil
}

type CInstruction struct {
	baseInstruction
	DestBits int // DEST_* for d1 d2 d3
	CompBits int // COMP_* for a c1 .. c6
	JumpBits int // JUMP_* for j1 j2 j3
}

func parseCInstruction(line *InputLine) (Instruction, error) {
	parts := cInstrRegexp.FindStringSubmatch(line.str)
	if parts == nil {
		return nil, nil
	}

	destStr := parts[1]
	compStr := parts[2]
	jumpStr := parts[3]

	instr := &CInstruction{
		baseInstruction: baseInstruction{lineNum: line.lineNum},
	}
	ok := false

	if instr.DestBits, ok = destStrToBits[destStr]; !ok {
		return nil, newLineError(line, fmt.Sprintf("invalid dest '%v'", destStr))
	}

	if instr.CompBits, ok = compStrToBits[compStr]; !ok {
		return nil, newLineError(line, fmt.Sprintf("invalid comp '%v'", compStr))
	}

	if instr.JumpBits, ok = jumpStrToBits[jumpStr]; !ok {
		return nil, newLineError(line, fmt.Sprintf("invalid jump '%v'", jumpStr))
	}

	return instr, nil
}

type parserFunc func(line *InputLine) (Instruction, error)

var (
	parserFuncs = []parserFunc{
		parseLabelInstruction,
		parseAInstruction,
		parseCInstruction,
	}
)

type Parser interface {
	NextInstruction() (Instruction, error)
}

type parserImpl struct {
	scanner *bufio.Scanner
	lineNum int
}

func NewParser(reader io.Reader) Parser {
	return &parserImpl{
		scanner: bufio.NewScanner(reader),
		lineNum: 0,
	}
}

func trimComment(str string) string {
	return strings.SplitN(str, "//", 2)[0]
}

func (p *parserImpl) readLine() (*InputLine, error) {
	if !p.scanner.Scan() {
		if err := p.scanner.Err(); err != nil {
			return nil, err
		}
		return nil, io.EOF
	}

	p.lineNum++
	return newInputLine(p.scanner.Text(), p.lineNum), nil
}

func (p *parserImpl) NextInstruction() (Instruction, error) {
	for {
		line, err := p.readLine()
		if err != nil {
			return nil, err
		}

		line.str = strings.TrimSpace(trimComment(line.str))
		if line.str == "" {
			continue
		}

		for _, f := range parserFuncs {
			if instr, err := f(line); instr != nil || err != nil {
				return instr, err
			}
		}

		return nil, newLineError(line, "unknown instruction type")
	}
}
