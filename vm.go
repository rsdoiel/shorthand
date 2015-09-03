//
// Package shorthand provides shorthand definition and expansion.
//
// shorthand.go - A simple definition and expansion notation to use
// as shorthand when a template language is too much.
//
// @author R. S. Doiel, <rsdoiel@gmail.com>
// copyright (c) 2015 all rights reserved.
// Released under the BSD 2-Clause license
// See: http://opensource.org/licenses/BSD-2-Clause
//
package shorthand

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/russross/blackfriday"
)

// The version nummber of library and utility
const Version = "v0.0.4"

//
// An Op is built from a multi character symbol
// Each element in the symbol has meaning
// " :" is the start of a glyph band the end ": " is a colon and trailing space
// = the source value is next, this is basic assignment of a string value to a symbol
// < is input from a file
// { expand (resolved label values)
// } assign a statement (i.e. label, op, value)
// ! input from a shell expression
// [ is a markdown expansion
// > is to output an to a file
// @ operate on whole symbol table
//

// SourceMap holds the source and value of an assignment
type SourceMap struct {
	Label    string // Label is the symbol to be replace based on Op and Value
	Op       string // Op is the type of assignment being made (if is an empty string if not an assignment)
	Value    string // Value is argument to the right of Op
	Expanded string // Expanded is the value calculated based on Label, Op and Value
	LineNo   int
}

// SymbolTable holds the exressions, values and other errata of parsing assignments making expansions
type SymbolTable struct {
	entries []SourceMap
	labels  map[string]int
}

type FunctionMap map[string]func(SourceMap) string

func New() *FunctionTable {
	fm := make(FunctionTable)
	return fm
}

func (fm *FunctionTable) RegisterOp(op string, run func(SourceMap) string) {
	_, ok := fm.operation[op]
	if ok == true {
		log.Fatalf("Cannot redefine function %s\n", op)
	}
	fm[op] = run
}

// Eval stores a shorthand assignment or expands and writes the content to stdout
func Eval(functionTable *FunctionTable, symbolTable *SymbolTable, s string, lineNo int) bool {
	sm, isAssignment := Parse(s, lineNo)
	if isAssignment == false {
		fmt.Printf("%s", Expand(symbolTable, s))
		return true
	}

	action, ok := functionTable[sm.Op]
	if ok == false {
		fmt.Fprintf(os.Stderr, "ERROR (%d): %s is not an expansion or valid assignment.\n", lineNo, s)
		return false
	}

	newSM, ok := action(sm)
	if ok == false {
		fmt.Fprintf(os.Stderr, "ERROR (%d): %s failed\n", lineNo, s)
		return false
	}

	if newSM.Label != "" {
		symbolTable.entries = append(table.entries, newSM)
		if symbolTable.labels == nil {
			symbolTable.labels = make(map[string]int)
		}
		symbolTable.labels[newSM.Label] = len(table.entries) - 1
	}
	return ok
}
