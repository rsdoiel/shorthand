//
// Package shorthand provides shorthand definition and expansion.
//
// shorthand.go - A simple definition and expansion notation to use
// as shorthand when a template language is too much.
//
// @author R. S. Doiel, <rsdoiel@gmail.com>
// copyright (c) 2015 all rights reserved.
// Released under the BSD 2-Clause license
// See: http://opensource.org/licenses/BSD-2-Clause
//
// shorthand_test.go - tests for short package for handling shorthand
// definition and expansion.
//
package shorthand

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	// 3rd Party packages
	"github.com/russross/blackfriday"
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
		"@now1 :set: $(date)": {
			Label:    "@now1",
			Op:       " :set: ",
			Source:   "$(date)",
			Expanded: "",
			LineNo:   1,
		},
		"this :set: a valid assignment": {
			Label:    "this",
			Op:       " :set: ",
			Source:   "a valid assignment",
			Expanded: "",
			LineNo:   2,
		},
		"this; :set: a valid assignment": {
			Label:    "this;",
			Op:       " :set: ",
			Source:   "a valid assignment",
			Expanded: "",
			LineNo:   3,
		},
		`now; :set: $(date +%H:%M);`: {
			Label:    "now;",
			Op:       " :set: ",
			Source:   `$(date +%H:%M);`,
			Expanded: "",
			LineNo:   4,
		},
		"@now2 :set: Fred\n": {
			Label:    "@now2",
			Op:       " :set: ",
			Source:   "Fred",
			Expanded: "",
			LineNo:   5,
		},
		"@file :import-text: file.txt": {
			Label:    "@file",
			Op:       " :import-text: ",
			Source:   "file.txt",
			Expanded: "",
			LineNo:   6,
		},
		"@now3 :bash: date": {
			Label:    "@now3",
			Op:       " :bash: ",
			Source:   "date",
			Expanded: "",
			LineNo:   7,
		},
		"@now4 :expand: @one @two": {
			Label:    "@now4",
			Op:       " :expand: ",
			Source:   "@one @two",
			Expanded: "",
			LineNo:   8,
		},
		"@now5 :import-shorthand: test.shorthand": {
			Label:    "@now5",
			Op:       " :import-shorthand: ",
			Source:   "test.shorthand",
			Expanded: "",
			LineNo:   9,
		},
		"@now6 :markdown: **strong words**": {
			Label:    "@now6",
			Op:       " :markdown: ",
			Source:   "**strong words**",
			Expanded: "",
			LineNo:   10,
		},
		"@now7 :import-markdown: test.md": {
			Label:    "@now7",
			Op:       " :import-markdown: ",
			Source:   "test.md",
			Expanded: "",
			LineNo:   11,
		},
		"@label0 :export: label0.txt": {
			Label:    "@label0",
			Op:       " :export: ",
			Source:   "label0.txt",
			Expanded: "",
			LineNo:   12,
		},
		"@label1 :export-all: label1.txt": {
			Label:    "@label1",
			Op:       " :export-all: ",
			Source:   "label1.txt",
			Expanded: "",
			LineNo:   13,
		},
		"@label2 :export-shorthand: label2.txt": {
			Label:    "@label2",
			Op:       " :export-shorthand: ",
			Source:   "label2.txt",
			Expanded: "",
			LineNo:   14,
		},
		"@label3 :export-all-shorthand: label3.txt": {
			Label:    "@label3",
			Op:       " :export-all-shorthand: ",
			Source:   "label3.txt",
			Expanded: "",
			LineNo:   15,
		},
		`@label4 :++ something`: {
			Label:    "",
			Op:       "",
			Source:   "@label4 :++ something",
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
		"{{pageTitle}} :set: Hello World": {
			Label:    "{{pageTitle}}",
			Op:       " :set: ",
			Source:   "Hello World",
			Expanded: "",
			LineNo:   18,
		},
		"{{year}} :bash: echo -n $(date +%Y)": {
			Label:    "{{year}}",
			Op:       " :bash: ",
			Source:   "echo -n $(date +%Y)",
			Expanded: "",
			LineNo:   19,
		},
		"{fred} :set: Fred": {
			Label:    "{fred}",
			Op:       " :set: ",
			Source:   "Fred",
			Expanded: "",
			LineNo:   20,
		},
		"{{strong}} :markdown: **strong words**": {
			Label:    "{{strong}}",
			Op:       " :markdown: ",
			Source:   "**strong words**",
			Expanded: "",
			LineNo:   21,
		},
		"{one} :set: 1": {
			Label:    "{one}",
			Op:       " :set: ",
			Source:   "1",
			Expanded: "",
			LineNo:   22,
		},
		"{two} :set: 2": {
			Label:    "{two}",
			Op:       " :set: ",
			Source:   "2",
			Expanded: "",
			LineNo:   23,
		},
		"{it} :expand: {one} {two}": {
			Label:    "{it}",
			Op:       " :expand: ",
			Source:   "{one} {two}",
			Expanded: "",
			LineNo:   24,
		},
		"{{html}} :import-markdown: testdata/test.md": {
			Label:    "{{html}}",
			Op:       " :import-markdown: ",
			Source:   "testdata/test.md",
			Expanded: "",
			LineNo:   25,
		},
		"{helloWorldTxT} :import-text: testdata/helloworld.txt": {
			Label:    "{helloWorldTxT}",
			Op:       " :import-text: ",
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
		"@now1 :set: $(date)": {
			Label:    "@now1",
			Op:       " :set: ",
			Source:   "$(date)",
			Expanded: "",
			LineNo:   1,
		},
		"this :set: a valid assignment": {
			Label:    "this",
			Op:       " :set: ",
			Source:   "a valid assignment",
			Expanded: "",
			LineNo:   2,
		},
		"this; :set: a valid assignment": {
			Label:    "this;",
			Op:       " :set: ",
			Source:   "a valid assignment",
			Expanded: "",
			LineNo:   3,
		},
		`now; :set: $(date +%H:%M);`: {
			Label:    "now;",
			Op:       " :set: ",
			Source:   `$(date +%H:%M);`,
			Expanded: "",
			LineNo:   4,
		},
		"@now2 :set: Fred\n": {
			Label:    "@now2",
			Op:       " :set: ",
			Source:   "Fred",
			Expanded: "",
			LineNo:   5,
		},
		"@file :import-text: file.txt": {
			Label:    "@file",
			Op:       " :import-text: ",
			Source:   "file.txt",
			Expanded: "",
			LineNo:   6,
		},
		"@now3 :bash: date": {
			Label:    "@now3",
			Op:       " :bash: ",
			Source:   "date",
			Expanded: "",
			LineNo:   7,
		},
		"@now4 :expand: @one @two": {
			Label:    "@now4",
			Op:       " :expand: ",
			Source:   "@one @two",
			Expanded: "",
			LineNo:   8,
		},
		"@now5 :import: test.shorthand": {
			Label:    "@now5",
			Op:       " :import: ",
			Source:   "test.shorthand",
			Expanded: "",
			LineNo:   9,
		},
		"@now6 :markdown: **strong words**": {
			Label:    "@now6",
			Op:       " :markdown: ",
			Source:   "**strong words**",
			Expanded: "",
			LineNo:   10,
		},
		"@now7 :import-markdown: test.md": {
			Label:    "@now7",
			Op:       " :import-markdown: ",
			Source:   "test.md",
			Expanded: "",
			LineNo:   11,
		},
		"@label0 :export: label0.txt": {
			Label:    "@label0",
			Op:       " :export: ",
			Source:   "label0.txt",
			Expanded: "",
			LineNo:   12,
		},
		"@label1 :export-all: label1.txt": {
			Label:    "@label1",
			Op:       " :export-all: ",
			Source:   "label1.txt",
			Expanded: "",
			LineNo:   13,
		},
		"@label2 :export: label2.txt": {
			Label:    "@label2",
			Op:       " :export: ",
			Source:   "label2.txt",
			Expanded: "",
			LineNo:   14,
		},
		"@label3 :export-all: label3.txt": {
			Label:    "@label3",
			Op:       " :export-all: ",
			Source:   "label3.txt",
			Expanded: "",
			LineNo:   15,
		},
		`@label4 :++ something`: {
			Label:    "",
			Op:       "",
			Source:   "@label4 :++ something",
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
		"{{pageTitle}} :set: Hello World": {
			Label:    "{{pageTitle}}",
			Op:       " :set: ",
			Source:   "Hello World",
			Expanded: "",
			LineNo:   18,
		},
		"{{year}} :bash: echo -n $(date +%Y)": {
			Label:    "{{year}}",
			Op:       " :bash: ",
			Source:   "echo -n $(date +%Y)",
			Expanded: "",
			LineNo:   19,
		},
		"{fred} :set: Fred": {
			Label:    "{fred}",
			Op:       " :set: ",
			Source:   "Fred",
			Expanded: "",
			LineNo:   20,
		},
		"{{strong}} :markdown: **strong words**": {
			Label:    "{{strong}}",
			Op:       " :markdown: ",
			Source:   "**strong words**",
			Expanded: "",
			LineNo:   21,
		},
		"{one} :set: 1": {
			Label:    "{one}",
			Op:       " :set: ",
			Source:   "1",
			Expanded: "",
			LineNo:   22,
		},
		"{two} :set: 2": {
			Label:    "{two}",
			Op:       " :set: ",
			Source:   "2",
			Expanded: "",
			LineNo:   23,
		},
		"{it} :expand: {one} {two}": {
			Label:    "{it}",
			Op:       " :expand: ",
			Source:   "{one} {two}",
			Expanded: "",
			LineNo:   24,
		},
		"{{html}} :import-markdown: testdata/test.md": {
			Label:    "{{html}}",
			Op:       " :import-markdown: ",
			Source:   "testdata/test.md",
			Expanded: "",
			LineNo:   25,
		},
		"{helloWorldTxT} :import-text: testdata/helloworld.txt": {
			Label:    "{helloWorldTxT}",
			Op:       " :import-text: ",
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
		"@now :set: $(date)",                        // 0
		"this :set: a valid assignment",             // 1
		"this; :set: is a valid assignment",         // 2
		"now; :set: $(date +\"%H:%M\");",            // 3
		"@here :import-text: testdata/testme.md",    // 4
		"@there :import: testdata/testme.md",        // 5
		"{here} :import-text: testdata/testme.md",   // 6
		"{{here}} :import-text: testdata/testme.md", // 7
		"This is not an assignment",                 // 8
		"this:set: is not a valid assignment",       // 9
		"nor :set:is this a valid assignment",       // 10
		"and not : =:is this a valid assignment",    // 11
		"also not := :is this a valid assignment",   // 12
	}
	vm := New()
	for i, src := range testData {
		eSM := vm.Parse(src, i)
		s, err := vm.Eval(src, i)
		sm := vm.Symbols.GetSymbol(eSM.Label)
		if notOk(err == nil) {
			fmt.Sprintf("err should be nill: %s", err)
		}
		if eSM.Label == "" && eSM.Op == "" {
			if notOk(s != "") {
				t.Errorf("not an assignment so should be non-empty: %q", s)
			}
			if notOk(s == vm.Expand(s)) {
				t.Errorf("expected %q, got %q", vm.Expand(s), s)
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
		"@now :set: $(date)",                          // 0
		"this :set: a valid assignment",               // 1
		"this; :set: is a valid assignment",           // 2
		"now; :set: $(date +\"%H:%M\");",              // 3
		"@here :import-text: testdata/testme.md",      // 4
		"@there :import: testdata/testme.md",          // 5
		"{here} :import-text: testdata/testme.md",     // 6
		"{{here}} :import-text: testdata/testme.md",   // 7
		"This is not an assignment",                   // 8
		"this:set: is not a valid assignment",         // 9
		"nor :set:is this a valid assignment",         // 10
		"and not : label:is this a valid assignment",  // 11
		"also not :label :is this a valid assignment", // 12
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
			if notOk(s == vm.Expand(s)) {
				t.Errorf("expected %q, got %q", vm.Expand(s), s)
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
	sm1 = vm.Parse("@now :set: This is now.", 1)
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

	vm.Eval("@now :set: This is now.", 1)
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

	vm.Eval("@me :set: Fred", 1)
	vm.Eval("@now :set: 9:00", 2)
	result := vm.Expand(text)
	if result != expected {
		t.Fatalf("Expected:\n\n" + expected + "\n\nReceived:\n\n" + result)
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

	buf, err := ioutil.ReadFile("testdata/testme.md")
	if notOk(err == nil) {
		t.Errorf("Should be able to read testdata/testme.md: %s", err)
	}

	text := string(buf)
	_, err = vm.Eval("@TESTME :import-text: testdata/testme.md", 1)
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
	s, err := vm.Eval("@ECHO :bash: echo -n 'Hello World!'", 1)
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
	vm.Eval(`@now :bash: date +%Y-%m-%d`, 1)
	vm.Eval("@title :expand: This is a title with date: @now", 2)
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

	vm.Eval("@one :set: 1", 1)
	vm.Eval("@two :set: 2", 2)
	vm.Eval("@three :set: 3", 3)
	vm.Eval("@text :set: "+text, 4)
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

	vm.Eval("@out :expand-expansion: @text", 5)
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
	_, err := vm.Eval("#T# :import-shorthand: testdata/test1.shorthand", 1)
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
		"@hello_world :set: Hello World",
		"@max :bash: echo -n 'Hello Max'",
		"@hello_world :export: testdata/expansion1.txt",
		"_ :export-all: testdata/expansion2.txt",
	}

	for i, data := range testData {
		_, err := vm.Eval(data, i)
		if notOk(err == nil) {
			t.Errorf("error for %s -> %s", data, err)
		}
	}

	for fname, expectedText := range testFiles {
		buf, err := ioutil.ReadFile(fname)
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
		"testdata/assigned1.txt": "@hello_world :set: Hello World",
		"testdata/assigned2.txt": "@hello_world :set: Hello World\n@max :bash: echo -n 'Hello Max'\n",
	}
	for fname := range testFiles {
		if _, err := os.Stat(fname); err != nil {
			os.Remove(fname)
		}
	}

	vm := New()

	testData := []string{
		`@hello_world :set: Hello World`,
		`@max :bash: echo -n 'Hello Max'`,
		`@hello_world :export-shorthand: testdata/assigned1.txt`,
		`_ :export-all-shorthand: testdata/assigned2.txt`,
	}

	for i, src := range testData {
		_, err := vm.Eval(src, i)
		if notOk(err == nil) {
			t.Errorf("%d %s error: %s", i, src, err)
		}
	}

	for fname, text := range testFiles {
		buf, err := ioutil.ReadFile(fname)
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

func TestMarkdownSupport(t *testing.T) {
	vm := New()

	testData := map[string]string{
		"[my link](http://example.org)": string(blackfriday.MarkdownCommon([]byte("[my link](http://example.org)"))),
		"**strong**":                    string(blackfriday.MarkdownCommon([]byte("**strong**"))),
	}

	i := 0
	for src, expected := range testData {
		vm.Eval(fmt.Sprintf("@test :markdown: %s", src), i)
		result := vm.Expand("@test")
		if notOk(expected == result) {
			t.Errorf("%s -> %s", expected, result)
		}
		i++
	}

	vm.Eval("@link :set: my link", i)
	i++
	vm.Eval("@url :set: http://example.com", i)
	i++
	vm.Eval("@html :expand-markdown: [@link](@url)", i)
	i++
	expected := "<p><a href=\"http://example.com\">my link</a></p>\n"
	result := vm.Expand("@html")
	if notOk(strings.Compare(expected, result) == 0) {
		t.Errorf("%s != %s", expected, result)
	}

	_, err := vm.Eval("@page :import-markdown: testdata/test.md", i)
	i++
	if notOk(err == nil) {
		t.Errorf("%d testdata/test.md error: %s", i, err)
	}
	result = vm.Expand("@page")
	if notOk(strings.Contains(result, "<h2>Another H2</h2>")) {
		t.Errorf("Should have a h2 from test.md,  %q", result)
	}

	_, err = vm.Eval("H2 :set: heading two element", i)
	i++
	if notOk(err == nil) {
		t.Errorf("Should be able to assign string to 'H2': %s", err)
	}
	result = vm.Expand("H2")
	if notOk(result == "heading two element") {
		t.Errorf("Should be able to expand 'H2': %s", result)
	}

	_, err = vm.Eval("@page :import-expand-markdown: testdata/test.md", i)
	i++
	if notOk(err == nil) {
		t.Errorf("%d testdata/test.md error: %s", i, err)
	}
	result = vm.Expand("@page")
	if notOk(strings.Contains(result, "<h2>Another heading two element</h2>")) {
		t.Errorf("Should have another heading two element from test.md, %q", result)
	}

	// Re-read testdata/test.md and process the maarkdown
	_, err = vm.Eval("H2 :set: heading two element", i)
	i++
	if notOk(err == nil) {
		t.Errorf("Should be able to assign string to 'H2': %s", err)
	}
	_, err = vm.Eval("@page :import-expand-markdown: testdata/test.md", i)
	i++
	if notOk(err == nil) {
		t.Errorf("%d testdata/test.md error: %s", i, err)
	}
	_, err = vm.Eval("@page :export: testdata/test.html", i)
	if notOk(err == nil) {
		t.Errorf("%d write testdata/test.html error: %s", i, err)
	}
	i++

	_, err = os.Stat("testdata/test.html")
	if notOk(err == nil) {
		t.Errorf("testdata/test.html should exist.")
	}

	buf, err := ioutil.ReadFile("testdata/test.html")
	if notOk(err == nil) {
		t.Errorf("should be able to read testdata/test.html.")
	}
	result2 := string(buf)
	if notOk(strings.Compare(result, result2) == 0) {
		t.Errorf("Should have same results: '%s' <> '%s'", result, result2)
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
	cnt := vm.Run(reader, false)
	if notOk(cnt == 3) {
		t.Errorf("Exited int wrong : %d", cnt)
	}
	fmt.Println("Success.")
}
