//
// Package stn is a library for processing Simple Timesheet Notation.
//
// shorthand_test.go - tests for short package for handling shorthand
// definition and expansion.
//
// @author R. S. Doiel, <rsdoiel@gmail.com>
// copyright (c) 2015 all rights reserved.
// Released under the BSD 2-Clause license
// See: http://opensource.org/licenses/BSD-2-Clause
//
package shorthand

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/rsdoiel/ok"
)

// Test IsAssignment
func TestIsAssignment(t *testing.T) {
	validAssignments := []string{
		"@now :=: $(date)",
		"this :=: a valid assignment",
		"this; :=: is a valid assignment",
		"now; :=: $(date +\"%H:%M\");",
		"@here :=<: testdata/testme.md",
		"@there :{<: testdata/testme.md",
		"{here} :=<: testdata/testme.md",
		"{{here}} :=<: testdata/testme.md",
	}

	invalidAssignments := []string{
		"This is not an assignment",
		"this:= is not a valid assignment",
		"nor =:is this a valid assignment",
		"and not : =:is this a valid assignment",
		"also not := :is this a valid assignment",
	}

	for i := range validAssignments {
		if IsAssignment(validAssignments[i]) == false {
			t.Fatalf(validAssignments[i] + " should be a valid assignment.")
		}
	}

	for i := range invalidAssignments {
		if IsAssignment(invalidAssignments[i]) == true {
			t.Fatalf(invalidAssignments[i] + " should be an invalid assignment.")
		}
	}
}

// Test Parse
func TestParse(t *testing.T) {
	validData := map[string]SourceMap{
		"@now1 :=: $(date)": SourceMap{
			Label: "@now1",
			Op:    AssignString,
			Value: "$(date)",
		},
		"this :=: a valid assignment": SourceMap{
			Label: "this",
			Op:    AssignString,
			Value: "a valid assignment",
		},
		"this; :=: a valid assignment": SourceMap{
			Label: "this;",
			Op:    AssignString,
			Value: "a valid assignment",
		},
		`now; :=: $(date +%H:%M);`: SourceMap{
			Label: "now;",
			Op:    AssignString,
			Value: `$(date +%H:%M);`,
		},
		"@now2 :=: Fred\n": SourceMap{
			Label: "@now2",
			Op:    AssignString,
			Value: "Fred",
		},
		"@file :=<: file.txt": SourceMap{
			Label: "@file",
			Op:    AssignInclude,
			Value: "file.txt",
		},
		"@now3 :!: date": SourceMap{
			Label: "@now3",
			Op:    AssignShell,
			Value: "date",
		},
		"@now4 :{: @one @two": SourceMap{
			Label: "@now4",
			Op:    AssignExpansion,
			Value: "@one @two",
		},
		"@now5 :}<: test.shorthand": SourceMap{
			Label: "@now5",
			Op:    IncludeAssignments,
			Value: "test.shorthand",
		},
		"@now6 :[: **strong words**": SourceMap{
			Label: "@now6",
			Op:    AssignMarkdown,
			Value: "**strong words**",
		},
		"@now7 :[<: test.md": SourceMap{
			Label: "@now7",
			Op:    IncludeMarkdown,
			Value: "test.md",
		},
		"@label0 :>: label0.txt": SourceMap{
			Label: "@label0",
			Op:    OutputAssignedExpansion,
			Value: "label0.txt",
		},
		"@label1 :@>: label1.txt": SourceMap{
			Label: "@label1",
			Op:    OutputAssignedExpansions,
			Value: "label1.txt",
		},
		"@label2 :}>: label2.txt": SourceMap{
			Label: "@label2",
			Op:    OutputAssignment,
			Value: "label2.txt",
		},
		"@label3 :@}>: label3.txt": SourceMap{
			Label: "@label3",
			Op:    OutputAssignments,
			Value: "label3.txt",
		},
	}

	for s, ex := range validData {
		sm, r := Parse(s, 1)
		ok.Ok(t, r == true, fmt.Sprintf("Expected Parse OK: label: %s, op: %s, val: %s, expanded: %s, %b", sm.Label, sm.Op, sm.Value, sm.Expanded, r))
		ok.Ok(t, sm.Label == ex.Label, "Label should match "+sm.Label+" ? "+ex.Label)
		ok.Ok(t, sm.Op == ex.Op, "Op should match "+sm.Op+" ? "+ex.Op)
		ok.Ok(t, sm.Value == ex.Value, "Value should match "+sm.Value+" ? "+ex.Value)
	}

	// Check an invalid assignment
	s := `@label4 :++ something`
	sm, r := Parse(s, 1)
	ok.Ok(t, r == false, "Should not parse "+s)
	ok.Ok(t, sm.Label == "", "Should not have a label "+sm.Label)
	ok.Ok(t, sm.Op == "", "Should not have an op "+sm.Op)
	ok.Ok(t, sm.Value == "", "Should have an empty value "+sm.Value)
	ok.Ok(t, sm.Expanded == s, "Expanded should have original s "+s)

	// Check an expansion
	s = "This should have @label4 and other things."
	sm, r = Parse(s, 1)
	ok.Ok(t, r == false, "Should not parse "+s)
	ok.Ok(t, sm.Label == "", "Should not have a label "+sm.Label)
	ok.Ok(t, sm.Op == "", "Should not have an op "+sm.Op)
	ok.Ok(t, sm.Value == "", "Should not have a value "+sm.Value)
	ok.Ok(t, sm.Expanded == s, "Expanded should have original s "+s)
}

