# go-docx Development Guide

Go rewrite of `python-docx` (v1.2.0). WIP — unit tests mostly green, acceptance tests partially passing (~200/650 scenarios).

## Dependencies

Two external test-only deps (`go.mod`): `cucumber/godog`, `stretchr/testify`. **Zero runtime deps.**

## Quick Commands

```bash
go vet ./...                     # must pass before commit
go test ./...                    # unit tests (all packages)
go test ./test/features/         # BDD acceptance tests (godog, 67 .feature files)
```

No Makefile, no CI, no golangci-lint config yet.

## Key Architecture

```
docx.go              # public API — type aliases re-exporting internal types
internal/oxml/       # low-level XML: dom/, ns/, xmodel/ (declarative registry), text/, table/, section/…
internal/opc/        # OPC zip package read/write
internal/odoc/       # Document open/save/paragraphs/tables/sections
internal/otext/      # Paragraph/Run/Font/ParagraphFormat/Hyperlink/TabStops
internal/otable/     # Table/Row/Cell/Column
internal/osect/      # Section/HeaderFooter
internal/parts/      # DocumentPart/StylesPart/…
internal/styles/     # Styles/Style/LatentStyles/LatentStyle
internal/shared/     # Length/Inches/Pt/Cm/Emu/Twips/RGBColor
internal/image/      # DPI extraction from PNG/JPEG headers
internal/enums/      # iota enumerations
internal/tpl/        # default.docx (go:embed)
internal/testutil/   # test data helpers
```

Cross-package dep direction: `public → internal → oxml/opc`. No circular deps.

## Testing Conventions

- Subtests named `it_*` / `and_*` / `but_*` per python-docx BDD style.
- `internal/*` tests use in-package access; root `docx` tests use `docx_test` (external).
- Test files live beside source; acceptance step defs in `test/features/steps/steps.go`.
- testdata `.docx` files in `test/features/steps/test_files/`.

## Step Definitions

All BDD step implementations are in the single file `test/features/steps/steps.go` (~4400 lines, growing). Step patterns and feature files follow the upstream python-docx Gherkin spec. When adding a new step:
1. Add Go code to the relevant `internal/` package (e.g., Font property → `internal/otext/font.go`)
2. Wire the step in `steps.go` with `ctx.Step(pattern, handler)`
3. The step handler extracts state from `s.document` / `s.paragraph` / `s.run` etc. (see `featureSuite` struct)

## Common Gotchas

- **`getRelOfType` returns nil, not panic**: OPC relationship lookup (e.g., for styles, core-properties) returns nil if not found. Callers must check nil before dereferencing.
- **All internal types have nil-receiver guards**: `Table`, `Run`, `Section`, `Paragraph`, `Font`, `ParagraphFormat`, etc. Methods check `t == nil || t.field == nil` before accessing inner fields.
- **LatentStyle uses `*bool` (tri-state)**, not plain bool: `nil` = unset, `&true`/`&false` = on/off.
- **`ParagraphFormat.KeepNext`/`KeepTogether`/`PageBreakBefore`/`WidowControl`** use `*bool` tri-state, not bool.
- **TabStops** exist as `internal/otext/tabstops.go` with `AddTabStop`/`ClearAll`/`Get`/`Remove`/`Len` methods.
- **Hyperlink.Address()** requires relationships to be loaded from the OPC package; it returns empty string if the external relationship is not resolved.
- **DPI not available from stdlib** `image/*` decoders — hand-parsed in `internal/image/`.
- **~130 CT_ classes** (XML element proxies in `internal/oxml/`) added incrementally on demand — not auto-generated.

## Project Status

| Metric | Count |
|--------|-------|
| Unit test packages | 19 pass, 0 fail |
| BDD feature files | 67 |
| BDD scenarios pass/fail | ~201 / ~407 |
| BDD undefined steps | ~52 |
| Go files | 200+ throughout `internal/` |

If docs conflict with executable source, trust the executable source.
