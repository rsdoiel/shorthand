---
title: shorthand
abstract: |-
  shorthand is a Go package and command line tool for doing simple text replacement and expansion. It works as a naive macro language.
authors:
  - family_name: Doiel
    given_name: R. S.
    id: https://orcid.org/0000-0003-0900-6903



repository_code: https://github.com/rsdoiel/shorthand
version: 0.3.0
license_url: https://www.gnu.org/licenses/agpl-3.0.txt

programming_language:
  - Go >= 1.26.4

keywords:
  - text markup
  - text expansion
  - text macros

date_released: 2026-10-20
---

About this software
===================

## shorthand 0.3.0

Maintenance release that brings shorthand up to current workspace practice. The command line no longer uses the third party `caltechlibrary/cli` module (v0.0.16, v0.0.18 was available) and is built on the standard `flag` package, so the module has no dependencies; it requires Go 1.26.4. The options `-h`, `-l`, `-v`, `-i`, `-o`, `-p`, `-n` and `-quiet` keep their long spellings, and `-generate-markdown`, a feature of the dropped module, is gone; the Markdown man page comes from `-help`. The command exits with the workspace codes: 2 for a bad command line, 65 when a line cannot be evaluated, 66 when an input file cannot be opened (it used to exit 0 after reporting it), 73 when the output file cannot be created. A file argument of `-` is standard input. With piped input the prompt and welcome text are no longer mixed into the output; they appear when standard input is a terminal or `-p` is given. Fixed: `-o` did nothing, because `Run` wrote to `os.Stdout` directly. `VirtualMachine` now has `Out` and `Err` writers (default `os.Stdout` and `os.Stderr`) and an `ErrorCount()`, so a program embedding it can send the output anywhere. Fixed: label expansion depended on the order the labels were defined in, which Go randomizes, so with `@a` and `@ab` defined `@ab` came out as the value of `@a` followed by `b` six runs in seven, and a value holding another label expanded in some runs and not in others. `Expand` now makes one pass in which the longest matching label wins and a value is not expanded again (`:expand-expansion:` expands twice), so the result is the same every run. The `_` label is now ignored as documented (a `:set: _ ...` used to replace every underscore in the text that followed), and a definition with no label no longer inserts its value between every character. It is also about twenty times faster on a document with hundreds of labels, because the lookup is built once per change to the definitions rather than once per line. Documentation: the README and tutorial no longer show a `-e` option the command has not had, `shorthand-syntax.md` now lists exactly the operators the code has (a test keeps it so), the stale copies of the help under `docs/` are gone, and the source headers and `copyright.md` say AGPL-3.0-or-later instead of BSD. Snap packaging is dropped; use the curl installer. The command line code moved to `internal/cli` with the help text as the constant `HelpText` in `helptext.go`, and test coverage went from 76% to 92% for the package and 100% for `internal/cli`. Output on the author's decade of timesheets is unchanged.

## Authors

- [R. S. Doiel](https://orcid.org/0000-0003-0900-6903)






shorthand is a Go package and command line tool for doing simple text replacement and expansion. It works as a naive macro language.

- [License](https://www.gnu.org/licenses/agpl-3.0.txt)
- [Code Repository](https://github.com/rsdoiel/shorthand)
  - [Issue Tracker](https://github.com/rsdoiel/shorthand/issues)

## Programming languages

- Go >= 1.26.4




## Software Requirements

- Go >= 1.26.4
- bash (runtime, for the :bash: operators)


## Software Suggestions

- CMTools >= 0.0.45b
- Pandoc >= 3.9
- GNU Make >= 3


