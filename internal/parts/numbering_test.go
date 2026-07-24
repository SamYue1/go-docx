package parts

import (
	"testing"

	"github.com/SamYue1/go-docx/internal/opc"
	"github.com/SamYue1/go-docx/internal/oxml"
	"github.com/stretchr/testify/assert"
)

func TestDescribeNumberingPart(t *testing.T) {
	t.Run("it_creates_default_numbering_part", func(t *testing.T) {
		np := DefaultNumberingPart(nil)
		assert.NotNil(t, np)
		assert.NotNil(t, np.Numbering())
		assert.Equal(t, "/word/numbering.xml", string(np.Partname()))
	})

	t.Run("it_creates_with_correct_content_type", func(t *testing.T) {
		np := DefaultNumberingPart(nil)
		assert.Equal(t, opc.CT_WML_NUMBERING, np.ContentType())
	})

	t.Run("it_provides_access_to_numbering_definitions", func(t *testing.T) {
		element := oxml.NewCT_Numbering()
		np := NewNumberingPart("/word/numbering.xml", opc.CT_WML_NUMBERING, element, nil)
		numbering := np.Numbering()
		assert.NotNil(t, numbering)
		assert.Equal(t, element, numbering)
	})

	t.Run("it_has_empty_num_list_by_default", func(t *testing.T) {
		np := DefaultNumberingPart(nil)
		assert.Empty(t, np.Numbering().Num_lst())
	})

	t.Run("it_constructs_default_part_with_correct_properties", func(t *testing.T) {
		pkg := opc.NewOpcPackage()
		np := DefaultNumberingPart(pkg)
		assert.Equal(t, opc.PackURI("/word/numbering.xml"), np.Partname())
		assert.Equal(t, opc.CT_WML_NUMBERING, np.ContentType())
		assert.Equal(t, pkg, np.Package())
		assert.NotNil(t, np.Numbering())
		assert.Empty(t, np.Numbering().Num_lst())
	})
}
