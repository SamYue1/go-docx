package parts

import (
	"testing"

	"github.com/SamYue1/go-docx/internal/opc"
	"github.com/SamYue1/go-docx/internal/oxml"
	"github.com/stretchr/testify/assert"
)

func TestDescribeCommentsPart(t *testing.T) {
	t.Run("it_creates_default_comments_part", func(t *testing.T) {
		cp := DefaultCommentsPart(nil)
		assert.NotNil(t, cp)
		assert.NotNil(t, cp.Comments())
		assert.Equal(t, "/word/comments.xml", string(cp.Partname()))
	})

	t.Run("it_creates_default_comments_xml", func(t *testing.T) {
		xml := DefaultCommentsXML()
		assert.NotEmpty(t, xml)
		assert.Contains(t, string(xml), "comments")
	})

	t.Run("it_creates_with_correct_content_type", func(t *testing.T) {
		cp := DefaultCommentsPart(nil)
		assert.Equal(t, opc.CT_WML_COMMENTS, cp.ContentType())
	})

	t.Run("it_provides_access_to_its_comments_collection", func(t *testing.T) {
		element := oxml.NewCT_Comments()
		cp := NewCommentsPart("/word/comments.xml", opc.CT_WML_COMMENTS, element, nil)
		comments := cp.Comments()
		assert.NotNil(t, comments)
		assert.Equal(t, element, comments)
	})

	t.Run("it_has_empty_comment_list_by_default", func(t *testing.T) {
		cp := DefaultCommentsPart(nil)
		assert.Empty(t, cp.Comments().Comment_lst())
	})

	t.Run("it_can_add_comments", func(t *testing.T) {
		cp := DefaultCommentsPart(nil)
		comment := oxml.NewCT_Comment(1, "Test Author")
		cp.Comments().Element.AddChild(comment.Element)
		lst := cp.Comments().Comment_lst()
		assert.Len(t, lst, 1)
		author, _ := lst[0].Author()
		assert.Equal(t, "Test Author", author)
		id, _ := lst[0].ID()
		assert.Equal(t, 1, id)
	})

	t.Run("it_constructs_default_part_with_correct_properties", func(t *testing.T) {
		pkg := opc.NewOpcPackage()
		cp := DefaultCommentsPart(pkg)
		assert.Equal(t, opc.PackURI("/word/comments.xml"), cp.Partname())
		assert.Equal(t, opc.CT_WML_COMMENTS, cp.ContentType())
		assert.Equal(t, pkg, cp.Package())
		assert.NotNil(t, cp.Comments())
		assert.Empty(t, cp.Comments().Comment_lst())
	})
}
