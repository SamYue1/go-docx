# python-docx → Go Complete Refactoring Plan

> Goal: fully refactor `python-docx` (v1.2.0, reading and writing Microsoft Word 2007+ `.docx`) into Go, with equivalent coverage of **features / tests / documentation**, replacing all underlying libraries with Go's **existing** libraries (stdlib first, battle-tested Go community libraries where necessary). This document does not fabricate APIs; all referenced third-party libraries are verified to exist.

---

## 1. Refactoring Goals and Principles

**Goals**

- Feature parity: preserve all of python-docx's public API capabilities (document/paragraph/Run/table/Section/header-footer/style/hyperlink/image/comment/numbering/settings/CoreProps, etc.).
- Faithful read/write: `.docx` (OOXML zip package) round-trip without losing unknown parts, altering element order, or breaking whitespace and namespaces.
- Dual-layer test migration: unit tests (pytest → `testing`) + acceptance tests (behave Gherkin → `cucumber/godog`).
- Documentation migration: Sphinx/rst doc site → Go doc + modern static site generator.

**Principles**

1. Replace all underlying libraries with Go's **existing** libraries; stdlib first, community libraries only if battle-tested and well-known.
2. No mechanical translation: Python's dynamic mechanisms (metaclasses, property descriptors, lxml custom element classes) are re-implemented in Go with static, explicit, statically-checkable code, **preserving semantics** rather than syntax.
3. Tests and documentation migrate alongside features; naming follows BDD style to retain the original project's intent.
4. Binary resources such as the default template (`default.docx`) are inlined with `go:embed`.

---

## 2. Python → Go Technology Stack Mapping

### 2.1 Runtime Dependencies (`requirements.txt` → Go)

| Python Dependency | Use in Original Project | Go Replacement (Existing Libraries) | Notes |
|---|---|---|---|
| `lxml` | XML DOM, xpath, **custom element classes (core of xmlchemy)** | `encoding/xml` (stdlib) + custom lightweight DOM; optional `github.com/antchfx/xml` for etree/xpath enhancements | lxml's `ElementBase` + namespace class lookup has no Go equivalent; need custom "declarative content model" framework (see §4.2). Prefer pure stdlib tokenizer-based custom DOM for round-trip fidelity. |
| `typing_extensions` | `Self` and other type extensions | Go type system built-in | Not needed. Go natively supports self-referencing types (though circular reference limitations exist, addressable with indirection). |

### 2.2 Image Handling (original project already **does not depend on Pillow** — pure handwritten header parsing)

| Format | Original Implementation | Go Replacement (stdlib) |
|---|---|---|
| PNG | `../python-docx/src/docx/image/png.py` — manual IHDR/pHYs parsing | `image/png` (stdlib); DPI requires reading `pHYs` chunk |
| JPEG | `../python-docx/src/docx/image/jpeg.py` — manual SOF/EXIF parsing | `image/jpeg` (stdlib); DPI requires reading `APP0`/EXIF resolution fields |
| GIF | `../python-docx/src/docx/image/gif.py` | `image/gif` (stdlib) |
| BMP | `../python-docx/src/docx/image/bmp.py` | `image/bmp` (stdlib, Go 1.9+) |
| TIFF | `../python-docx/src/docx/image/tiff.py` | `image/tiff` (stdlib) |

> DPI parsing: python-docx reads PNG `pHYs` and JPEG `JFIF`/`EXIF` for DPI; Go stdlib image decoders do not return DPI, so these two parsing routines must be **replicated** in `internal/image` (small amount of code, pure stdlib), without introducing third-party libraries.

### 2.3 Test Dependencies

| Python Tool | Use Case | Go Replacement (Existing Libraries) | Notes |
|---|---|---|---|
| `pytest` | Unit tests + BDD collection (`Describe`/`it_`) | `testing` (stdlib) + `github.com/stretchr/testify/assert` | BDD naming uses subtests `t.Run("it_can_add_a_row", ...)` to preserve intent. |
| `behave` | Gherkin `.feature` acceptance tests | `github.com/cucumber/godog` | Verified: Cucumber for Go, supports Gherkin, runs in-process with `go test` via `TestMain` (see §6.2). |
| `pyparsing` | CXEL test DSL (`../python-docx/tests/unitutil/cxml.py`) | **Not ported** | CXEL is a project-specific DSL with no Go equivalent; replace with `go:embed` XML literal fragments + table-driven tests (see §6.3). |
| `unittest.mock` | Mock dependencies | Prefer: interfaces + handwritten fakes; supplement: `github.com/stretchr/testify/mock`, `go.uber.org/mock` (generator) | Go style favors explicit fakes; dynamic mocking uses testify/mock or moq. |
| `pytest-coverage` | Coverage | `go test -cover -coverprofile` (stdlib go toolchain) | Built-in, no third-party needed. |
| `pytest-xdist` | Parallel execution | `go test -parallel` + `t.Parallel()` (stdlib) | Built-in. |

### 2.4 Documentation and Tooling Dependencies

