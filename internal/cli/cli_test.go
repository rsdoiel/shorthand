// Copyright (C) 2015-2026 R. S. Doiel
// SPDX-License-Identifier: AGPL-3.0-or-later
// See the LICENSE file.

package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func run(t *testing.T, stdin string, args ...string) (string, string, int) {
	t.Helper()
	var out, errOut strings.Builder
	code := Run("shorthand", args, strings.NewReader(stdin), &out, &errOut)
	return out.String(), errOut.String(), code
}

func write(t *testing.T, dir, name, text string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(text), 0o666); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestStandardOptionsInEverySpelling(t *testing.T) {
	for _, opt := range []string{"-h", "-help", "--help"} {
		out, _, code := run(t, "", opt)
		if code != 0 || !strings.Contains(out, "shorthand(1) user manual") || !strings.Contains(out, "# EXIT STATUS") {
			t.Errorf("%s: exit %d\n%s", opt, code, out)
		}
		for _, token := range []string{"{app_name}", "{version}", "{release_date}", "{release_hash}"} {
			if strings.Contains(out, token) {
				t.Errorf("%s: %s was not replaced", opt, token)
			}
		}
	}
	for _, opt := range []string{"-v", "-version", "--version"} {
		if out, _, code := run(t, "", opt); code != 0 || !strings.HasPrefix(out, "shorthand 0.") {
			t.Errorf("%s: exit %d %q", opt, code, out)
		}
	}
	for _, opt := range []string{"-l", "-license", "--license"} {
		if out, _, code := run(t, "", opt); code != 0 || !strings.Contains(out, "GNU AFFERO GENERAL PUBLIC LICENSE") && !strings.Contains(out, "GNU Affero General Public License") {
			t.Errorf("%s: exit %d %q", opt, code, out[:min(len(out), 80)])
		}
	}
	if out, _, code := run(t, "", "-examples"); code != 0 || !strings.Contains(out, ":import-text:") {
		t.Errorf("-examples: exit %d", code)
	}
}

func TestFilesShareDefinitionsAndStdinIsDash(t *testing.T) {
	dir := t.TempDir()
	defs := write(t, dir, "defs.shorthand", ":set: @who Freda\n")
	page := write(t, dir, "page.txt", "Hello @who\n")
	out, errOut, code := run(t, "From stdin: @who\n", defs, page, "-")
	if code != 0 || errOut != "" {
		t.Fatalf("exit %d, stderr %q", code, errOut)
	}
	if want := "Hello Freda\nFrom stdin: Freda\n"; out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}

func TestInputOptionComesBeforeTheFiles(t *testing.T) {
	dir := t.TempDir()
	defs := write(t, dir, "defs.shorthand", ":set: @who Freda\n")
	page := write(t, dir, "page.txt", "Hello @who\n")
	for _, opt := range []string{"-i", "-input"} {
		out, _, code := run(t, "", opt, defs, page)
		if code != 0 || out != "Hello Freda\n" {
			t.Errorf("%s: exit %d, got %q", opt, code, out)
		}
	}
}

func TestOutputOptionWritesTheFile(t *testing.T) {
	dir := t.TempDir()
	page := write(t, dir, "page.txt", ":set: @who Freda\nHello @who\n")
	target := filepath.Join(dir, "out.txt")
	for _, opt := range []string{"-o", "-output"} {
		out, _, code := run(t, "", opt, target, page)
		if code != 0 || out != "" {
			t.Fatalf("%s: exit %d, stdout %q", opt, code, out)
		}
		if got, _ := os.ReadFile(target); string(got) != "Hello Freda\n" {
			t.Errorf("%s: file holds %q", opt, got)
		}
	}
}

func TestBashOperatorAsUsedForTimesheets(t *testing.T) {
	// The header of a timesheet defines @now; stnparse then reads the result.
	dir := t.TempDir()
	sheet := write(t, dir, "Time_Sheet.txt", ":bash: @now echo -n 3:30\n\nTimesheet\n\n2026-10-16\n\n11:45 - @now; Acme; Widgets; notes\n")
	out, errOut, code := run(t, "", sheet)
	if code != 0 || errOut != "" {
		t.Fatalf("exit %d, stderr %q", code, errOut)
	}
	if !strings.Contains(out, "11:45 - 3:30; Acme; Widgets; notes\n") || strings.Contains(out, ":bash:") {
		t.Errorf("got %q", out)
	}
}

func TestPromptIsOffForPipedInputUnlessAsked(t *testing.T) {
	out, _, _ := run(t, "hello\n")
	if out != "hello\n" {
		t.Errorf("piped input should not be prompted: %q", out)
	}
	out, _, _ = run(t, "hello\n", "-p", "> ")
	if !strings.Contains(out, "> hello\n") || !strings.Contains(out, "Welcome") {
		t.Errorf("-p should turn the prompt on: %q", out)
	}
	out, _, _ = run(t, "hello\n", "-p", "> ", "-n")
	if out != "hello\n" {
		t.Errorf("-n should turn it off: %q", out)
	}
	out, _, _ = run(t, ":help:\n", "-p", "> ", "-n")
	if !strings.Contains(out, ":set:") {
		t.Errorf(":help: lists the operators: %q", out)
	}
}

func TestExitCodes(t *testing.T) {
	dir := t.TempDir()
	missing := filepath.Join(dir, "missing.txt")
	good := write(t, dir, "good.txt", "fine\n")
	bad := write(t, dir, "bad.txt", ":import-text: @x /no/such/file\nafter the error\n")
	cases := []struct {
		name string
		args []string
		want int
	}{
		{"bad flag", []string{"-bogus"}, 2},
		{"missing input", []string{missing}, 66},
		{"missing input among good ones", []string{good, missing}, 66},
		{"quiet still fails", []string{"-quiet", missing}, 66},
		{"a line that cannot be evaluated", []string{bad}, 65},
		{"cannot create the output", []string{"-o", filepath.Join(dir, "no", "dir", "out"), good}, 73},
		{"good file", []string{good}, 0},
	}
	for _, c := range cases {
		if _, _, code := run(t, "", c.args...); code != c.want {
			t.Errorf("%s: exit %d, want %d", c.name, code, c.want)
		}
	}
	out, errOut, _ := run(t, "", bad, good)
	if !strings.Contains(out, "after the error") || !strings.Contains(out, "fine") || !strings.Contains(errOut, "ERROR (1)") {
		t.Errorf("the rest should still be processed\nstdout %q\nstderr %q", out, errOut)
	}
	if _, errOut, _ := run(t, "", "-quiet", missing); errOut != "" {
		t.Errorf("-quiet should not report the file: %q", errOut)
	}
}

func TestIsTerminal(t *testing.T) {
	if isTerminal(strings.NewReader("x")) {
		t.Error("a reader that is not a file is not a terminal")
	}
	f, err := os.Open(filepath.Join(t.TempDir(), "x"))
	if err == nil {
		f.Close()
	}
	dir := t.TempDir()
	plain, err := os.Create(filepath.Join(dir, "plain"))
	if err != nil {
		t.Fatal(err)
	}
	defer plain.Close()
	if isTerminal(plain) {
		t.Error("a regular file is not a terminal")
	}
	if tty, err := os.Open(os.DevNull); err == nil {
		defer tty.Close()
		if !isTerminal(tty) {
			t.Skip("/dev/null is not a character device here")
		}
	}
}
