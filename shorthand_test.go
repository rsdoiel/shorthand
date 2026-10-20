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
// shorthand_test.go - tests for short package for handling shorthand
// definition and expansion.
package shorthand

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"
)

// notOk is similar to assertError true
func notOk(expression bool) bool {
	if expression == true {
		return false
	}
	return true
}

// TestParseOps
func TestParseOps(t *testing.T) {
	validData := map[string]SourceMap{
		":set: @now1 $(date)": {
			Label:    "@now1",
			Op:       ":set:",
			Source:   "$(date)",
			Expanded: "",
			LineNo:   1,
		},
		":set: this a valid assignment": {
			Label:    "this",
			Op:       ":set:",
			Source:   "a valid assignment",
			Expanded: "",
			LineNo:   2,
		},
		":set: this; a valid assignment": {
			Label:    "this;",
			Op:       ":set:",
			Source:   "a valid assignment",
			Expanded: "",
			LineNo:   3,
		},
		`:set: now; $(date +%H:%M);`: {
			Label:    "now;",
			Op:       ":set:",
			Source:   `$(date +%H:%M);`,
			Expanded: "",
			LineNo:   4,
		},
		":set: @now2 Fred\n": {
			Label:    "@now2",
			Op:       ":set:",
			Source:   "Fred",
			Expanded: "",
			LineNo:   5,
		},
		":import-text: @file file.txt": {
			Label:    "@file",
			Op:       ":import-text:",
			Source:   "file.txt",
			Expanded: "",
			LineNo:   6,
		},
		":bash: @now3 date": {
			Label:    "@now3",
			Op:       ":bash:",
			Source:   "date",
			Expanded: "",
			LineNo:   7,
		},
		":expand: @now4 @one @two": {
			Label:    "@now4",
			Op:       ":expand:",
			Source:   "@one @two",
			Expanded: "",
			LineNo:   8,
		},
		":import-shorthand: @now5 test.shorthand": {
			Label:    "@now5",
			Op:       ":import-shorthand:",
			Source:   "test.shorthand",
			Expanded: "",
			LineNo:   9,
		},
		":export: @label0 label0.txt": {
			Label:    "@label0",
			Op:       ":export:",
			Source:   "label0.txt",
			Expanded: "",
			LineNo:   12,
		},
		":export-all: @label1 label1.txt": {
			Label:    "@label1",
			Op:       ":export-all:",
			Source:   "label1.txt",
			Expanded: "",
			LineNo:   13,
		},
		":export-shorthand: @label2 label2.txt": {
			Label:    "@label2",
			Op:       ":export-shorthand:",
			Source:   "label2.txt",
			Expanded: "",
			LineNo:   14,
		},
		":export-all-shorthand: @label3 label3.txt": {
			Label:    "@label3",
			Op:       ":export-all-shorthand:",
			Source:   "label3.txt",
			Expanded: "",
			LineNo:   15,
		},
		`:++ @label4 something`: {
			Label:    "",
			Op:       "",
			Source:   ":++ @label4 something",
			Expanded: "",
			LineNo:   16,
		},
		"This should have @label4 and other things.": {
			Label:    "",
			Op:       "",
			Source:   "This should have @label4 and other things.",
			Expanded: "",
			LineNo:   17,
		},
		":set: {{pageTitle}} Hello World": {
			Label:    "{{pageTitle}}",
			Op:       ":set:",
			Source:   "Hello World",
			Expanded: "",
			LineNo:   18,
		},
		":bash: {{year}} echo -n $(date +%Y)": {
			Label:    "{{year}}",
			Op:       ":bash:",
			Source:   "echo -n $(date +%Y)",
			Expanded: "",
			LineNo:   19,
		},
		":set: {fred} Fred": {
			Label:    "{fred}",
			Op:       ":set:",
			Source:   "Fred",
			Expanded: "",
			LineNo:   20,
		},
		":set: {one} 1": {
			Label:    "{one}",
			Op:       ":set:",
			Source:   "1",
			Expanded: "",
			LineNo:   22,
		},
		":set: {two} 2": {
			Label:    "{two}",
			Op:       ":set:",
			Source:   "2",
			Expanded: "",
			LineNo:   23,
		},
		":expand: {it} {one} {two}": {
			Label:    "{it}",
			Op:       ":expand:",
			Source:   "{one} {two}",
			Expanded: "",
			LineNo:   24,
		},
		":import-text: {helloWorldTxT} testdata/helloworld.txt": {
			Label:    "{helloWorldTxT}",
			Op:       ":import-text:",
			Source:   "testdata/helloworld.txt",
			Expanded: "",
			LineNo:   26,
		},
	}

	vm := New()

	i := 1
	for s, ex := range validData {
		sm := vm.Parse(s, i)
		if notOk(strings.Compare(sm.Label, ex.Label) == 0) {
			t.Errorf("%d %q Label should match %q ? %q", i, s, sm.Label, ex.Label)
		}
		if notOk(strings.Compare(sm.Op, ex.Op) == 0) {
			t.Errorf("%d %q Op should match %q ? %q", i, s, sm.Op, ex.Op)
		}
		if notOk(strings.Compare(sm.Source, ex.Source) == 0) {
			t.Errorf("%d %q Source should match %q ? %q", i, s, sm.Source, ex.Source)
		}
		if notOk(strings.Compare(sm.Expanded, ex.Expanded) == 0) {
			t.Errorf("%d %q Source should match %q ? %q", i, s, sm.Source, ex.Source)
		}
		i++
	}
}