| Python Tool | Use Case | Go Replacement (Existing Libraries/Tools) | Notes |
|---|---|---|---|
| `setuptools` / `build` | Package sdist/wheel | Go module (`go.mod`) + `GoReleaser` (community de facto standard) | `go build` produces binaries; library distribution is via `go get`. |
| `twine` | PyPI upload | GoReleaser publishes to GitHub Release / proxy | No PyPI equivalent; Go libraries are distributed via module proxy. |
| `tox` | Multi-version matrix | GoRepro / GitHub matrix (`setup-go` multi-version) + `go vet`/`go test` per version; optional `act` | No single equivalent command, but Go cross-version matrix testing is done directly in CI. |
| `ruff` | Lint/format | `gofmt` + `go vet` (stdlib) + `github.com/golangci/golangci-lint` | golangci-lint is the de facto standard lint aggregator in the Go ecosystem. |
| `pyright` (strict) | Type checking | `go vet` (stdlib) + optional `golang.org/x/tools/go/analysis`; types are statically guaranteed by the compiler | Go compilation = type checking; strict semantics are naturally achieved. |
| `Sphinx` + `Jinja2` + `MarkupSafe` + `alabaster` | Documentation site | **`mkdocs-material`** (mature community) or **Hugo**; rst→md via `pandoc` one-time migration | Python doc stack locked to ancient versions; Go ecosystem-agnostic; migrate to Markdown + any static site generator. |
| `types-lxml` | Type stubs | Not needed | Go has native types. |

### 2.5 Stdlib Mapping (Python stdlib → Go stdlib) Quick Reference

| Capability | Python | Go (existing) |
|---|---|---|
| Zip read/write | `zipfile` | `archive/zip` |
| XML | `xml.etree`, `lxml` | `encoding/xml` |
| Hashing | `hashlib.sha1` | `crypto/sha1` |
| Time | `datetime` | `time` |
| Enums | `enum` | Explicit `iota` constants + custom types (see §4.4) |
| Embedded resources | (`importlib.resources`) | `embed` (`go:embed`) |
| IO | `typing.IO[bytes]` | `io.Reader`/`io.Writer` |
| Paths | `os.path` | `path`/`filepath` |

---

## 3. Architecture Mapping (Two Layers → Go Package Structure)

python-docx has a two-layer architecture (see `AGENTS.md` and architecture diagram); the Go refactoring maintains isomorphic structure, with the following package mapping:

| Python Package | Responsibility | Go Package | Visibility |
|---|---|---|---|
| `docx` (top-level `api.py`, `__init__.py`) | Public API entry point | `github.com/SamYue1/go-docx` (root package, maps to `docx.Document`) | public |
| `docx.oxml` | Low-level OpenXML (lxml custom elements + xmlchemy) | `internal/oxml` | internal |
| `docx.oxml.ns`/`parser`/`xmlchemy`/`simpletypes` | Namespace/parsing/declarative framework/simple types | `internal/oxml/{ns,parser,xmodel,stypes}` | internal |
| `docx.opc` | OOXML package (zip) loading/serialization, Part, Rel | `internal/opc` | internal |
| `docx.image` | Pure handwritten image header parsing | `internal/image` | internal |
| `docx.parts` | DocumentPart/ImagePart/StylesPart, etc. | `internal/parts` | internal |
| `docx.styles` | Style collection/factory | `internal/styles` | internal |
| `docx.text`, `docx.table`, `docx.section`, `docx.document`, `docx.comments` etc. (object layer) | User-facing object wrappers | `internal/otext`/`internal/otable`/`internal/osect`/`internal/odoc`… → re-exported through root `docx.*` | internal→public re-export |
| `docx.enum` | Enums (`BaseEnum`/`BaseXmlEnum`) | `internal/enums` | internal |
| `docx.shared` (`Length`/`Pt`/`Inches`/`Emu`/`Twips`/`RGBColor`/`lazyproperty`) | Shared types and utilities | `internal/shared` | internal |
| `docx.templates/default.docx` | Built-in default template | `internal/tpl/*.docx` + `go:embed` | internal |
| `../python-docx/tests/` (pytest) | Unit tests | `*_test.go` alongside source packages | test |
| `../python-docx/tests/unitutil` (cxml/file/mock) | Test utilities | `internal/testutil` + `testdata` (`go:embed`) | internal (test build tag) |
| `../python-docx/features/` (behave) | Acceptance test specs | `test/features/*.feature` (reuse all 67 files as-is) | test |
| `../python-docx/features/steps/` | behave step implementations | `test/features/steps/*.go` (godog steps) | test |

> Cross-package dependencies are strictly top-down unidirectional (public → internal/oxml/opc, etc.), maintaining zero circular dependencies (the original project was also verified to have no circular dependencies).

---

## 4. Key Mechanism Refactoring

This section addresses areas where "Python uses dynamic language features with no direct Go equivalent" — this is where the greatest refactoring risk and effort lies.

### 4.1 Namespaces and QN (`docx.oxml.ns`)

**Original**: `nsmap` (prefix→uri), `qn("w:p")→"{uri}p"` (lxml Clark notation), `NamespacePrefixedTag`.

**Go approach**: pure stdlib.

```go
// internal/oxml/ns/ns.go
var Map = map[string]string{
    "w":  "http://schemas.openxmlformats.org/wordprocessingml/2006/main",
    "a":  "http://schemas.openxmlformats.org/drawingml/2006/main",
    // … identical to python-docx ../python-docx/src/docx/oxml/ns.py:nsmap
}
var Prefix = map[string]string{} // uri->prefix, reversed at init

func Qn(tag string) (uri, local string) { // equivalent to qn()
    idx := strings.IndexByte(tag, ':')
    return Map[tag[:idx]], tag[idx+1:]
}
func Clark(tag string) string { u, l := Qn(tag); return "{" + u + "}" + l }
```

Semantics identical to Python's `qn`, zero dependencies.

### 4.2 xmlchemy Declarative System (Core Challenge)

