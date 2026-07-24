# python-docx → Go 完整重构方案

> 目标：把 `python-docx`（v1.2.0，读写 Microsoft Word 2007+ `.docx`）完整重构到 Go 语言，等价覆盖**特性 / 测试 / 文档**，并将原项目所依赖的全部基础类库替换为 Go **既有**库（标准库优先，必要时采用 Go 社区成熟库）。本文档不臆造 API，凡引用第三方库均以已验证存在为准。

---

## 1. 重构目标与原则

**目标**

- 功能等价：保留 python-docx 全部公开 API 能力（文档/段落/Run/表格/Section/页眉页脚/样式/超链接/图片/批注/编号/设置/CoreProps 等）。
- 保真读写：`.docx`（OOXML zip 包）往返不丢未知部件、不改元素顺序、不破坏空白与命名空间。
- 双层测试等价迁移：单元测试（pytest → `testing`）+ 验收测试（behave Gherkin → `cucumber/godog`）。
- 文档等价迁移：Sphinx/rst 文档站点 → Go doc + 现代静态文档生成器。

**原则**

1. 基础类库全部替换为 Go **已经存在**的库；标准库优先，社区库仅取成熟、知名者。
2. 不机械翻译：Python 的动态机制（元类、属性描述符、lxml 自定义元素类）在 Go 里用静态、显式、可静态检查的方式重新实现，**保留语义**而非保留语法。
3. 测试与文档随特性一起迁移，命名沿用 BDD 风格以保留原项目意图。
4. 默认模板（`default.docx`）等二进制资源用 `go:embed` 内联。

---

## 2. Python → Go 技术栈映射总表

### 2.1 运行时依赖（`requirements.txt` → Go）

| Python 依赖 | 在原项目中的用途 | Go 替代（既有库） | 说明 |
|---|---|---|---|
| `lxml` | XML DOM、xpath、**自定义元素类（xmlchemy 的核心基石）** | `encoding/xml`（标准库）+ 自研轻量 DOM；可选 `github.com/antchfx/xml` 的 etree/xpath 增强 | lxml 的 `ElementBase` + namespace class lookup 在 Go 无对应，需自研等价"声明式 content model"框架（见 §4.2）。优先纯标准库 tokenizer 自研 DOM 以保证 round-trip 保真。 |
| `typing_extensions` | `Self` 等类型扩展 | Go 类型系统内建 | 不需要。Go 原生支持自引用类型（虽有循环引用限制，可加间接层）。 |

### 2.2 图像处理（原项目已**不依赖 Pillow**，纯手写 header 解析）

| 格式 | 原实现 | Go 替代（标准库） |
|---|---|---|
| PNG | `../python-docx/src/docx/image/png.py` 手解析 IHDR/pHYs | `image/png`（标准库）；DPI 需补读 `pHYs` chunk |
| JPEG | `../python-docx/src/docx/image/jpeg.py` 手解析 SOF/EXIF | `image/jpeg`（标准库）；DPI 需补读 `APP0`/EXIF 分辨率字段 |
| GIF | `../python-docx/src/docx/image/gif.py` | `image/gif`（标准库） |
| BMP | `../python-docx/src/docx/image/bmp.py` | `image/bmp`（标准库，Go 1.9+） |
| TIFF | `../python-docx/src/docx/image/tiff.py` | `image/tiff`（标准库） |

> DPI 解析：python-docx 自行读 PNG `pHYs`、JPEG `JFIF`/`EXIF` 得 DPI；Go 标准库 host 解码不返回 DPI，需在 `internal/image` 内**复刻**这两段解析逻辑（小量代码，纯标准库），不引入第三方。

### 2.3 测试依赖

| Python 工具 | 用途 | Go 替代（既有库） | 说明 |
|---|---|---|---|
| `pytest` | 单元测试 + BDD 收集（`Describe`/`it_`） | `testing`（标准库） + `github.com/stretchr/testify/assert` | BDD 命名用子测试 `t.Run("it_can_add_a_row", ...)` 保留意图。 |
| `behave` | Gherkin `.feature` 验收测试 | `github.com/cucumber/godog` | 已验证：Cucumber for Go，支持 Gherkin，可通过 `TestMain` 与 `go test` 同进程跑（见 §6.2）。 |
| `pyparsing` | CXEL 测试 DSL（`../python-docx/tests/unitutil/cxml.py`） | **不移植** | CXEL 是项目自造 DSL，无 Go 对应；改为 `go:embed` XML 字面片段 + 表驱动（见 §6.3）。 |
| `unittest.mock` | mock 依赖 | 优先：接口 + 手写 fake；补充：`github.com/stretchr/testify/mock`、`go.uber.org/mock`（生成器） | Go 风格偏好显式 fake；动态 mock 用 testify/mock 或 moq。 |
| `pytest-coverage` | 覆盖率 | `go test -cover -coverprofile`（标准库 go 工具链） | 内建，无需第三方。 |
| `pytest-xdist` | 并行 | `go test -parallel` + `t.Parallel()`（标准库） | 内建。 |

### 2.4 文档与工具依赖

