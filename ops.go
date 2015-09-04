//
// Shorthand package operators - assign a function with the func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) signature
// and use RegisterOp (e.g. in the New() function) to add support to Shorthand.
//
package shorthand

import (
	"fmt"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

//ExitShorthand - call os.Exit() with appropriate value and exit the repl
var ExitShorthand = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	if sm.Source == "" {
		os.Exit(0)
	}
	fmt.Fprintf(os.Stderr, sm.Source)
	os.Exit(1)
	return SourceMap{Label: "", Op: ":exit:", Source: "", Expanded: ""}, nil
}

//AssignStringCallback take the Source and copy to Expanded
var AssignStringCallback = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	expanded := sm.Source
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded}, nil
}

//AssignIncludeCallback read a file using Source as filename and put the results in Expanded
var AssignIncludeCallback = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	buf, err := ioutil.ReadFile(sm.Source)
	if err != nil {
		return SourceMap{Label: "", Op: ":exit:", Source: "", Expanded: ""},
			fmt.Errorf("Cannot read %s: %s\n", sm.Source, err)
	}
	expanded := string(buf)
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded}, nil
}

// IncludeAssignmentsCallback evaluates the file for assignment operations
var IncludeAssignmentsCallback = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	var output []string
	buf, err := ioutil.ReadFile(sm.Source)
	if err != nil {
		return SourceMap{Label: "", Op: ":exit:", Source: "", Expanded: ""}, err
	}
	lineNo := 1
	for _, src := range strings.Split(string(buf), "\n") {
		s, err := vm.Eval(src, lineNo)
		if err != nil {
			return SourceMap{Label: "", Op: ":exit:", Source: "", Expanded: ""},
				fmt.Errorf("ERROR (%s %d): %s", sm.Source, lineNo, err)
		}
		if s != "" {
			output = append(output, s)
		}
	}
	expanded := strings.Join(output, "\n")
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded}, nil
}

// AssignExpanshionCallback expands Source and copy to Expanded
var AssignExpansionCallback = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	expanded := Expand(vm.Symbols, sm.Source)
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded}, nil
}

// AssignExpandExpansionsCallback expand an expanded Source and copy to Expanded
var AssignExpandExpansionCallback = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	tmp := Expand(vm.Symbols, sm.Source)
	expanded := Expand(vm.Symbols, tmp)
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded}, nil
}

// IncludeExpansionCallback include the filename from Source, expand and copy to Expanded
var IncludeExpansionCallback = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	buf, err := ioutil.ReadFile(sm.Source)
	if err != nil {
		return sm, err
	}
	expanded := Expand(vm.Symbols, string(buf))
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded}, nil
}

// AssignShellCallback pass Source to shell and copy stdout to Expanded
var AssignShellCallback = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	buf, err := exec.Command("bash", "-c", sm.Source).Output()
	if err != nil {
		return sm, err
	}
	expanded := string(buf)
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded}, nil
}

var AssignExpandShellCallback = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	buf, err := exec.Command("bash", "-c", Expand(vm.Symbols, sm.Source)).Output()
	if err != nil {
		return sm, err
	}
	expanded := string(buf)
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded}, nil
}

var AssignMarkdownCallback = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	expanded := string(blackfriday.MarkdownCommon([]byte(sm.Source)))
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded}, nil
}

var AssignExpandMarkdownCallback = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	tmp := Expand(vm.Symbols, sm.Source)
	expanded := string(blackfriday.MarkdownCommon([]byte(tmp)))
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded}, nil
}

var IncludeMarkdownCallback = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	buf, err := ioutil.ReadFile(sm.Source)
	if err != nil {
		return sm, err
	}
	tmp := string(buf)
	expanded := string(blackfriday.MarkdownCommon([]byte(tmp)))
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded}, nil
}

var IncludeExpandMarkdownCallback = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	buf, err := ioutil.ReadFile(sm.Source)
	if err != nil {
		return sm, err
	}
	tmp := Expand(vm.Symbols, string(buf))
	expanded := string(blackfriday.MarkdownCommon([]byte(tmp)))
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded}, nil
}

var OutputAssignedExpansionCallback = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	return sm, fmt.Errorf("OutputAssignedExpansionCallback() not implemented.")
}

var OutputAssignedExpansionsCallback = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	return sm, fmt.Errorf("OutputAssignedExpansionsCallback() not implemented.")
}

var OutputAssignmentCallback = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	return sm, fmt.Errorf("OutputAssignmentCallback() not implemented.")
}

var OutputAssignmentsCallback = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	return sm, fmt.Errorf("OutputAssignmentsCallback() not implemented.")
}
