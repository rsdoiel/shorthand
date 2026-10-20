// Copyright (C) 2015-2026 R. S. Doiel
// SPDX-License-Identifier: AGPL-3.0-or-later
// See the LICENSE file.

// Package shorthand provides shorthand definition and expansion.
//
// shorthand.go - A simple definition and expansion notation to use
// as shorthand when a template language is too much.
//
// @author R. S. Doiel, <rsdoiel@gmail.com>
package shorthand

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

/** HowItWorks is the help text describing the operators of shorthand, shown by
 * shorthand -examples.
 *
 * Example:
 *   fmt.Print(shorthand.HowItWorks)
 */
var (
	HowItWorks = `

ASSIGNMENTS AND EXPANSIONS

Shorthand is a simple label expansion utility. It is based on a simple key value substitution.  It supports this following types of definitions

+ Assign a string to a LABEL
+ Assign the contents of a file to a LABEL
+ Assign the output of a Bash shell expression to a LABEL
+ Assign the output of a shorthand expansion to a LABEL
+ Read a file of shorthand assignments and assign any expansions to the LABEL
+ Output a LABEL value to a file
+ Output all LABEL values to a file
+ Output a LABEL assignment statement to a file
+ Output all assignment statements to a file

*shorthand* replaces the LABEL with the value assigned to it whereever it is encountered in the text being passed. The assignment statement is 
not written to stdout output.

operator                    | meaning                                  | example
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :set:                      | Assign String                            | :set: {{name}} Freda
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :import-text:              | Assign the contents of a file            | :import-text: {{content}} myfile.txt
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :import-shorthand:         | Get assignments from a file              | :import-shorthand: _ myfile.shorthand
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :expand:                   | Assign an expansion                      | :expand: {{reportTitle}} Report: @title for @date
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :expand-expansion:         | Assign expanded expansion                | :expand-expansion: {{reportHeading}} @reportTitle
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :import:                   | Include a file, procesisng the shorthand | :import: {{nav}} mynav.shorthand
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :bash:                     | Assign Shell output                      | :bash: {{date}} date +%Y-%m-%%d
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :expand-and-bash:          | Assign Expand then gete Shell output     | :expand-and-bash: {{entry}} cat header.txt @filename footer.txt
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :export:                   | Output a label's value to a file         | :export: {{content}} content.txt
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :export-all:               | Output all assigned Expansions           | :export-all: _ contents.txt
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :export-shorthand:             | Output Assignment                        | :export-shorthand: {{content}} content.shorthand
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :export-all-shorthand:        | Output all shorthand assignments      | :export-all-shorthand: _ contents.shorthand
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :exit:                     | Exit the shorthand repl                  | :exit:
----------------------------|------------------------------------------|---------------------------------------------------------------------



Notes: Using an underscore as a LABEL means the label will be ignored. There are no guarantees of order when writing values or assignment 
statements to a file.

The spaces following surrounding ":set:", ":import-text:", ":bash:", ":expand:", ":export:", etc. are required.


EXAMPLE

In this example a file containing the text of pre-amble is assigned to the label @PREAMBLE, the time 3:30 is assigned to the label {{NOW}}.

    :import-text: {{PREAMBLE}} /home/me/preamble.text
	:set: {{NOW}} 3:30

    At {{NOW}} I will be reading the {{PREAMBLE}} until everyone falls asleep.


If the file preamble.txt contained the phrase "Hello World" (including the quotes but without any carriage return or line feed) the output 
after processing the shorthand would look like - 

    At 3:30 I will be reading the "Hello World" until everyone falls asleep.

Notice the lines containing the assignments are not included in the output and that no carriage returns or line feeds are added the the 
substituted labels.

+ Assign shorthand expansions to a LABEL
    + :expand: LABEL SHORTHAND_TO_BE_EXPANDED
	+ :expand: @content@ @report_name@ @report_date@
        + this would concatenate report name and date


PROCESSING MARKDOWN PAGES

_shorthand_ is a label expander or light weight macro expande. It
can be combined with programs like *pandoc* as a pre-processor.

In this example we'll build a HTML page with shorthand labels from 
a couple markdown documents. Then we will use the render HTML as a 
template for a blog page entry.

Our markdown file serving as a template will be call 
"post-template.md".  It should contain the outline of the 
structure of the page plus some shorthand labels we'll expand 
later.


    # @blogTitle

    ## @pageTitle

    ### @dateString

    @contentBlocks


For the purposes of this exercise we'll use _shorthand_ as a repl 
and just enter the assignments sequencly.  Also rather than use 
the output of shorthand directly we'll build up the content for the 
page in a label and use shorthand itself to write the final page out.

The steps we'll follow will be to 

1. Read in our markdown file page.md and turn it into an HTML with embedded shorthand labels
2. Assign some values to the labels
3. Expand the labels in the HTML and assign to a new label
4. Write the new label out to are page call "page.html"

Start the repl with this version of the shorthand command:

    shorthand -p "? "

The _-p_ option tells _shorthand_ to use the value "? " as the prompt. When _shorthand_ starts it will display "? " to indicate it is 
ready for an assignment or expansion.

The following assumes you are in the _shorthand_ repl.

Load the mardkown file and transform it into HTML with embedded shorthand labels

	:import: @pageTemplate post-template.md
	:bash: @dateString date
	:set: @blogTitle My Blog
	:set: @pageTitle A Post
	:import: @contentBlock a-post.md
	:export: @output post.md
	:exit:

Then we can run the following command to pipe it through Pandoc.

    pandoc -s -t html post.md > post.html
`
)