| Python 工具 | 用途 | Go 替代（既有库/工具） | 说明 |
|---|---|---|---|
| `setuptools` / `build` | 打包 sdist/wheel | Go module（`go.mod`）+ `GoReleaser`（社区事实标准） | `go build` 产出二进制；库发布即 `go get`。 |
| `twine` | PyPI 上传 | GoReleaser 发布到 GitHub Release / proxy | 无 PyPI 等价，Go 库经 module proxy 发布。 |
| `tox` | 多版本矩阵 | GoRepro / GitHub matrix（`setup-go` 多版本） + `go vet`/`go test` per version；可选 `act` | 无单一对应命令，但 Go 跨版本矩阵通过 CI 直接做。 |
| `ruff` | Lint/格式 | `gofmt` + `go vet`（标准库） + `github.com/golangci/golangci-lint` | golangci-lint 为 Go 生态事实标准 lint 聚合器。 |
| `pyright`（strict） | 类型检查 | `go vet`（标准库） + 可选 `golang.org/x/tools/go/analysis`；类型由编译器静态保证 | Go 编译即类型检查，strict 语义天然达成。 |
| `Sphinx` + `Jinja2` + `MarkupSafe` + `alabaster` | 文档站点 | **`mkdocs-material`**（社区成熟）或 **Hugo**；rst→md 用 `pandoc` 一次性迁移 | Python 文档栈古早版本锁死，Go 生态无关；迁移到 Markdown + 任意静态生成器。 |
| `types-lxml` | 类型 stub | 不需要 | Go 原生类型。 |

### 2.5 标准库映射（Python stdlib → Go stdlib）速查

| 能力 | Python | Go（既有） |
|---|---|---|
| zip 读写 | `zipfile` | `archive/zip` |
| XML | `xml.etree`、`lxml` | `encoding/xml` |
| 哈希 | `hashlib.sha1` | `crypto/sha1` |
| 时间 | `datetime` | `time` |
| 枚举 | `enum` | 显式 `iota` 常量 + 自定义类型（见 §4.4） |
| 嵌入资源 | (`importlib.resources`) | `embed`（`go:embed`） |
| IO | `typing.IO[bytes]` | `io.Reader`/`io.Writer` |
| 路径 | `os.path` | `path`/`filepath` |

---

## 3. 架构映射（两层 → Go 包结构）

python-docx 是两层架构（见 `AGENTS.md` 与架构图），Go 重构保持同构，包结构映射如下：

| Python 包 | 职责 | Go 包 | 可见性 |
|---|---|---|---|
| `docx`（顶层 `api.py`、`__init__.py`） | 公开 API 入口 | `github.com/SamYue1/go-docx`（根包，`docx.Document` 对照） | public |
| `docx.oxml` | 低层 OpenXML（lxml 自定义元素 + xmlchemy） | `internal/oxml` | internal |
| `docx.oxml.ns`/`parser`/`xmlchemy`/`simpletypes` | 命名空间/解析/声明式框架/简单类型 | `internal/oxml/{ns,parser,xmodel,stypes}` | internal |
| `docx.opc` | OOXML 包（zip）加载/序列化、Part、Rel | `internal/opc` | internal |
| `docx.image` | 纯手写图像 header 解析 | `internal/image` | internal |
| `docx.parts` | DocumentPart/ImagePart/StylesPart 等 | `internal/parts` | internal |
| `docx.styles` | 样式集合/工厂 | `internal/styles` | internal |
| `docx.text`、`docx.table`、`docx.section`、`docx.document`、`docx.comments` 等对象层 | 面向用户的对象包装 | `internal/otext`/`internal/otable`/`internal/osect`/`internal/odoc`… → 通过根包 `docx.*` 再导出 | internal→public re-export |
| `docx.enum` | 枚举（`BaseEnum`/`BaseXmlEnum`） | `internal/enums` | internal |
| `docx.shared`（`Length`/`Pt`/`Inches`/`Emu`/`Twips`/`RGBColor`/`lazyproperty`） | 共享类型与工具 | `internal/shared` | internal |
| `docx.templates/default.docx` | 内置默认模板 | `internal/tpl/*.docx` + `go:embed` | internal |
| `../python-docx/tests/`（pytest） | 单元测试 | 与源码同包的 `*_test.go` | test |
| `../python-docx/tests/unitutil`（cxml/file/mock） | 测试工具 | `internal/testutil` + `testdata`（`go:embed`） | internal(test build tag) |
| `../python-docx/features/`（behave） | 验收测试规格 | `test/features/*.feature`（原样复用 67 文件） | test |
| `../python-docx/features/steps/` | behave step 实现 | `test/features/steps/*.go`（godog step） | test |

> 跨包依赖一律自顶向下单向（公开包 → internal/oxml/opc 等），保持无循环依赖（原项目经检测亦无循环依赖）。

---

## 4. 关键机制重构方案

本节处理"Python 用了动态语言特性、Go 无可直接对应"的地方——这是重构最大风险与工作量所在。

### 4.1 命名空间与 QN（`docx.oxml.ns`）

**原项目**：`nsmap`（prefix→uri）、`qn("w:p")→"{uri}p"`（lxml clark notation）、`NamespacePrefixedTag`.

**Go 方案**：纯标准库。

```go
// internal/oxml/ns/ns.go
var Map = map[string]string{
    "w":  "http://schemas.openxmlformats.org/wordprocessingml/2006/main",
    "a":  "http://schemas.openxmlformats.org/drawingml/2006/main",
    // …与 python-docx ../python-docx/src/docx/oxml/ns.py:nsmap 逐项一致
}
var Prefix = map[string]string{} // uri->prefix，init 时反转

func Qn(tag string) (uri, local string) { // 等价 qn()
    idx := strings.IndexByte(tag, ':')
    return Map[tag[:idx]], tag[idx+1:]
}
func Clark(tag string) string { u, l := Qn(tag); return "{" + u + "}" + l }
```

