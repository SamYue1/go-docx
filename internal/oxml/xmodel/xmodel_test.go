package xmodel

import (
	"testing"

	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
	"github.com/stretchr/testify/assert"
)

func wURI() string {
	return ns.NsMap["w"]
}

func wQn(local string) string {
	return ns.Qn("w:" + local)
}

func wEl(local string) *dom.Element {
	uri := wURI()
	return dom.NewElement(uri, local)
}

func TestDescribeRegistry(t *testing.T) {
	t.Run("it_stores_and_retrieves_child_defs", func(t *testing.T) {
		reg := NewRegistry()
		reg.Add("w:p", ChildDef{Tag: "w:pPr", Kind: ZeroOrOne})
		reg.Add("w:p", ChildDef{Tag: "w:r", Kind: ZeroOrMore})

		defs := reg.Get("w:p")
		assert.Len(t, defs, 2)
		assert.Equal(t, "w:pPr", defs[0].Tag)
		assert.Equal(t, "w:r", defs[1].Tag)
	})

	t.Run("it_returns_nil_for_unknown_tag", func(t *testing.T) {
		reg := NewRegistry()
		assert.Nil(t, reg.Get("w:unknown"))
	})
}

func TestDescribeGetChild(t *testing.T) {
	t.Run("it_finds_first_matching_child", func(t *testing.T) {
		parent := wEl("p")
		r1 := wEl("r")
		r2 := wEl("r")
		pPr := wEl("pPr")
		parent.AddChild(r1)
		parent.AddChild(pPr)
		parent.AddChild(r2)

		result := GetChild(parent, wQn("pPr"))
		assert.Same(t, pPr, result)
	})

	t.Run("it_returns_nil_when_not_found", func(t *testing.T) {
		parent := wEl("p")
		parent.AddChild(wEl("r"))

		result := GetChild(parent, wQn("pPr"))
		assert.Nil(t, result)
	})

	t.Run("it_returns_nil_on_empty_parent", func(t *testing.T) {
		parent := wEl("p")
		result := GetChild(parent, wQn("r"))
		assert.Nil(t, result)
	})
}

func TestDescribeGetChildren(t *testing.T) {
	t.Run("it_finds_all_matching_children", func(t *testing.T) {
		parent := wEl("p")
		r1 := wEl("r")
		r2 := wEl("r")
		pPr := wEl("pPr")
		parent.AddChild(r1)
		parent.AddChild(pPr)
		parent.AddChild(r2)

		results := GetChildren(parent, wQn("r"))
		assert.Len(t, results, 2)
		assert.Same(t, r1, results[0])
		assert.Same(t, r2, results[1])
	})

	t.Run("it_returns_empty_slice_when_not_found", func(t *testing.T) {
		parent := wEl("p")
		parent.AddChild(wEl("pPr"))

		results := GetChildren(parent, wQn("r"))
		assert.Empty(t, results)
	})
}

func TestDescribeInsertChild(t *testing.T) {
	t.Run("it_appends_when_no_successors", func(t *testing.T) {
		parent := wEl("p")
		r1 := wEl("r")
		parent.AddChild(r1)

		child := wEl("pPr")
		InsertChild(parent, child, nil)

		assert.Len(t, parent.Children(), 2)
		assert.Same(t, r1, parent.Children()[0])
		assert.Same(t, child, parent.Children()[1])
		assert.Same(t, parent, child.Parent())
	})

	t.Run("it_inserts_before_first_successor", func(t *testing.T) {
		parent := wEl("p")
		r1 := wEl("r")
		hlink := wEl("hyperlink")
		parent.AddChild(r1)
		parent.AddChild(hlink)

		pPr := wEl("pPr")
		InsertChild(parent, pPr, []string{wQn("r"), wQn("hyperlink")})

		assert.Len(t, parent.Children(), 3)
		assert.Same(t, pPr, parent.Children()[0])
		assert.Same(t, r1, parent.Children()[1])
	})

	t.Run("it_appends_when_no_successor_found", func(t *testing.T) {
		parent := wEl("p")
		r1 := wEl("r")
		parent.AddChild(r1)

		pPr := wEl("pPr")
		InsertChild(parent, pPr, []string{wQn("hyperlink")})

		assert.Len(t, parent.Children(), 2)
		assert.Same(t, r1, parent.Children()[0])
		assert.Same(t, pPr, parent.Children()[1])
	})
}

