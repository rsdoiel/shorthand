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

	// My packages
	"github.com/rsdoiel/ok"
)

// TestParse
func TestParse(t *testing.T) {
	validData := map[string]SourceMap{
		"@now1 :=: $(date)": SourceMap{
			Label:    "@now1",
			Op:       " :=: ",
			Source:   "$(date)",
			Expanded: "",
			LineNo:   1,
		},
		"this :=: a valid assignment": SourceMap{
			Label:    "this",
			Op:       " :=: ",
			Source:   "a valid assignment",
			Expanded: "",
			LineNo:   2,
		},
		"this; :=: a valid assignment": SourceMap{
			Label:    "this;",
			Op:       " :=: ",
			Source:   "a valid assignment",
			Expanded: "",
			LineNo:   3,
		},
		`now; :=: $(date +%H:%M);`: SourceMap{
			Label:    "now;",
			Op:       " :=: ",
			Source:   `$(date +%H:%M);`,
			Expanded: "",
			LineNo:   4,
		},
		"@now2 :=: Fred\n": SourceMap{
			Label:    "@now2",
			Op:       " :=: ",
			Source:   "Fred",
			Expanded: "",
			LineNo:   5,
		},
		"@file :=<: file.txt": SourceMap{
			Label:    "@file",
			Op:       " :=<: ",
			Source:   "file.txt",
			Expanded: "",
			LineNo:   6,
		},
		"@now3 :!: date": SourceMap{
			Label:    "@now3",
			Op:       " :!: ",
			Source:   "date",
			Expanded: "",
			LineNo:   7,
		},
		"@now4 :{: @one @two": SourceMap{
			Label:    "@now4",
			Op:       " :{: ",
			Source:   "@one @two",
			Expanded: "",
			LineNo:   8,
		},
		"@now5 :}<: test.shorthand": SourceMap{
			Label:    "@now5",
			Op:       " :}<: ",
			Source:   "test.shorthand",
			Expanded: "",
			LineNo:   9,
		},
		"@now6 :[: **strong words**": SourceMap{
			Label:    "@now6",
			Op:       " :[: ",
			Source:   "**strong words**",
			Expanded: "",
			LineNo:   10,
		},
		"@now7 :[<: test.md": SourceMap{
			Label:    "@now7",
			Op:       " :[<: ",
			Source:   "test.md",
			Expanded: "",
			LineNo:   11,
		},
		"@label0 :>: label0.txt": SourceMap{
			Label:    "@label0",
			Op:       " :>: ",
			Source:   "label0.txt",
			Expanded: "",
			LineNo:   12,
		},
		"@label1 :@>: label1.txt": SourceMap{
			Label:    "@label1",
			Op:       " :@>: ",
			Source:   "label1.txt",
			Expanded: "",
			LineNo:   13,
		},
		"@label2 :}>: label2.txt": SourceMap{
			Label:    "@label2",
			Op:       " :}>: ",
			Source:   "label2.txt",
			Expanded: "",
			LineNo:   14,
		},
		"@label3 :@}>: label3.txt": SourceMap{
			Label:    "@label3",
			Op:       " :@}>: ",
			Source:   "label3.txt",
			Expanded: "",
			LineNo:   15,
		},
		`@label4 :++ something`: SourceMap{
			Label:    "",
			Op:       "",
			Source:   "@label4 :++ something",
			Expanded: "",
			LineNo:   16,
		},
		"This should have @label4 and other things.": SourceMap{
			Label:    "",
			Op:       "",
			Source:   "This should have @label4 and other things.",
			Expanded: "",
			LineNo:   17,
		},
		"{{pageTitle}} :=: Hello World": SourceMap{
			Label:    "{{pageTitle}}",
			Op:       " :=: ",
			Source:   "Hello World",
			Expanded: "",
			LineNo:   18,
		},
		"{{year}} :!: echo -n $(date +%Y)": SourceMap{
			Label:    "{{year}}",
			Op:       " :!: ",
			Source:   "echo -n $(date +%Y)",
			Expanded: "",
			LineNo:   19,
		},
		"{fred} :=: Fred": SourceMap{
			Label:    "{fred}",
			Op:       " :=: ",
			Source:   "Fred",
			Expanded: "",
			LineNo:   20,
		},
		"{{strong}} :[: **strong words**": SourceMap{
			Label:    "{{strong}}",
			Op:       " :[: ",
			Source:   "**strong words**",
			Expanded: "",
			LineNo:   21,
		},
		"{one} :=: 1": SourceMap{
			Label:    "{one}",
			Op:       " :=: ",
			Source:   "1",
			Expanded: "",
			LineNo:   22,
		},
		"{two} :=: 2": SourceMap{
			Label:    "{two}",
			Op:       " :=: ",
			Source:   "2",
			Expanded: "",
			LineNo:   23,
		},
		"{it} :{: {one} {two}": SourceMap{
			Label:    "{it}",
			Op:       " :{: ",
			Source:   "{one} {two}",
			Expanded: "",
			LineNo:   24,
		},
		"{{html}} :[<: testdata/test.md": SourceMap{
			Label:    "{{html}}",
			Op:       " :[<: ",
			Source:   "testdata/test.md",
			Expanded: "",
			LineNo:   25,
		},
		"{helloWorldTxT} :=<: testdata/helloworld.txt": SourceMap{
			Label:    "{helloWorldTxT}",
			Op:       " :=<: ",
			Source:   "testdata/helloworld.txt",
			Expanded: "",
			LineNo:   26,
		},
	}

	vm := New()

	i := 1
	for s, ex := range validData {
		sm := vm.Parse(s, i)
		ok.Ok(t, strings.Compare(sm.Label, ex.Label) == 0, "Label should match "+sm.Label+" ? "+ex.Label)
		ok.Ok(t, strings.Compare(sm.Op, ex.Op) == 0, "Op should match "+sm.Op+" ? "+ex.Op)
		ok.Ok(t, strings.Compare(sm.Source, ex.Source) == 0, "Source should match "+sm.Source+" ? "+ex.Source)
		ok.Ok(t, strings.Compare(sm.Expanded, ex.Expanded) == 0, "Source should match "+sm.Source+" ? "+ex.Source)
		i++
	}
}

