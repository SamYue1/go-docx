# go-docx Development Guide

## Project Status

A Go-equivalent rewrite of `python-docx` (v1.2.0), **no Go source code exists yet** — before §9 Phase 1. The sole design document is `docs/REFACTORING.md`; must read before modifying any code.

## Core Constraints

- **Zero third-party runtime dependencies** — stdlib first. `encoding/xml` for DOM/serialization, `archive/zip` for OPC, `image/*` for image decoding.
- **Go 1.21+** — for `go:embed` (1.16) and per-iteration loop variable scoping (explicit capture needed before 1.22).

## Architecture At a Glance

```
go-docx/                     # public API re-export entry
internal/oxml/{dom,ns,parser,xmodel,stypes,text,table,section,...}
internal/opc/                # OPC zip package read/write
internal/image/              # image header parsing + DPI supplement
internal/parts/              # DocumentPart/ImagePart/...
internal/otext/otable/osect/odoc/  # user-facing object layer
internal/shared/             # Length/Pt/Inches/Emu/RGBColor
internal/enums/              # iota + ToXML/FromXML
internal/tpl/                # default.docx (go:embed)
internal/testutil/           # test data / snippetSeq / mock
test/features/*.feature      # 67 Gherkin files (from upstream)
test/features/steps/*.go     # godog steps
```

Cross-package dependencies are strictly unidirectional (public → internal → oxml/opc), zero circular deps.

## Commands

```bash
go vet ./...                              # must pass
go test ./...                             # unit tests
go test ./test/features/...               # acceptance tests (godog Gherkin)
# All three must be green before commit:
go vet ./... && go test ./... && go test ./test/features/...
```

Currently no `Makefile`, no CI, no `golangci-lint` config. Prioritize adding them when needed.

## Testing Conventions

- **TDD**: red → green → refactor; no "implement first, test later".
- **BDD subtest names**: `TestDescribeXxx` + `t.Run("it_*"/"and_*"/"but_*", ...)`.
- **Table-driven first**: `cases := []struct{...}` + `t.Run(c.name, ...)`.
- **Package boundary**: root `docx` package uses external test package `docx_test`; `internal/*` uses in-package tests to access private members.
- **Assertions**: `testify/assert`, avoid handwritten `if + t.Errorf`.
- **Fakes over mocks**: dependency injection with interfaces + handwritten fakes; `testify/mock` only for large interfaces.
- **Round-trip testing**: `Open→Save` uses byte-level golden files in `testdata/golden/`, updated via `go test -update` (custom flag), reviewed manually.
- **CXEL not ported**: build DOM with `oxml.Elem("w:p", ...)` or parse with `dom.Parse(xmlStr)`.
- **testdata** loaded via `go:embed`; `snippetSeq(name)` reads `testdata/snippets/<name>.txt`.
- Acceptance tests run via `test/features/features_test.go`'s `TestMain` + `godog.TestSuite`.

## DPI Caveat

Go `image/*` decoders **do not return DPI**. Must hand-parse PNG `pHYs` and JPEG `JFIF`/`EXIF` for DPI (`internal/image`), using `encoding/binary` only.

## xmlchemy Equivalent

Python's metaclass auto-generated methods are the core Go challenge. Adopting **Plan A**: thin DOM + declarative registry (`internal/oxml/xmodel`) + explicit CT_ methods. ~130 CT_ classes added incrementally on demand.

## Revision Notes

Extracted from `README.md`, `docs/REFACTORING.md`, and repository state. If docs conflict with executable sources, trust the executable source.