// TestParseReadable
func TestParseReadable(t *testing.T) {
	validData := map[string]SourceMap{
		":set: @now1 $(date)": {
			Label:    "@now1",
			Op:       ":set:",
			Source:   "$(date)",
			Expanded: "",
			LineNo:   1,
		},
		":set: this a valid assignment": {
			Label:    "this",
			Op:       ":set:",
			Source:   "a valid assignment",
			Expanded: "",
			LineNo:   2,
		},
		":set: this; a valid assignment": {
			Label:    "this;",
			Op:       ":set:",
			Source:   "a valid assignment",
			Expanded: "",
			LineNo:   3,
		},
		`:set: now; $(date +%H:%M);`: {
			Label:    "now;",
			Op:       ":set:",
			Source:   `$(date +%H:%M);`,
			Expanded: "",
			LineNo:   4,
		},
		":set: @now2 Fred\n": {
			Label:    "@now2",
			Op:       ":set:",
			Source:   "Fred",
			Expanded: "",
			LineNo:   5,
		},
		":import-text: @file file.txt": {
			Label:    "@file",
			Op:       ":import-text:",
			Source:   "file.txt",
			Expanded: "",
			LineNo:   6,
		},
		":bash: @now3 date": {
			Label:    "@now3",
			Op:       ":bash:",
			Source:   "date",
			Expanded: "",
			LineNo:   7,
		},
		":expand: @now4 @one @two": {
			Label:    "@now4",
			Op:       ":expand:",
			Source:   "@one @two",
			Expanded: "",
			LineNo:   8,
		},
		":import: @now5 test.shorthand": {
			Label:    "@now5",
			Op:       ":import:",
			Source:   "test.shorthand",
			Expanded: "",
			LineNo:   9,
		},
		":export: @label0 label0.txt": {
			Label:    "@label0",
			Op:       ":export:",
			Source:   "label0.txt",
			Expanded: "",
			LineNo:   12,
		},
		":export-all: @label1 label1.txt": {
			Label:    "@label1",
			Op:       ":export-all:",
			Source:   "label1.txt",
			Expanded: "",
			LineNo:   13,
		},
		":export: @label2 label2.txt": {
			Label:    "@label2",
			Op:       ":export:",
			Source:   "label2.txt",
			Expanded: "",
			LineNo:   14,
		},
		":export-all: @label3 label3.txt": {
			Label:    "@label3",
			Op:       ":export-all:",
			Source:   "label3.txt",
			Expanded: "",
			LineNo:   15,
		},
		`:++ @label4 something`: {
			Label:    "",
			Op:       "",
			Source:   ":++ @label4 something",
			Expanded: "",
			LineNo:   16,
		},
		"This should have @label4 and other things.": {
			Label:    "",
			Op:       "",
			Source:   "This should have @label4 and other things.",
			Expanded: "",
			LineNo:   17,
		},
		":set: {{pageTitle}} Hello World": {
			Label:    "{{pageTitle}}",
			Op:       ":set:",
			Source:   "Hello World",
			Expanded: "",
			LineNo:   18,
		},
		":bash: {{year}} echo -n $(date +%Y)": {
			Label:    "{{year}}",
			Op:       ":bash:",
			Source:   "echo -n $(date +%Y)",
			Expanded: "",
			LineNo:   19,
		},
		":set: {fred} Fred": {
			Label:    "{fred}",
			Op:       ":set:",
			Source:   "Fred",
			Expanded: "",
			LineNo:   20,
		},
		":set: {one} 1": {
			Label:    "{one}",
			Op:       ":set:",
			Source:   "1",
			Expanded: "",
			LineNo:   22,
		},
		":set: {two} 2": {
			Label:    "{two}",
			Op:       ":set:",
			Source:   "2",
			Expanded: "",
			LineNo:   23,
		},
		":expand: {it} {one} {two}": {
			Label:    "{it}",
			Op:       ":expand:",
			Source:   "{one} {two}",
			Expanded: "",
			LineNo:   24,
		},
		":import-text: {helloWorldTxT} testdata/helloworld.txt": {
			Label:    "{helloWorldTxT}",
			Op:       ":import-text:",
			Source:   "testdata/helloworld.txt",
			Expanded: "",
			LineNo:   26,
		},
	}

	vm := New()

	i := 1
	for s, ex := range validData {
		sm := vm.Parse(s, i)
		if notOk(strings.Compare(sm.Label, ex.Label) == 0) {
			t.Errorf("%d %q Label should match %q ? %q", i, s, sm.Label, ex.Label)
		}
		if notOk(strings.Compare(sm.Op, ex.Op) == 0) {
			t.Errorf("%d %q Op should match %q ? %q", i, s, sm.Op, ex.Op)
		}
		if notOk(strings.Compare(sm.Source, ex.Source) == 0) {
			t.Errorf("%d %q Source should match %q ? %q", i, s, sm.Source, ex.Source)
		}
		if notOk(strings.Compare(sm.Expanded, ex.Expanded) == 0) {
			t.Errorf("%d %q Source should match %q ? %q", i, s, sm.Source, ex.Source)
		}
		i++
	}
}