// TestEval
func TestEval(t *testing.T) {
	testData := []string{
		"@now :=: $(date)",                        // 0
		"this :=: a valid assignment",             // 1
		"this; :=: is a valid assignment",         // 2
		"now; :=: $(date +\"%H:%M\");",            // 3
		"@here :=<: testdata/testme.md",           // 4
		"@there :{<: testdata/testme.md",          // 5
		"{here} :=<: testdata/testme.md",          // 6
		"{{here}} :=<: testdata/testme.md",        // 7
		"This is not an assignment",               // 8
		"this:=: is not a valid assignment",       // 9
		"nor :=:is this a valid assignment",       // 10
		"and not : =:is this a valid assignment",  // 11
		"also not := :is this a valid assignment", // 12
	}
	vm := New()
	for i, src := range testData {
		eSM := vm.Parse(src, i)
		s, err := vm.Eval(src, i)
		sm := vm.Symbols.GetSymbol(eSM.Label)
		ok.Ok(t, err == nil, fmt.Sprintf("err should be nill: %s", err))
		if eSM.Label == "" && eSM.Op == "" {
			ok.Ok(t, s != "", "Not an assignment so should be non-empty: "+s)
			ok.Ok(t, s == vm.Expand(s), "expected '"+vm.Expand(s)+"' got "+s)
		} else {
			ok.Ok(t, eSM.LineNo == sm.LineNo, fmt.Sprintf("expected line no. %d got %d", eSM.LineNo, sm.LineNo))
			ok.Ok(t, eSM.Source == sm.Source, "expected source '"+eSM.Source+"' got "+sm.Source)
			ok.Ok(t, eSM.Label == sm.Label, "expected label '"+eSM.Label+"' got "+sm.Label)
			ok.Ok(t, eSM.Op == sm.Op, "expected op '"+eSM.Op+"' got "+sm.Op)
			ok.Ok(t, sm.Expanded != "", "expected some expansion, got "+sm.Expanded)
		}
		i++
	}
}

