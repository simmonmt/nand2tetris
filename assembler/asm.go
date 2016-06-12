package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"simmonmt.org/nand2tetris/assembler/hack"
)

type passInstrHandler func(instr hack.Instruction) error

func lookupOrAddLabel(label string, symTab *hack.SymbolTable) (uint16, error) {
	if addr, found := symTab.Lookup(label); found {
		return addr, nil
	}

	addr, err := symTab.AddSymbol(label)
	return addr, err
}

func firstPass(reader io.Reader, symTab *hack.SymbolTable) error {
	instNum := uint16(0)

	handler := func(instr hack.Instruction) error {
		switch instr := instr.(type) {
		case *hack.LabelInstruction:
			if err := symTab.AddLabel(instr.LabelName, instNum); err != nil {
				return fmt.Errorf("failed to add label %v to symbol table: %v",
					instr.LabelName, err)
			}
		default:
			instNum++
		}

		return nil
	}

	return doPass(reader, handler)
}

func secondPass(reader io.Reader, symTab *hack.SymbolTable, codeGen *hack.CodeGen) error {
	handler := func(instr hack.Instruction) error {
		switch instr := instr.(type) {
		case *hack.AInstruction:
			addr := instr.Address
			if instr.Label != "" {
				var err error
				addr, err = lookupOrAddLabel(instr.Label, symTab)
				if err != nil {
					return fmt.Errorf("%d: %v", instr.LineNum(), err)
				}
			}

			codeGen.EmitAInstruction(addr)
			return nil

		case *hack.CInstruction:
			codeGen.EmitCInstruction(instr.DestBits, instr.CompBits, instr.JumpBits)
			return nil
		}

		return nil
	}

	return doPass(reader, handler)
}

func doPass(reader io.Reader, handler passInstrHandler) error {
	p := hack.NewParser(reader)

	for {
		instr, err := p.NextInstruction()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if err := handler(instr); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("missing source file parameter")
	}

	inputFile := os.Args[1]
	inputBase := path.Base(inputFile)

	fp, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("failed to open %v: %v", inputFile, err)
	}
	defer fp.Close()

	symTab := hack.NewSymbolTable()

	if err := firstPass(fp, symTab); err != nil {
		log.Fatalf("%v: error: %v", inputBase, err)
	}

	if off, err := fp.Seek(0, 0); off != 0 || err != nil {
		log.Fatalf("%v: failed to seek to 0 for pass 2; off now %v err %v", off, err)
	}

	codeGen := &hack.CodeGen{}

	if err := secondPass(fp, symTab, codeGen); err != nil {
		log.Fatalf("%v: error: %v", inputBase, err)
	}
}
