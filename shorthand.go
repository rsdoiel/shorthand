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
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

// The version nummber of library and utility
const Version = "v0.0.2"

// Assignment Ops
const (
	AssignString         string = " := "
	AssignInclude        string = " :< "
	AssignShell          string = " :! "
	AssignExpansion      string = " :{ "
	IncludeAssignments   string = " :={ "
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
	OutputAssignedValue,
	OutputAssignedValues,
	OutputAssignment,
	OutputAssignments,
}

// SourceMap holds the source and value of an assignment
type SourceMap struct {
	src   string
	value string
}

// Abbrevations holds the shorthand and translation
var Abbreviations = make(map[string]SourceMap)

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
func HasAssignment(key string) bool {
	_, ok := Abbreviations[key]
	return ok
}

// ParseDefinition returns LABEL, ASSIGNMENT, VALUE
func ParseDefinition(s string) (string, string, string) {
	for _, op := range ops {
		if strings.Index(s, op) != -1 {
			parts := strings.SplitN(strings.TrimSpace(s), op, 2)
			return parts[0], op, parts[1]
		}
	}
	return "", "", ""
}

// WriteAssignment writes a single assignment statement to filename
func WriteAssignment(fname string, sm SourceMap, writeSourceCode bool) bool {
	fp, err := os.Create(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	if writeSourceCode == true {
		fmt.Fprintf(fp, "%s", sm.src)
	} else {
		fmt.Fprintf(fp, "%s", sm.value)
	}
	return true
}

// WriteAssignments write all assignment statements to filename
func WriteAssignments(fname string, assigned map[string]SourceMap, writeSourceCode bool) bool {
	fp, err := os.Create(fname)
	if err != nil {
		log.Printf("Cannot write to %s, error: %s\n", fname, err)
		return false
	}
	defer fp.Close()
	for k := range assigned {
		//FIXME: Probably should not write out the value of _
		if writeSourceCode == true {
			fmt.Fprintf(fp, "%s\n", Abbreviations[k].src)
		} else {
			fmt.Fprintf(fp, "%s\n", Abbreviations[k].value)
		}
	}
	return true
}

// ReadAssignments read in all of the lines of fname and
// add any assignment statements found to Abbreviations
// and expand any assignments and return as a string.
func ReadAssignments(fname string) (string, error) {
	var (
		out []string
	)
	buf, err := ioutil.ReadFile(fname)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Cannot read %s: %s\n", fname, err))
	}
	for i, line := range strings.Split(string(buf), "\n") {
		if IsAssignment(line) {
			ok := Assign(line)
			if ok == false {
				return "", errors.New(fmt.Sprintf("Error at line %d in %s\n", i+1, fname))
			}
		} else {
			out = append(out, Expand(line))
		}
	}
	return strings.Join(out, "\n"), nil
}

// Assign stores a shorthand and its expansion or writes and assignment or assignments
func Assign(s string) bool {
	var (
		key   string
		op    string
		value string
	)

	if IsAssignment(s) {
		key, op, value = ParseDefinition(s)
	} else {
		log.Fatalf("[%s] is an invalid assignment.\n", s)
	}

	if op == OutputAssignedValue {
		return WriteAssignment(value, Abbreviations[key], false)
	} else if op == OutputAssignedValues {
		return WriteAssignments(value, Abbreviations, false)
	} else if op == OutputAssignment {
		return WriteAssignment(value, Abbreviations[key], true)
	} else if op == OutputAssignments {
		return WriteAssignments(value, Abbreviations, true)
	} else if op == AssignExpansion {
		value = Expand(value)
	} else if op == IncludeAssignments {
		value, err := ReadAssignments(value)
		if err != nil {
			log.Fatalf("Error processing %s: %s\n", value, err)
		}
	} else if op == AssignInclude {
		buf, err := ioutil.ReadFile(value)
		if err != nil {
			log.Fatalf("Cannot read %s: %v\n", value, err)
		}
		value = string(buf)
	} else if op == AssignShell {
		buf, err := exec.Command("bash", "-c", value).Output()
		if err != nil {
			log.Fatalf("Shell command returned error: %s\n", err)
		}
		value = string(buf)
	}

	if op == "" || key == "" || value == "" {
		return false
	}
	if key == "_" {
		return true
	}
	Abbreviations[key] = SourceMap{src: s, value: value}
	_, ok := Abbreviations[key]
	return ok
}

// Expand takes a text and expands all shorthands
func Expand(text string) string {
	// Iterate through the list of key/SourceMaps in abbreviations
	for key, sm := range Abbreviations {
		text = strings.Replace(text, key, sm.value, -1)
	}
	//fmt.Printf("DEBUG Expand(%s)\n", text)
	return text
}

// Clear remove all the elements of a map.
func Clear() {
	for key := range Abbreviations {
		delete(Abbreviations, key)
	}
}
