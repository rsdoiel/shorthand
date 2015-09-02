package shorthand

import (
	"fmt"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"os"
	"strings"
)

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
