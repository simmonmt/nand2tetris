package hack

import (
	"fmt"
	"sort"
)

var (
	predefinedSymbols = map[string]uint16{
		"R0":     0,
		"R1":     1,
		"R2":     2,
		"R3":     3,
		"R4":     4,
		"R5":     5,
		"R6":     6,
		"R7":     7,
		"R8":     8,
		"R9":     9,
		"R10":    10,
		"R11":    11,
		"R12":    12,
		"R13":    13,
		"R14":    14,
		"R15":    15,
		"SCREEN": 16384,
		"KBD":    24576,
		"SP":     0,
		"LCL":    1,
		"ARG":    2,
		"THIS":   3,
		"THAT":   4,
	}
)

type SymbolTable struct {
	entries     map[string]uint16
	nextSymAddr uint16
}

func (t *SymbolTable) Dump() {
	syms := make([]string, len(t.entries))
	i := 0
	for k, _ := range t.entries {
		syms[i] = k
		i++
	}
	sort.Strings(syms)

	for _, k := range syms {
		fmt.Printf("%10s = %d\n", k, t.entries[k])
	}
}

func (t *SymbolTable) AddLabel(name string, val uint16) error {
	if _, found := t.entries[name]; found {
		return fmt.Errorf("symbol already exists")
	}

	t.entries[name] = val
	return nil
}

func (t *SymbolTable) AddSymbol(name string) (uint16, error) {
	addr := t.nextSymAddr
	t.nextSymAddr++

	err := t.AddLabel(name, addr)
	return addr, err
}

func (t *SymbolTable) Lookup(name string) (uint16, bool) {
	val, found := t.entries[name]
	return val, found
}

func NewSymbolTable() *SymbolTable {
	t := &SymbolTable{
		entries:     map[string]uint16{},
		nextSymAddr: 16,
	}

	for k, v := range predefinedSymbols {
		t.entries[k] = v
	}

	return t
}