func TestSymbolTable(t *testing.T) {
	st := new(SymbolTable)
	ok.Ok(t, len(st.entries) == 0, "st.entries should be zero")
	ok.Ok(t, len(st.labels) == 0, "st.labels should be zero too")

	r := HasAssignment(st, "@missing")
	ok.Ok(t, r == false, "Should fail with an empty symbol table")

	ok.Ok(t, Assign(st, "@now :=: This is now.", 1), "Should have a successful assignment ")
	r = HasAssignment(st, "@now")
	ok.Ok(t, r == true, "Should have a @now symbol")

	s := "This is '@now'"
	r = HasAssignment(st, s)
	ok.Ok(t, r == false, "s is an expansion "+s)

	resultText := Expand(st, s)
	ok.Ok(t, resultText == "This is 'This is now.'", "Should have an expansion. ["+resultText+"]")
}

// Test Expand
func TestExpand(t *testing.T) {
	st := new(SymbolTable)
	ok.Ok(t, len(st.entries) == 0, "st.entries should be zero")
	ok.Ok(t, len(st.labels) == 0, "st.labels should be zero too")

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

	Assign(st, "@me :=: Fred\n", 1)
	Assign(st, "@now :=: 9:00", 2)
	result := Expand(st, text)
	if result != expected {
		t.Fatalf("Expected:\n\n" + expected + "\n\nReceived:\n\n" + result)
	}
}

// Test include file
func TestInclude(t *testing.T) {
	st := new(SymbolTable)
	ok.Ok(t, len(st.entries) == 0, "st.entries should be zero")
	ok.Ok(t, len(st.labels) == 0, "st.labels should be zero too")

	text := `
Today is @NOW.

Now add the testme.md to this.
-------------------------------
@TESTME
-------------------------------
Did it work?
`

	Assign(st, "@NOW :=: 2015-07-04", 1)
	expected := true
	results := HasAssignment(st, "@NOW")
	ok.Ok(t, results == expected, "Should have @NOW assignment")
	Assign(st, "@TESTME :=<: testdata/testme.md", 2)
	results = HasAssignment(st, "@TESTME")
	ok.Ok(t, results == expected, "Should have @TESTME assignment")
	resultText := Expand(st, text)
	l := len(text)
	ok.Ok(t, len(resultText) > l, "Should have more results: "+resultText)
	ok.Ok(t, strings.Contains(resultText, "A nimble webserver"), fmt.Sprintf("Should have 'A nimble webserver' in %s", resultText))
	ok.Ok(t, strings.Contains(resultText, "JSON"), fmt.Sprintf("Should have 'JSON' in %s", resultText))
}

func TestShellAssignment(t *testing.T) {
	st := new(SymbolTable)
	ok.Ok(t, len(st.entries) == 0, "st.entries should be zero")
	ok.Ok(t, len(st.labels) == 0, "st.labels should be zero too")

	expected := true
	expectedText := "Hello World!"
	Assign(st, "@ECHO :!: echo -n 'Hello World!'", 1)
	results := HasAssignment(st, "@ECHO")
	ok.Ok(t, results == expected, "Should have @ECHO assignment")
	resultText := Expand(st, "@ECHO")
	l := len(strings.Trim(resultText, "\n"))
	ok.Ok(t, l == len(expectedText), fmt.Sprintf("Expected length %d got %d for @ECHO", len(expectedText), l))
	ok.Ok(t, strings.Contains(strings.Trim(resultText, "\n"), expectedText), "Should have matching text for @ECHO")
}

