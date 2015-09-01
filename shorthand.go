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
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/russross/blackfriday"
)

// The version nummber of library and utility
const Version = "v0.0.3"

// Assignment Ops
const (
	AssignString         string = " := "
	AssignInclude        string = " :< "
	AssignShell          string = " :! "
	AssignExpansion      string = " :{ "
	IncludeAssignments   string = " :={ "
	AssignMarkdown       string = " :[ "
	IncludeMarkdown      string = " :=[ "
	OutputAssignedValue  string = " :> "
	OutputAssignedValues string = " :=> "
	OutputAssignment     string = " :} "
	OutputAssignments    string = " :=} "
)

var ops = []string{
	AssignString,
	AssignInclude,
	AssignShell,
	AssignExpansion,
	IncludeAssignments,
	AssignMarkdown,
	IncludeMarkdown,
	OutputAssignedValue,
	OutputAssignedValues,
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
func HasAssignment(table SymbolTable, label string) bool {
	for _, sm := range table.entries {
		if sm.Label == label {
			return true
		}

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
func WriteAssignment(fname string, sm SourceMap, writeSourceCode bool) bool {
	fp, err := os.Create(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	if writeSourceCode == true {
		fmt.Fprintf(fp, "%s%s%s", sm.Label, sm.Op, sm.Value)
	} else {
		fmt.Fprintf(fp, "%s", sm.Expanded)
	}
	return true
}

// WriteAssignments write all assignment statements to filename
func WriteAssignments(fname string, table SymbolTable, writeSourceCode bool) bool {
	fp, err := os.Create(fname)
	if err != nil {
		log.Printf("Cannot write to %s, error: %s\n", fname, err)
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
func ReadAssignments(fname string, table SymbolTable) error {
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
		log.Fatalf(fmt.Sprintf("Cannot read %s: %s\n", fname, err))
	}
	return string(blackfriday.MarkdownCommon(buf))
}

// Expand takes some text and expands all labels to their values
func Expand(table SymbolTable, text string) string {
	// Iterate through the list of key/SourceMaps in abbreviations
	for _, sm := range table.entries {
		text = strings.Replace(text, sm.Label, sm.Expanded, -1)
	}
	return text
}

// Assign stores a shorthand and its expansion or writes and assignment or assignments
func Assign(table SymbolTable, s string, lineNo int) bool {
	sm, ok := Parse(s, lineNo)
	if ok == false {
		return ok
	}

	// These functions do not change the symbol table
	if sm.Op == OutputAssignedValue {
		return WriteAssignment(sm.Value, sm, false)
	}

	if sm.Op == OutputAssignedValues {
		return WriteAssignments(sm.Value, table, false)
	}

	if sm.Op == OutputAssignment {
		return WriteAssignment(sm.Value, sm, true)
	}

	if sm.Op == OutputAssignments {
		return WriteAssignments(sm.Value, table, true)
	}

	// These functions change the symbol table
	if sm.Op == AssignExpansion {
		sm.Expanded = Expand(table, sm.Value)
	} else if sm.Op == IncludeAssignments {
		err := ReadAssignments(sm.Value, table)
		if err != nil {
			log.Fatalf("Error processing %s: %s\n", sm.Value, err)
		}
	} else if sm.Op == AssignInclude {
		buf, err := ioutil.ReadFile(sm.Value)
		if err != nil {
			log.Fatalf("Cannot read %s: %s\n", sm.Value, err)
		}
		sm.Expanded = string(buf)
	} else if sm.Op == AssignShell {
		buf, err := exec.Command("bash", "-c", sm.Value).Output()
		if err != nil {
			log.Fatalf("Shell command returned error: %s\n", err)
		}
		sm.Expanded = string(buf)
	} else if sm.Op == AssignMarkdown {
		sm.Expanded = strings.TrimSpace(string(blackfriday.MarkdownCommon([]byte(sm.Value))))
	} else if sm.Op == IncludeMarkdown {
		sm.Expanded = ReadMarkdown(sm.Value)
	}
	if sm.Label != "" {
		table.entries = append(table.entries, sm)
	}
	return ok
}
