// Copyright (C) 2015-2026 R. S. Doiel
// SPDX-License-Identifier: AGPL-3.0-or-later
// See the LICENSE file.

// Package cli holds the command line behavior of shorthand so that plain
// go test covers it. The program under cmd/shorthand only calls Run.
package cli

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	shorthand "github.com/rsdoiel/shorthand"
)

// The exit codes of every command line tool in the workspace (DR-0003).
const (
	ExitOK         = 0
	ExitUsage      = 2
	ExitData       = 65
	ExitNoInput    = 66
	ExitCantCreate = 73
)

const welcome = `
  Welcome to shorthand the simple label expander.
  Use ':exit:' to quit the repl, ':help:' to get a list of supported operators.
`

/** Run is the shorthand command: it expands the labels in the files named in
 * args, or in stdin when there are none.
 *
 * Parameters:
 *   appName (string) — the program name, for help and usage
 *   args ([]string) — the command line without the program name
 *   stdin (io.Reader) — standard input
 *   stdout (io.Writer) — standard output
 *   stderr (io.Writer) — standard error
 *
 * Returns:
 *   int — the exit code
 *
 * Example:
 *   os.Exit(cli.Run("shorthand", os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
 */
func Run(appName string, args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	var (
		showHelp, showLicense, showVersion, showExamples bool
		noPrompt, quiet                                  bool
		input, output, prompt                            string
	)
	fs := flag.NewFlagSet(appName, flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprintf(stderr, "Usage: %s [OPTIONS] [SHORTHAND_FILES]\nTry %s -help for the manual.\n", appName, appName)
	}
	fs.BoolVar(&showHelp, "h", false, "display help")
	fs.BoolVar(&showHelp, "help", false, "display help")
	fs.BoolVar(&showLicense, "l", false, "display license")
	fs.BoolVar(&showLicense, "license", false, "display license")
	fs.BoolVar(&showVersion, "v", false, "display version")
	fs.BoolVar(&showVersion, "version", false, "display version")
	fs.BoolVar(&showExamples, "examples", false, "display examples")
	fs.StringVar(&input, "i", "", "input filename")
	fs.StringVar(&input, "input", "", "input filename")
	fs.StringVar(&output, "o", "", "output filename")
	fs.StringVar(&output, "output", "", "output filename")
	fs.StringVar(&prompt, "p", "=> ", "prompt of the interactive interpreter")
	fs.StringVar(&prompt, "prompt", "=> ", "prompt of the interactive interpreter")
	fs.BoolVar(&noPrompt, "n", false, "turn off the prompt")
	fs.BoolVar(&noPrompt, "no-prompt", false, "turn off the prompt")
	fs.BoolVar(&quiet, "quiet", false, "do not report a file that cannot be opened")
	if err := fs.Parse(args); err != nil {
		return ExitUsage
	}
	promptGiven := false
	fs.Visit(func(f *flag.Flag) {
		if f.Name == "p" || f.Name == "prompt" {
			promptGiven = true
		}
	})

	switch {
	case showHelp:
		fmt.Fprintln(stdout, shorthand.FmtHelp(HelpText, appName, shorthand.Version, shorthand.ReleaseDate, shorthand.ReleaseHash))
		return ExitOK
	case showExamples:
		fmt.Fprint(stdout, shorthand.HowItWorks)
		return ExitOK
	case showLicense:
		fmt.Fprintln(stdout, strings.TrimSpace(shorthand.LicenseText))
		return ExitOK
	case showVersion:
		fmt.Fprintf(stdout, "%s %s %s\n", appName, shorthand.Version, shorthand.ReleaseHash)
		return ExitOK
	}

	out := stdout
	if output != "" && output != "-" {
		f, err := os.Create(output)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return ExitCantCreate
		}
		defer f.Close()
		out = f
	}
	vm := shorthand.New()
	vm.Out, vm.Err = out, stderr

	files := fs.Args()
	if input != "" {
		files = append([]string{input}, files...)
	}
	if len(files) == 0 {
		interpret(vm, stdin, out, prompt, noPrompt, promptGiven)
	} else {
		vm.SetPrompt("")
	}
	noInput := 0
	for _, name := range files {
		var in io.Reader = stdin
		if name != "-" {
			f, err := os.Open(name)
			if err != nil {
				if !quiet {
					fmt.Fprintln(stderr, err)
				}
				noInput++
				continue
			}
			defer f.Close()
			in = f
		}
		vm.Run(bufio.NewReader(in))
	}
	switch {
	case noInput > 0:
		return ExitNoInput
	case vm.ErrorCount() > 0:
		return ExitData
	}
	return ExitOK
}

// interpret runs the interpreter on stdin. A prompt, the welcome text and the
// :help: operator are for a person at a terminal; with input piped in they would
// be mixed into the output, so they are used only when stdin is a terminal or the
// prompt was asked for with -p.
func interpret(vm *shorthand.VirtualMachine, stdin io.Reader, out io.Writer, prompt string, noPrompt bool, promptGiven bool) {
	if noPrompt || !(promptGiven || isTerminal(stdin)) {
		prompt = ""
	}
	vm.SetPrompt(prompt)
	vm.RegisterOp(":help:", func(vm *shorthand.VirtualMachine, sm shorthand.SourceMap) (shorthand.SourceMap, error) {
		fmt.Fprintf(out, "\nThe following operators are supported in shorthand:\n\n")
		for _, op := range vm.Ops {
			fmt.Fprintf(out, "\t%s\t%s\n", op, vm.Help[op])
		}
		fmt.Fprintf(out, "\nshorthand %s\n\n", shorthand.Version)
		return shorthand.SourceMap{Label: "", Op: ":help:", Source: "", Expanded: ""}, nil
	}, "This help message")
	if prompt != "" {
		fmt.Fprint(out, welcome)
	}
	vm.Run(bufio.NewReader(stdin))
}

// isTerminal says r is a terminal rather than a file or a pipe.
func isTerminal(r io.Reader) bool {
	f, ok := r.(*os.File)
	if !ok {
		return false
	}
	info, err := f.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}
