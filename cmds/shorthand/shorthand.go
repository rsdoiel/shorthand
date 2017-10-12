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
	"path"

	// my packages
	shorthand "github.com/rsdoiel/shorthand"

	// Caltech Library packages
	"github.com/caltechlibrary/cli"
)

type expressionList []string

var (
	usage = `USAGE: %s [OPTIONS] [FILES_TO_PROCESS]`

	description = `%s is a command line utility to expand labels based on their
assigned definitions. The render output is the transformed text 
and without the shorthand definitions themselves. %s reads 
from standard input and writes to standard output.`
	license = `%s %s

copyright (c) 2015 all rights reserved.
Released under the BSD 2-Clause license.
See: http://opensource.org/licenses/BSD-2-Clause`

	welcome = `
  Welcome to shorthand the simple label expander and markdown processor.
  Use ':exit:' to quit the repl, ':help:' to get a list of supported operators.
`

	// Standard Options
	showHelp     bool
	showLicense  bool
	showVersion  bool
	showExamples bool

	// Application Options
	expression              expressionList
	prompt                  string
	noprompt                bool
	vm                      *shorthand.VirtualMachine
	lineNo                  int
	postProcessWithMarkdown bool
)

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
	// Standard Options
	flag.BoolVar(&showHelp, "h", false, "display help")
	flag.BoolVar(&showHelp, "help", false, "display help")
	flag.BoolVar(&showLicense, "l", false, "display license")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.BoolVar(&showVersion, "v", false, "Version information")
	flag.BoolVar(&showVersion, "version", false, "Version information")
	flag.BoolVar(&showExamples, "example", false, "display example(s)")

	// Application Options
	flag.Var(&expression, "e", "The shorthand notation(s) you wish at add")
	flag.StringVar(&prompt, "p", "=> ", "Output a prompt for interactive processing")
	flag.BoolVar(&noprompt, "n", false, "Turn off the prompt for interactive processing")
	flag.BoolVar(&postProcessWithMarkdown, "m", false, "Run final output through markdown processor")
	flag.BoolVar(&postProcessWithMarkdown, "markdown", false, "Run final output through markdown processor")
}

func main() {
	appName := path.Base(os.Args[0])
	flag.Parse()
	args := flag.Args()
	cfg := cli.New(appName, "", shorthand.Version)
	cfg.UsageText = fmt.Sprintf(usage, appName)
	cfg.DescriptionText = fmt.Sprintf(description, appName, appName)
	cfg.OptionText = "OPTIONS\n\n"
	cfg.ExampleText = shorthand.HowItWorks
	cfg.LicenseText = fmt.Sprintf(license, appName)

	if showHelp == true {
		if len(args) > 0 {
			fmt.Println(cfg.Help(args...))
		} else {
			fmt.Println(cfg.Usage())
		}
		os.Exit(0)
	}
	if showExamples == true {
		fmt.Println(cfg.ExampleText)
		os.Exit(0)
	}
	if showLicense == true {
		fmt.Println(cfg.License())
		os.Exit(0)
	}
	if showVersion == true {
		fmt.Println(cfg.Version())
		os.Exit(0)
	}

	vm = shorthand.New()
	vm.RegisterOp(":exit:", exitShorthand, "Exit shorthand repl")

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
			vm.Run(reader, postProcessWithMarkdown)
		}
	} else {
		// Run as repl
		vm.RegisterOp(":help:", helpShorthand, "This help message")
		if prompt != "" {
			fmt.Println(welcome)
		}
		reader := bufio.NewReader(os.Stdin)
		vm.Run(reader, false)
	}
}
