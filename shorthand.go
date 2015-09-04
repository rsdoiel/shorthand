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
	"github.com/russross/blackfriday"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// The version nummber of library and utility
const Version = "v0.0.4-next"

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

// OperatorMap is a map of operator testings and their related functions
type OperatorMap map[string]func(*VirtualMachine, SourceMap) (SourceMap, error)

// Assignment Ops
const (
	AssignString             string = " :=: "
	AssignInclude            string = " :=<: "
	IncludeAssignments       string = " :}<: "
	AssignExpansion          string = " :{: "
	AssignExpandExpansion    string = " :{{: "
	IncludeExpansion         string = " :{<: "
	AssignShell              string = " :!: "
	AssignExpandShell        string = " :{!: "
	AssignMarkdown           string = " :[: "
	AssignExpandMarkdown     string = " :{[: "
	IncludeMarkdown          string = " :[<: "
	IncludeExpandMarkdown    string = " :{[<: "
	OutputAssignedExpansion  string = " :>: "
	OutputAssignedExpansions string = " :@>: "
	OutputAssignment         string = " :}>: "
	OutputAssignments        string = " :@}>: "
)

var ops = []string{
	AssignString,
	AssignInclude,
	IncludeAssignments,
	AssignExpansion,
	AssignExpandExpansion,
	IncludeExpansion,
	AssignShell,
	AssignExpandShell,
	AssignMarkdown,
	AssignExpandMarkdown,
	IncludeMarkdown,
	IncludeExpandMarkdown,
	OutputAssignedExpansion,
	OutputAssignedExpansions,
	OutputAssignment,
	OutputAssignments,
}

func warning(msgs string) {
	fmt.Fprintln(os.Stderr, msgs)
}

// IsAssignment checks to see if a string contains an assignment (e.g. has a ' := ' in the string.)
func IsAssignment(text string) bool {
	for _, op := range ops {
		if strings.Index(text, op) != -1 {
			return true
		}
	}
	return false
}

// HasAssignment checks to see if a shortcut has already been assigned.
func HasAssignment(table *SymbolTable, label string) bool {
	_, ok := table.labels[label]
	if ok == true {
		return true
	}
	return false
}

// Parse returns SourceMap object and bool with true if assignment or false if not
func Parse(s string, lineNo int) (SourceMap, bool) {
	for _, op := range ops {
		if strings.Index(s, op) != -1 {
			parts := strings.SplitN(strings.TrimSpace(s), op, 2)
			return SourceMap{Label: parts[0], Op: op, Source: parts[1], LineNo: lineNo, Expanded: ""}, true
		}
	}
	return SourceMap{Label: "", Op: "", Source: "", LineNo: lineNo, Expanded: s}, false
}

// Expand takes some text and expands all labels to their values
func Expand(table *SymbolTable, text string) string {
	// labels hash should also point at the last known state of
	// the label
	result := text
	for _, i := range table.labels {
		sm := table.entries[i]
		if strings.Contains(text, sm.Label) {
			tmp := strings.Replace(result, sm.Label, sm.Expanded, -1)
			result = tmp
		}
	}
	return result
}