func TestDescribeRemoveAllChildren(t *testing.T) {
	t.Run("it_removes_matching_children", func(t *testing.T) {
		parent := wEl("p")
		pPr := wEl("pPr")
		r1 := wEl("r")
		r2 := wEl("r")
		parent.AddChild(pPr)
		parent.AddChild(r1)
		parent.AddChild(r2)

		RemoveAllChildren(parent, wQn("r"))

		assert.Len(t, parent.Children(), 1)
		assert.Same(t, pPr, parent.Children()[0])
	})

	t.Run("it_removes_multiple_tags", func(t *testing.T) {
		parent := wEl("p")
		pPr := wEl("pPr")
		r1 := wEl("r")
		hlink := wEl("hyperlink")
		parent.AddChild(pPr)
		parent.AddChild(r1)
		parent.AddChild(hlink)

		RemoveAllChildren(parent, wQn("r"), wQn("hyperlink"))

		assert.Len(t, parent.Children(), 1)
		assert.Same(t, pPr, parent.Children()[0])
	})

	t.Run("it_handles_no_matches", func(t *testing.T) {
		parent := wEl("p")
		r1 := wEl("r")
		parent.AddChild(r1)

		RemoveAllChildren(parent, "nonexistent")
		assert.Len(t, parent.Children(), 1)
	})

	t.Run("it_handles_empty_children", func(t *testing.T) {
		parent := wEl("p")
		RemoveAllChildren(parent, wQn("r"))
		assert.Empty(t, parent.Children())
	})
}

func TestDescribeAddChild(t *testing.T) {
	t.Run("it_creates_and_inserts_child_with_successors", func(t *testing.T) {
		reg := NewRegistry()
		pTag := "w:p"
		rTag := "w:r"
		pPrTag := "w:pPr"
		reg.Add(pTag, ChildDef{Tag: pPrTag, Kind: ZeroOrOne, Successors: []string{rTag}})

		parent := wEl("p")
		r1 := wEl("r")
		parent.AddChild(r1)

		child := AddChild(parent, reg, pTag, pPrTag)
		assert.NotNil(t, child)
		assert.Equal(t, wQn("pPr"), child.ClarkTag())
		assert.Len(t, parent.Children(), 2)
		assert.Same(t, child, parent.Children()[0])
		assert.Same(t, parent, child.Parent())
	})

	t.Run("it_creates_child_without_successors", func(t *testing.T) {
		reg := NewRegistry()
		pTag := "w:p"
		rTag := "w:r"
		reg.Add(pTag, ChildDef{Tag: rTag, Kind: ZeroOrMore})

		parent := wEl("p")
		child := AddChild(parent, reg, pTag, rTag)
		assert.NotNil(t, child)
		assert.Equal(t, wQn("r"), child.ClarkTag())
		assert.Len(t, parent.Children(), 1)
	})
}

func TestDescribeGetOrAddChild(t *testing.T) {
	t.Run("it_returns_existing_child", func(t *testing.T) {
		reg := NewRegistry()
		pTag := "w:p"
		pPrTag := "w:pPr"
		reg.Add(pTag, ChildDef{Tag: pPrTag, Kind: ZeroOrOne})

		parent := wEl("p")
		existing := wEl("pPr")
		parent.AddChild(existing)

		result := GetOrAddChild(parent, reg, pTag, pPrTag)
		assert.Same(t, existing, result)
		assert.Len(t, parent.Children(), 1)
	})

	t.Run("it_creates_child_when_missing", func(t *testing.T) {
		reg := NewRegistry()
		pTag := "w:p"
		pPrTag := "w:pPr"
		reg.Add(pTag, ChildDef{Tag: pPrTag, Kind: ZeroOrOne})

		parent := wEl("p")
		result := GetOrAddChild(parent, reg, pTag, pPrTag)
		assert.NotNil(t, result)
		assert.Equal(t, wQn("pPr"), result.ClarkTag())
		assert.Len(t, parent.Children(), 1)
	})
}
