package opc

import (
	"sort"
)

type PackageWriter struct{}

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

func writeContentTypesStream(physWriter PhysPkgWriter, parts []*Part) error {
	cti := newContentTypesItemFromParts(parts)
	return physWriter.Write(CONTENT_TYPES_URI, cti.blob())
}

func writePkgRels(physWriter PhysPkgWriter, pkgRels *Relationships) error {
	return physWriter.Write(PACKAGE_URI.RelsURI(), pkgRels.XML())
}

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

type ContentTypesItem struct {
	defaults  CaseInsensitiveDict
	overrides map[string]string
}

func newContentTypesItem() *ContentTypesItem {
	return &ContentTypesItem{
		defaults:  NewCaseInsensitiveDict(),
		overrides: make(map[string]string),
	}
}

func newContentTypesItemFromParts(parts []*Part) *ContentTypesItem {
	cti := newContentTypesItem()
	cti.defaults.Set("rels", CT_OPC_RELATIONSHIPS)
	cti.defaults.Set("xml", CT_XML)
	for _, part := range parts {
		cti.addContentType(part.Partname(), part.ContentType())
	}
	return cti
}

func (cti *ContentTypesItem) addContentType(partname PackURI, contentType string) {
	ext := partname.Ext()
	if IsDefaultContentType(ext, contentType) {
		cti.defaults.Set(ext, contentType)
	} else {
		cti.overrides[string(partname)] = contentType
	}
}

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
