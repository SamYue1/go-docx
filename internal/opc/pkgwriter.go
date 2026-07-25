package opc

import (
	"sort"
)

// PackageWriter serialises an in-memory OpcPackage (its relationships and
// parts) into a physical OPC zip archive via a PhysPkgWriter.
type PackageWriter struct{}

// Write writes the complete OPC package: content types stream, package-level
// relationships, and all parts with their per-part relationships.
func (pw *PackageWriter) Write(physWriter PhysPkgWriter, pkgRels *Relationships, parts []*Part) error {
	if err := writeContentTypesStream(physWriter, parts); err != nil {
		return err
	}
	if err := writePkgRels(physWriter, pkgRels); err != nil {
		return err
	}
	if err := writeParts(physWriter, parts); err != nil {
		return err
	}
	return nil
}

// writeContentTypesStream builds the [Content_Types].xml from the given
// parts and writes it to the physical package.
func writeContentTypesStream(physWriter PhysPkgWriter, parts []*Part) error {
	cti := newContentTypesItemFromParts(parts)
	return physWriter.Write(CONTENT_TYPES_URI, cti.blob())
}

// writePkgRels serialises the package-level relationships to XML and writes
// it to the physical package at /_rels/.rels.
func writePkgRels(physWriter PhysPkgWriter, pkgRels *Relationships) error {
	return physWriter.Write(PACKAGE_URI.RelsURI(), pkgRels.XML())
}

// writeParts writes every part's blob and (if non-empty) its relationships
// XML to the physical package.
func writeParts(physWriter PhysPkgWriter, parts []*Part) error {
	for _, part := range parts {
		if err := physWriter.Write(part.Partname(), part.Blob()); err != nil {
			return err
		}
		if part.Rels().Len() > 0 {
			if err := physWriter.Write(part.Partname().RelsURI(), part.Rels().XML()); err != nil {
				return err
			}
		}
	}
	return nil
}

// ContentTypesItem builds the [Content_Types].xml content: a collection of
// Default and Override entries derived from the parts in the package.
type ContentTypesItem struct {
	defaults  CaseInsensitiveDict
	overrides map[string]string
}

// newContentTypesItem creates an empty ContentTypesItem.
func newContentTypesItem() *ContentTypesItem {
	return &ContentTypesItem{
		defaults:  NewCaseInsensitiveDict(),
		overrides: make(map[string]string),
	}
}

// newContentTypesItemFromParts creates a ContentTypesItem and populates it
// with Default entries for .rels and .xml, then adds each part's content type.
func newContentTypesItemFromParts(parts []*Part) *ContentTypesItem {
	cti := newContentTypesItem()
	cti.defaults.Set("rels", CT_OPC_RELATIONSHIPS)
	cti.defaults.Set("xml", CT_XML)
	for _, part := range parts {
		cti.addContentType(part.Partname(), part.ContentType())
	}
	return cti
}

// addContentType registers the content type for a part. If the extension
// and content type match a known default, it is stored as a Default entry;
// otherwise it is stored as an Override.
func (cti *ContentTypesItem) addContentType(partname PackURI, contentType string) {
	ext := partname.Ext()
	if IsDefaultContentType(ext, contentType) {
		cti.defaults.Set(ext, contentType)
	} else {
		cti.overrides[string(partname)] = contentType
	}
}

// blob serialises the content types data into a [Content_Types].xml byte
// slice, with Default and Override elements sorted alphabetically.
func (cti *ContentTypesItem) blob() []byte {
	typesEl := NewTypesElement()
	exts := make([]string, 0, len(cti.defaults))
	for ext := range cti.defaults {
		exts = append(exts, ext)
	}
	sort.Strings(exts)
	for _, ext := range exts {
		ct, _ := cti.defaults.Get(ext)
		child := NewDefaultElement(ext, ct)
		typesEl.AddChild(child)
	}
	overrideKeys := make([]string, 0, len(cti.overrides))
	for k := range cti.overrides {
		overrideKeys = append(overrideKeys, k)
	}
	sort.Strings(overrideKeys)
	for _, key := range overrideKeys {
		child := NewOverrideElement(key, cti.overrides[key])
		typesEl.AddChild(child)
	}
	return serializePartXML(typesEl)
}