语义与 Python `qn` 一致，无依赖。

### 4.2 xmlchemy 声明式系统（核心难点）

**原项目机制**（`../python-docx/src/docx/oxml/xmlchemy.py`）：
- `BaseOxmlElement` 继承 `lxml.etree.ElementBase`，元类 `MetaOxmlElement` 在类创建时扫描类属性。
- 类属性是 `ZeroOrOne`、`ZeroOrMore`、`OneAndOnlyOne`、`OneOrMore`、`OptionalAttribute`、`RequiredAttribute`、`Choice`、`ZeroOrOneChoice` 之一，描述 content model。
- 元类据此**自动生成**方法：`xxx`（getter）、`xxx_lst`（列表）、`_add_xxx`/`add_xxx`/`_new_xxx`/`_insert_xxx`/`_remove_xxx`/`get_or_add_xxx`/`get_or_change_to_xxx`。
- `successors` 控制插入顺序以遵守 OOXML schema 的元素序列规则。
- 属性经 `simpletypes` 验证后写入。

**Go 无元类、无属性描述符、无 lxml 自定义元素类查找**——必须重新设计，保留语义。

**设计：薄层 DOM + 声明式注册 + 显式方法（推荐方案 A）**

1. **底层 DOM**（保证 round-trip 保真，纯标准库 `encoding/xml` tokenizer）：

```go
// internal/oxml/dom/element.go
type Element struct {
    uri, local string          // 限定名拆分
    attrs  []Attr             // 顺序保留，含命名空间声明
    children []*Element        // 顺序保留，未知元素一样保留
    text   string             // 文本节点（mixed content）
    parent *Element
}
```

解析：用 `xml.Decoder` token 流自建树（保留所有节点、顺序、空白、NS 声明）；序列化：手写 emitter（按 OOXML 习惯 pretty-print）。这取代 lxml，且**无信息丢失**。

2. **CT_ 元素为 DOM 的强类型薄包装**：

```go
// internal/oxml/xmodel/element.go
type CTParagraph struct { *dom.Element }   // 内嵌 DOM 节点
func (p *CTParagraph) PPr() *CTPPr          // getter: findfirst("w:pPr")
func (p *CTParagraph) GetOrAddPPr() *CTPPr  // 对应 get_or_add_pPr
func (p *CTParagraph) RList() []*CTR        // 对应 r_lst
func (p *CTParagraph) AddR() *CTR           // 对应 _add_r / add_r
func (p *CTParagraph) InsertPPr(v *CTPPr)  // 对应 _insert_pPr（successors 控序）
```

3. **声明式注册**（保留 OOXML 序列规则，集中维护，比 Python 元类更显式但减少手写错漏）：

```go
// internal/oxml/xmodel/registry.go
type Child struct {
    Tag        string        // "w:pPr"
    Kind       Kind          // ZeroOrOne|ZeroOrMore|OneAndOnlyOne|OneOrMore
    Successors []string      // ["w:r", "w:hyperlink"] 插入顺序约束
}
var Reg = newRegistry()
func init() {
    Reg.Add("w:p", Child{"w:pPr", ZeroOrOne, nil})
    Reg.Add("w:p", Child{"w:hyperlink", ZeroOrMore, nil})
    Reg.Add("w:p", Child{"w:r", ZeroOrMore, nil})
    // … 与 ../python-docx/src/docx/oxml/__init__.py 的 register_element_cls 一一对齐
}
```

各元素类型在 `internal/oxml/<域>/` 下用代码生成或手写方法；`GetOrAdd/Add/Insert/Remove/GetList` 由 `xmodel` 提供通用实现，CT_ 包装类型转发。**这是 Python 元类"自动生成方法"的等价物**——由通用框架按注册表实现，CT_ 只声明 schema 关系。

**备选方案 B**（元素量大时再评估）：从 `ref/xsd` schema 用 `xgen`/自研代码生成器生成 CT_ 强类型结构 + `MarshalXML/UnmarshalXML`，放弃薄层 DOM 转为完全结构化 struct。代价：对未知/顺序节点保真需要额外 `Any` 字段与原始保留策略，复杂度更高。**首轮重构选 A**。

### 4.3 simpletypes（`docx.oxml.simpletypes`）

**原项目**：`BaseSimpleType` + `from_xml/to_xml/validate/convert_*`，每 `ST_*`（`ST_OnOff`、`ST_HexColor`、`ST_Coordinate`、`ST_UniversalMeasure`…）一个类。

**Go 方案**（接口 + 类型适配）：

```go
// internal/oxml/stypes/stype.go
type Simple interface {
    FromXML(string) (any, error)   // 对应 from_xml
    ToXML(any) (string, error)      // 对应 to_xml（内部调 Validate）
    Validate(any) error             // 对应 validate
}
type OnOff struct{}  // 对应 ST_OnOff
func (OnOff) FromXML(s string) (bool, error) { … }
// ST_HexColor → type RGBColor struct{R,G,B uint8} + FromXML("auto")=AUTO
```

单位换算（`Emu`/`Pt`/`Twips`/通用度量 `ST_UniversalMeasure`）落到 §4.5 的 `internal/shared`。