// TestEval
func TestEval(t *testing.T) {
	testData := []string{
		":set: @now $(date)",                        // 0
		":set: this a valid assignment",             // 1
		":set: this; is a valid assignment",         // 2
		":set: now; $(date +\"%H:%M\");",            // 3
		":import-text: @here testdata/testme.md",    // 4
		":import: @there testdata/testme.md",        // 5
		":import-text: {here} testdata/testme.md",   // 6
		":import-text: {{here}} testdata/testme.md", // 7
		"This is not an assignment",                 // 8
		":set:this is not a valid assignment",       // 9
		":set: this is a valid assignment",          // 10
		": set: is not this a valid assignment",     // 11
		":set : is not this a valid assignment",     // 12
	}
	vm := New()
	for i, src := range testData {
		eSM := vm.Parse(src, i)
		s, err := vm.Eval(src, i)
		sm := vm.Symbols.GetSymbol(eSM.Label)
		if notOk(err == nil) {
			t.Errorf("err should be nil: %s", err)
		}
		if eSM.Label == "" && eSM.Op == "" {
			if notOk(s != "") {
				t.Errorf("not an assignment so should be non-empty: %q", s)
			}
			// Eval of text is Expand of that text. (It used to compare with
			// Expand(s), which only held because the old Expand sometimes
			// expanded a value that holds another label, depending on the
			// order of the labels.)
			if notOk(s == vm.Expand(src)) {
				t.Errorf("expected %q, got %q", vm.Expand(src), s)
			}
		} else {
			if notOk(eSM.LineNo == sm.LineNo) {
				t.Errorf("expected line no. %d got %d", eSM.LineNo, sm.LineNo)
			}
			if notOk(eSM.Source == sm.Source) {
				t.Errorf("expected source %q, got %q", eSM.Source, sm.Source)
			}
			if notOk(eSM.Label == sm.Label) {
				t.Errorf("expected label %q, got %q", eSM.Label, sm.Label)
			}
			if notOk(eSM.Op == sm.Op) {
				t.Errorf("expected op %q, %q", eSM.Op, sm.Op)
			}
			if notOk(sm.Expanded != "") {
				t.Errorf("expected some expansion, got %q", sm.Expanded)
			}
		}
		i++
	}
}