func TestExpandedAssignment(t *testing.T) {
	st := new(SymbolTable)
	ok.Ok(t, len(st.entries) == 0, "st.entries should be zero")
	ok.Ok(t, len(st.labels) == 0, "st.labels should be zero too")

	dateFormat := "2006-01-02"
	now := time.Now()
	// Date will generate a LF so the text will also contain it. So we'll test against a Trim later.
	Assign(st, `@now :!: date +%Y-%m-%d`, 1)
	Assign(st, "@title :{: This is a title with date: @now", 2)
	text := `@title`
	expected := true
	results := HasAssignment(st, "@now")
	ok.Ok(t, results == expected, "Should have @now")
	results = HasAssignment(st, "@title")
	ok.Ok(t, results == expected, "Should have @title")
	expectedText := fmt.Sprintf("This is a title with date: %s", now.Format(dateFormat))
	resultText := Expand(st, text)
	l := len(strings.Trim(resultText, "\n"))
	ok.Ok(t, l == len(expectedText), "Should have expected length for @title")
	ok.Ok(t, strings.Contains(strings.Trim(resultText, "\n"), expectedText), "Should have matching text for @title")

	// Now test a label that holds multiple lines that need expanding.
	st2 := new(SymbolTable)
	ok.Ok(t, len(st2.entries) == 0, "st.entries should be zero")
	ok.Ok(t, len(st2.labels) == 0, "st.labels should be zero too")

	text = `
@one this is a line
@two this is also a line
@three this is the last line
`

	Assign(st2, "@one :=: 1", 1)
	Assign(st2, "@two :=: 2", 2)
	Assign(st2, "@three :=: 3", 3)
	Assign(st2, "@text :=: "+text, 4)
	resultText = Expand(st2, "@text")

	ok.Ok(t, strings.Contains(resultText, "@one this is a line"), "Should have line @one ["+resultText+"]")
	ok.Ok(t, strings.Contains(resultText, "@two this is also a line"), "Should have line @two "+resultText+"]")
	ok.Ok(t, strings.Contains(resultText, "@three this is the last line"), "Should have line @three "+resultText+"]")

	Assign(st2, "@out :{{: @text", 5)
	resultText = Expand(st2, "@out")

	ok.Ok(t, strings.Contains(resultText, "1 this is a line"), "Should have line 1 "+resultText)
	ok.Ok(t, strings.Contains(resultText, "2 this is also a line"), "Should have line 2 "+resultText)
	ok.Ok(t, strings.Contains(resultText, "3 this is the last line"), "Should have line 3 "+resultText)

	// Now test evaluating a shorthand file
	st3 := new(SymbolTable)
	ok.Ok(t, len(st3.entries) == 0, "st.entries should be zero")
	ok.Ok(t, len(st3.labels) == 0, "st.labels should be zero too")

	Assign(st3, "#T# :}<: testdata/test1.shorthand", 1)
	expected = true
	results = HasAssignment(st3, "@now")
	ok.Ok(t, results == expected, "Should have @now from :}<:")
	results = HasAssignment(st3, "@title")
	ok.Ok(t, results == expected, "Should have @title from :}<:")
	results = HasAssignment(st3, "@greeting")
	ok.Ok(t, results == expected, "Should have @greeting")
	titleExpansion := Expand(st3, "@title")
	nowExpansion := Expand(st3, "@now")
	greetingExpansion := Expand(st3, "@greeting")

	shorthandText := `
@title @now
@greeting
`

	resultText = Expand(st3, shorthandText)
	ok.Ok(t, strings.Contains(resultText, titleExpansion), "Should have title: "+titleExpansion)
	ok.Ok(t, strings.Contains(resultText, nowExpansion), "Should have name: "+nowExpansion)
	ok.Ok(t, strings.Contains(resultText, greetingExpansion), "Should have greeting: "+greetingExpansion)
}