**Original mechanism** (`../python-docx/src/docx/oxml/xmlchemy.py`):
- `BaseOxmlElement` inherits from `lxml.etree.ElementBase`; metaclass `MetaOxmlElement` scans class attributes at class creation time.
- Class attributes are `ZeroOrOne`, `ZeroOrMore`, `OneAndOnlyOne`, `OneOrMore`, `OptionalAttribute`, `RequiredAttribute`, `Choice`, `ZeroOrOneChoice` — describing the content model.
- The metaclass **auto-generates** methods: `xxx` (getter), `xxx_lst` (list), `_add_xxx`/`add_xxx`/`_new_xxx`/`_insert_xxx`/`_remove_xxx`/`get_or_add_xxx`/`get_or_change_to_xxx`.
- `successors` controls insertion order to comply with OOXML schema element sequence rules.
- Attributes are validated through `simpletypes` before writing.

**Go has no metaclasses, no property descriptors, no lxml custom element class lookup** — must be redesigned while preserving semantics.

**Design: Thin DOM + Declarative Registry + Explicit Methods (Recommended Plan A)**

1. **Underlying DOM** (round-trip fidelity, pure stdlib `encoding/xml` tokenizer):

```go
// internal/oxml/dom/element.go
type Element struct {
    uri, local string          // qualified name parts
    attrs  []Attr              // order-preserving, including namespace declarations
    children []*Element        // order-preserving, unknown elements also preserved
    text   string              // text node (mixed content)
    parent *Element
}
```

Parsing: build tree from `xml.Decoder` token stream (preserve all nodes, order, whitespace, NS declarations); serialization: hand-written emitter (pretty-print per OOXML conventions). This replaces lxml with **no information loss**.

2. **CT_ elements as strongly-typed thin wrappers over DOM**:

```go
// internal/oxml/xmodel/element.go
type CTParagraph struct { *dom.Element }   // embeds DOM node
func (p *CTParagraph) PPr() *CTPPr          // getter: findfirst("w:pPr")
func (p *CTParagraph) GetOrAddPPr() *CTPPr  // corresponds to get_or_add_pPr
func (p *CTParagraph) RList() []*CTR        // corresponds to r_lst
func (p *CTParagraph) AddR() *CTR           // corresponds to _add_r / add_r
func (p *CTParagraph) InsertPPr(v *CTPPr)  // corresponds to _insert_pPr (controlled by successors)
```

3. **Declarative Registry** (preserves OOXML sequence rules, centrally maintained, more explicit than Python metaclasses but reduces manual errors):

```go
// internal/oxml/xmodel/registry.go
type Child struct {
    Tag        string        // "w:pPr"
    Kind       Kind          // ZeroOrOne|ZeroOrMore|OneAndOnlyOne|OneOrMore
    Successors []string      // ["w:r", "w:hyperlink"] insertion order constraints
}
var Reg = newRegistry()
func init() {
    Reg.Add("w:p", Child{"w:pPr", ZeroOrOne, nil})
    Reg.Add("w:p", Child{"w:hyperlink", ZeroOrMore, nil})
    Reg.Add("w:p", Child{"w:r", ZeroOrMore, nil})
    // … aligned one-to-one with register_element_cls in ../python-docx/src/docx/oxml/__init__.py
}
```

Each element type lives in `internal/oxml/<domain>/` with either code-generated or handwritten methods; `GetOrAdd/Add/Insert/Remove/GetList` are provided by the `xmodel` generic implementation, with CT_ wrapper types forwarding calls. **This is the equivalent of Python's metaclass "auto-generated methods"** — implemented by a generic framework driven by the registry, with CT_ types only declaring schema relationships.

**Fallback Plan B** (evaluate when the number of elements grows large): generate CT_ strongly-typed structs + `MarshalXML/UnmarshalXML` from `ref/xsd` schemas using `xgen`/custom code generator, abandoning the thin DOM approach for fully structured structs. Cost: unknown/out-of-order node fidelity requires extra `Any` fields and raw preservation strategies, adding complexity. **Plan A is chosen for the initial refactoring**.

### 4.3 Simple Types (`docx.oxml.simpletypes`)

**Original**: `BaseSimpleType` + `from_xml/to_xml/validate/convert_*`, one class per `ST_*` (`ST_OnOff`, `ST_HexColor`, `ST_Coordinate`, `ST_UniversalMeasure`…).

**Go approach** (interface + type adapters):

```go
// internal/oxml/stypes/stype.go
type Simple interface {
    FromXML(string) (any, error)   // corresponds to from_xml
    ToXML(any) (string, error)      // corresponds to to_xml (calls Validate internally)
    Validate(any) error             // corresponds to validate
}
type OnOff struct{}  // corresponds to ST_OnOff
func (OnOff) FromXML(s string) (bool, error) { … }
// ST_HexColor → type RGBColor struct{R,G,B uint8} + FromXML("auto")=AUTO
```

Unit conversion (`Emu`/`Pt`/`Twips`/universal measure `ST_UniversalMeasure`) goes to §4.5's `internal/shared`.

### 4.4 Enums (`docx.enum.base`)

**Original**: `BaseEnum`(int)/`BaseXmlEnum`(int), each member has `ms_api_value`+`xml_value`+docstr, provides `from_xml`/`to_xml`.

**Go approach** (`iota` types + methods + lookup tables):

