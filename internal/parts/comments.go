package parts

import (
	"github.com/SamYue1/go-docx/internal/opc"
	"github.com/SamYue1/go-docx/internal/oxml"
	"github.com/SamYue1/go-docx/internal/oxml/dom"
	"github.com/SamYue1/go-docx/internal/oxml/ns"
)

type CommentsPart struct {
	*StoryPart
	comments *oxml.CT_Comments
}

func NewCommentsPart(partname opc.PackURI, contentType string, element *oxml.CT_Comments, pkg *opc.OpcPackage) *CommentsPart {
	sp := NewStoryPart(partname, contentType, element.Element, pkg)
	return &CommentsPart{
		StoryPart: sp,
		comments:  element,
	}
}

func DefaultCommentsPart(pkg *opc.OpcPackage) *CommentsPart {
	partname, _ := opc.NewPackURI("/word/comments.xml")
	contentType := opc.CT_WML_COMMENTS
	element := oxml.NewCT_Comments()
	return NewCommentsPart(partname, contentType, element, pkg)
}

func (cp *CommentsPart) Comments() *oxml.CT_Comments {
	return cp.comments
}

func DefaultCommentsXML() []byte {
	e := dom.NewElement(ns.NsMap["w"], "comments")
	return e.Bytes()
}