// Assign stores a shorthand and its expansion or writes and assignment or assignments
func Assign(table *SymbolTable, s string, lineNo int) bool {
	sm, ok := Parse(s, lineNo)
	if ok == false {
		return ok
	}

	// These functions do not change the symbol table
	if sm.Op == OutputAssignedExpansion {
		return WriteAssignment(sm.Source, sm.Label, table, false)
	}

	if sm.Op == OutputAssignedExpansions {
		return WriteAssignments(sm.Source, table, false)
	}

	if sm.Op == OutputAssignment {
		return WriteAssignment(sm.Source, sm.Label, table, true)
	}

	if sm.Op == OutputAssignments {
		return WriteAssignments(sm.Source, table, true)
	}

	// These functions change the symbol table
	if sm.Op == AssignString {
		sm.Expanded = sm.Source
	} else if sm.Op == AssignInclude {
		buf, err := ioutil.ReadFile(sm.Source)
		if err != nil {
			warning(fmt.Sprintf("Cannot read %s: %s\n", sm.Source, err))
			return false
		}
		sm.Expanded = string(buf)
	} else if sm.Op == IncludeExpansion {
		buf, err := ioutil.ReadFile(sm.Source)
		if err != nil {
			warning(fmt.Sprintf("Cannot read %s: %s\n", sm.Source, err))
			return false
		}
		sm.Expanded = Expand(table, string(buf))
	} else if sm.Op == AssignShell {
		buf, err := exec.Command("bash", "-c", sm.Source).Output()
		if err != nil {
			warning(fmt.Sprintf("Shell command returned error: %s\n", err))
			return false
		}
		sm.Expanded = string(buf)
	} else if sm.Op == AssignExpandShell {
		buf, err := exec.Command("bash", "-c", Expand(table, sm.Source)).Output()
		if err != nil {
			warning(fmt.Sprintf("Shell command returned error: %s\n", err))
			return false
		}
		sm.Expanded = string(buf)
	} else if sm.Op == AssignExpansion {
		sm.Expanded = Expand(table, sm.Source)
	} else if sm.Op == AssignExpandExpansion {
		tmp := Expand(table, sm.Source)
		sm.Expanded = Expand(table, tmp)
	} else if sm.Op == IncludeAssignments {
		err := ReadAssignments(sm.Source, table)
		if err != nil {
			warning(fmt.Sprintf("Error processing %s: %s\n", sm.Source, err))
			return false
		}
	} else if sm.Op == AssignMarkdown {
		sm.Expanded = strings.TrimSpace(string(blackfriday.MarkdownCommon([]byte(sm.Source))))
	} else if sm.Op == AssignExpandMarkdown {
		sm.Expanded = strings.TrimSpace(string(blackfriday.MarkdownCommon([]byte(Expand(table, sm.Source)))))
	} else if sm.Op == IncludeMarkdown {
		sm.Expanded = ReadMarkdown(sm.Source)
	} else if sm.Op == IncludeExpandMarkdown {
		buf, err := ioutil.ReadFile(sm.Source)
		if err != nil {
			warning(fmt.Sprintf("Cannot read %s: %s\n", sm.Source, err))
			return false
		}
		sm.Expanded = strings.TrimSpace(string(blackfriday.MarkdownCommon([]byte(Expand(table, string(buf))))))
	}
	if sm.Label != "" {
		table.entries = append(table.entries, sm)
		if table.labels == nil {
			table.labels = make(map[string]int)
		}
		table.labels[sm.Label] = len(table.entries) - 1
	}
	return ok
}

// WriteAssignment writes a single assignment statement to filename
func WriteAssignment(fname string, label string, table *SymbolTable, writeSourceCode bool) bool {
	fp, err := os.Create(fname)
	if err != nil {
		warning(fmt.Sprintf("%s", err))
		return false
	}
	defer fp.Close()

	i, ok := table.labels[label]
	if ok == true {
		sm := table.entries[i]
		if writeSourceCode == true {
			fmt.Fprintf(fp, "%s%s%s", sm.Label, sm.Op, sm.Source)
		} else {
			fmt.Fprintf(fp, "%s", sm.Expanded)
		}
		return true
	}
	return false
}

// WriteAssignments write all assignment statements to filename
func WriteAssignments(fname string, table *SymbolTable, writeSourceCode bool) bool {
	fp, err := os.Create(fname)
	if err != nil {
		warning(fmt.Sprintf("Cannot write to %s, error: %s\n", fname, err))
		return false
	}
	defer fp.Close()
	for _, sm := range table.entries {
		if writeSourceCode == true {
			fmt.Fprintf(fp, "%s%s%s\n", sm.Label, sm.Op, sm.Source)
		} else {
			fmt.Fprintf(fp, "%s\n", sm.Expanded)
		}
	}
	return true
}

// ReadAssignments read in all of the lines of fname and
// add any assignment statements found to Abbreviations
// and expand any assignments and return as a string.
func ReadAssignments(fname string, table *SymbolTable) error {
	buf, err := ioutil.ReadFile(fname)
	if err != nil {
		return fmt.Errorf("Cannot read %s: %s\n", fname, err)
	}
	lines := strings.Split(string(buf), "\n")
	for i, text := range lines {
		if IsAssignment(text) {
			ok := Assign(table, text, i)
			if ok == false {
				return fmt.Errorf("Error at line %d in %s\n", i+1, fname)
			}
		}
	}
	return nil
}

// ReadMarkdown reads in a file and processes it with Blackfriday MarkdownCommon
func ReadMarkdown(fname string) string {
	buf, err := ioutil.ReadFile(fname)
	if err != nil {
		warning(fmt.Sprintf("Cannot read %s: %s\n", fname, err))
		return ""
	}
	return string(blackfriday.MarkdownCommon(buf))
}

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