### 4.4 枚举（`docx.enum.base`）

**原项目**：`BaseEnum`(int)/`BaseXmlEnum`(int)，每成员带 `ms_api_value`+`xml_value`+docstr，提供 `from_xml`/`to_xml`。

**Go 方案**（`iota` 类型 + 方法 + 表）：

```go
// internal/enums/align.go
type WdParagraphAlignment int
const (
    WdAlignLeft   WdParagraphAlignment = iota // 0 (MS API 值)
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

> Go 无内置枚举类型，iota + 自定义类型 + 表是 Go 既有惯例，等价 `BaseXmlEnum`。MS API 整数值与 `xml_value` 用注释或表保留（docstr 转 Go doc 注释）。

### 4.5 shared：长度与单位（`docx.shared`）

**原项目**：`Length(int)` 是 EMU 计数；`EMU=914400/in`、`Pt`/`Inches`/`Cm`/`Mm`/`Twips`/`Emu`/`Px`、`RGBColor`、`lazyproperty`.

**Go 方案**（纯标准库）：

```go
// internal/shared/length.go
type Length int64                // EMU 计数，等价 python Length(int)
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
func RGBFromString(s string) (RGBColor, error) // 等价 RGBColor.from_string，含 "auto" 占位
```

`lazyproperty`（首次访问惰性计算）在 Go 用字段 + `sync.Once` 直接等价，无需独立类型。

### 4.6 OPC 包读写（`docx.opc` + `package.py`）

**原项目**：`PackageReader.from_file` 解 zip → 遍历部件与 rels → `Unmarshaller` 建关系图 → `OpcPackage`；`PackageWriter.write` 反向序列化；`PackURI` 管理部件名；`Relationships` 管理关系；`Package.open` 维护 `default.docx` 模板。

**Go 方案**（标准库 `archive/zip` + `encoding/xml`）：

```go
// internal/opc/package.go
func Open(r io.ReaderAt, size int64) (*Package, error) { // 对照 Package.open
    zr, _ := zip.NewReader(r, size)
    // 1) 找 [Content_Types].xml；2) 找 _rels/.rels；3) 按 rels 遍历加载 parts
    // 4) Unmarshal: 建 part map + Relationships 图（DFS，处理 external）
}
func (p *Package) Save(w io.Writer) error { // 对照 Package.save
    // before_marshal → 按 partname 写 zip：含 [Content_Types].xml、各 part、各 .rels
}
```

- `Part` 抽象：`XmlPart`/`ImagePart`/`CorePropsPart`/`NumberingPart`… 实现 `Marshaler`/`Unmarshaler`。
- `PartFactory` → Go 的 `func(partname, contentType, reltype string, blob []byte, pkg *Package) Part` 工厂函数 + 按 contentType 路由（对应 `PartFactory`）。
- 触发点：读 `[Content_Types].xml` 后，未识别部件**保留原始 blob 走 round-trip**（保真关键）。
- 默认模板：`//go:embed tpl/default.docx` 嵌入。

### 4.7 image 部件（`docx.image`）

**原项目**：`Image.from_file/_from_stream` → 工厂（`_ImageHeaderFactory`）按 magic bytes 选 PNG/JPEG/GIF/BMP/TIFF header 解析 dim+DPI → `ImagePart`（含 SHA1 去重）。

**Go 方案**（标准库 `image/*` + 补 DPI）：

```go
// internal/image/image.go
func FromStream(r io.ReadSeeker) (*Image, error) {
    cfg, fmt, err := image.DecodeConfig(r) // stdlib: 返回 Width/Height + 格式
    dim := Dim{cfg.Width, cfg.Height}
    dpi, err := readDPI(r, fmt)             // 见下，补 PNG pHYs / JPEG JFIF+EXIF
    return &Image{Dim: dim, DPI: dpi, Ext: ext, Sha1: sha1.Blob(blob)}, nil
}
```

- DPI：`readDPI` 仅复刻 python-docx 两段小逻辑（PNG `pHYs` chunk、JPEG `APP0`/`APP1 EXIF` 的 XResolution/YResolution/Unit），纯标准库 `encoding/binary`。
- SHA1：`crypto/sha1`（标准库），用于 `ImageParts._get_by_sha1` 去重。
- Partname 生成：`/word/media/image%d.{ext}`，复刻 `_next_image_partname` 复用空洞编号逻辑。

### 4.8 对象层（document/paragraph/run/table/section/…）

逐类平移：每个对象持有其 CT_ 元素指针，转发语义。例如：

```go
// internal/otext/paragraph.go
type Paragraph struct{ p *oxml.CTParagraph; parent BlockItem }
func (p *Paragraph) Text() string          { /* p.p.Text(): join r/hyperlink text */ }
func (p *Paragraph) AddRun(s string) *Run   { r := p.p.AddR(); return &Run{r, p} }
func (p *Paragraph) Style() (string, bool) { return p.p.Style() }
// Alignment / InnerContent / Hyperlinks / RenderedPageBreaks / IterInnerContent …
```

所有公开对象经根包 `docx` re-export（如 `docx.Paragraph = otext.Paragraph` 或类型别名），保持用户 API `docx.Document(...)` / `.AddParagraph` / `.Tables` 等。

---

## 5. 特性清单（67 个验收 `.feature` 全保留）