// TestEvalReadable
func TestEvalReadable(t *testing.T) {
	testData := []string{
		":set: @now $(date)",                         // 0
		":set: this a valid assignment",              // 1
		":set: this; is a valid assignment",          // 2
		":set: now; $(date +\"%H:%M\");",             // 3
		":import-text: @here testdata/testme.md",     // 4
		":import: @there testdata/testme.md",         // 5
		":import-text: {here} testdata/testme.md",    // 6
		":import-text: {{here}} testdata/testme.md",  // 7
		"This is not an assignment",                  // 8
		"this:set: is not a valid assignment",        // 9
		":set:nor is this a valid assignment",        // 10
		": set: this is not a valid assignment",      // 11
		":set : also not is this a valid assignment", // 12
	}
	vm := New()
	for i, src := range testData {
		eSM := vm.Parse(src, i)
		s, err := vm.Eval(src, i)
		sm := vm.Symbols.GetSymbol(eSM.Label)
		if notOk(err == nil) {
			t.Errorf("err should be nill: %s", err)
		}
		if eSM.Label == "" && eSM.Op == "" {
			if notOk(s != "") {
				t.Errorf("Not an assignment so should be non-empty: %q ", s)
			}
			/*
				if notOk(s == vm.Expand(s)) {
					t.Errorf("expected %q, got %q", vm.Expand(s), s)
				}
			*/
		} else {
			if notOk(eSM.LineNo == sm.LineNo) {
				t.Errorf("expected line no. %d got %d", eSM.LineNo, sm.LineNo)
			}
			if notOk(eSM.Source == sm.Source) {
				t.Errorf("expected source %q, got %q", eSM.Source, sm.Source)
			}
			if notOk(eSM.Label == sm.Label) {
				t.Errorf("expected label %q, got %q", eSM.Label, sm.Label)
			}
			if notOk(eSM.Op == sm.Op) {
				t.Errorf("expected op %q, got %q", eSM.Op, sm.Op)
			}
			if notOk(sm.Expanded != "") {
				t.Errorf("expected some expansion, got %q", sm.Expanded)
			}
		}
		i++
	}
}

func TestSymbolTable(t *testing.T) {
	vm := New()
	st := new(SymbolTable)
	if notOk(len(st.entries) == 0) {
		t.Errorf("st.entries should be zero")
	}
	if notOk(len(st.labels) == 0) {
		t.Errorf("st.labels should be zero too")
	}

	sm1 := st.GetSymbol("@missing")
	if notOk(sm1.LineNo == -1) {
		t.Errorf("Should fail with an empty symbol table")
	}
	sm1 = vm.Parse(":set: @now This is now.", 1)
	i := st.SetSymbol(sm1)
	if notOk(i == 0) {
		t.Errorf("Expected i to be zero as first element in symbol table")
	}
	sm2 := st.GetSymbol("@now")
	if notOk(sm1.Label == sm2.Label) {
		t.Errorf("expected label %q, got %q", sm1.Label, sm2.Label)
	}
	if notOk(sm1.Op == sm2.Op) {
		t.Errorf("expected op %q, got %q", sm1.Op, sm2.Op)
	}
	if notOk(sm1.Source == sm2.Source) {
		t.Errorf("expected source %q, got %q", sm1.Source, sm2.Source)
	}
	if notOk(sm1.Expanded == sm2.Expanded) {
		t.Errorf("expected expanded %q, got %q'", sm1.Expanded, sm2.Expanded)
	}
	if notOk(sm1.LineNo == sm2.LineNo) {
		t.Errorf("expected expanded %q, got %q", sm1.Expanded, sm2.Expanded)
	}

	vm.Eval(":set: @now This is now.", 1)
	resultText := vm.Expand("This is '@now'")
	if notOk(resultText == "This is 'This is now.'") {
		t.Errorf("Should have an expansion. [%s]", resultText)
	}
}

// Test Expand
func TestExpand(t *testing.T) {
	vm := New()
	if notOk(len(vm.Symbols.entries) == 0) {
		t.Errorf("vm.Symbols.entries should be zero")
	}
	if notOk(len(vm.Symbols.labels) == 0) {
		t.Errorf("vm.Symbols.labels should be zero too")
	}

	text := `
		   @me

		   This is some line that should not change.

		   8:00 - @now; some stuff

		   This "now" should not change. This "me" should not change.`

	expected := `
		   Fred

		   This is some line that should not change.

		   8:00 - 9:00; some stuff

		   This "now" should not change. This "me" should not change.`

	vm.Eval(":set: @me Fred", 1)
	vm.Eval(":set: @now 9:00", 2)
	result := vm.Expand(text)
	if result != expected {
		t.Fatalf("Expected:\n\n%s\n\nReceived:\n\n%s", expected, result)
	}
}