/** SourceMap holds the source and value of an assignment.
 *
 * Parameters:
 *   Label (string) — the symbol to replace
 *   Op (string) — the kind of assignment, empty if the line is not one
 *   Source (string) — the argument to the right of Op
 *   Expanded (string) — the value worked out from Label, Op and Source
 *   LineNo (int) — the line the assignment came from
 *
 * Example:
 *   sm := shorthand.SourceMap{Label: "@who", Op: ":set:", Source: "Freda"}
 */
type SourceMap struct {
	Label    string // Label is the symbol to be replace based on Op and Source
	Op       string // Op is the type of assignment being made (if is an empty string if not an assignment)
	Source   string // Source is argument to the right of Op
	Expanded string // Expanded is the value calculated based on Label, Op and Source
	LineNo   int
}

/** SymbolTable holds the assignments made so far, and which one each label
 * currently points at.
 *
 * Example:
 *   table := new(shorthand.SymbolTable)
 *   table.SetSymbol(shorthand.SourceMap{Label: "@who", Expanded: "Freda"})
 */
type SymbolTable struct {
	entries []SourceMap
	labels  map[string]int
}

/** GetSymbol finds a label's current definition.
 *
 * Parameters:
 *   sym (string) — the label
 *
 * Returns:
 *   SourceMap — the definition, or one with LineNo -1 if the label is not defined
 *
 * Example:
 *   sm := table.GetSymbol("@who")
 */
func (st *SymbolTable) GetSymbol(sym string) SourceMap {
	i, ok := st.labels[sym]
	if ok == true {
		return st.entries[i]
	}
	return SourceMap{Label: "", Op: "", Source: "", Expanded: "", LineNo: -1}
}

/** GetSymbols lists the current definition of every label, in no particular order.
 *
 * Returns:
 *   []SourceMap — one entry for each label
 *
 * Example:
 *   for _, sm := range table.GetSymbols() { fmt.Println(sm.Label) }
 */
func (st *SymbolTable) GetSymbols() []SourceMap {
	var symbols []SourceMap

	for _, i := range st.labels {
		symbols = append(symbols, st.entries[i])
	}
	return symbols
}

/** SetSymbol adds a definition and points its label at it, replacing an earlier one.
 *
 * Parameters:
 *   sm (SourceMap) — the definition
 *
 * Returns:
 *   int — the position of the definition in the table
 *
 * Example:
 *   table.SetSymbol(shorthand.SourceMap{Label: "@who", Expanded: "Freda"})
 */
func (st *SymbolTable) SetSymbol(sm SourceMap) int {
	st.entries = append(st.entries, sm)
	if st.labels == nil {
		st.labels = make(map[string]int)
	}
	i := len(st.entries) - 1
	st.labels[sm.Label] = i
	return st.labels[sm.Label]
}

/** OperatorMap maps an operator such as :set: to the function that carries it out.
 *
 * Example:
 *   var ops shorthand.OperatorMap = vm.Operators
 */
type OperatorMap map[string]func(*VirtualMachine, SourceMap) (SourceMap, error)

/** VirtualMachine holds the symbols, the operators and the prompt of one shorthand
 * session. Out and Err are where Run writes; they default to os.Stdout and os.Stderr.
 *
 * Example:
 *   vm := shorthand.New()
 *   vm.Out = &buffer
 */
