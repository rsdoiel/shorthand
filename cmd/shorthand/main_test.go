// Copyright (C) 2015-2026 R. S. Doiel
// SPDX-License-Identifier: AGPL-3.0-or-later
// See the LICENSE file.

package main

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

/** TestProgram builds shorthand and runs it, so main and the wiring to
 * internal/cli are covered.
 *
 * Example:
 *   go test ./cmd/shorthand
 */
func TestProgram(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "shorthand")
	if out, err := exec.Command("go", "build", "-o", bin, ".").CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(":set: @who Freda\nHello @who\n")
	out, err := cmd.Output()
	if err != nil || string(out) != "Hello Freda\n" {
		t.Errorf("piped expansion: %q, %v", out, err)
	}
	if out, err := exec.Command(bin, "-v").Output(); err != nil || !strings.HasPrefix(string(out), "shorthand 0.") {
		t.Errorf("-v: %q, %v", out, err)
	}
	err = exec.Command(bin, "-bogus").Run()
	if ee, ok := err.(*exec.ExitError); !ok || ee.ExitCode() != 2 {
		t.Errorf("a bad flag should exit 2, got %v", err)
	}
}
