//
// shorthand.go - command line utility to process shorthand definitions
// and render output with the transformed text and without any
// shorthand definitions.
//
// @author R. S. Doiel, <rsdoiel@gmail.com>
// copyright (c) 2015 all rights reserved.
// Released under the BSD 2-Clause license.
// See: http://opensource.org/licenses/BSD-2-Clause
//
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	// my packages
	shorthand "github.com/rsdoiel/shorthand"
)

type expressionList []string

var (
	help       bool
	version    bool
	expression expressionList
	prompt     string
	noprompt   bool
	vm         *shorthand.VirtualMachine
	lineNo     int
)

var usage = func(exit_code int, msg string) {
	var fh = os.Stderr
	if exit_code == 0 {
		fh = os.Stdout
	}
	cmdName := os.Args[0]

	fmt.Fprintf(fh, `%s
USAGE %s [options] [FILES_TO_PROCESS]

%s is a command line utility to expand labels based on their
assigned definitions. The render output is the transformed text 
and without the shorthand definitions themselves. %s reads 
from standard input and writes to standard output. The basic 
definision form is

    LABEL := VALUE

When the label is encountered in the text outside of a definition
the label will be replaced with its value.

To create a shortand for the label "ACME" with the value
"the point at which someone or something is best"
would be done with the following line

    ACME := the point at which someone or something is best

Now each time the shorthand "ACME" is encountered the phrase
"the point at which someone or something is best" will replace it. Thus

    My, ACME, will come

would become

    My, the point at which someone or something is best, will come

These are supported types of assignments (followed by their forms). 

	+ Assign a string to a label

		LABEL := STRING

	+ Assign the contents of a file to a label

		LABEL :< FILENAME

	+ Assign the output of a Bash shell expression to a label

		LABEL :! BASH_EXPRESSION

	+ Assign the contents of a Shorthand expression to a label

		LABEL :{ SHORTHAND_EXPRESSION

	+ Read a file of shorthand assignments and assign any expansions to the label

		LABEL :={ SHORTHAND_FILENAME

	+ Write out the value for a label

		LABEL :> FILENAME

	+ Write out the values for all labels (order is not guaranteed)

		LABEL :=> FILENAME
		
	+ Write out the assignment statement for a label

		LABEL :} FILENAME

	+ Write out all the assignment statements for all labels

		LABEL :=} FILENAME

	+ Render markdown to a label

		LABEL :[ MARKDOWN_EXPRESSION

	+ Render markdown file to a label

		LABEL :=[ FILENAME


If the label is an underscore then it gets ignored.

You can evaluate shorthand expression in the command line much like you do
with the Unix command sed.  But you can also easily embed them in the file
or include from a file.

EXAMPLE

Pass the current date and time as shorthands transform the file "input.txt"
into "output.txt" with shorthands converted.

    %s -e "@now :! date +%%H:%%M" \
	   -e "@today :! date +%%Y-%%m-%%d" < input.txt > output.txt

You can also embed the shorthand definitions directly in your text file.

OPTIONS
`, msg, cmdName, cmdName, cmdName, cmdName)

	flag.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(fh, "\t-%s\t\t%s\n", f.Name, f.Usage)
	})

	fmt.Fprintf(fh, `
copyright (c) 2015 all rights reserved.
Released under the BSD 2-Clause license.
See: http://opensource.org/licenses/BSD-2-Clause
`)
	os.Exit(exit_code)
}

var welcome = `
  Welcome to shorthand the simple label expander and markdown processor.
  Use ':exit:' to quit the repl, ':help:' to get a list of supported operators.

`

var helpShorthand = func(vm *shorthand.VirtualMachine, sm shorthand.SourceMap) (shorthand.SourceMap, error) {
	fmt.Printf(`
The following operators are supported in shorthand:

`)
	for op, msg := range vm.Help {
		fmt.Printf("\t%s\t%s\n", op, msg)
	}
	fmt.Printf("\nshorthand %s\n\n", shorthand.Version)
	return shorthand.SourceMap{Label: "", Op: ":help:", Source: "", Expanded: ""}, nil
}

//exitShorthand - call os.Exit() with appropriate value and exit the repl
var exitShorthand = func(vm *shorthand.VirtualMachine, sm shorthand.SourceMap) (shorthand.SourceMap, error) {
	if sm.Source == "" {
		os.Exit(0)
	}
	fmt.Fprintf(os.Stderr, sm.Source)
	os.Exit(1)
	return shorthand.SourceMap{Label: "", Op: ":exit:", Source: "", Expanded: ""}, nil
}

func revision() {
	fmt.Printf("%s %s\n", filepath.Base(os.Args[0]), shorthand.Version)
	os.Exit(0)
}

func (e *expressionList) String() string {
	return fmt.Sprintf("%s", *e)
}

func (e *expressionList) Set(value string) error {
	lineNo++
	out, err := vm.Eval(value, lineNo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR (%d): %s\n", lineNo, err)
		return err
	}
	if out != "" {
		fmt.Fprintf(os.Stdout, "%s\n", out)
	}
	return nil
}

func init() {
	vm = shorthand.New()
	vm.RegisterOp(":exit:", exitShorthand, "Exit shorthand repl")
}

func main() {
	flag.Var(&expression, "e", "The shorthand notation(s) you wish at add")
	flag.StringVar(&prompt, "p", "=> ", "Output a prompt for interactive processing")
	flag.BoolVar(&noprompt, "n", false, "Turn off the prompt for interactive processing")
	flag.BoolVar(&help, "h", false, "Display this help document")
	flag.BoolVar(&help, "help", false, "Display this help document")
	flag.BoolVar(&version, "v", false, "Version information")
	flag.BoolVar(&version, "version", false, "Version information")
	flag.Parse()
	args := flag.Args()
	if help == true {
		usage(0, "")
	}
	if version == true {
		revision()
	}

	if noprompt == true {
		prompt = ""
	}
	vm.SetPrompt(prompt)

	// If a filename is provided on the command line use it instead of standard input.
	if len(args) > 0 {
		vm.SetPrompt("")
		for _, arg := range args {
			fp, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
			}
			defer fp.Close()
			reader := bufio.NewReader(fp)
			vm.Run(reader)
		}
	} else {
		vm.RegisterOp(":help:", helpShorthand, "This help message")
		if prompt != "" {
			fmt.Println(welcome)
		}
		reader := bufio.NewReader(os.Stdin)
		vm.Run(reader)
	}
}
