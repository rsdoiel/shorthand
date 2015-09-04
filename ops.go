//
// Shorthand package operators - assign a function with the func(sm SourceMap) (SourceMap, error) signature
// and use RegisterOp (e.g. in the New() function) to add support to Shourthand.
//
package shorthand

import (
	"fmt"
	"os"
)

var ExitShorthand = func(sm SourceMap) (SourceMap, error) {
	if sm.Source == "" {
		os.Exit(0)
	}
	fmt.Fprintf(os.Stderr, sm.Source)
	os.Exit(1)
	return sm, nil
}

var AssignStringCallback = func(sm SourceMap) (SourceMap, error) {
	return sm, fmt.Errorf("AssignString() not implemented.")
}

var AssignIncludeCallback = func(sm SourceMap) (SourceMap, error) {
	return sm, fmt.Errorf("AssignInclude() not implemented.")
}

var IncludeAssignmentsCallback = func(sm SourceMap) (SourceMap, error) {
	return sm, fmt.Errorf("IncludeAssignments() not implemented.")
}

var AssignExpansionCallback = func(sm SourceMap) (SourceMap, error) {
	return sm, fmt.Errorf("AssignExpansion() not implemented.")
}

var AssignExpandExpansionCallback = func(sm SourceMap) (SourceMap, error) {
	return sm, fmt.Errorf("AssignExpandExpansion() not implemented.")
}

var IncludeExpansionCallback = func(sm SourceMap) (SourceMap, error) {
	return sm, fmt.Errorf("IncludeExpansion() not implemented.")
}

var AssignShellCallback = func(sm SourceMap) (SourceMap, error) {
	return sm, fmt.Errorf("AssignShell() not implemented.")
}

var AssignExpandShellCallback = func(sm SourceMap) (SourceMap, error) {
	return sm, fmt.Errorf("AssignExpandShell() not implemented.")
}

var AssignMarkdownCallback = func(sm SourceMap) (SourceMap, error) {
	return sm, fmt.Errorf("AssignMarkdown() not implemented.")
}

var AssignExpandMarkdownCallback = func(sm SourceMap) (SourceMap, error) {
	return sm, fmt.Errorf("AssignExpandMarkdown() not implemented.")
}

var IncludeMarkdownCallback = func(sm SourceMap) (SourceMap, error) {
	return sm, fmt.Errorf("IncludeMarkdown() not implemented.")
}

var IncludeExpandMarkdownCallback = func(sm SourceMap) (SourceMap, error) {
	return sm, fmt.Errorf("IncludeExpandMarkdown() not implemented.")
}

var OutputAssignedExpansionCallback = func(sm SourceMap) (SourceMap, error) {
	return sm, fmt.Errorf("OutputAssignedExpansion() not implemented.")
}

var OutputAssignedExpansionsCallback = func(sm SourceMap) (SourceMap, error) {
	return sm, fmt.Errorf("OutputAssignedExpansions() not implemented.")
}

var OutputAssignmentCallback = func(sm SourceMap) (SourceMap, error) {
	return sm, fmt.Errorf("OutputAssignment() not implemented.")
}

var OutputAssignmentsCallback = func(sm SourceMap) (SourceMap, error) {
	return sm, fmt.Errorf("OutputAssignments() not implemented.")
}