```go
// internal/enums/align.go
type WdParagraphAlignment int
const (
    WdAlignLeft   WdParagraphAlignment = iota // 0 (MS API value)
    WdAlignCenter                              // 1
    WdAlignRight                              // 2
    // …
)
var alignmentXML = map[WdParagraphAlignment]string{
    WdAlignLeft: "left", WdAlignCenter: "center", WdAlignRight: "right",
}
func (v WdParagraphAlignment) ToXML() string  { return alignmentXML[v] }
func AlignmentFromXML(s string) (WdParagraphAlignment, error) { … }
```

> Go has no built-in enum type; iota + custom type + lookup tables is the established Go convention, equivalent to `BaseXmlEnum`. MS API integer values and `xml_value` are preserved via comments or tables (docstr → Go doc comments).

### 4.5 Shared: Lengths and Units (`docx.shared`)

**Original**: `Length(int)` is an EMU count; `EMU=914400/in`, `Pt`/`Inches`/`Cm`/`Mm`/`Twips`/`Emu`/`Px`, `RGBColor`, `lazyproperty`.

**Go approach** (pure stdlib):

```go
// internal/shared/length.go
type Length int64                // EMU count, equivalent to Python Length(int)
const (
  EMUsPerInch Length = 914400
  EMUsPerCm          = 360000
  EMUsPerMm          = 36000
  EMUsPerPt          = 12700
  EMUsPerTwip        = 635
)
func Inches(v float64) Length { return Length(int(v*float64(EMUsPerInch)+0.5)) }
func Pt(v float64) Length    { return Length(int(v*float64(EMUsPerPt)+0.5)) }
// Twips()/Cm()/Mm()/Px()/Emu()…
func (l Length) Inches() float64 { return float64(l) / float64(EMUsPerInch) }
func (l Length) Pt() float64     { return float64(l) / float64(EMUsPerPt) }
// Twips()/Cm()/Mm()/Emu()/Px()…
type RGBColor struct{ R, G, B uint8 }
func RGBFromString(s string) (RGBColor, error) // equivalent to RGBColor.from_string, includes "auto" placeholder
```

`lazyproperty` (lazy evaluation on first access) in Go is directly equivalent to a field + `sync.Once`, without requiring a standalone type.

### 4.6 OPC Package Read/Write (`docx.opc` + `package.py`)

**Original**: `PackageReader.from_file` unzips → iterates parts and rels → `Unmarshaller` builds relationship graph → `OpcPackage`; `PackageWriter.write` reverse-serializes; `PackURI` manages part names; `Relationships` manages relationships; `Package.open` maintains the `default.docx` template.

**Go approach** (stdlib `archive/zip` + `encoding/xml`):

```go
// internal/opc/package.go
func Open(r io.ReaderAt, size int64) (*Package, error) { // maps to Package.open
    zr, _ := zip.NewReader(r, size)
    // 1) find [Content_Types].xml; 2) find _rels/.rels; 3) traverse rels to load parts
    // 4) Unmarshal: build part map + Relationships graph (DFS, handling external)
}
func (p *Package) Save(w io.Writer) error { // maps to Package.save
    // before_marshal → write zip by partname: includes [Content_Types].xml, each part, each .rels
}
```

- `Part` abstraction: `XmlPart`/`ImagePart`/`CorePropsPart`/`NumberingPart`… implement `Marshaler`/`Unmarshaler`.
- `PartFactory` → Go function `func(partname, contentType, reltype string, blob []byte, pkg *Package) Part` with contentType routing (maps to `PartFactory`).
- Trigger: unrecognized parts after reading `[Content_Types].xml` **preserve raw blob for round-trip** (key to fidelity).
- Default template: `//go:embed tpl/default.docx`.

### 4.7 Image Part (`docx.image`)

**Original**: `Image.from_file/_from_stream` → factory (`_ImageHeaderFactory`) selects PNG/JPEG/GIF/BMP/TIFF header parser by magic bytes for dim+DPI → `ImagePart` (with SHA1 dedup).

**Go approach** (stdlib `image/*` + DPI supplement):

```go
// internal/image/image.go
func FromStream(r io.ReadSeeker) (*Image, error) {
    cfg, fmt, err := image.DecodeConfig(r) // stdlib: returns Width/Height + format
    dim := Dim{cfg.Width, cfg.Height}
    dpi, err := readDPI(r, fmt)             // see below: PNG pHYs / JPEG JFIF+EXIF
    return &Image{Dim: dim, DPI: dpi, Ext: ext, Sha1: sha1.Blob(blob)}, nil
}
```

- DPI: `readDPI` replicates only the two small parsing routines from python-docx (PNG `pHYs` chunk, JPEG `APP0`/`APP1 EXIF` XResolution/YResolution/Unit), pure stdlib `encoding/binary`.
- SHA1: `crypto/sha1` (stdlib), used for `ImageParts._get_by_sha1` dedup.
- Partname generation: `/word/media/image%d.{ext}`, replicating `_next_image_partname` gap-reuse logic.

### 4.8 Object Layer (document/paragraph/run/table/section/…)

Each class is a direct translation: each object holds a pointer to its CT_ element and forwards semantic methods. For example:

```go
// internal/otext/paragraph.go
type Paragraph struct{ p *oxml.CTParagraph; parent BlockItem }
func (p *Paragraph) Text() string          { /* p.p.Text(): join r/hyperlink text */ }
func (p *Paragraph) AddRun(s string) *Run   { r := p.p.AddR(); return &Run{r, p} }
func (p *Paragraph) Style() (string, bool) { return p.p.Style() }
// Alignment / InnerContent / Hyperlinks / RenderedPageBreaks / IterInnerContent …
```

All public objects are re-exported through the root `docx` package (e.g., `docx.Paragraph = otext.Paragraph` or type alias), maintaining the user API of `docx.Document(...)` / `.AddParagraph` / `.Tables`, etc.

