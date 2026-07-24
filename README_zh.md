# go-docx

> 读写 Microsoft Word 2007+（`.docx`，OOXML）的 Go 库——从 Python 的 [python-docx](https://github.com/python-openxml/python-docx) 完整重构而来。

[![Go Reference](https://pkg.go.dev/badge/github.com/SamYue1/go-docx.svg)](https://pkg.go.dev/github.com/SamYue1/go-docx)
[![Go Report Card](https://goreportcard.com/badge/github.com/SamYue1/go-docx)](https://goreportcard.com/report/github.com/SamYue1/go-docx)
[![build](https://img.shields.io/badge/build-go%20test-brightgreen)](#测试)

`go-docx` 是 [python-docx](https://github.com/python-openxml/python-docx)（v1.2.0）的 Go 语言等价重写：在保留原项目**全部公开 API 能力与 `.docx` 往返保真**的前提下，把运行时依赖（`lxml`、Pillow 等）替换为 Go **既有**库（标准库优先），测试与文档一并迁移。重构依据见 [`docs/REFACTORING.md`](docs/REFACTORING.md)。

> **仓库地址**：`https://github.com/SamYue1/go-docx`。
> **上游来源**：`https://github.com/python-openxml/python-docx`。

## 状态

进行中（WIP）。按 [`docs/REFACTORING.md`](docs/REFACTORING.md) §9 分阶段推进：骨架与基础 → 声明式 XML 框架 → OPC 包 → CT 元素类 → 对象层 → 子域补全 → 验收测试对齐 → 文档与 CI。当前尚未发布稳定版。

## 特性（目标等价覆盖 python-docx）

- 打开/保存 `.docx`（含内置默认模板），支持路径与 `io.Reader`/`io.Writer`
- 段落与 Run：文本、加粗/斜体/上下标、字体、颜色、对齐、缩进、行距、break
- 表格：增删行列、合并单元格、样式、对齐、行高规则
- Section：页面尺寸/方向、页边距、分节
- 页眉页脚：首页/奇偶页
- 样式：段落/字符/表样式、潜在样式（latent styles）
- 超链接：地址、文本、runs、break
- 图片：嵌入图片、尺寸与 DPI 自适应（PNG/JPEG/GIF/BMP/TIFF）
- 批注：新增、作者、首字母、`Comment.text`、区间标记
- 核心属性、文档设置、编号部件

完整特性映射见 [`docs/REFACTORING.md`](docs/REFACTORING.md) §5（67 个验收 `.feature` 全保留）。

## 安装

```bash
go get github.com/SamYue1/go-docx
```

要求 Go 1.21+（`go:embed` 已在 1.16 引入；选用 1.21 以获得循环变量作用域等测试便利）。

## 快速开始

对照 python-docx 的 README 示例：

```go
package main

import (
	"fmt"
	"os"

	"github.com/SamYue1/go-docx"
)

func main() {
	doc := docx.NewDocument()
	doc.AddParagraph("It was a dark and stormy night.")
	if err := doc.Save("dark-and-stormy.docx"); err != nil {
		panic(err)
	}

	f, _ := os.Open("dark-and-stormy.docx")
	defer f.Close()
	doc2, err := docx.Open(f)
	if err != nil {
		panic(err)
	}
	fmt.Println(doc2.Paragraphs[0].Text())
	// It was a dark and stormy night.
}
```

## 测试

本库用**测试驱动开发（TDD）**推进重构，风格向 Go 靠拢：

- **单元测试**：`testing` + `github.com/stretchr/testify/assert`，命名沿用原项目 BDD 风格（`it_`/`and_`/`but_` 作为子测试名），表驱动优先。
- **验收测试**：`github.com/cucumber/godog`，原 python-docx 的 67 个 Gherkin `.feature` 原样复用，steps 用 Go 重写。
- **纪律**：每处改动走 红→绿→重构；提交前三者全绿。

```bash
go vet ./...
go test ./...                 # 单元测试
go test ./test/features/...   # 验收测试（godo gherkin）
```

详见 [`docs/REFACTORING.md`](docs/REFACTORING.md) §6（含 TDD 工作流与 Go 测试风格公约）。

## 架构

保持 python-docx 的两层架构：低层 OpenXML（DOM + 声明式 content model，对应原 `xmlchemy`）与高层面向用户的对象层；累积层用 `internal/` 限定，公开 API 经根包 `docx` 再导出。包结构与映射见 [`docs/REFACTORING.md`](docs/REFACTORING.md) §3。

## 文档

文档站点由 Sphinx/rst 迁移到 Markdown + 静态生成器；API 文档由 Go doc 自动生成。迁移方案见 [`docs/REFACTORING.md`](docs/REFACTORING.md) §7。

## 贡献

按 TDD 提交：`go vet ./... && go test ./... && go test ./test/features/...` 全绿为门槛。Lint 用 `golangci-lint`。

## 许可

MIT，与上游 [python-docx](https://github.com/python-openxml/python-docx) 一致。

## 致谢

`go-docx` 源自 Steve Canny 及 [python-openxml/python-docx](https://github.com/python-openxml/python-docx) 维护者的多年工作。本仓库为其 Go 语言重构与等价再造，特此致谢。