// Test include file
func TestInclude(t *testing.T) {
	vm := New()
	if notOk(len(vm.Symbols.entries) == 0) {
		t.Errorf("vm.Symbols.entries should be zero")
	}
	if notOk(len(vm.Symbols.labels) == 0) {
		t.Errorf("vm.Symbols.labels should be zero too")
	}

	buf, err := os.ReadFile("testdata/testme.md")
	if notOk(err == nil) {
		t.Errorf("Should be able to read testdata/testme.md: %s", err)
	}

	text := string(buf)
	_, err = vm.Eval(":import-text: @TESTME testdata/testme.md", 1)
	if notOk(err == nil) {
		t.Errorf("Should not get error of Eval assignment")
	}
	resultText, err := vm.Eval("@TESTME", 1)
	if notOk(err == nil) {
		t.Errorf("Should not get error on eval expand")
	}
	if notOk(strings.Compare(text, resultText) == 0) {
		t.Errorf("Should get same text for @TESTME")
	}

	l := len(text)
	if notOk(len(resultText) >= l) {
		t.Errorf("Should %d have got %d results: %s", l, len(resultText), resultText)
	}
	if notOk(strings.Contains(resultText, "A nimble webserver")) {
		t.Errorf("Should have 'A nimble webserver' in %s", resultText)
	}
	if notOk(strings.Contains(resultText, "JSON")) {
		t.Errorf("Should have 'JSON' in %s", resultText)
	}
}

func TestShellAssignment(t *testing.T) {
	vm := New()
	s, err := vm.Eval(":bash: @ECHO echo -n 'Hello World!'", 1)
	if notOk(err == nil) {
		t.Errorf("assignment should not have an error: %s", err)
	}
	if notOk(s == "") {
		t.Errorf("Assignment should yield an empty string %q", s)
	}

	s, err = vm.Eval("@ECHO", 2)
	if notOk(err == nil) {
		t.Errorf("Expansion should not have an error: %s", err)
	}
	if notOk(s == "Hello World!") {
		t.Errorf("Should have @ECHO assignment: %q", s)
	}
}

func TestExpandedAssignment(t *testing.T) {
	vm := New()

	dateFormat := "2006-01-02"
	now := time.Now()
	// Date will generate a LF so the text will also contain it. So we'll test against a Trim later.
	vm.Eval(`:bash: @now date +%Y-%m-%d`, 1)
	vm.Eval(":expand: @title This is a title with date: @now", 2)
	resultText, err := vm.Eval("@title", 3)
	if notOk(err == nil) {
		t.Errorf("Expanded title should not have an error %s\n", err)
	}
	expectedText := fmt.Sprintf("This is a title with date: %s\n", now.Format(dateFormat))
	if notOk(resultText == expectedText) {
		t.Errorf("expected %q, got %q", expectedText, resultText)
	}

	// Now test a label that holds multiple lines that need expanding.
	text := `
			@one this is a line
			@two this is also a line
			@three this is the last line`

	vm.Eval(":set: @one 1", 1)
	vm.Eval(":set: @two 2", 2)
	vm.Eval(":set: @three 3", 3)
	vm.Eval(":set: @text "+text, 4)
	resultText = vm.Expand("@text")

	if notOk(strings.Contains(resultText, "@one this is a line")) {
		t.Errorf("Should have line @one [%s]", resultText)
	}
	if notOk(strings.Contains(resultText, "@two this is also a line")) {
		t.Errorf("Should have line @two %q", resultText)
	}
	if notOk(strings.Contains(resultText, "@three this is the last line")) {
		t.Errorf("Should have line @three %q", resultText)
	}

	vm.Eval(":expand-expansion: @out @text", 5)
	resultText = vm.Expand("@out")

	if notOk(strings.Contains(resultText, "1 this is a line")) {
		t.Errorf("Should have line 1 %q", resultText)
	}
	if notOk(strings.Contains(resultText, "2 this is also a line")) {
		t.Errorf("Should have line 2 %q", resultText)
	}
	if notOk(strings.Contains(resultText, "3 this is the last line")) {
		t.Errorf("Should have line 3 %q", resultText)
	}
}

