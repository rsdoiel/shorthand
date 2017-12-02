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
	"fmt"
	"os"

	// my packages
	shorthand "github.com/rsdoiel/shorthand"

	// Caltech Library packages
	"github.com/caltechlibrary/cli"
)

var (
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
	showHelp             bool
	showLicense          bool
	showVersion          bool
	showExamples         bool
	inputFName           string
	outputFName          string
	quiet                bool
	generateMarkdownDocs bool

	// Application Options
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

func main() {
	// Create app to hold the CLUI
	app := cli.NewCli(shorthand.Version)
	appName := app.AppName()

	// Describe expexted non-option parameters
	app.AddParams("[SHORTHAND_FILES]")

	// Add some help texts
	app.AddHelp("welcome", []byte(welcome))
	app.AddHelp("description", []byte(fmt.Sprintf(description, appName, appName)))
	app.AddHelp("examples", []byte(shorthand.HowItWorks))
	app.AddHelp("license", []byte(fmt.Sprintf(license, appName)))

	// Standard Options
	app.BoolVar(&showHelp, "h,help", false, "display help")
	app.BoolVar(&showLicense, "l,license", false, "display license")
	app.BoolVar(&showVersion, "v,version", false, "diplsay version")
	app.BoolVar(&showExamples, "examples", false, "display examples")
	app.StringVar(&inputFName, "i,input", "", "input filename")
	app.StringVar(&outputFName, "o,output", "", "output filename")
	app.BoolVar(&quiet, "quiet", false, "suppress error messages")
	app.BoolVar(&generateMarkdownDocs, "generate-markdown-docs", false, "output documentation in Markdown")

	// Application Options
	app.StringVar(&prompt, "p,prompt", "=> ", "Output a prompt for interactive processing")
	app.BoolVar(&noprompt, "n,no-prompt", false, "Turn off the prompt for interactive processing")
	app.BoolVar(&postProcessWithMarkdown, "m,markdown", false, "Run final output through markdown processor")

	app.Parse()
	args := app.Args()

	// Setup IO
	var err error

	app.Eout = os.Stderr

	app.In, err = cli.Open(inputFName, os.Stdin)
	cli.ExitOnError(app.Eout, err, quiet)
	defer cli.CloseFile(inputFName, app.In)

	app.Out, err = cli.Create(outputFName, os.Stdout)
	cli.ExitOnError(app.Eout, err, quiet)
	defer cli.CloseFile(outputFName, app.Out)

	// Handle options
	if generateMarkdownDocs {
		app.GenerateMarkdownDocs(app.Out)
		os.Exit(0)
	}
	if showHelp == true {
		if len(args) > 0 {
			fmt.Fprintf(app.Out, app.Help(args...))
		} else {
			app.Usage(app.Out)
		}
		os.Exit(0)
	}
	if showExamples == true {
		fmt.Fprintf(app.Out, app.Help("examples"))
		os.Exit(0)
	}
	if showLicense == true {
		fmt.Fprintf(app.Out, "%s\n", app.License())
		os.Exit(0)
	}
	if showVersion == true {
		fmt.Fprintf(app.Out, "%s\n", app.Version())
		os.Exit(0)
	}

	vm = shorthand.New()
	vm.RegisterOp(":exit:", exitShorthand, "Exit shorthand repl")

	if noprompt == true {
		prompt = ""
	}
	vm.SetPrompt(prompt)

	// if inputFName
	if inputFName != "" {
		vm.SetPrompt("")
		fp, err := os.Open(inputFName)
		if err != nil {
			fmt.Fprintf(app.Eout, "%s\n", err)
		}
		defer fp.Close()
		reader := bufio.NewReader(fp)
		vm.Run(reader, postProcessWithMarkdown)
	}

	// If a filename is provided on the command line use it instead of standard input.
	if len(args) > 0 {
		vm.SetPrompt("")
		for _, arg := range args {
			fp, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(app.Eout, "%s\n", err)
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
		reader := bufio.NewReader(app.In)
		vm.Run(reader, false)
	}
}