---

## 5. Feature Checklist (All 67 Acceptance `.feature` Files Preserved)

There are **67** `../python-docx/features/*.feature` files (22 domain step files, 9 enum modules). The Go refactoring copies `.feature` files as-is into `test/features/` (godog is Gherkin-compatible); steps are rewritten in Go. Organized by domain to ensure no features are lost:

| Domain | Involved `.feature` (excerpt) | Corresponding Public API |
|---|---|---|
| API/Open | `api-open-document` | `docx.Document(path/reader/nil)` |
| Paragraph-Create/Modify | `par-add-run`, `par-set-text`, `par-clear-paragraph`, `par-insert-paragraph`, `par-alignment-prop`, `par-style-prop`, `par-access-parfmt`, `par-access-inner-content` | `Paragraph.{AddRun,Text,Alignment,Style,ParagraphFormat,IterInnerContent}` |
| Run | `run-add-content`, `run-add-picture`, `run-access-font`, `run-access-inner-content`, `run-clear-run`, `run-char-style`, `run-enum-props`, `txt-add-break`, `txt-font-color`, `txt-font-props` | `Run.{AddText,AddBreak,AddPicture,Font,Style,…}` |
| Table | `tbl-add-row-or-col`, `tbl-cell-access/add-table/props/text`, `tbl-col-props`, `tbl-item-access`, `tbl-merge-cells`, `tbl-props`, `tbl-row-props`, `tbl-style`, `blk-add-table`, `doc-add-table` | `Table.{AddRow,AddColumn,Cell,Rows,Columns,Style,Alignment,…}`, `Document.AddTable` |
| Section/Page | `doc-access-sections`, `doc-add-section`, `sct-section` | `Document.Sections`, `Section.{PageSize,Orientation,Margins,Header,Footer,…}` |
| Headers/Footers | `hdr-header-footer` | `Section.{Header,Footer,FirstPageHeader,EvenPageHeader,…}` |
| Styles | `doc-styles`, `sty-access-font`, `sty-access-latent-styles`, `sty-access-parfmt`, `sty-add-style`, `sty-delete-style`, `sty-latent-add-del`, `sty-latent-props`, `sty-style-props` | `Document.Styles`, `Styles.{Add,Delete,…}`, `LatentStyles` |
| Hyperlinks | `hlk-props`, (paragraph hyperlink features in `par-*`) | `Hyperlink.{Address,Fragment,URL,Runs,Text,…}` |
| Inner Content/Iteration | `blk-iter-inner-content`, `par-access-inner-content`, `run-access-inner-content` | `BlockItemContainer.IterInnerContent`, `Paragraph/Run/Section.IterInnerContent` |
| Page Breaks | `par-add-paragraph`, `doc-add-page-break`, `pbk-split-para`, `run-add-*` (page break related) | `Document.AddPageBreak`, `Paragraph/Run.ContainsPageBreak`, `RenderedPageBreak` |
| Images/Shapes | `img-characterize-image`, `doc-add-picture`, `shp-inline-shape-access/size` | `Document.AddPicture`, `InlineShape` |
| Comments | `doc-add-comment`, `doc-comments`, `cmt-mutations`, `cmt-props` | `Document.AddComment`, `Comment.{Author,Initials,Text}`, `Run.MarkCommentRange` |
| Numbering | `num-access-numbering-part` | `DocumentPart.NumberingPart` (internal API) |
| Document Settings | `doc-settings` | `Document.Settings` (`odd_and_even_pages_header_footer`, etc.) |
| Core Properties | `doc-coreprops` | `Document.CoreProperties` (author/title/modified/last_modified_by…) |
| Headings | `doc-add-heading` | `Document.AddHeading` |
| Collection Access | `doc-access-collections` | `Document.{Paragraphs,Tables,…}` |

**Coverage method**: Go acceptance tests verify feature parity by running through each item in this checklist via `test/features/*.feature`.

---

## 6. Test Migration Plan (TDD-Driven, Go-Idiomatic Style)

> **Default workflow is test-driven development (TDD)**: every behavioral change starts with a test (red) → minimal implementation to pass (green) → refactor while staying green. The subsections below describe how python-docx's two test layers are migrated and how refactoring proceeds with TDD. §6.4 provides the red-green-refactor workflow across both layers; §6.5 defines Go test style conventions.

### 6.1 Unit Tests (pytest → `testing`)

**BDD naming preserved** (original project uses `python_classes=["Describe"]`, `python_functions=["it_","its_","they_","and_","but_"]` in `../python-docx/pyproject.toml`, default `test_*` not collected):

Go uses subtest names to preserve semantics:

```go
// internal/otext/paragraph_test.go
func TestDescribeParagraph(t *testing.T) {
    t.Run("it_can_add_a_run", func(t *testing.T) { … })
    t.Run("it_can_iterate_its_inner_content", func(t *testing.T) { … })
    t.Run("and_can_clear_content", func(t *testing.T) { … })
}
```

File naming mirrors the source: `internal/otext/paragraph.go ↔ internal/otext/paragraph_test.go` (following `.projections.json` alternate rules).

Assertions: `testify/assert` (existing library), not bare `if`.

### 6.2 Acceptance Tests (behave → godog)

`../python-docx/features/*.feature` (67 files) are **copied as-is** to `test/features/`; godog is fully Gherkin-compatible. Step implementations go from `../python-docx/features/steps/*.py` → `test/features/steps/*.go`:

```go
// test/features/features_test.go
func TestMain(m *testing.M) {
    opts := godog.Options{Format: "progress", Paths: []string{"features"}}
    status := godog.TestSuite{
        Name: "go-docx", TestSuiteInitializer: initTS,
        ScenarioInitializer:  initScenario, Options: &opts,
    }.Run()
    if st := m.Run(); st > status { status = st }
    os.Exit(status)
}
// initScenario: godog.T(*ctx) registers Given/When/Then, mapping to ../python-docx/features/steps/*.py
```

> behave's Python step expressions (regex/parse) must be adapted one by one using godog `ctx.Step(regexp, fn)`; most map directly, a few with Python-friendly syntax (e.g., `text=`) are rewritten as Go regex capture groups. `../python-docx/features/environment.py`'s `before_all` (creating `_scratch`) is handled by godog `TestSuiteInitializer`'s `BeforeSuite`.

Running: `go test ./test/features/...` in the same process as unit tests.

### 6.3 Test Data and CXEL Handling

python-docx uses `../python-docx/tests/unitutil/cxml.py` (CXEL, based on `pyparsing`) to parse compact strings like `"w:p/w:r"` into oxml element trees, and reads `../python-docx/tests/test_files/snippets/<name>.txt` (blank-line delimited) via `snippet_seq("name")`.

The Go refactoring **does not port CXEL**, using instead standard Go approaches:

1. **`testdata` with `go:embed`**:

```
test/testdata/snippets/add-row-col.txt   ← copied as-is from ../python-docx/tests/test_files/snippets/*
test/testdata/*.docx                     ← copied as-is from ../python-docx/tests/test_files/*.docx / png / jpg...
```

```go
//go:embed testdata/snippets/*.txt
var snippetsFS embed.FS
func snippetSeq(name string) []string { b,_ := snippetsFS.ReadFile("testdata/snippets/"+name+".txt"); return strings.Split(string(b), "\n\n") }
```