type VirtualMachine struct {
	prompt    string
	errors    int
	Out       io.Writer // where Run writes expansions and the prompt, os.Stdout if nil
	Err       io.Writer // where Run writes errors, os.Stderr if nil
	Symbols   *SymbolTable
	Operators OperatorMap
	Ops       []string
	Help      map[string]string
}

/** New returns a VirtualMachine with the built-in operators registered, writing to
 * os.Stdout and os.Stderr.
 *
 * Returns:
 *   *VirtualMachine — a machine with no definitions
 *
 * Example:
 *   vm := shorthand.New()
 */
func New() *VirtualMachine {
	vm := new(VirtualMachine)
	vm.Out = os.Stdout
	vm.Err = os.Stderr
	vm.Symbols = new(SymbolTable)
	vm.Operators = make(OperatorMap)
	vm.Help = make(map[string]string)

	// Register the built-in operators (readable versions)
	vm.RegisterOp(":set:", AssignString, "Assign a string to label")
	vm.RegisterOp(":import-text:", AssignInclude, "Include content and assign to label")
	vm.RegisterOp(":import-shorthand:", ImportAssignments, "Import assignments from a shorthand file")

	vm.RegisterOp(":expand:", AssignExpansion, "Expand and assign to label")
	vm.RegisterOp(":expand-expansion:", AssignExpandExpansion, "Expand and expansion and assign to label")
	vm.RegisterOp(":import:", IncludeExpansion, "Include a file, evaluate the shorthand")

	vm.RegisterOp(":bash:", AssignShell, "Assign the output of a Bash command to label")
	vm.RegisterOp(":expand-and-bash:", AssignExpandShell, "Expand and then assign the results of a Bash command to label")

	vm.RegisterOp(":export:", OutputExpansion, "Write an the contents of an label to a file")
	vm.RegisterOp(":export-all:", OutputExpansions, "Write all label contents to a (order not guaranteed)")

	vm.RegisterOp(":export-shorthand:", ExportAssignment, "Export assignment to a file")
	vm.RegisterOp(":export-all-shorthand:", ExportAssignments, "Expand all assignments (order not guaranteed)")

	return vm
}

// out is where Run writes expansions and the prompt.
func (vm *VirtualMachine) out() io.Writer {
	if vm.Out == nil {
		return os.Stdout
	}
	return vm.Out
}

// err is where Run writes errors.
func (vm *VirtualMachine) err() io.Writer {
	if vm.Err == nil {
		return os.Stderr
	}
	return vm.Err
}

/** ErrorCount returns how many lines Run has failed to evaluate, so a caller
 * can exit non-zero after it has processed everything.
 *
 * Returns:
 *   int — the number of errors reported by Run so far
 *
 * Example:
 *   vm.Run(reader)
 *   if vm.ErrorCount() > 0 { os.Exit(65) }
 */
func (vm *VirtualMachine) ErrorCount() int {
	return vm.errors
}

/** SetPrompt sets the prompt Run writes before reading each line. Empty is no prompt.
 *
 * Parameters:
 *   s (string) — the prompt
 *
 * Example:
 *   vm.SetPrompt("=> ")
 */
func (vm *VirtualMachine) SetPrompt(s string) {
	vm.prompt = s
}

/** RegisterOp adds an operator.
 *
 * Parameters:
 *   op (string) — the operator, e.g. :set:
 *   callback (func(*VirtualMachine, SourceMap) (SourceMap, error)) — what carries it out
 *   help (string) — one line for the :help: listing
 *
 * Returns:
 *   error — if the operator is already registered
 *
 * Example:
 *   err := vm.RegisterOp(":upper:", upper, "Assign upper case text")
 */
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

/** Parse splits a line into a SourceMap using the registered operators. A line with no
 * operator comes back with Label and Op empty and Source set to the line.
 *
 * Parameters:
 *   s (string) — the line
 *   lineNo (int) — its line number
 *
 * Returns:
 *   SourceMap — the parsed line, Expanded empty
 *
 * Example:
 *   sm := vm.Parse(":set: @who Freda", 1) // Op ":set:", Label "@who"
 */
