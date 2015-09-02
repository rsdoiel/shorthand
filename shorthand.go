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
const Version = "v0.0.3"

//
// An Op is built from a multi character glyph
// Each element in the glyph has meaning
// " :" is the start of a glyb and the end is a trailing space
// = the source value is next, this is basic assignment of a string value to a symbol
// < is input from a file
// { expand (resolved label values)
// } assign a statement (i.e. label, op, value)
// ! input from a shell expression
// [ is a markdown expansion
// > is to output an to a file
// * operate on symbol table

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
	OutputAssignedExpansion  string = " :>: "
	OutputAssignedExpansions string = " :*>: "
	OutputAssignment         string = " :}>: "
	OutputAssignments        string = " :*}>: "
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
	OutputAssignedExpansion,
	OutputAssignedExpansions,
	OutputAssignment,
	OutputAssignments,
}

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

// warning writes a message to stderr
func warning(msg string) {
	fmt.Fprintf(os.Stderr, "%s\n", msg)
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
			return SourceMap{Label: parts[0], Op: op, Value: parts[1], LineNo: lineNo, Expanded: ""}, true
		}
	}
	return SourceMap{Label: "", Op: "", Value: "", LineNo: lineNo, Expanded: s}, false
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
			fmt.Fprintf(fp, "%s%s%s", sm.Label, sm.Op, sm.Value)
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
			fmt.Fprintf(fp, "%s%s%s\n", sm.Label, sm.Op, sm.Value)
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
	for i, text := range strings.Split(string(buf), "\n") {
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
		return WriteAssignment(sm.Value, sm.Label, table, false)
	}

	if sm.Op == OutputAssignedExpansions {
		return WriteAssignments(sm.Value, table, false)
	}

	if sm.Op == OutputAssignment {
		return WriteAssignment(sm.Value, sm.Label, table, true)
	}

	if sm.Op == OutputAssignments {
		return WriteAssignments(sm.Value, table, true)
	}

	// These functions change the symbol table
	if sm.Op == AssignString {
		sm.Expanded = sm.Value
	} else if sm.Op == AssignInclude {
		buf, err := ioutil.ReadFile(sm.Value)
		if err != nil {
			warning(fmt.Sprintf("Cannot read %s: %s\n", sm.Value, err))
			return false
		}
		sm.Expanded = string(buf)
	} else if sm.Op == IncludeExpansion {
		buf, err := ioutil.ReadFile(sm.Value)
		if err != nil {
			warning(fmt.Sprintf("Cannot read %s: %s\n", sm.Value, err))
			return false
		}
		sm.Expanded = Expand(table, string(buf))
	} else if sm.Op == AssignShell {
		buf, err := exec.Command("bash", "-c", sm.Value).Output()
		if err != nil {
			warning(fmt.Sprintf("Shell command returned error: %s\n", err))
			return false
		}
		sm.Expanded = string(buf)
	} else if sm.Op == AssignExpandShell {
		buf, err := exec.Command("bash", "-c", Expand(table, sm.Value)).Output()
		if err != nil {
			warning(fmt.Sprintf("Shell command returned error: %s\n", err))
			return false
		}
		sm.Expanded = string(buf)
	} else if sm.Op == AssignExpansion {
		sm.Expanded = Expand(table, sm.Value)
	} else if sm.Op == AssignExpandExpansion {
		tmp := Expand(table, sm.Value)
		sm.Expanded = Expand(table, tmp)
	} else if sm.Op == IncludeAssignments {
		err := ReadAssignments(sm.Value, table)
		if err != nil {
			warning(fmt.Sprintf("Error processing %s: %s\n", sm.Value, err))
			return false
		}
	} else if sm.Op == AssignMarkdown {
		sm.Expanded = strings.TrimSpace(string(blackfriday.MarkdownCommon([]byte(sm.Value))))
	} else if sm.Op == AssignExpandMarkdown {
		sm.Expanded = strings.TrimSpace(string(blackfriday.MarkdownCommon([]byte(Expand(table, sm.Value)))))
	} else if sm.Op == IncludeMarkdown {
		sm.Expanded = ReadMarkdown(sm.Value)
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