2. **DOM construction helper** (replacing CXEL's `element(...)`): use `oxml.Elem("w:p", oxml.Elem("w:r", ...))` to build trees directly, or use `dom.Parse(string)` for assertions expressed as XML text. No need for parser combinator libraries.

**Mock** (replacing `unittest.mock`): `../python-docx/tests/unitutil/mock.py` provides `class_mock`/`instance_mock`/`property_mock`. For Go:

- Prefer defining interfaces + handwritten fakes (Go convention).
- Dynamic stubbing with `testify/mock` (existing library): `m := new(MockFoo); m.On("Bar").Return(baz)`. The generator `go.uber.org/mock` (existing) can generate mocks from interfaces.

### 6.4 TDD Workflow (Red → Green → Refactor)

Refactoring progresses per §9's phases, **every phase and every change** is test-driven. Both test layers have independent red-green cycles:

**Unit layer (`testing`, before implementation)**

1. **Red**: in the mirror file `xxx_test.go`, write expected behavior using table-driven subtests with BDD naming (`it_`/`and_`/`but_`…). `go test` should fail (compile error or assertion failure).
2. **Green**: write the **minimal** implementation in `xxx.go` to pass the test, nothing more.
3. **Refactor**: while staying green, clean up code (extract helpers, register in `xmodel`, strengthen fidelity).

```go
// internal/shared/length_test.go —— Red: write expectations first
func TestDescribeLength(t *testing.T) {
    t.Run("it_converts_inches_to_emu", func(t *testing.T) {
        cases := []struct{ in float64; want shared.Length }{
            {1.0, shared.EMUsPerInch},
            {0.5, shared.EMUsPerInch / 2},
        }
        for _, c := range cases {
            got := shared.Inches(c.in)
            assert.Equal(t, c.want, got, "Inches(%v)", c.in)
        }
    })
}
// Green: then write Inches(); after green, extract EMUsPerInch as constant (refactor).
```

**Acceptance layer (godog, spec before implementation)**

1. **Red**: copy the corresponding `.feature` (from original `../python-docx/features/*.feature`, Gherkin unchanged) into `test/features/` and provide step signatures; step implementations start with `t.Fatal("not implemented")`. `go test ./test/features/...` should be red.
2. **Green**: in `test/features/steps/<domain>.go`, call the implemented object-layer APIs to make steps pass.
3. **Refactor**: consolidate construction/assertion logic in `internal/testutil` helpers, keeping scenarios readable and semantically aligned with the Python steps.

> This follows python-docx's existing practice of "write `.feature` first, then run red with `xfail`" (see git log `xfail: acceptance test for ...`) — the Go side uses godog's failed steps directly, without `xfail` markers.

### 6.5 Go Test Style Conventions

Style leans toward Go conventions rather than copying Python testing habits:

- **Table-driven first**: multiple inputs/branches use `cases := []struct{name string; ...}{...}` + `t.Run(c.name, …)`, carrying BDD names (`it_can_add_a_row`) as subtest names; minimize one-function-per-`Test*` patterns.
- **Package boundary**: public API (root `docx` package) uses **external test package** `package docx_test`, verifying only through exported APIs to prevent depending on internals; `internal/*` uses in-package tests `package oxml` for direct access to private members (Go allows same-package access).
- **Parallelism**: subtests with no shared state add `t.Parallel()`; **note** that `t.Parallel()` interacts with closure capture of loop variables — variables must be captured explicitly or use Go 1.22+'s per-iteration scoping.
- **Helpers and construction**: repeated construction goes in `t.Helper()`-annotated helpers; DOM construction uses `oxml.Elem("w:p", …)`, assertions against snippets use `snippetSeq(name)` — **do not** port pyparsing/CXEL.
- **Assertions**: use `testify/assert` (existing library) for equality/containment assertions, avoiding large handwritten `if + t.Errorf` blocks; but **do not import a third-party library for every comparison** — when semantics are clear, zero-dependency `cmp`-style approaches are fine.
- **Fakes over mocks**: dependencies are injected via **interfaces**; prefer handwritten fakes (Go convention, zero dependencies) first; use `testify/mock` or `go.uber.org/mock` for larger interfaces or dynamic behavior — **do not** mock value types, structs, or named package types.
- **Golden file review**: OPC `Open→Save` round-trip uses byte-level golden files; update snapshots via `go test -update` (custom flag, reading `os.Getenv`/`flag.Bool`, no third-party) with explicit human review; golden artifacts committed to `testdata/golden/`.
- **TDD discipline**: no "implement first, test later"; `go vet ./... && go test ./... && go test ./test/features/...` must all pass before commit, serving as the acceptance gate for each §9 phase.

---

## 7. Documentation Migration Plan

| Original (Python/Sphinx) | Target (Go/Markdown ecosystem) |
|---|---|
| `docs/` (Sphinx 1.8.6 / Jinja2 2.11.3 / MarkupSafe 0.23 / alabaster) rst sources | `pandoc` one-time rst→md, migrate to `docs/` mkdocs-material project |
| API docs (`DocsPageFormatter` generates rst for enums) | Go doc comments (`// Func ...`) generated per package; enum docs via Go doc |
| `make docs` (`sphinx-build`) | `mkdocs serve` / `mkdocs build`; or Hugo |
| `.readthedocs.yaml` | mkdocs material site config (GitHub Pages / Netlify deployment) |
| `HISTORY.rst` | Convert to `CHANGELOG.md` (Keep a Changelog format) |
| `README.md` | Translated, examples changed to Go calls |

> The original project's documentation stack is locked to ancient versions (documented in AGENTS.md); migrating to the Markdown ecosystem removes this burden entirely. `DocsPageFormatter` (in `enum/base.py`) — which generates rst for enums — is no longer needed in Go: `go doc` automatically renders constants and their comments.

---

## 8. Go Module Package Structure (Proposed)

```
go-docx/
  go.mod                         # module github.com/SamYue1/go-docx
  docx.go                        # Public entry: Document()/Document type
  docx_test.go
  cmd/
    go-docx/                     # (Optional) CLI: read/write/convert text/validate
  internal/
    oxml/
      dom/                       # Lightweight DOM (encoding/xml tokenizer, round-trip fidelity)
      ns/                        # Namespace + Qn/Clark
      parser/                    # parse_xml / OxmlElement equivalent
      xmodel/                    # xmlchemy equivalent: declarative registry + GetOrAdd/Add/Insert/...
      stypes/                    # simpletypes equivalent
      text/  table/  section/    # CT_P / CT_Tbl / CT_SectPr ... per-domain elements
      styles/ comments/ settings/ coreprops/ ...
    opc/                         # archive/zip package loading/serialization, Part/Rel/PackURI
    image/                       # image/* + DPI supplement + SHA1 dedup
    parts/                       # DocumentPart/ImagePart/StylesPart...
    styles/                      # Style collection/factory
    otext/ otable/ osect/ odoc/  # Object layer
    shared/                      # Length/Pt/Inches/Twips/Emu/RGBColor
    enums/                       # iota enums + ToXML/FromXML
    tpl/  default.docx (go:embed)
    testutil/                    # testdata embed / snippetSeq / mock helpers
  test/
    features/*.feature           # 67 files, as-is
    features/steps/*.go          # godog steps
    features/features_test.go    # TestMain + TestSuite
    testdata/                    # snippets/*.txt + *.docx + images
  docs/                          # mkdocs-material (migrated from Sphinx)
  .github/workflows/             # Go matrix (trims 1.21–latest) + golangci-lint
```

---

## 9. Phased Migration Plan

Dependencies are bottom-up; each phase produces independently verifiable output.

1. **Skeleton and Foundations** (§4.1, 4.5): `go.mod`, `internal/oxml/ns`, `internal/shared` (Length/units/RGBColor), `internal/oxml/dom` (parsing + faithful serialization), `internal/enums` framework. Verification: unit tests covering `qn`/length conversion/round-trip XML byte fidelity.
2. **Declarative Framework** (§4.2, 4.3): `internal/oxml/xmodel` (registry + generic GetOrAdd/Add/Insert/Remove/List/Choice) + `internal/oxml/stypes` (ST_OnOff/HexColor/Coordinate/UniversalMeasure…). Verification: use 1–2 sample elements (CT_P/CT_R) to exercise declarative CRUD.
3. **OPC Package Layer** (§4.6): `internal/opc` + `internal/parts` base classes + `go:embed tpl/default.docx`. Verification: `Open` default template → `Save` byte-level faithful write-back (diff against original `default.docx`).
4. **Core CT_ Element Class Migration** (§4.2 table): following `oxml/__init__.py` registration order — document/body, text(paragraph/run/parfmt/font/hyperlink), table, section, styles, settings, core classes. Verification: unit tests per domain (BDD naming).
5. **Object Layer + Public API** (§4.8): `internal/otext/otable/osect/odoc` + root `docx` re-export. Verification: `docx.Document(path)` → `.Paragraphs/.Tables/.Styles/...` reads correctly; `.AddParagraph/.AddTable` writes back correctly.
6. **Image / Hyperlink / Comments / Numbering / Settings / Header-Footer** (§4.7 + §5 corresponding domains): complete all sub-domains covered by the 67 features.
7. **Acceptance Test Alignment** (§6.2): 67 `.feature` files copied in + steps rewritten in Go; `go test ./test/features/...` all green. Supplement `internal/testutil` mocks as needed.
8. **Documentation and CI** (§7 + §2.4): rst→md migration to mkdocs; golangci-lint + Go matrix workflow; `CHANGELOG.md`.

Each phase defines acceptance commands: `go vet ./... && go test ./...` (unit tests) and `go test ./test/features/...` (acceptance).

---

## 10. Risks and Trade-offs

| Risk | Impact | Mitigation |
|---|---|---|
| **xmlchemy equivalent framework is high-effort** | §4.2 is the core challenge, involving ~130 CT_ classes | Generic `xmodel` framework absorbs repetition; centralized CT_ declaration registry improves readability; initial round implements only elements needed by the 67 features, rest added on demand. Evaluate Plan B (schema codegen) if needed. |
| **No standard xpath in Go** | python-docx uses lxml `xpath` (see `CT_P.text`/`clear_content`/`lastRenderedPageBreaks`) | 99% of xpath queries are simple axes like `child::*`/`descendant::w:r` — replaced by explicit `Children/FindAll(descend)` methods; for the rare complex query, evaluate `antchfx/xml` (includes xpath). |
| **DOM round-trip fidelity** | OS/Word is sensitive to unrecognized nodes, element order, whitespace | Custom DOM (§4.2-1) preserves all node order and raw blob; byte-level `Open→Save` comparison tests lock down `default.docx` and multiple real samples. |
| **DPI parsing not directly available from stdlib** | Image size conversion depends on DPI | Replicate python-docx's two small parsing routines in `internal/image` (PNG pHYs / JPEG JFIF+EXIF), pure `encoding/binary`. |
| **behave → godog step expression differences** | 67 features' Given/When/Then patterns need individual adaptation | godog supports regex steps; most behave patterns map directly; a few with Python-friendly syntax (e.g., `text=`) rewritten as capture groups. |
| **lxml `resolve_entities=False` security semantics** | XXE/entity expansion protection | Custom DOM disables external entities and DTD by default (consistent with original parser). |
| **Default template binary embedding** | `default.docx` must be inlined | `//go:embed tpl/default.docx` directly inlines, cross-platform. |
| **Minimal third-party library set** | User requires "replace all with existing Go base libraries" | **Zero third-party** at runtime (stdlib only). Test/tooling third-party only: `cucumber/godog`, `stretchr/testify`, `golangci/golangci-lint`, `go.uber.org/mock`(optional), `GoReleaser`(CI); all are mature Go ecosystem libraries, not "reinventing base libraries." |

---

## 11. Quick Reference: Original vs. Refactored

| Original (python-docx) | Refactored (go-docx) |
|---|---|
| `from docx import Document` | `import "github.com/SamYue1/go-docx"` using `docx.Document(path)` |
| `docx.oxml.parser.parse_xml/OxmlElement` | `internal/oxml/parser` |
| `docx.oxml.ns.qn` | `internal/oxml/ns.Qn` |
| `docx.oxml.xmlchemy.*` | `internal/oxml/xmodel` (declarative registry + generic methods) |
| `docx.oxml.simpletypes.ST_*` | `internal/oxml/stypes` |
| `docx.opc.package.OpcPackage` | `internal/opc.Package` |
| `docx.opc.pkgreader/pkgwriter` | `internal/opc` (reader/writer based on archive/zip) |
| `docx.image.image.Image` | `internal/image.Image` |
| `docx.shared.{Emu,Pt,Inches,Twips,RGBColor}` | `internal/shared` |
| `docx.enum.*.{BaseEnum,BaseXmlEnum}` | `internal/enums` (iota + ToXML/FromXML) |
| `Document/Paragraph/Run/Table/Section/Styles/...` | Root `docx` package re-exports |
| `make test` / `make accept` | `go test ./...` / `go test ./test/features/...` |
| `uv run ruff` / `uv run pyright` | `golangci-lint run` / `go vet` + compiler |
| `make docs` (sphinx) | `mkdocs build` |

---

## 12. References (Verifiable)

- Project structure: read from `../python-docx/src/docx/`, `../python-docx/tests/`, `../python-docx/features/`, `../python-docx/pyproject.toml`, `../python-docx/Makefile`, `../python-docx/tox.ini`, `../python-docx/requirements*.txt`, `../python-docx/AGENTS.md`, `../python-docx/.projections.json`.
- Core mechanisms: read from `../python-docx/src/docx/oxml/{ns,parser,xmlchemy,simpletypes}.py`, `../python-docx/src/docx/oxml/__init__.py`, `../python-docx/src/docx/opc/package.py`, `../python-docx/src/docx/image/*`.
- Feature checklist: read from `../python-docx/features/*.feature` (67 files), `../python-docx/features/steps/` (22 files).
- godog: `/cucumber/godog` (verified via Context7; `godog.TestSuite`+`TestMain` runs in-process with `go test`).
- Go image stdlib: `image/png`, `image/jpeg`, `image/gif`, `image/bmp`, `image/tiff`.
- Go container/hash/XML: `archive/zip`, `crypto/sha1`, `encoding/xml`, `embed`.

> This document, `AGENTS.md`, and the `codegraph` architecture artifacts together serve as the basis for subsequent implementation sessions (it is recommended to start with §9 phases 1–2: skeleton and foundations).