func TestImport(t *testing.T) {
	vm := New()

	dateFormat := "2006-01-02"
	now := time.Now()
	atNow := fmt.Sprintf("%s\n", now.Format(dateFormat))
	atTitle := "Today's date report."
	atGreeting := "Hello World"

	// Now test evaluating a shorthand file
	_, err := vm.Eval(":import-shorthand: #T# testdata/test1.shorthand", 1)
	if notOk(err == nil) {
		t.Errorf("Should be able to import testdata/test1.shorthand")
	}

	result := vm.Expand("@now")
	if notOk(result == atNow) {
		t.Errorf("Should have @now %q", result)
	}
	result = vm.Expand("@title")
	if notOk(result == atTitle) {
		t.Errorf("Should have @title %q", result)
	}
	result = vm.Expand("@greeting")
	if notOk(result == atGreeting) {
		t.Errorf("Should have @greeting %q", result)
	}

	shorthandText := `
		@title @now
		@greeting
		`

	result = vm.Expand(shorthandText)
	if notOk(strings.Contains(result, atTitle)) {
		t.Errorf("Should have title: %q", atTitle)
	}
	if notOk(strings.Contains(result, atNow)) {
		t.Errorf("Should have name: %q", atNow)
	}
	if notOk(strings.Contains(result, atGreeting)) {
		t.Errorf("Should have greeting: %q", atGreeting)
	}
}

func TestExpandingSourcesToFile(t *testing.T) {
	testFiles := map[string]string{
		"testdata/expansion1.txt": "Hello World",
		"testdata/expansion2.txt": "Hello World\nHello Max\n",
	}

	// Clean up testdata area
	for fname := range testFiles {
		if _, err := os.Stat(fname); err != nil {
			os.Remove(fname)
		}
	}

	//Generate test data and verify output
	vm := New()

	testData := []string{
		":set: @hello_world Hello World",
		":bash: @max echo -n 'Hello Max'",
		":export: @hello_world testdata/expansion1.txt",
		":export-all: _ testdata/expansion2.txt",
	}

	for i, data := range testData {
		_, err := vm.Eval(data, i)
		if notOk(err == nil) {
			t.Errorf("error for %s -> %s", data, err)
		}
	}

	for fname, expectedText := range testFiles {
		buf, err := os.ReadFile(fname)
		if notOk(err == nil) {
			t.Errorf("%s error: %s", fname, err)
		}
		resultText := string(buf)
		terms := strings.Split(expectedText, "\n")
		for _, term := range terms {
			if notOk(strings.Contains(resultText, term)) {
				t.Errorf("%s expected '%s' got '%s'", fname, term, resultText)
			}
		}
	}
}

func TestExportAssignments(t *testing.T) {
	testFiles := map[string]string{
		"testdata/assigned1.txt": ":set: @hello_world Hello World",
		"testdata/assigned2.txt": ":set: @hello_world Hello World\n:bash: @max echo -n 'Hello Max'\n",
	}
	for fname := range testFiles {
		if _, err := os.Stat(fname); err != nil {
			os.Remove(fname)
		}
	}

	vm := New()

	testData := []string{
		`:set: @hello_world Hello World`,
		`:bash: @max echo -n 'Hello Max'`,
		`:export-shorthand: @hello_world testdata/assigned1.txt`,
		`:export-all-shorthand: _ testdata/assigned2.txt`,
	}

	for i, src := range testData {
		_, err := vm.Eval(src, i)
		if notOk(err == nil) {
			t.Errorf("%d %s error: %s", i, src, err)
		}
	}

	for fname, text := range testFiles {
		buf, err := os.ReadFile(fname)
		if notOk(err == nil) {
			t.Errorf("Should beable to read %q", fname)
		}
		resultText := string(buf)
		terms := strings.Split(text, "\n")
		for _, term := range terms {
			if notOk(strings.Contains(resultText, term)) {
				t.Errorf("%s expected '%s' got '%s'", fname, text, resultText)
			}
		}
	}
}

func TestRun(t *testing.T) {
	vm := New()
	if vm == nil {
		t.Error("vm was not created by New()")
	}
	fp, err := os.Open("testdata/run1.shorthand")
	if notOk(err == nil) {
		t.Errorf("Should be able to open testdata/run1.shorthand")
	}

	reader := bufio.NewReader(fp)

	fmt.Println("Starting vm.Run()")
	cnt := vm.Run(reader)
	if notOk(cnt == 3) {
		t.Errorf("Exited int wrong : %d", cnt)
	}
	fmt.Println("Success.")
}

/** TestRunWritesToTheMachinesWriters checks that Run writes the expansion to
 * vm.Out and the errors to vm.Err, so a caller can send them anywhere. Run
 * used to write to os.Stdout, which is why the command's -o option did nothing.
 *
 * Example:
 *   go test -run TestRunWritesToTheMachinesWriters
 */