func TestSymbolTable(t *testing.T) {
	vm := New()
	st := new(SymbolTable)
	ok.Ok(t, len(st.entries) == 0, "st.entries should be zero")
	ok.Ok(t, len(st.labels) == 0, "st.labels should be zero too")

	sm1 := st.GetSymbol("@missing")
	ok.Ok(t, sm1.LineNo == -1, "Should fail with an empty symbol table")
	sm1 = vm.Parse("@now :=: This is now.", 1)
	i := st.SetSymbol(sm1)
	ok.Ok(t, i == 0, "Expected i to be zero as first element in symbol table")
	sm2 := st.GetSymbol("@now")
	ok.Ok(t, sm1.Label == sm2.Label, "expected label '"+sm1.Label+"' got "+sm2.Label)
	ok.Ok(t, sm1.Op == sm2.Op, "expected op '"+sm1.Op+"' got "+sm2.Op)
	ok.Ok(t, sm1.Source == sm2.Source, "expected source '"+sm1.Source+"' got "+sm2.Source)
	ok.Ok(t, sm1.Expanded == sm2.Expanded, "expected expanded '"+sm1.Expanded+"' got "+sm2.Expanded)
	ok.Ok(t, sm1.LineNo == sm2.LineNo, "expected expanded '"+sm1.Expanded+"' got "+sm2.Expanded)

	vm.Eval("@now :=: This is now.", 1)
	resultText := vm.Expand("This is '@now'")
	ok.Ok(t, resultText == "This is 'This is now.'", "Should have an expansion. ["+resultText+"]")
}

// Test Expand
func TestExpand(t *testing.T) {
	vm := New()
	ok.Ok(t, len(vm.Symbols.entries) == 0, "vm.Symbols.entries should be zero")
	ok.Ok(t, len(vm.Symbols.labels) == 0, "vm.Symbols.labels should be zero too")

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

	vm.Eval("@me :=: Fred", 1)
	vm.Eval("@now :=: 9:00", 2)
	result := vm.Expand(text)
	if result != expected {
		t.Fatalf("Expected:\n\n" + expected + "\n\nReceived:\n\n" + result)
	}
}

// Test include file
func TestInclude(t *testing.T) {
	vm := New()
	ok.Ok(t, len(vm.Symbols.entries) == 0, "vm.Symbols.entries should be zero")
	ok.Ok(t, len(vm.Symbols.labels) == 0, "vm.Symbols.labels should be zero too")

	buf, err := ioutil.ReadFile("testdata/testme.md")
	ok.Ok(t, err == nil, fmt.Sprintf("Should be able to read testdata/testme.md: %s", err))

	text := string(buf)
	_, err = vm.Eval("@TESTME :=<: testdata/testme.md", 1)
	ok.Ok(t, err == nil, "Should not get error of Eval assignment")
	resultText, err := vm.Eval("@TESTME", 1)
	ok.Ok(t, err == nil, "Should not get error on eval expand")
	ok.Ok(t, strings.Compare(text, resultText) == 0, "Should get same text for @TESTME")

	l := len(text)
	ok.Ok(t, len(resultText) >= l, fmt.Sprintf("Should %d have got %d results: %s", l, len(resultText), resultText))
	ok.Ok(t, strings.Contains(resultText, "A nimble webserver"), fmt.Sprintf("Should have 'A nimble webserver' in %s", resultText))
	ok.Ok(t, strings.Contains(resultText, "JSON"), fmt.Sprintf("Should have 'JSON' in %s", resultText))
}

func TestShellAssignment(t *testing.T) {
	vm := New()
	s, err := vm.Eval("@ECHO :!: echo -n 'Hello World!'", 1)
	ok.Ok(t, err == nil, fmt.Sprintf("assignment should not have an error: %s", err))
	ok.Ok(t, s == "", "Assignment should yield an empty string "+s)

	s, err = vm.Eval("@ECHO", 2)
	ok.Ok(t, err == nil, fmt.Sprintf("Expansion should not have an error: %s", err))
	ok.Ok(t, s == "Hello World!", "Should have @ECHO assignment: "+s)
}

