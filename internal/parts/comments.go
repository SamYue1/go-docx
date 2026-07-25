package parts

import (
	"github.com/SamYue1/go-docx/internal/opc"
	"github.com/SamYue1/go-docx/internal/oxml"
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
)

// CommentsPart wraps a StoryPart and the CT_Comments XML element for document comments.
type CommentsPart struct {
	*StoryPart
	comments *oxml.CT_Comments
}

// NewCommentsPart creates a new CommentsPart with the given partname, content type, comments element, and package.
func NewCommentsPart(partname opc.PackURI, contentType string, element *oxml.CT_Comments, pkg *opc.OpcPackage) *CommentsPart {
	sp := NewStoryPart(partname, contentType, element.Element, pkg)
	return &CommentsPart{
		StoryPart: sp,
		comments:  element,
	}
}

// DefaultCommentsPart creates a CommentsPart with default settings at /word/comments.xml.
func DefaultCommentsPart(pkg *opc.OpcPackage) *CommentsPart {
	partname, _ := opc.NewPackURI("/word/comments.xml")
	contentType := opc.CT_WML_COMMENTS
	element := oxml.NewCT_Comments()
	return NewCommentsPart(partname, contentType, element, pkg)
}

// Comments returns the CT_Comments XML element.
func (cp *CommentsPart) Comments() *oxml.CT_Comments {
	return cp.comments
}

// DefaultCommentsXML returns the default minimal XML bytes for an empty comments part.
func DefaultCommentsXML() []byte {
	e := dom.NewElement(ns.NsMap["w"], "comments")
	return e.Bytes()
}