func (vm *VirtualMachine) Parse(s string, lineNo int) SourceMap {
	for _, op := range vm.Ops {
		if strings.Contains(s, op) {
			// NOTE: I've changed to a VERB SUBJECT OBJECT from SUBJECT VERB OBJECT form
			parts := strings.SplitN(strings.TrimSpace(s), " ", 3)
			if len(parts) == 3 {
				return SourceMap{Label: parts[1], Op: op, Source: parts[2], LineNo: lineNo, Expanded: ""}
			}
			if len(parts) == 2 {
				return SourceMap{Label: parts[1], Op: op, Source: "", LineNo: lineNo, Expanded: ""}
			}
			return SourceMap{Label: "", Op: op, Source: "", LineNo: lineNo, Expanded: ""}
		}
	}
	return SourceMap{Label: "", Op: "", Source: s, LineNo: lineNo, Expanded: ""}
}

/** Expand replaces every defined label in text with its value.
 *
 * Parameters:
 *   text (string) — the text to expand
 *
 * Returns:
 *   string — the expanded text
 *
 * Example:
 *   vm.Eval(":set: @who Freda", 1)
 *   vm.Expand("Hello @who") // "Hello Freda"
 */
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

/** Eval carries out an assignment line, or expands any other line.
 *
 * Parameters:
 *   s (string) — the line
 *   lineNo (int) — its line number
 *
 * Returns:
 *   string — the expanded text, empty for an assignment
 *   error — why the assignment failed
 *
 * Example:
 *   text, err := vm.Eval("Hello @who", 2)
 */
func (vm *VirtualMachine) Eval(s string, lineNo int) (string, error) {
	sm := vm.Parse(s, lineNo)
	// If not an assignment Expand and return the expansion
	if sm.Label == "" && sm.Op == "" {
		return vm.Expand(s), nil
	}
	return "", vm.EvalSymbol(sm)
}

/** EvalSymbol applies the operator of a SourceMap you built yourself and records
 * the result, like Eval without the parse step.
 *
 * Parameters:
 *   sm (SourceMap) — the assignment
 *
 * Returns:
 *   error — if the operator is not supported or fails
 *
 * Example:
 *   err := vm.EvalSymbol(shorthand.SourceMap{Label: "@who", Op: ":set:", Source: "Freda"})
 */
func (vm *VirtualMachine) EvalSymbol(sm SourceMap) error {
	callback, ok := vm.Operators[sm.Op]
	if ok == false {
		return fmt.Errorf("ERROR (%d): `%s %s %s` is not a supported assignment.\n", sm.LineNo, sm.Op, sm.Label, sm.Source)
	}

	// Make the associated assignment and save the symbol to the symbol table.
	newSM, err := callback(vm, sm)
	if err != nil {
		return err
	}

	vm.Symbols.SetSymbol(newSM)
	return nil
}

/** Apply processes a byte slice like Run does a reader, for embedded use. It stops at
 * a line with :exit: or :quit: in it.
 *
 * Parameters:
 *   src ([]byte) — the shorthand text
 *
 * Returns:
 *   []byte — the expanded text
 *   error — the first line that failed
 *
 * Example:
 *   out, err := vm.Apply([]byte(":set: @who Freda\nHello @who\n"))
 */
func (vm *VirtualMachine) Apply(src []byte) ([]byte, error) {
	vm.SetPrompt("")

	out := []string{}
	for lineNo, line := range strings.Split(string(src), "\n") {
		if strings.Contains(line, ":exit:") || strings.Contains(line, ":quit:") {
			break
		}
		r, err := vm.Eval(line, lineNo)
		if err != nil {
			return nil, fmt.Errorf("line (%d): %s\n", lineNo, err)
		}
		out = append(out, r)
	}
	return []byte(strings.Join(out, "\n")), nil
}

/** Run reads until end of input or a line with :exit: or :quit: in it, writing
 * expansions and the prompt to vm.Out and errors to vm.Err.
 *
 * Parameters:
 *   in (*bufio.Reader) — where to read from
 *
 * Returns:
 *   int — the number of lines read
 *
 * Example:
 *   lines := vm.Run(bufio.NewReader(os.Stdin))
 */
func (vm *VirtualMachine) Run(in *bufio.Reader) int {
	lineNo := 0
	for {
		if vm.prompt != "" {
			fmt.Fprint(vm.out(), vm.prompt)
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
			vm.errors++
			fmt.Fprintf(vm.err(), "ERROR (%d): %s\n", lineNo, err)
		}
		if out != "" {
			fmt.Fprint(vm.out(), out)
		}
	}
	return lineNo
}
