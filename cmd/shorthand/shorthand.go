// Copyright (C) 2015-2026 R. S. Doiel
// SPDX-License-Identifier: AGPL-3.0-or-later
// See the LICENSE file.

// shorthand is a simple label expander. Run shorthand -help.
package main

import (
	"os"
	"path"

	"github.com/rsdoiel/shorthand/internal/cli"
)

func main() {
	os.Exit(cli.Run(path.Base(os.Args[0]), os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
