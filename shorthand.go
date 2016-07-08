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
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Version nummber of library and utility
const Version = "v0.0.6"

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
	Label    string // Label is the symbol to be replace based on Op and Source
	Op       string // Op is the type of assignment being made (if is an empty string if not an assignment)
	Source   string // Source is argument to the right of Op
	Expanded string // Expanded is the value calculated based on Label, Op and Source
	LineNo   int
}

// SymbolTable holds the exressions, values and other errata of parsing assignments making expansions
type SymbolTable struct {
	entries []SourceMap
	labels  map[string]int
}

// GetSymbol finds the symbol entry and returns the SourceMap
func (st *SymbolTable) GetSymbol(sym string) SourceMap {
	i, ok := st.labels[sym]
	if ok == true {
		return st.entries[i]
	}
	return SourceMap{Label: "", Op: "", Source: "", Expanded: "", LineNo: -1}
}

// GetSymbols returns a list of all symbols defined by labels as an array of SourceMaps
func (st *SymbolTable) GetSymbols() []SourceMap {
	var symbols []SourceMap

	for _, i := range st.labels {
		symbols = append(symbols, st.entries[i])
	}
	return symbols
}

// SetSymbol adds a SourceMap to entries and points the labels at the most recent definition.
func (st *SymbolTable) SetSymbol(sm SourceMap) int {
	st.entries = append(st.entries, sm)
	if st.labels == nil {
		st.labels = make(map[string]int)
	}
	i := len(st.entries) - 1
	st.labels[sm.Label] = i
	return st.labels[sm.Label]
}

// OperatorMap is a map of operator testings and their related functions
type OperatorMap map[string]func(*VirtualMachine, SourceMap) (SourceMap, error)

// VirtualMachine defines the structure which holds symbols, operator map,
// ops and current prompt setting for a shorthand instance.
type VirtualMachine struct {
	prompt    string
	Symbols   *SymbolTable
	Operators OperatorMap
	Ops       []string
	Help      map[string]string
}

// New returns a VirtualMachine struct and registers all Operators
func New() *VirtualMachine {
	vm := new(VirtualMachine)
	vm.Symbols = new(SymbolTable)
	vm.Operators = make(OperatorMap)
	vm.Help = make(map[string]string)

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

	// Register the built-in operators (glyph versions, these should be depreciated at somepoint)
	vm.RegisterOp(" :=: ", AssignString, "Assign a string to label")
	vm.RegisterOp(" :=<: ", AssignInclude, "Include content and assign to label")
	vm.RegisterOp(" :}<: ", ImportAssignments, "Import assignments from a shorthand file")

	vm.RegisterOp(" :{: ", AssignExpansion, "Expand and assign to label")
	vm.RegisterOp(" :{{: ", AssignExpandExpansion, "Expand and expansion and assign to label")
	vm.RegisterOp(" :{<: ", IncludeExpansion, "Include a file, expand and assign to label")

	vm.RegisterOp(" :!: ", AssignShell, "Assign the output of a Bash command to label")
	vm.RegisterOp(" :{!: ", AssignExpandShell, "Expand and then assign the results of a Bash command to label")

	vm.RegisterOp(" :[: ", AssignMarkdown, "Convert markdown and assign to label")
	vm.RegisterOp(" :{[: ", AssignExpandMarkdown, "Expand and convert markdown and assign to label")
	vm.RegisterOp(" :[<: ", IncludeMarkdown, "Include and convert markdown and assign to label")
	vm.RegisterOp(" :{[<: ", IncludeExpandMarkdown, "Include an expansion, convert with Markdown and assign to label")

	vm.RegisterOp(" :>: ", OutputExpansion, "Write an expansion to a file")
	vm.RegisterOp(" :@>: ", OutputExpansions, "Write all expansions to a file (order not guaranteed)")

	vm.RegisterOp(" :}>: ", ExportAssignment, "Export assignment to a file")
	vm.RegisterOp(" :@}>: ", ExportAssignments, "Expand all assignments (order not guaranteed)")

	// Register the built-in operators (readable versions)
	vm.RegisterOp(" :label: ", AssignString, "Assign a string to label")
	vm.RegisterOp(" :import-text: ", AssignInclude, "Include content and assign to label")
	vm.RegisterOp(" :import-shorthand: ", ImportAssignments, "Import assignments from a shorthand file")

	vm.RegisterOp(" :expand: ", AssignExpansion, "Expand and assign to label")
	vm.RegisterOp(" :expand-expansion: ", AssignExpandExpansion, "Expand and expansion and assign to label")
	vm.RegisterOp(" :import-expansion: ", IncludeExpansion, "Include a file, expand and assign to label")

	vm.RegisterOp(" :bash: ", AssignShell, "Assign the output of a Bash command to label")
	vm.RegisterOp(" :expand-and-bash: ", AssignExpandShell, "Expand and then assign the results of a Bash command to label")

	vm.RegisterOp(" :markdown: ", AssignMarkdown, "Convert markdown and assign to label")
	vm.RegisterOp(" :expand-markdown: ", AssignExpandMarkdown, "Expand and convert markdown and assign to label")
	vm.RegisterOp(" :import-markdown: ", IncludeMarkdown, "Include and convert markdown and assign to label")
	vm.RegisterOp(" :import-expand-markdown: ", IncludeExpandMarkdown, "Include an expansion, convert with Markdown and assign to label")

	vm.RegisterOp(" :export-expansion: ", OutputExpansion, "Write an expansion to a file")
	vm.RegisterOp(" :export-all-expansions: ", OutputExpansions, "Write all expansions to a file (order not guaranteed)")

	vm.RegisterOp(" :export-shorthand: ", ExportAssignment, "Export assignment to a file")
	vm.RegisterOp(" :export-all-shorthand: ", ExportAssignments, "Expand all assignments (order not guaranteed)")
	return vm
}

