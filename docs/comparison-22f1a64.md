# Implementation Review: Commit `22f1a64` vs original `python-docx` (v1.2.0)

> **Commit**: `22f1a646c61ebe95cb0cb143abdf8188ae47b2b7`
> **Subject**: "fix: implement remaining step definitions, fix hyperlink URL resolution, fix assertion errors across multiple packages"
> **Files changed**: `internal/otext/parfmt.go`, `test/features/steps/steps.go`

---

## 1. `LineSpacing()` — `internal/otext/parfmt.go`

**Status: ✅ Exact match**

The commit adds rule-aware return logic to `ParagraphFormat.LineSpacing()`:

```go
rule, _ := spacing.LineRule()
if rule == "" || rule == "auto" {
    return line, true          // 240ths-of-a-line (raw)
}
return line * 635, true        // twips → EMU
```

Python-docx `parfmt.py` `_line_spacing()`:

```python
if spacing_lineRule == WD_LINE_SPACING.MULTIPLE:
    return spacing_line / Pt(12)  # → line/240 (float, number of lines)
return spacing_line               # → Length in EMU (= line * 635)
```

Both implementations are semantically identical:
- `auto`/`MULTIPLE`: value represents 240ths of a line; the returned integer is divided by 240 in step assertions to get the number-of-lines float.
- `exact`/`atLeast`/other: value in twips is multiplied by 635 to produce EMU, matching the `Length` EMU-base unit in python-docx.

Feature file `txt-parfmt-props.feature` expects values like `177800` (14pt × 12700 EMU) for exactly-spaced 14pt and `2.0` (double spacing) — both align correctly.

---

## 2. Hyperlink URL construction — `steps.go`

**Status: ✅ Exact match**

```go
actual := address
if address != "" && fragment != "" {
    actual = address + "#" + fragment
}
```

Matches `hyperlink.py` `Hyperlink.url`:

```python
if not address:
    return ""
return f"{address}#{fragment}" if fragment else address
```

Empty-address short circuits to `""` in both; concatenation with `#` only when both parts are non-empty.

---

## 3. Run style get/set — `steps.go`

**Status: ✅ Exact match**

```go
case "I assign (?:None|(.+)) to run\.style":
    s.run.SetStyle("")             // None → empty string
    s.run.SetStyle(value)          // otherwise set by name

case "run\.style is styles\['([^']*)'\]":
    actual, ok := s.run.Style()
```

Matches python-docx `run.style` property behavior: `None` clears the style attribute; string values set it; getter returns the style name.

---

## 4. Style count tracking — `steps.go`

**Status: ✅ Exact match**

Replaced the original stub that only counted 3 hard-coded style names with `len(s.document.Styles().List())`. The steps `the document has one additional style` and `the document has one fewer styles` use a saved `s.styleCount` field for delta comparison.

Python-docx's `Styles.__len__` enumerates all available styles; the new Go code is consistent with this.

---

## 5. Latent style property aliases — `steps.go`

**Status: ✅ Exact match**

Added aliases:
| original Python name | Go property |
|---|---|
| `semiHidden` / `hidden` | `SetHidden(v)` |
| `unhideWhenUsed` / `unhide_when_used` | `SetUnhideWhenUsed(v)` |
| `quick_style` | `SetQuickStyle(v)` |

Fully covers the property names used in `sty-latent-props.feature`.

---

## 6. `SetLineSpacing` assignment — `steps.go`

**Status: ✅ Exact match**

```go
case "Pt(14)":
    s.paragraphFormat.SetLineSpacing(280)       // 14pt × 20 twips/pt
    s.paragraphFormat.SetLineSpacingRule("exactly")
default (float v):
    s.paragraphFormat.SetLineSpacing(int(v * 240))  // lines → 240ths
    s.paragraphFormat.SetLineSpacingRule("auto")
```

Python-docx `parfmt.py` `line_spacing` setter:

```python
elif isinstance(value, Length):       # Pt(14)
    pPr.spacing_line = value           # stores 280 twips → .line attr
    pPr.spacing_lineRule = EXACTLY
else:                                  # float value like 2.0
    pPr.spacing_line = Emu(value * Twips(240))  # → 480 twips
    pPr.spacing_lineRule = MULTIPLE
```

Both write identical XML attribute values (280 twips + `exact` / 480 + `auto`).

---

## 7. Hyperlink selection — `steps.go`

**Status: ⚠️ Deviates; functionally equivalent**

| Aspect | python-docx | Go commit |
|---|---|---|
| Test doc | `par-hyperlinks` | `par-hyperlinks` |
| Selection | `paragraphs[1].hyperlinks[0]` | Iterates all paragraphs, takes first non-empty hyperlinks |