`../python-docx/features/*.feature` 共 **67** 文件（22 个域 step 文件，9 个 enum 模块）。Go 重构将 `.feature` 原样搬入 `test/features/`（godog 兼容 Gherkin），steps 用 Go 重写。按域归纳保证不丢特性：

| 域 | 涉及 `.feature`（节选） | 对应公开 API |
|---|---|---|
| API/打开 | `api-open-document` | `docx.Document(path/reader/nil)` |
| 段落-创建/修改 | `par-add-run`、`par-set-text`、`par-clear-paragraph`、`par-insert-paragraph`、`par-alignment-prop`、`par-style-prop`、`par-access-parfmt`、`par-access-inner-content` | `Paragraph.{AddRun,Text,Alignment,Style,ParagraphFormat,IterInnerContent}` |
| Run | `run-add-content`、`run-add-picture`、`run-access-font`、`run-access-inner-content`、`run-clear-run`、`run-char-style`、`run-enum-props`、`txt-add-break`、`txt-font-color`、`txt-font-props` | `Run.{AddText,AddBreak,AddPicture,Font,Style,…}` |
| 表格 | `tbl-add-row-or-col`、`tbl-cell-access/add-table/props/text`、`tbl-col-props`、`tbl-item-access`、`tbl-merge-cells`、`tbl-props`、`tbl-row-props`、`tbl-style`、`blk-add-table`、`doc-add-table` | `Table.{AddRow,AddColumn,Cell,Rows,Columns,Style,Alignment,…}`、`Document.AddTable` |
| Section/页面 | `doc-access-sections`、`doc-add-section`、`sct-section` | `Document.Sections`、`Section.{PageSize,Orientation,Margins,Header,Footer,…}` |
| 页眉页脚 | `hdr-header-footer` | `Section.{Header,Footer,FirstPageHeader,EvenPageHeader,…}` |
| 样式 | `doc-styles`、`sty-access-font`、`sty-access-latent-styles`、`sty-access-parfmt`、`sty-add-style`、`sty-delete-style`、`sty-latent-add-del`、`sty-latent-props`、`sty-style-props` | `Document.Styles`、`Styles.{Add,Delete,…}`、`LatentStyles` |
| 超链接 | `hlk-props`、（段落超链接特性见 `par-*`） | `Hyperlink.{Address,Fragment,URL,Runs,Text,…}` |
| 内部内容/迭代 | `blk-iter-inner-content`、`par-access-inner-content`、`run-access-inner-content` | `BlockItemContainer.IterInnerContent`、`Paragraph/Run/Section.IterInnerContent` |
| 分页 | `par-add-paragraph`、`doc-add-page-break`、`pbk-split-para`、`run-add-*`（page break 相关） | `Document.AddPageBreak`、`Paragraph/Run.ContainsPageBreak`、`RenderedPageBreak` |
| 图片/形状 | `img-characterize-image`、`doc-add-picture`、`shp-inline-shape-access/size` | `Document.AddPicture`、`InlineShape` |
| 批注 | `doc-add-comment`、`doc-comments`、`cmt-mutations`、`cmt-props` | `Document.AddComment`、`Comment.{Author,Initials,Text}`、`Run.MarkCommentRange` |
| 编号 | `num-access-numbering-part` | `DocumentPart.NumberingPart`（内部 API） |
| 文档设置 | `doc-settings` | `Document.Settings`（`odd_and_even_pages_header_footer` 等） |
| 核心属性 | `doc-coreprops` | `Document.CoreProperties`（author/title/modified/last_modified_by…） |
| 标题 | `doc-add-heading` | `Document.AddHeading` |
| 集合访问 | `doc-access-collections` | `Document.{Paragraphs,Tables,…}` |

**覆盖方法**：Go 验收测试从该清单逐项在 `test/features/*.feature` 跑通即可视为特性等价。

---

## 6. 测试迁移方案（TDD 驱动，风格向 Go 靠拢）

> **默认工作流是测试驱动开发（TDD）**：每一处行为改动都先写测试（红）→ 写最小实现使其通过（绿）→ 重构并保持绿色（重构）。下面各小节描述如何把 python-docx 的两层测试迁过来，并据此以 TDD 推进重构。下文 §6.4 给出贯穿两层测试的红绿循环工作流，§6.5 给出 Go 测试的风格公约。

### 6.1 单元测试（pytest → `testing`）

**BDD 命名保留**（原项目 `python_classes=["Describe"]`、`python_functions=["it_","its_","they_","and_","but_"]` 在 `../python-docx/pyproject.toml`，默认 `test_*` 不被收集）：

Go 用子测试名保留语义：

```go
// internal/otext/paragraph_test.go
func TestDescribeParagraph(t *testing.T) {
    t.Run("it_can_add_a_run", func(t *testing.T) { … })
    t.Run("it_can_iterate_its_inner_content", func(t *testing.T) { … })
    t.Run("and_can_clear_content", func(t *testing.T) { … })
}
```

文件命名保持镜像：`internal/otext/paragraph.go ↔ internal/otext/paragraph_test.go`（对照 `.projections.json` 的 alternate 规则）。

断言：`testify/assert`（既有库），不用裸 `if`。

### 6.2 验收测试（behave → godog）

