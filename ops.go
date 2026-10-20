// Copyright (C) 2015-2026 R. S. Doiel
// SPDX-License-Identifier: AGPL-3.0-or-later
// See the LICENSE file.

// Package shorthand provides shorthand definition and expansion.
//
// shorthand.go - A simple definition and expansion notation to use
// as shorthand when a template language is too much.
//
// @author R. S. Doiel, <rsdoiel@gmail.com>
//
// operators - assign a function with the func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) signature
// and use RegisterOp (e.g. in the New() function) to add support to Shorthand.
package shorthand

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

/** AssignString is the :set: operator: it copies Source to Expanded.
 *
 * Parameters:
 *   vm (*VirtualMachine) — the machine holding the symbols
 *   sm (SourceMap) — the parsed line: Label, Op and Source
 *
 * Returns:
 *   SourceMap — the same line with Expanded filled in
 *   error — why the operator failed, e.g. a file that cannot be read
 *
 * Example:
 *   vm.RegisterOp(":set:", shorthand.AssignString, "Assign a string to label")
 */
var AssignString = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	expanded := sm.Source
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded, LineNo: sm.LineNo}, nil
}

/** AssignInclude is the :import-text: operator: it reads the file named by Source into Expanded.
 *
 * Parameters:
 *   vm (*VirtualMachine) — the machine holding the symbols
 *   sm (SourceMap) — the parsed line: Label, Op and Source
 *
 * Returns:
 *   SourceMap — the same line with Expanded filled in
 *   error — why the operator failed, e.g. a file that cannot be read
 *
 * Example:
 *   vm.Eval(":import-text: {{nav}} nav.html", 1)
 */
var AssignInclude = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	buf, err := os.ReadFile(sm.Source)
	if err != nil {
		return SourceMap{Label: "", Op: ":exit:", Source: "", Expanded: "", LineNo: sm.LineNo},
			fmt.Errorf("Cannot read %s: %s\n", sm.Source, err)
	}
	expanded := string(buf)
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded, LineNo: sm.LineNo}, nil
}

/** ImportAssignments is the :import-shorthand: operator: it evaluates the file named by
 * Source for assignments, so its definitions are available afterward.
 *
 * Parameters:
 *   vm (*VirtualMachine) — the machine holding the symbols
 *   sm (SourceMap) — the parsed line: Label, Op and Source
 *
 * Returns:
 *   SourceMap — the same line with Expanded filled in
 *   error — why the operator failed, e.g. a file that cannot be read
 *
 * Example:
 *   vm.Eval(":import-shorthand: _ definitions.shorthand", 1)
 */
var ImportAssignments = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	var output []string
	buf, err := os.ReadFile(sm.Source)
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

/** AssignExpansion is the :expand: operator: it expands the labels in Source and
 * copies the result to Expanded.
 *
 * Parameters:
 *   vm (*VirtualMachine) — the machine holding the symbols
 *   sm (SourceMap) — the parsed line: Label, Op and Source
 *
 * Returns:
 *   SourceMap — the same line with Expanded filled in
 *   error — why the operator failed, e.g. a file that cannot be read
 *
 * Example:
 *   vm.Eval(":expand: {{title}} Report for @date", 1)
 */
var AssignExpansion = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	expanded := vm.Expand(sm.Source)
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded, LineNo: sm.LineNo}, nil
}

/** AssignExpandExpansion is the :expand-expansion: operator: it expands Source twice,
 * for a label whose value contains another label.
 *
 * Parameters:
 *   vm (*VirtualMachine) — the machine holding the symbols
 *   sm (SourceMap) — the parsed line: Label, Op and Source
 *
 * Returns:
 *   SourceMap — the same line with Expanded filled in
 *   error — why the operator failed, e.g. a file that cannot be read
 *
 * Example:
 *   vm.Eval(":expand-expansion: {{heading}} @title", 1)
 */
var AssignExpandExpansion = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	tmp := vm.Expand(sm.Source)
	expanded := vm.Expand(tmp)
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded, LineNo: sm.LineNo}, nil
}

/** IncludeExpansion is the :import: operator: it reads the file named by Source,
 * evaluates the shorthand in it and copies the result to Expanded.
 *
 * Parameters:
 *   vm (*VirtualMachine) — the machine holding the symbols
 *   sm (SourceMap) — the parsed line: Label, Op and Source
 *
 * Returns:
 *   SourceMap — the same line with Expanded filled in
 *   error — why the operator failed, e.g. a file that cannot be read
 *
 * Example:
 *   vm.Eval(":import: {{page}} page.shorthand", 1)
 */
var IncludeExpansion = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	buf, err := os.ReadFile(sm.Source)
	if err != nil {
		return sm, err
	}
	expanded := vm.Expand(string(buf))
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded, LineNo: sm.LineNo}, nil
}