func TestRunWritesToTheMachinesWriters(t *testing.T) {
	vm := New()
	var out, errOut strings.Builder
	vm.Out, vm.Err = &out, &errOut
	src := ":set: @who Freda\nHello @who\n:import-text: @x /no/such/file\nAgain @who\n"
	lines := vm.Run(bufio.NewReader(strings.NewReader(src)))
	if lines != 4 {
		t.Errorf("Run processed %d lines, want 4", lines)
	}
	if want := "Hello Freda\nAgain Freda\n"; out.String() != want {
		t.Errorf("Out = %q, want %q", out.String(), want)
	}
	if !strings.Contains(errOut.String(), "ERROR (3)") || !strings.Contains(errOut.String(), "/no/such/file") {
		t.Errorf("Err = %q, want the error for line 3", errOut.String())
	}
	if vm.ErrorCount() != 1 {
		t.Errorf("ErrorCount = %d, want 1", vm.ErrorCount())
	}
}

func TestRunPromptGoesToOut(t *testing.T) {
	vm := New()
	var out strings.Builder
	vm.Out = &out
	vm.SetPrompt("=> ")
	vm.Run(bufio.NewReader(strings.NewReader("hello\n")))
	if got := out.String(); got != "=> hello\n=> " {
		t.Errorf("Out = %q", got)
	}
}

func TestRunStopsAtExit(t *testing.T) {
	vm := New()
	var out strings.Builder
	vm.Out = &out
	vm.Run(bufio.NewReader(strings.NewReader("one\n:exit:\ntwo\n")))
	if out.String() != "one\n" {
		t.Errorf("Out = %q, want only the line before :exit:", out.String())
	}
}

func TestApplyExpandsBytesAndStopsAtExit(t *testing.T) {
	vm := New()
	out, err := vm.Apply([]byte(":set: @who Freda\nHello @who\n:exit:\nnot reached"))
	if err != nil {
		t.Fatal(err)
	}
	if want := "\nHello Freda"; string(out) != want {
		t.Errorf("Apply = %q, want %q", out, want)
	}
	if _, err := vm.Apply([]byte(":import-text: @x /no/such/file")); err == nil || !strings.Contains(err.Error(), "line (0)") {
		t.Errorf("Apply should report the line that failed, got %v", err)
	}
}

func TestParseShortForms(t *testing.T) {
	vm := New()
	cases := map[string]SourceMap{
		":set:":              {Op: ":set:", LineNo: 3},
		":set: @who":         {Op: ":set:", Label: "@who", LineNo: 3},
		":set: @who a b  c ": {Op: ":set:", Label: "@who", Source: "a b  c", LineNo: 3},
		"plain text":         {Source: "plain text", LineNo: 3},
	}
	for src, want := range cases {
		if got := vm.Parse(src, 3); got != want {
			t.Errorf("Parse(%q) = %+v, want %+v", src, got, want)
		}
	}
}

func TestRegisterOpRefusesADuplicate(t *testing.T) {
	vm := New()
	if err := vm.RegisterOp(":set:", AssignString, "again"); err == nil {
		t.Error("registering :set: twice should fail")
	}
	upper := func(vm *VirtualMachine, sm SourceMap) (SourceMap, error) {
		sm.Expanded = strings.ToUpper(sm.Source)
		return sm, nil
	}
	if err := vm.RegisterOp(":upper:", upper, "Assign upper case text"); err != nil {
		t.Fatal(err)
	}
	vm.Eval(":upper: @loud quiet words", 1)
	if got := vm.Expand("@loud"); got != "QUIET WORDS" {
		t.Errorf("a registered operator should work, got %q", got)
	}
}

func TestEvalSymbolUnsupportedOperator(t *testing.T) {
	vm := New()
	err := vm.EvalSymbol(SourceMap{Label: "@x", Op: ":nope:", Source: "y", LineNo: 7})
	if err == nil || !strings.Contains(err.Error(), "ERROR (7)") {
		t.Errorf("got %v", err)
	}
}

func TestZeroValueMachineWritesToTheOSStreams(t *testing.T) {
	vm := &VirtualMachine{}
	if vm.out() != os.Stdout || vm.err() != os.Stderr {
		t.Error("a machine with no writers should use os.Stdout and os.Stderr")
	}
	if New().ErrorCount() != 0 {
		t.Error("a new machine has no errors")
	}
}

func TestFmtHelpReplacesTheTokens(t *testing.T) {
	got := FmtHelp("{app_name} {version} {release_date} {release_hash}", "x", "1.2.3", "2026-01-02", "abc")
	if got != "x 1.2.3 2026-01-02 abc" {
		t.Errorf("FmtHelp = %q", got)
	}
}