func TestExpandedAssignment(t *testing.T) {
	vm := New()

	dateFormat := "2006-01-02"
	now := time.Now()
	// Date will generate a LF so the text will also contain it. So we'll test against a Trim later.
	vm.Eval(`@now :!: date +%Y-%m-%d`, 1)
	vm.Eval("@title :{: This is a title with date: @now", 2)
	resultText, err := vm.Eval("@title", 3)
	ok.Ok(t, err == nil, fmt.Sprintf("Expanded title should not have an error %s\n", err))
	expectedText := fmt.Sprintf("This is a title with date: %s\n", now.Format(dateFormat))
	ok.Ok(t, resultText == expectedText, "expected '"+expectedText+"' got '"+resultText+"'")

	// Now test a label that holds multiple lines that need expanding.
	text := `
			@one this is a line
			@two this is also a line
			@three this is the last line`

	vm.Eval("@one :=: 1", 1)
	vm.Eval("@two :=: 2", 2)
	vm.Eval("@three :=: 3", 3)
	vm.Eval("@text :=: "+text, 4)
	resultText = vm.Expand("@text")

	ok.Ok(t, strings.Contains(resultText, "@one this is a line"), "Should have line @one ["+resultText+"]")
	ok.Ok(t, strings.Contains(resultText, "@two this is also a line"), "Should have line @two "+resultText+"]")
	ok.Ok(t, strings.Contains(resultText, "@three this is the last line"), "Should have line @three "+resultText+"]")

	vm.Eval("@out :{{: @text", 5)
	resultText = vm.Expand("@out")

	ok.Ok(t, strings.Contains(resultText, "1 this is a line"), "Should have line 1 "+resultText)
	ok.Ok(t, strings.Contains(resultText, "2 this is also a line"), "Should have line 2 "+resultText)
	ok.Ok(t, strings.Contains(resultText, "3 this is the last line"), "Should have line 3 "+resultText)
}

func TestImport(t *testing.T) {
	vm := New()

	dateFormat := "2006-01-02"
	now := time.Now()
	atNow := fmt.Sprintf("%s\n", now.Format(dateFormat))
	atTitle := "Today's date report."
	atGreeting := "Hello World"

	// Now test evaluating a shorthand file
	_, err := vm.Eval("#T# :}<: testdata/test1.shorthand", 1)
	ok.Ok(t, err == nil, "Should be able to import testdata/test1.shorthand")

	result := vm.Expand("@now")
	ok.Ok(t, result == atNow, "Should have @now "+result)
	result = vm.Expand("@title")
	ok.Ok(t, result == atTitle, "Should have @title "+result)
	result = vm.Expand("@greeting")
	ok.Ok(t, result == atGreeting, "Should have @greeting "+result)

	shorthandText := `
		@title @now
		@greeting
		`

	result = vm.Expand(shorthandText)
	ok.Ok(t, strings.Contains(result, atTitle), "Should have title: "+atTitle)
	ok.Ok(t, strings.Contains(result, atNow), "Should have name: "+atNow)
	ok.Ok(t, strings.Contains(result, atGreeting), "Should have greeting: "+atGreeting)
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
		"@hello_world :=: Hello World",
		"@max :!: echo -n 'Hello Max'",
		"@hello_world :>: testdata/expansion1.txt",
		"_ :@>: testdata/expansion2.txt",
	}

	for i, data := range testData {
		_, err := vm.Eval(data, i)
		ok.Ok(t, err == nil, fmt.Sprintf("error for %s -> %s", data, err))
	}

	for fname, expectedText := range testFiles {
		buf, err := ioutil.ReadFile(fname)
		ok.Ok(t, err == nil, fmt.Sprintf("%s error: %s", fname, err))
		resultText := string(buf)
		terms := strings.Split(expectedText, "\n")
		for _, term := range terms {
			ok.Ok(t, strings.Contains(resultText, term),
				fmt.Sprintf("%s expected '%s' got '%s'", fname, term, resultText))
		}
	}
}