/** AssignShell is the :bash: operator: it runs Source with Bash and copies its
 * standard output to Expanded.
 *
 * Parameters:
 *   vm (*VirtualMachine) — the machine holding the symbols
 *   sm (SourceMap) — the parsed line: Label, Op and Source
 *
 * Returns:
 *   SourceMap — the same line with Expanded filled in
 *   error — why the operator failed, e.g. a file that cannot be read
 *
 * Example:
 *   vm.Eval(":bash: @now date", 1)
 */
var AssignShell = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	buf, err := exec.Command("bash", "-c", sm.Source).Output()
	if err != nil {
		return sm, err
	}
	expanded := string(buf)
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded, LineNo: sm.LineNo}, nil
}

/** AssignExpandShell is the :expand-and-bash: operator: it expands Source, runs the
 * result with Bash and copies its output to Expanded.
 *
 * Parameters:
 *   vm (*VirtualMachine) — the machine holding the symbols
 *   sm (SourceMap) — the parsed line: Label, Op and Source
 *
 * Returns:
 *   SourceMap — the same line with Expanded filled in
 *   error — why the operator failed, e.g. a file that cannot be read
 *
 * Example:
 *   vm.Eval(":expand-and-bash: {{entry}} cat @filename", 1)
 */
var AssignExpandShell = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	buf, err := exec.Command("bash", "-c", vm.Expand(sm.Source)).Output()
	if err != nil {
		return sm, err
	}
	expanded := string(buf)
	return SourceMap{Label: sm.Label, Op: sm.Op, Source: sm.Source, Expanded: expanded, LineNo: sm.LineNo}, nil
}

/** OutputExpansion is the :export: operator: it writes the expansion of Label to the
 * file named by Source.
 *
 * Parameters:
 *   vm (*VirtualMachine) — the machine holding the symbols
 *   sm (SourceMap) — the parsed line: Label, Op and Source
 *
 * Returns:
 *   SourceMap — the same line with Expanded filled in
 *   error — why the operator failed, e.g. a file that cannot be read
 *
 * Example:
 *   vm.Eval(":export: {{content}} content.txt", 1)
 */
var OutputExpansion = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	oSM := vm.Symbols.GetSymbol(sm.Label)
	out := oSM.Expanded
	fname := sm.Source
	err := os.WriteFile(fname, []byte(out), 0666)
	if err != nil {
		return sm, fmt.Errorf("%d Write error %s: %s", sm.LineNo, fname, err)
	}
	return oSM, nil
}

/** OutputExpansions is the :export-all: operator: it writes every expansion to the
 * file named by Source, in no particular order.
 *
 * Parameters:
 *   vm (*VirtualMachine) — the machine holding the symbols
 *   sm (SourceMap) — the parsed line: Label, Op and Source
 *
 * Returns:
 *   SourceMap — the same line with Expanded filled in
 *   error — why the operator failed, e.g. a file that cannot be read
 *
 * Example:
 *   vm.Eval(":export-all: _ contents.txt", 1)
 */
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

/** ExportAssignment is the :export-shorthand: operator: it writes the assignment of
 * Label to the file named by Source.
 *
 * Parameters:
 *   vm (*VirtualMachine) — the machine holding the symbols
 *   sm (SourceMap) — the parsed line: Label, Op and Source
 *
 * Returns:
 *   SourceMap — the same line with Expanded filled in
 *   error — why the operator failed, e.g. a file that cannot be read
 *
 * Example:
 *   vm.Eval(":export-shorthand: {{content}} content.shorthand", 1)
 */
var ExportAssignment = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	oSM := vm.Symbols.GetSymbol(sm.Label)
	out := fmt.Sprintf("%s %s %s", oSM.Op, oSM.Label, oSM.Source)
	fname := sm.Source
	err := os.WriteFile(fname, []byte(out), 0666)
	if err != nil {
		return sm, fmt.Errorf("%d Write error %s: %s", sm.LineNo, fname, err)
	}
	return oSM, nil
}

/** ExportAssignments is the :export-all-shorthand: operator: it writes every
 * assignment to the file named by Source, in no particular order.
 *
 * Parameters:
 *   vm (*VirtualMachine) — the machine holding the symbols
 *   sm (SourceMap) — the parsed line: Label, Op and Source
 *
 * Returns:
 *   SourceMap — the same line with Expanded filled in
 *   error — why the operator failed, e.g. a file that cannot be read
 *
 * Example:
 *   vm.Eval(":export-all-shorthand: _ contents.shorthand", 1)
 */
var ExportAssignments = func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
	fp, err := os.Create(sm.Source)
	if err != nil {
		return sm, fmt.Errorf("%d Create error %s: %s", sm.LineNo, sm.Source, err)
	}
	defer fp.Close()
	symbols := vm.Symbols.GetSymbols()
	for _, oSM := range symbols {
		fmt.Fprintf(fp, "%s %s %s\n", oSM.Op, oSM.Label, oSM.Source)
	}
	return sm, nil
}
