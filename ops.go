//
// Shorthand package operators - assign a function with the func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) signature
// and use RegisterOp (e.g. in the New() function) to add support to Shorthand.
//
package shorthand

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	// 3rd party packages
	"github.com/russross/blackfriday"
)

//AssignString take the Source and copy to Expanded
var AssignString = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	expanded := sm.Source
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded, LineNo: sm.LineNo}, nil
}

//AssignInclude read a file using Source as filename and put the results in Expanded
var AssignInclude = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	buf, err := ioutil.ReadFile(sm.Source)
	if err != nil {
		return SourceMap{Label: "", Op: ":exit:", Source: "", Expanded: "", LineNo: sm.LineNo},
			fmt.Errorf("Cannot read %s: %s\n", sm.Source, err)
	}
	expanded := string(buf)
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded, LineNo: sm.LineNo}, nil
}

// ImportAssignments evaluates the file for assignment operations
var ImportAssignments = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	var output []string
	buf, err := ioutil.ReadFile(sm.Source)
	if err != nil {
		return SourceMap{Label: "", Op: ":exit:", Source: "", Expanded: "", LineNo: sm.LineNo}, err
	}
	lineNo := 1
	for _, src := range strings.Split(string(buf), "\n") {
		s, err := vm.Eval(src, lineNo)
		if err != nil {
			return SourceMap{Label: "", Op: ":exit:", Source: "", Expanded: "", LineNo: sm.LineNo},
				fmt.Errorf("ERROR (%s %d): %s", sm.Source, lineNo, err)
		}
		if s != "" {
			output = append(output, s)
		}
	}
	expanded := strings.Join(output, "\n")
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded, LineNo: sm.LineNo}, nil
}

// AssignExpanshion expands Source and copy to Expanded
var AssignExpansion = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	expanded := vm.Expand(sm.Source)
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded, LineNo: sm.LineNo}, nil
}

// AssignExpandExpansions expand an expanded Source and copy to Expanded
var AssignExpandExpansion = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	tmp := vm.Expand(sm.Source)
	expanded := vm.Expand(tmp)
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded, LineNo: sm.LineNo}, nil
}

// IncludeExpansion include the filename from Source, expand and copy to Expanded
var IncludeExpansion = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	buf, err := ioutil.ReadFile(sm.Source)
	if err != nil {
		return sm, err
	}
	expanded := vm.Expand(string(buf))
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded, LineNo: sm.LineNo}, nil
}

// AssignShell pass Source to shell and copy stdout to Expanded
var AssignShell = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	buf, err := exec.Command("bash", "-c", sm.Source).Output()
	if err != nil {
		return sm, err
	}
	expanded := string(buf)
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded, LineNo: sm.LineNo}, nil
}

// AssignExpandShell expand Source, pass to Bash and assign output to Expanded
var AssignExpandShell = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	buf, err := exec.Command("bash", "-c", vm.Expand(sm.Source)).Output()
	if err != nil {
		return sm, err
	}
	expanded := string(buf)
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded, LineNo: sm.LineNo}, nil
}

// AssignMarkdown process Source with Blackfriday and copy
var AssignMarkdown = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	expanded := string(blackfriday.MarkdownCommon([]byte(sm.Source)))
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded, LineNo: sm.LineNo}, nil
}

var AssignExpandMarkdown = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	tmp := vm.Expand(sm.Source)
	expanded := string(blackfriday.MarkdownCommon([]byte(tmp)))
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded, LineNo: sm.LineNo}, nil
}

var IncludeMarkdown = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	buf, err := ioutil.ReadFile(sm.Source)
	if err != nil {
		return sm, err
	}
	tmp := string(buf)
	expanded := string(blackfriday.MarkdownCommon([]byte(tmp)))
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded, LineNo: sm.LineNo}, nil
}

var IncludeExpandMarkdown = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	buf, err := ioutil.ReadFile(sm.Source)
	if err != nil {
		return sm, err
	}
	tmp := vm.Expand(string(buf))
	expanded := string(blackfriday.MarkdownCommon([]byte(tmp)))
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded, LineNo: sm.LineNo}, nil
}

var OutputExpansion = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	oSM := vm.Symbols.GetSymbol(sm.Label)
	out := oSM.Expanded
	fname := sm.Source
	err := ioutil.WriteFile(fname, []byte(out), 0666)
	if err != nil {
		return sm, fmt.Errorf("%d Write error %s: %s", sm.LineNo, fname, err)
	}
	return oSM, nil
}

var OutputExpansions = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	fp, err := os.Create(sm.Source)
	if err != nil {
		return sm, fmt.Errorf("%d Create error %s: %s", sm.LineNo, sm.Source, err)
	}
	defer fp.Close()
	symbols := vm.Symbols.GetSymbols()
	for _, oSM := range symbols {
		fmt.Fprintln(fp, vm.Expand(oSM.Expanded))
	}
	return sm, nil
}

var ExportAssignment = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	oSM := vm.Symbols.GetSymbol(sm.Label)
	out := fmt.Sprintf("%s%s%s", oSM.Label, oSM.Op, oSM.Source)
	fname := sm.Source
	err := ioutil.WriteFile(fname, []byte(out), 0666)
	if err != nil {
		return sm, fmt.Errorf("%d Write error %s: %s", sm.LineNo, fname, err)
	}
	return oSM, nil
}

var ExportAssignments = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	fp, err := os.Create(sm.Source)
	if err != nil {
		return sm, fmt.Errorf("%d Create error %s: %s", sm.LineNo, sm.Source, err)
	}
	defer fp.Close()
	symbols := vm.Symbols.GetSymbols()
	for _, oSM := range symbols {
		fmt.Fprintf(fp, "%s%s%s\n", oSM.Label, oSM.Op, oSM.Source)
	}
	return sm, nil
}