`../python-docx/features/*.feature`（67 文件）**原样**复制到 `test/features/`，godog 完全兼容 Gherkin。step 实现由 `../python-docx/features/steps/*.py` → `test/features/steps/*.go`：

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
// initScenario: godog.T(*ctx) 注册 Given/When/Then，对照 ../python-docx/features/steps/*.py
```

> behave 的 Python step 表达式（正则/parse）需逐条用 godog `ctx.Step(regexp, fn)` 适配；多数可直接套用，少量含 Python 友好语法（如 `text=`）改写为 Go 正则捕获组。`../python-docx/features/environment.py` 的 `before_all`（建 `_scratch`）用 godog `TestSuiteInitializer` 的 `BeforeSuite` 等价。

运行：`go test ./test/features/...` 与单元测试同进程。

### 6.3 测试数据与 CXEL 处理

python-docx 用 `../python-docx/tests/unitutil/cxml.py`（CXEL，基于 `pyparsing`）把 `"w:p/w:r"` 等紧凑串解析成 oxml 元素树，并用 `snippet_seq("name")` 读 `../python-docx/tests/test_files/snippets/<name>.txt`（空行分段）。

Go 重构**不移植 CXEL**，改用熟标准库方案：

1. **`testdata` 用 `go:embed`**：

```
test/testdata/snippets/add-row-col.txt   ← 原样搬 ../python-docx/tests/test_files/snippets/*
test/testdata/*.docx                     ← 原样搬 ../python-docx/tests/test_files/*.docx / png / jpg...
```

```go
//go:embed testdata/snippets/*.txt
var snippetsFS embed.FS
func snippetSeq(name string) []string { b,_ := snippetsFS.ReadFile("testdata/snippets/"+name+".txt"); return strings.Split(string(b), "\n\n") }
```

2. **DOM 构造 helper**（取代 CXEL 的 `element(...)`）：用 `oxml.Elem("w:p", oxml.Elem("w:r", ...))` 直接建树，或对仍以 XML 文本表达的断言用 `dom.Parse(string)`。无需引入解析器组合子库。

**mock**（取代 `unittest.mock`）：`../python-docx/tests/unitutil/mock.py` 提供 `class_mock`/`instance_mock`/`property_mock`。Go：

- 优先定义接口 + 手写 fake（Go 惯例）。
- 动态 stub 用 `testify/mock`（既有库）：`m := new(MockFoo); m.On("Bar").Return(baz)`；生成器 `go.uber.org/mock`（既有）可按接口生成 mock。

### 6.4 TDD 工作流（红 → 绿 → 重构）

重构按 §9 的阶段推进，**每个阶段、每个改动**都走测试驱动。两层测试各自有红绿循环：

**单元层（`testing`，先于实现）**

1. **红**：在镜像文件 `xxx_test.go` 里用表驱动子测试写出期望行为，命名沿用 BDD（`it_`/`and_`/`but_`…）。`go test` 应失败（编译错误或断言不通过）。
2. **绿**：在 `xxx.go` 写**最小**实现让测试通过，不多做。
3. **重构**：在保持绿色前提下整理代码（提取 helper、注册到 `xmodel`、强化保真）。

```go
// internal/shared/length_test.go —— 红：先写期望
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
// 绿：再写 Inches()；绿后把 EMUsPerInch 提为常量（重构）。
```

**验收层（godog，规格先于实现）**

1. **红**：把对应 `.feature`（来自原 `../python-docx/features/*.feature`，原样不改 Gherkin）放入 `test/features/` 并给出 step 签名；步骤实现先 `t.Fatal("not implemented")`。`go test ./test/features/...` 应红。
2. **绿**：在 `test/features/steps/<域>.go` 调用已实现的对象层 API 让步骤通过。
3. **重构**：把 step 里的构造/断言收敛到 `internal/testutil` 的 helper，保持场景可读、与 Python step 语义一致。

> 这与 python-docx 既有"先写 `.feature` 再跑红的 `xfail`"实践（见 git log 中 `xfail: acceptance test for ...`）一脉相承——Go 侧用 godog 的失败步骤直接表达，不用 `xfail` 标记。

### 6.5 Go 测试风格公约

风格向 Go 靠拢，而非照搬 Python 测试习惯：

- **表驱动优先**：多输入/多分支用 `cases := []struct{name string; ...}{...}` + `t.Run(c.name, …)`，顺承 BDD 名（`it_can_add_a_row`）做子测试名；少写一函数一 `Test*`。
- **包边界**：公开 API（根包 `docx`）用**外部测试包** `package docx_test`，只能通过导出 API 验证，防止依赖内部细节；`internal/*` 用同包测试 `package oxml`，需访问私有成员时直接用（Go 允许同包）。
- **并行**：无共享状态的子测试加 `t.Parallel()`；**注意** `t.Parallel()` 会与闭包捕获的循环变量交互，循环变量需显式捕获或用 Go 1.22+ 的 per-iteration 作用域。
- **helper 与构造**：重复构造放 `t.Helper()` 标注的 helper；DOM 构造用 `oxml.Elem("w:p", …)`，断言对照片段用 `snippetSeq(name)`——**不**回退移植 pyparsing/CXEL。
- **断言**：用 `testify/assert`（既有库）做相等/包含类断言，避免大段手写 `if + t.Errorf`；但**不为每个比较都引第三方**，语义清楚处直接用 `cmp` 思路的零依赖写法也可。
- **fake 优先于 mock**：依赖通过**接口**注入，先写手写 fake（Go 惯例、零依赖）；接口面较大或需动态行为时再用 `testify/mock` 或 `go.uber.org/mock` 生成——值类型、struct、具名包类型**不要** mock。
- **golden files 评审**：OPC `Open→Save` 的 round-trip 用字节级 golden 文件卡住，需更新快照时用 `go test -update`（自实现 flag，读 `os.Getenv`/`flag.Bool`，不引第三方）显式触发并人眼复核；golden 产物提交进 `testdata/golden/`。
- **TDD 纪律**：禁止"先实现再补测试"；提交前 `go vet ./... && go test ./... && go test ./test/features/...` 三者全绿，作为 §9 每阶段的验收门槛。

---

## 7. 文档迁移方案

| 原（Python/Sphinx） | 目标（Go/Markdown 生态） |
|---|---|
| `docs/`（Sphinx 1.8.6 / Jinja2 2.11.3 / MarkupSafe 0.23 / alabaster）rst 源 | `pandoc` 一次性 rst→md，迁到 `docs/` mkdocs-material 项目 |
| API 文档（`DocsPageFormatter` 为枚举生成 rst） | Go doc 注释（`// Func ...`）天然随包生成；枚举文档用 Go doc |
| `make docs`（`sphinx-build`） | `mkdocs serve` / `mkdocs build`；或 Hugo |
| `.readthedocs.yaml` | mkdocs material site config（GitHub Pages / Netlify 部署） |
| `HISTORY.rst` | 改 `CHANGELOG.md`（Keep a Changelog 格式） |
| `README.md` | 平移，示例改为 Go 调用 |

> 原项目文档栈依赖被锁死在古早版本（AGENTS.md 已记），迁移到 Markdown 生态同时彻底摆脱该负担。`DocsPageFormatter`（`enum/base.py`）这种给枚举生成 rst 的逻辑在 Go 里不再需要——`go doc` 自动渲染常量与其注释。

---

## 8. Go 模块包结构（建议）

```
go-docx/
  go.mod                         # module github.com/SamYue1/go-docx
  docx.go                        # 公开入口：Document()/Document type
  docx_test.go
  cmd/
    go-docx/                     # （可选）CLI：读写/转文本/校验
  internal/
    oxml/
      dom/                       # 轻量 DOM（encoding/xml tokenizer，round-trip 保真）
      ns/                        # 命名空间 + Qn/Clark
      parser/                    # parse_xml / OxmlElement 等价
      xmodel/                    # xmlchemy 等价：声明式注册 + GetOrAdd/Add/Insert/...
      stypes/                    # simpletypes 等价
      text/  table/  section/    # CT_P / CT_Tbl / CT_SectPr ... 各域元素
      styles/ comments/ settings/ coreprops/ ...
    opc/                         # archive/zip 包加载/序列化、Part/Rel/PackURI
    image/                       # image/* + DPI 补齐 + SHA1 去重
    parts/                       # DocumentPart/ImagePart/StylesPart...
    styles/                     # 样式集合/工厂
    otext/ otable/ osect/ odoc/  # 对象层
    shared/                      # Length/Pt/Inches/Twips/Emu/RGBColor
    enums/                       # iota 枚举 + ToXML/FromXML
    tpl/  default.docx (go:embed)
    testutil/                    # testdata embed / snippetSeq / mock helpers
  test/
    features/*.feature           # 原样 67 个
    features/steps/*.go          # godog step
    features/features_test.go    # TestMain + TestSuite
    testdata/                    # snippets/*.txt + *.docx + 图片
  docs/                          # mkdocs-material（迁移自 Sphinx）
  .github/workflows/             # Go matrix(trims 1.21–latest) + golangci-lint
```

---

## 9. 分阶段迁移计划

依赖自底向上，每阶段产出可独立验证。

1. **骨架与基础**（§4.1、4.5）：`go.mod`、`internal/oxml/ns`、`internal/shared`（Length/单位/RGBColor）、`internal/oxml/dom`（解析+保真序列化）、`internal/enums` 框架。验证：单测覆盖 `qn`/长度换算/round-trip XML 不丢字节。
2. **声明式框架**（§4.2、4.3）：`internal/oxml/xmodel`（注册表 + 通用 GetOrAdd/Add/Insert/Remove/List/Choice）+ `internal/oxml/stypes`（ST_OnOff/HexColor/Coordinate/UniversalMeasure…）。验证：用 1~2 个示例元素（CT_P/CT_R）跑通声明式增删改。
3. **OPC 包层**（§4.6）：`internal/opc` + `internal/parts` 基类 + `go:embed tpl/default.docx`。验证：`Open` 默认模板 → `Save` 字节级保真回写（对照原 `default.docx`）。
4. **核心 CT_ 元素类迁移**（§4.2 表）：按 `oxml/__init__.py` 注册顺序——document/body、text(paragraph/run/parfmt/font/hyperlink)、table、section、styles、settings 核心类。验证：每域配套 unit tests（BDD 命名）。
5. **对象层 + 公开 API**（§4.8）：`internal/otext/otable/osect/odoc` + 根包 `docx` re-export。验证：`docx.Document(path)` → `.Paragraphs/.Tables/.Styles/...` 读通；`.AddParagraph/.AddTable` 写回正确。
6. **image / hyperlink / comments / numbering / settings / header-footer**（§4.7 + §5 对应域）：补全 67 feature 覆盖的所有子域。
7. **验收测试对齐**（§6.2）：67 `.feature` 搬入 + steps Go 重写，`go test ./test/features/...` 全绿。期间补 `internal/testutil` mocks。
8. **文档与 CI**（§7 + §2.4）：rst→md 迁 mkdocs；golangci-lint + Go matrix workflow；`CHANGELOG.md`。

每阶段定义验收命令：`go vet ./... && go test ./...`（单测）与 `go test ./test/features/...`（验收）。

---

## 10. 风险与取舍

| 风险 | 影响 | 缓解 |
|---|---|---|
| **xmlchemy 等价框架工作量大** | §4.2 是核心难点，涉及 ~130 个 CT_ 类 | 通用 `xmodel` 框架吸收重复；CT_ 声明注册集中可读；首轮仅实现 67 feature 涉及的元素，其余按需补。必要时评估方案 B（schema codegen）。 |
| **Go 无标准 xpath** | python-docx 用 lxml `xpath`（见 `CT_P.text`/`clear_content`/`lastRenderedPageBreaks`） | 99% xpath 是 `child::*`/`descendant::w:r` 等简单轴，用显式 `Children/FindAll(descend)` 方法替代；极少数复杂查询评估引入 `antchfx/xml`（含 xpath）。 |
| **DOM round-trip 保真** | OS/Word 对未识别节点、元素顺序、空白敏感 | 自研 DOM（§4.2-1）保留全部节点顺序与原始 blob；`Open→Save` 字节级对比测试卡住 default.docx 与多个真实样本。 |
| **DPI 解析不直接由标准库提供** | image 尺寸换算依赖 DPI | 在 `internal/image` 复刻 python-docx 两段小解析（PNG pHYs / JPEG JFIF+EXIF），纯 `encoding/binary`。 |
| **behave → godog step 表达差异** | 67 feature 的 Given/When/Then 句式需逐条适配 | godog 支持正则 step；多数 behave 句式直接可映射；少量含 `text=` 等 Python 友好语法改捕获组。 |
| **lxml `resolve_entities=False` 安全语义** | XXE/实体扩展防护 | 自研 DOM 默认禁用外部实体与 DTD（与原 parser 一致）。 |
| **默认模板二进制嵌入** | `default.docx` 须内联 | `//go:embed tpl/default.docx` 直接内联，跨平台。 |
| **第三方库最小集** | 用户要求"全部换 Go 已有基础库" | 运行时**零第三方**（仅标准库）。测试/工具第三方仅：`cucumber/godog`、`stretchr/testify`、`golangci/golangci-lint`、`go.uber.org/mock`(可选)、`GoReleaser`(CI)；这些均为 Go 生态成熟库，非"重新发明基础库"。 |

---

## 11. 与原项目的对照速查

| 原（python-docx） | 重构（go-docx） |
|---|---|
| `from docx import Document` | `import "github.com/SamYue1/go-docx"` 用 `docx.Document(path)` |
| `docx.oxml.parser.parse_xml/OxmlElement` | `internal/oxml/parser` |
| `docx.oxml.ns.qn` | `internal/oxml/ns.Qn` |
| `docx.oxml.xmlchemy.*` | `internal/oxml/xmodel`（声明式注册 + 通用方法） |
| `docx.oxml.simpletypes.ST_*` | `internal/oxml/stypes` |
| `docx.opc.package.OpcPackage` | `internal/opc.Package` |
| `docx.opc.pkgreader/pkgwriter` | `internal/opc`（reader/writer 基于 archive/zip） |
| `docx.image.image.Image` | `internal/image.Image` |
| `docx.shared.{Emu,Pt,Inches,Twips,RGBColor}` | `internal/shared` |
| `docx.enum.*.{BaseEnum,BaseXmlEnum}` | `internal/enums`（iota + ToXML/FromXML） |
| `Document/Paragraph/Run/Table/Section/Styles/...` | 根包 `docx` 再导出 |
| `make test` / `make accept` | `go test ./...` / `go test ./test/features/...` |
| `uv run ruff` / `uv run pyright` | `golangci-lint run` / `go vet` + 编译器 |
| `make docs`（sphinx） | `mkdocs build` |

---

## 12. 参考来源（可验证）

- 项目结构：`../python-docx/src/docx/`、`../python-docx/tests/`、`../python-docx/features/`、`../python-docx/pyproject.toml`、`../python-docx/Makefile`、`../python-docx/tox.ini`、`../python-docx/requirements*.txt`、`../python-docx/AGENTS.md`、`../python-docx/.projections.json` 实读。
- 核心机制：`../python-docx/src/docx/oxml/{ns,parser,xmlchemy,simpletypes}.py`、`../python-docx/src/docx/oxml/__init__.py`、`../python-docx/src/docx/opc/package.py`、`../python-docx/src/docx/image/*` 实读。
- 特性清单：`../python-docx/features/*.feature`（67 文件）、`../python-docx/features/steps/`（22 文件）实读。
- godog：`/cucumber/godog`（Context7 已验证，`godog.TestSuite`+`TestMain` 与 `go test` 同进程）。
- Go 图像标准库：`image/png`、`image/jpeg`、`image/gif`、`image/bmp`、`image/tiff`。
- Go 容器/哈希/XML：`archive/zip`、`crypto/sha1`、`encoding/xml`、`embed`。

> 本文档、`AGENTS.md` 与 `codegraph` 架构产物三者可共同作为后续实现会话（建议先执行 §9 阶段 1–2 的骨架）的依据。