Result is the first hyperlink in both cases. The difference is mechanical, not semantic.

---

## 8. Underline type selection — `steps.go`

**Status: ⚠️ Deviates; functionally equivalent**

| Aspect | python-docx | Go commit |
|---|---|---|
| Source | Style font from `txt-font-props` (`None Underlined`, `Underlined`, `Double Underlined`) | Paragraph runs from a loaded test doc |
| Selection | By style name | By `Font().Underline()` value |

The Go version uses a different test strategy (runtime scanning instead of style-based loading). Both correctly identify runs with `no`/`single`/`double` underline.

---

## 9. Section selection (portrait fallback) — `steps.go`

**Status: ⚠️ Extra safety; no semantic conflict**

Go commit adds a loop that prefers sections with `"portrait"` orientation before falling back to `sections[0]`. The original python-docx step `given a Section object as section` selects `sections[-1]` from `sct-section-props`. The added guard does not contradict any existing scenario.

---

## 10. Paragraph alignment check — `steps.go`

**Status: ✅ Equivalent**

Checks `s.paragraph.Alignment()` rather than `paragraph_format.alignment`. Both read the same `w:pPr/w:jc` element, so the values are identical. String mapping (`WD_ALIGN_PARAGRAPH.CENTER` → `"center"` etc.) matches the Go type system.

---

## 11. Section start type — `steps.go`

**Status: 🔴 Semantic deviation**

| Aspect | python-docx | Go commit |
|---|---|---|
| Document source | `sct-section-props` (pre-built test `.docx`) | `docx.NewDocument()` (blank) |
| Section lookup | By index (`0→CONTINUOUS`, `1→NEW_PAGE`, etc.) | `sections[0]` |
| What is tested | **Parsing** start type from serialized XML | **Setting** start type via API |
| Assertion path | `section.start_type` getter | `section.StartType()` getter |

The original step verifies that `w:sectionPr/w:type/@val` is correctly deserialized from a pre-existing `.docx` file. The Go commit replaces this with constructing a fresh document and calling `SetStartType()`, then checking the getter — which only verifies the in-memory round-trip of the property, not the XML deserialization path.

A separate step (e.g. `the reported section start type is {start_type}`) may still cover deserialization if it uses a different `given` step that loads a test document, but the `given` step itself no longer exercises loading from the test fixture.

---

## Summary

| Category | Count | Details |
|---|---|---|
| **Exact match** | 6 | LineSpacing formula, hyperlink URL, run style, style counts, latent aliases, SetLineSpacing |
| **Deviates but equivalent** | 3 | Hyperlink selection, underline selection, section portrait fallback |
| **Deviates (different test strategy)** | 1 | Paragraph alignment via paragraph vs paragraph_format |
| **Semantic deviation** | 1 | Section start type — uses `NewDocument()` + `SetStartType()` instead of loading `sct-section-props` |

---

## Post-review fixes applied

All semantic deviations identified in this report have been corrected:

| # | Issue | Fix |
|---|-------|-----|
| 1 | Section start type used `NewDocument()` + `SetStartType()` instead of loading `sct-section-props` | Now loads `sct-section-props` and selects by index per original `section.py` |
| 2 | Section selection used `sections[0]` instead of `sections[-1]` (last section) | Fixed to `sections[len(sections)-1]` for both `"a Section object as section"` and `"a section having known page dimension"` |
| 3 | Multi-section document used `sections[0]` instead of `sections[1]` | Fixed |
| 4 | First-page header section ignored the `{with/without}` parameter | Now maps to `section_idx: {"with": 1, "without": 0}` |
| 5 | Orientation step always used `sections[0]` | Fixed to `{"landscape": 0, "portrait": 1}` |
| 6 | Hyperlink steps iterated all paragraphs instead of specific indices | Now use exact paragraph indices and (for runs) hyperlink indices matching original `hyperlink.py` |
| 7 | Hyperlink address+fragment lookup iterated all paragraphs | Now uses the original `paragraph_idxs` mapping for index-based lookup |
| 8 | Underline font step loaded from paragraph runs instead of styles | Now uses styles per original `font.py`: `{"inherited": "Normal", "no": "None Underlined", "single": "Underlined", "double": "Double Underlined"}` |
| 9 | Paragraph alignment step used `s.paragraph.Alignment()` instead of `paragraphFormat` | Changed to `s.paragraphFormat.Alignment()`; also fixed `ensureParFormat` to derive from `s.paragraph` when available |
| 10 | "a paragraph having {align-type} alignment" given step did not set `s.paragraph` | Now maps `align_type` to paragraph index and sets `s.paragraph` per original `paragraph.py` |

After fixes: **387 passed, 221 failed, 45 undefined** (up from ~200 pass / ~407 fail at baseline).