// The operator tables in the docs must name the operators the code has, and
// all of them. They had drifted: shorthand-syntax.md listed :markdown:,
// :export-label: and others that do not exist.
func TestDocsListExactlyTheRegisteredOperators(t *testing.T) {
	registered := map[string]bool{}
	for _, op := range New().Ops {
		registered[op] = true
	}
	re := regexp.MustCompile(`:[a-z]+(?:-[a-z]+)*:`)
	docs := map[string]string{"shorthand-syntax.md": "", "HowItWorks": HowItWorks}
	if src, err := os.ReadFile("shorthand-syntax.md"); err == nil {
		docs["shorthand-syntax.md"] = string(src)
	} else {
		t.Fatalf("shorthand-syntax.md: %v", err)
	}
	for name, text := range docs {
		named := map[string]bool{}
		for _, op := range re.FindAllString(text, -1) {
			if op == ":exit:" || op == ":quit:" || op == ":help:" {
				continue
			}
			named[op] = true
			if !registered[op] {
				t.Errorf("%s names %s, which is not an operator", name, op)
			}
		}
		for op := range registered {
			if !named[op] {
				t.Errorf("%s does not mention %s", name, op)
			}
		}
	}
}

/** TestExpandOverlappingLabels checks that where two labels match at the same
 * place the longer one wins, every time. Expand used to walk the labels in map
 * order, so with @a and @ab defined "@ab" was "two" in one run in seven.
 *
 * Example:
 *   go test -run TestExpandOverlappingLabels
 */
func TestExpandOverlappingLabels(t *testing.T) {
	for i := 0; i < 300; i++ {
		vm := New()
		vm.Eval(":set: @a one", 1)
		vm.Eval(":set: @ab two", 2)
		vm.Eval(":set: @abc three", 3)
		for text, want := range map[string]string{
			"@a":           "one",
			"@ab":          "two",
			"@abc":         "three",
			"@abcd":        "threed",
			"@ax @ab @a.":  "onex two one.",
			"no labels":    "no labels",
			"@@ab@a@abc@b": "@two" + "one" + "three" + "@b",
		} {
			if got := vm.Expand(text); got != want {
				t.Fatalf("run %d: Expand(%q) = %q, want %q", i, text, got, want)
			}
		}
	}
}

// A value is not expanded again: that is what :expand-expansion: is for, so
// Expand has the same answer whatever order the labels were defined in.
func TestExpandIsOneLevel(t *testing.T) {
	for i := 0; i < 300; i++ {
		vm := New()
		vm.Eval(":set: @a one", 1)
		vm.Eval(":set: @b @a", 2)
		if got := vm.Expand("@a @b"); got != "one @a" {
			t.Fatalf("run %d: Expand = %q, want %q", i, got, "one @a")
		}
		vm.Eval(":expand-expansion: @c @b", 3)
		if got := vm.Expand("@c"); got != "one" {
			t.Fatalf("run %d: :expand-expansion: gave %q, want %q", i, got, "one")
		}
		vm.Eval(":expand: @d x @b y", 4)
		if got := vm.Expand("@d"); got != "x @a y" {
			t.Fatalf("run %d: :expand: gave %q, want %q", i, got, "x @a y")
		}
	}
}

// The underscore label is ignored, as the docs say, and a definition with no
// label does not touch the text. Both used to be replaced everywhere.
func TestUnderscoreAndEmptyLabelsAreIgnored(t *testing.T) {
	vm := New()
	vm.Eval(":set: _ ignored", 1)
	vm.EvalSymbol(SourceMap{Op: ":help:", Source: "no label", Expanded: "XX", LineNo: 2})
	vm.Eval(":set: @x hello", 3)
	if got := vm.Expand("snake_case abc @x"); got != "snake_case abc hello" {
		t.Errorf("Expand = %q", got)
	}
}

// A label defined again, or a new one, is seen by the next Expand.
func TestExpandSeesRedefinitions(t *testing.T) {
	vm := New()
	if got := vm.Expand("@a"); got != "@a" {
		t.Errorf("before any definition: %q", got)
	}
	vm.Eval(":set: @a one", 1)
	if got := vm.Expand("@a"); got != "one" {
		t.Errorf("after the first definition: %q", got)
	}
	vm.Eval(":set: @a two", 2)
	if got := vm.Expand("@a"); got != "two" {
		t.Errorf("after redefining: %q", got)
	}
	vm.Eval(":set: @ab three", 3)
	if got := vm.Expand("@a @ab"); got != "two three" {
		t.Errorf("after a longer label: %q", got)
	}
}