func TestExportAssignments(t *testing.T) {
	testFiles := map[string]string{
		"testdata/assigned1.txt": "@hello_world :=: Hello World",
		"testdata/assigned2.txt": "@hello_world :=: Hello World\n@max :!: echo -n 'Hello Max'\n",
	}
	for fname := range testFiles {
		if _, err := os.Stat(fname); err != nil {
			os.Remove(fname)
		}
	}

	vm := New()

	testData := []string{
		`@hello_world :=: Hello World`,
		`@max :!: echo -n 'Hello Max'`,
		`@hello_world :}>: testdata/assigned1.txt`,
		`_ :@}>: testdata/assigned2.txt`,
	}

	for i, src := range testData {
		_, err := vm.Eval(src, i)
		ok.Ok(t, err == nil, fmt.Sprintf("%d %s error: %s", i, src, err))
	}

	for fname, text := range testFiles {
		buf, err := ioutil.ReadFile(fname)
		ok.Ok(t, err == nil, "Should beable to read "+fname)
		resultText := string(buf)
		terms := strings.Split(text, "\n")
		for _, term := range terms {
			ok.Ok(t, strings.Contains(resultText, term),
				fmt.Sprintf("%s expected '%s' got '%s'", fname, text, resultText))
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
		vm.Eval(fmt.Sprintf("@test :[: %s", src), i)
		result := vm.Expand("@test")
		ok.Ok(t, expected == result, fmt.Sprintf("%s -> %s", expected, result))
		i++
	}

	vm.Eval("@link :=: my link", i)
	i++
	vm.Eval("@url :=: http://example.com", i)
	i++
	vm.Eval("@html :{[: [@link](@url)", i)
	i++
	expected := "<p><a href=\"http://example.com\">my link</a></p>\n"
	result := vm.Expand("@html")
	ok.Ok(t, strings.Compare(expected, result) == 0, fmt.Sprintf("%s != %s", expected, result))

	_, err := vm.Eval("@page :[<: testdata/test.md", i)
	i++
	ok.Ok(t, err == nil, fmt.Sprintf("%d testdata/test.md error: %s", i, err))
	result = vm.Expand("@page")
	ok.Ok(t, strings.Contains(result, "<h2>Another H2</h2>"), "Should have a h2 from test.md"+result)

	_, err = vm.Eval("H2 :=: heading two element", i)
	i++
	ok.Ok(t, err == nil, fmt.Sprintf("Should be able to assign string to 'H2': %s", err))
	result = vm.Expand("H2")
	ok.Ok(t, result == "heading two element", fmt.Sprintf("Should be able to expand 'H2': %s", result))

	_, err = vm.Eval("@page :{[<: testdata/test.md", i)
	i++
	ok.Ok(t, err == nil, fmt.Sprintf("%d testdata/test.md error: %s", i, err))
	result = vm.Expand("@page")
	ok.Ok(t, strings.Contains(result, "<h2>Another heading two element</h2>"), "Should have another heading two element from test.md"+result)

	// Re-read testdata/test.md and process the maarkdown
	_, err = vm.Eval("H2 :=: heading two element", i)
	i++
	ok.Ok(t, err == nil, fmt.Sprintf("Should be able to assign string to 'H2': %s", err))
	_, err = vm.Eval("@page :{[<: testdata/test.md", i)
	i++
	ok.Ok(t, err == nil, fmt.Sprintf("%d testdata/test.md error: %s", i, err))
	_, err = vm.Eval("@page :>: testdata/test.html", i)
	ok.Ok(t, err == nil, fmt.Sprintf("%d write testdata/test.html error: %s", i, err))
	i++

	_, err = os.Stat("testdata/test.html")
	ok.Ok(t, err == nil, "testdata/test.html should exist.")

	buf, err := ioutil.ReadFile("testdata/test.html")
	ok.Ok(t, err == nil, "should be able to read testdata/test.html.")
	result2 := string(buf)
	ok.Ok(t, strings.Compare(result, result2) == 0, fmt.Sprintf("Should have same results: '%s' <> '%s'", result, result2))

}

func TestRun(t *testing.T) {
	vm := New()
	if vm == nil {
		t.Error("vm was not created by New()")
	}
	fp, err := os.Open("testdata/run1.shorthand")
	ok.Ok(t, err == nil, "Should be able to open testdata/run1.shorthand")

	reader := bufio.NewReader(fp)

	fmt.Println("Starting vm.Run()")
	cnt := vm.Run(reader)
	ok.Ok(t, cnt == 3, fmt.Sprintf("Exited int wrong : %d", cnt))
	fmt.Println("Success.")
}