func TestExpandingValuesToFile(t *testing.T) {
	if _, err := os.Stat("testdata/helloworld1.txt"); err != nil {
		os.Remove("testdata/helloworld1.txt")
	}
	if _, err := os.Stat("testdata/helloworld2.txt"); err != nil {
		os.Remove("testdata/helloworld2.txt")
	}
	st := new(SymbolTable)
	ok.Ok(t, len(st.entries) == 0, "st.entries should be zero")
	ok.Ok(t, len(st.labels) == 0, "st.labels should be zero too")

	a1 := `@hello_world :=: Hello World`
	a2 := `@max :!: echo -n 'Hello Max'`
	e1 := "Hello World"
	e2 := "Hello Max"
	Assign(st, a1, 1)
	Assign(st, `@hello_world :>: testdata/helloworld1.txt`, 2)
	b, err := ioutil.ReadFile("testdata/helloworld1.txt")
	ok.Ok(t, err == nil, "Should beable to read testdata/helloworld1.txt")
	resultText := string(b)
	ok.Ok(t, strings.Contains(resultText, e1), "Should find "+e1+" in "+resultText)
	Assign(st, a2, 3)
	Assign(st, `@hello_world :@>: testdata/helloworld2.txt`, 4)
	b, err = ioutil.ReadFile("testdata/helloworld2.txt")
	ok.Ok(t, err == nil, "Should be able to read testdata/helloworld2.txt")
	resultText = string(b)
	ok.Ok(t, strings.Contains(resultText, e1), "Should find "+e1+" in "+resultText)
	ok.Ok(t, strings.Contains(resultText, e2), "Should find "+e2+" in "+resultText)
}

func TestExpandingAssignmentsToFile(t *testing.T) {
	if _, err := os.Stat("testdata/assigned1.txt"); err != nil {
		os.Remove("testdata/assigned1.txt")
	}
	if _, err := os.Stat("testdata/assigned2.txt"); err != nil {
		os.Remove("testdata/assigned2.txt")
	}
	st := new(SymbolTable)
	ok.Ok(t, len(st.entries) == 0, "st.entries should be zero")
	ok.Ok(t, len(st.labels) == 0, "st.labels should be zero too")

	a1 := `@hello_world :=: Hello World`
	a2 := `@max :!: echo -n 'Hello Max'`
	Assign(st, a1, 1)
	Assign(st, `@hello_world :}>: testdata/assigned1.txt`, 2)
	b, err := ioutil.ReadFile("testdata/assigned1.txt")
	ok.Ok(t, err == nil, "Should beable to read testdata/assigned1.txt")
	resultText := string(b)
	ok.Ok(t, strings.Contains(resultText, a1), "Should have @hello_world assignment in file.")
	Assign(st, a2, 3)
	Assign(st, `_ :@}>: testdata/assigned2.txt`, 4)
	b, err = ioutil.ReadFile("testdata/assigned2.txt")
	ok.Ok(t, err == nil, "Should have all assigments in file.")
	resultText = string(b)
	ok.Ok(t, strings.Contains(resultText, a1), "Should find "+a1)
	ok.Ok(t, strings.Contains(resultText, a2), "Should find "+a2)
}

func TestMarkdownSupport(t *testing.T) {
	st := new(SymbolTable)
	ok.Ok(t, len(st.entries) == 0, "st.entries should be zero")
	ok.Ok(t, len(st.labels) == 0, "st.labels should be zero too")

	s1 := "[my link](http://example.org)"
	Assign(st, fmt.Sprintf("@link :[: %s", s1), 1)
	ok.Ok(t, HasAssignment(st, "@link"), "Should have @link assignment")
	e1 := `<p><a href="http://example.org">my link</a></p>`
	r1 := Expand(st, "@link")
	ok.Ok(t, r1 == e1, fmt.Sprintf("@link shourl render as %s, got %s", e1, r1))

	st2 := new(SymbolTable)
	ok.Ok(t, len(st2.entries) == 0, "st.entries should be zero")
	ok.Ok(t, len(st2.labels) == 0, "st.labels should be zero too")

	s2 := "[@link](@url) is a shorthand link in Markdown."
	a2 := "@link :=: My Link"
	Assign(st2, a2, 1)
	a3 := "@url :=: http://www.example.org"
	Assign(st2, a3, 2)
	e2 := `<p><a href="http://www.example.org">My Link</a> is a shorthand link in Markdown.</p>`
	a4 := `@html :{[: ` + s2
	Assign(st2, a4, 3)
	r2 := Expand(st2, "@html")
	ok.Ok(t, r2 == e2, "Expected ["+e2+"] found ["+r2+"]")
}