// SetPrompt sets the string value of the prompt for a VirtualMachine instance
func (vm *VirtualMachine) SetPrompt(s string) {
	vm.prompt = s
}

// RegisterOp associate a operation and function
func (vm *VirtualMachine) RegisterOp(op string, callback func(*VirtualMachine, SourceMap) (SourceMap, error), help string) error {
	_, ok := vm.Operators[op]
	if ok == true {
		return fmt.Errorf("Cannot redefine function %s\n", op)
	}
	vm.Operators[op] = callback
	vm.Ops = append(vm.Ops, op)
	vm.Help[op] = help
	return nil
}

// Parse a string, return a source map. Takes advantage of the internal ops list.
// If no valid op is found then return a source map with Label and Op set to an empty string
// while Source is set the the string that was parsed.  Expanded should always be an empty string
// at the parse stage.
func (vm *VirtualMachine) Parse(s string, lineNo int) SourceMap {
	for _, op := range vm.Ops {
		if strings.Contains(s, op) {
			parts := strings.SplitN(strings.TrimSpace(s), op, 2)
			if len(parts) == 2 {
				return SourceMap{Label: parts[0], Op: op, Source: parts[1], LineNo: lineNo, Expanded: ""}
			}
			if len(parts) == 1 {
				return SourceMap{Label: parts[0], Op: op, Source: "", LineNo: lineNo, Expanded: ""}
			}
			return SourceMap{Label: "", Op: op, Source: "", LineNo: lineNo, Expanded: ""}
		}
	}
	return SourceMap{Label: "", Op: "", Source: s, LineNo: lineNo, Expanded: ""}
}

// Expand takes some text and expands all labels to their values
func (vm *VirtualMachine) Expand(text string) string {
	// labels hash should also point at the last known state of
	// the label
	result := text
	symbols := vm.Symbols.GetSymbols()
	for _, sm := range symbols {
		if strings.Contains(text, sm.Label) {
			tmp := strings.Replace(result, sm.Label, sm.Expanded, -1)
			result = tmp
		}
	}
	return result
}

// Eval stores a shorthand assignment or expands and writes the content to stdout
// Returns the expanded  and any error
func (vm *VirtualMachine) Eval(s string, lineNo int) (string, error) {
	sm := vm.Parse(s, lineNo)
	// If not an assignment Expand and return the expansion
	if sm.Label == "" && sm.Op == "" {
		return fmt.Sprintf("%s", vm.Expand(s)), nil
	}

	callback, ok := vm.Operators[sm.Op]
	if ok == false {
		return "", fmt.Errorf("ERROR (%d): %s is not a supported assignment.\n", lineNo, s)
	}

	// Make the associated assignment and save the symbol to the symbol table.
	newSM, err := callback(vm, sm)
	if err != nil {
		return "", err
	}

	vm.Symbols.SetSymbol(newSM)
	return "", nil
}

// Run takes a reader (e.g. os.Stdin), and two writers (e.g. os.Stdout and os.Stderr)
// It reads until EOF, :exit:, or :quit: operation is encountered
// returns the number of lines processed.
func (vm *VirtualMachine) Run(in *bufio.Reader) int {
	lineNo := 0
	for {
		if vm.prompt != "" {
			fmt.Fprint(os.Stdout, vm.prompt)
		}
		src, rErr := in.ReadString('\n')
		if rErr != nil {
			break
		}
		lineNo++
		if strings.Contains(src, ":exit:") || strings.Contains(src, ":quit:") {
			break
		}
		out, err := vm.Eval(src, lineNo)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR (%d): %s\n", lineNo, err)
		}
		if out != "" {
			fmt.Fprint(os.Stdout, out)
		}
	}
	return lineNo
}