type VirtualMachine struct {
	Symbols   *SymbolTable
	Operators OperatorMap
	Ops       []string
}

func New() *VirtualMachine {
	vm := new(VirtualMachine)
	vm.Symbols = new(SymbolTable)
	vm.Operators = make(OperatorMap)
	// Now register the built-in operators
	/*
		vm.RegisterOp(" :=: ", AssignStringCallback)
		vm.RegisterOp(" :=<: ", AssignIncludeCallback)
		vm.RegisterOp(" :}<: ", IncludeAssignmentsCallback)
		vm.RegisterOp(" :{: ", AssignExpansionCallback)
		vm.RegisterOp(" :{{: ", AssignExpandExpansionCallback)
		vm.RegisterOp(" :{<: ", IncludeExpansionCallback)
		vm.RegisterOp(" :!: ", AssignShellCallback)
		vm.RegisterOp(" :{!: ", AssignExpandShellCallback)
		vm.RegisterOp(" :[: ", AssignMarkdownCallback)
		vm.RegisterOp(" :{[: ", AssignExpandMarkdownCallback)
		vm.RegisterOp(" :[<: ", IncludeMarkdownCallback)
		vm.RegisterOp(" :{[<: ", IncludeExpandMarkdownCallback)
		vm.RegisterOp(" :>: ", OutputAssignedExpansionCallback)
		vm.RegisterOp(" :@>: ", OutputAssignedExpansionsCallback)
		vm.RegisterOp(" :}>: ", OutputAssignmentCallback)
		vm.RegisterOp(" :@}>: ", OutputAssignmentsCallback)
		vm.RegisterOp(" :exit: ", ExitShorthand)
		vm.RegisterOp(" :quit: ", ExitShorthand)
	*/
	return vm
}

// RegisterOp associate a operations with a function
func (vm *VirtualMachine) RegisterOp(op string, callback func(*VirtualMachine, SourceMap) (SourceMap, error)) error {
	_, ok := vm.Operators[op]
	if ok == true {
		return fmt.Errorf("Cannot redefine function %s\n", op)
	}
	vm.Operators[op] = callback
	vm.Ops = append(vm.Ops, op)
	return nil
}

// Parse vm method. Takes advantage of the internal ops list.
func (vm *VirtualMachine) Parse(s string, lineNo int) (SourceMap, error) {
	for _, op := range vm.Ops {
		if strings.Index(s, op) != -1 {
			parts := strings.SplitN(strings.TrimSpace(s), op, 2)
			return SourceMap{Label: parts[0], Op: op, Source: parts[1], LineNo: lineNo, Expanded: ""}, nil
		}
	}
	return SourceMap{Label: "", Op: "", Source: "", LineNo: lineNo, Expanded: s}, fmt.Errorf("%s\n", s)
}

// Eval stores a shorthand assignment or expands and writes the content to stdout
// Returns the expanded  and any error
func (vm *VirtualMachine) Eval(s string, lineNo int) (string, error) {
	sm, isAssignment := vm.Parse(s, lineNo)
	if isAssignment != nil {
		return fmt.Sprintf("%s", Expand(vm.Symbols, s)), nil
	}

	callback, ok := vm.Operators[sm.Op]
	if ok == false {
		return "", fmt.Errorf("ERROR (%d): %s is not an expansion or valid assignment.\n", lineNo, s)
	}

	newSM, err := callback(vm, sm)
	if err != nil {
		return "", err
	}

	if newSM.Label != "" {
		vm.Symbols.entries = append(vm.Symbols.entries, newSM)
		if vm.Symbols.labels == nil {
			vm.Symbols.labels = make(map[string]int)
		}
		vm.Symbols.labels[newSM.Label] = len(vm.Symbols.entries) - 1
	}
	return "", nil
}

// Run takes a reader (e.g. os.Stdin), and two writers (e.g. os.Stdout and os.Stderr)
// It reads until EOF, :exit:, or :quit: operation is encountered
// returns the number of lines processed.
func (vm *VirtualMachine) Run(in *bufio.Reader) int {
	lineNo := 0
	for {
		src, rErr := in.ReadString('\n')
		if rErr != nil {
			break
		}
		lineNo += 1
		if strings.Contains(src, ":exit:") || strings.Contains(src, ":quit:") {
			break
		}
		out, err := vm.Eval(src, lineNo)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR (%d): %s\n", lineNo, err)
		}
		if out != "" {
			fmt.Fprintf(os.Stdout, out)
		}
	}
	return lineNo
}
