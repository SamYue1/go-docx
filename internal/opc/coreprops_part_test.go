package opc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDescribeCorePropertiesPart(t *testing.T) {
	t.Run("it_provides_access_to_its_core_props_object", func(t *testing.T) {
		el := NewDefaultCorePropertiesElement()
		partname, _ := NewPackURI("/docProps/core.xml")
		part := NewPart(partname, CT_OPC_CORE_PROPERTIES, serializePartXML(el), NewOpcPackage())

		blob := part.Blob()
		parsedEl, err := parseXML(blob)
		assert.NoError(t, err)
		assert.NotNil(t, parsedEl)

		cp := NewCoreProperties(parsedEl)
		assert.Equal(t, "Word Document", cp.Title())
		assert.Equal(t, "go-docx", cp.LastModifiedBy())
		assert.Equal(t, "1", cp.Revision())
	})

	t.Run("it_can_create_a_default_core_properties_part", func(t *testing.T) {
		pkg := NewOpcPackage()
		cpPart := pkg.createDefaultCorePropertiesPart()
		assert.NotNil(t, cpPart)
		assert.Equal(t, CT_OPC_CORE_PROPERTIES, cpPart.ContentType())

		cp, err := parseXML(cpPart.Blob())
		assert.NoError(t, err)
		coreProps := NewCoreProperties(cp)
		assert.Equal(t, "Word Document", coreProps.Title())
		assert.Equal(t, "go-docx", coreProps.LastModifiedBy())
		assert.Equal(t, "1", coreProps.Revision())
		assert.False(t, coreProps.Modified().IsZero())
		delta := time.Now().UTC().Sub(coreProps.Modified())
		assert.Less(t, delta, 5*time.Second)
	})
}
