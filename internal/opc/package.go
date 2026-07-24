package opc

import (
	"fmt"
	"io"
)

type OpcPackage struct {
	rels *Relationships
	partFactory func(PackURI, string, string, []byte) *Part
}

func NewOpcPackage() *OpcPackage {
	return &OpcPackage{
		rels: NewRelationships(PACKAGE_URI.BaseURI()),
		partFactory: func(partname PackURI, contentType, reltype string, blob []byte) *Part {
			return NewPartFromFactory(partname, contentType, reltype, blob, nil)
		},
	}
}

func (pkg *OpcPackage) AfterUnmarshal() {}

func (pkg *OpcPackage) CoreProperties() *CoreProperties {
	cpPart := pkg.CorePropertiesPart()
	if cpPart == nil {
		element := NewDefaultCorePropertiesElement()
		return NewCoreProperties(element)
	}
	blob := cpPart.Blob()
	element, err := parseXML(blob)
	if err != nil || element == nil {
		return NewCoreProperties(NewDefaultCorePropertiesElement())
	}
	return NewCoreProperties(element)
}

func (pkg *OpcPackage) CorePropertiesPart() *Part {
	cpPart := pkg.PartRelatedBy(RT_CORE_PROPERTIES)
	if cpPart == nil {
		return nil
	}
	return cpPart
}

func (pkg *OpcPackage) createDefaultCorePropertiesPart() *Part {
	partname, _ := NewPackURI("/docProps/core.xml")
	element := NewDefaultCorePropertiesElement()
	blob := serializePartXML(element)
	part := NewPart(partname, CT_OPC_CORE_PROPERTIES, blob, pkg)
	pkg.RelateTo(part, RT_CORE_PROPERTIES)
	return part
}

func (pkg *OpcPackage) IterParts() []*Part {
	return pkg.walkParts()
}

func (pkg *OpcPackage) walkParts() []*Part {
	var result []*Part
	visited := make(map[*Part]bool)
	pkg.walkPartsFrom(pkg.rels, visited, &result)
	return result
}

func (pkg *OpcPackage) walkPartsFrom(rels *Relationships, visited map[*Part]bool, result *[]*Part) {
	for _, rel := range rels.rels {
		if rel.isExternal {
			continue
		}
		part := rel.targetPart
		if part == nil || visited[part] {
			continue
		}
		visited[part] = true
		*result = append(*result, part)
		pkg.walkPartsFrom(part.Rels(), visited, result)
	}
}

func (pkg *OpcPackage) LoadRel(relType string, target interface{}, rID string, isExternal bool) *Relationship {
	return pkg.rels.AddRelationship(relType, target, rID, isExternal)
}

func (pkg *OpcPackage) MainDocumentPart() *Part {
	return pkg.PartRelatedBy(RT_OFFICE_DOCUMENT)
}

func (pkg *OpcPackage) NextPartname(template string) PackURI {
	existing := make(map[string]bool)
	for _, part := range pkg.walkParts() {
		existing[string(part.Partname())] = true
	}
	for n := 1; n <= len(existing)+1; n++ {
		candidate := fmt.Sprintf(template, n)
		if !existing[candidate] {
			pu, err := NewPackURI(candidate)
			if err == nil {
				return pu
			}
		}
	}
	pu, _ := NewPackURI(fmt.Sprintf(template, len(existing)+1))
	return pu
}

func (pkg *OpcPackage) PartRelatedBy(relType string) *Part {
	rel := pkg.rels.getRelOfType(relType)
	if rel == nil || rel.isExternal {
		return nil
	}
	return rel.targetPart
}

func (pkg *OpcPackage) Parts() []*Part {
	return pkg.walkParts()
}

func (pkg *OpcPackage) PartsCount() int {
	return len(pkg.walkParts())
}

func (pkg *OpcPackage) RelateTo(part *Part, relType string) string {
	return pkg.rels.GetOrAdd(relType, part).RID()
}

func (pkg *OpcPackage) Rels() *Relationships {
	return pkg.rels
}

func (pkg *OpcPackage) SetRels(rels *Relationships) {
	pkg.rels = rels
}

func Open(r io.ReaderAt, size int64) (*OpcPackage, error) {
	physReader, err := NewPhysPkgReaderFromReaderAt(r, size)
	if err != nil {
		return nil, fmt.Errorf("opc: failed to open package: %w", err)
	}
	defer physReader.Close()

	pkgReader, err := PackageReaderFromFile(physReader)
	if err != nil {
		return nil, fmt.Errorf("opc: failed to read package: %w", err)
	}

	pkg := NewOpcPackage()
	Unmarshal(pkgReader, pkg)
	return pkg, nil
}

func OpenFromPath(path string) (*OpcPackage, error) {
	physReader, err := NewPhysPkgReader(path)
	if err != nil {
		return nil, fmt.Errorf("opc: failed to open package at '%s': %w", path, err)
	}
	defer physReader.Close()

	pkgReader, err := PackageReaderFromFile(physReader)
	if err != nil {
		return nil, fmt.Errorf("opc: failed to read package: %w", err)
	}

	pkg := NewOpcPackage()
	Unmarshal(pkgReader, pkg)
	return pkg, nil
}

func (pkg *OpcPackage) Save(w io.Writer) error {
	for _, part := range pkg.Parts() {
		part.BeforeMarshal()
	}
	physPath, ok := w.(interface{ Name() string })
	var physWriter PhysPkgWriter
	var err error
	if ok {
		physWriter, err = NewPhysPkgWriter(physPath.Name())
	} else {
		return fmt.Errorf("opc: Save requires a file path or file-like writer")
	}
	if err != nil {
		return fmt.Errorf("opc: failed to create package writer: %w", err)
	}

	writer := &PackageWriter{}
	if err := writer.Write(physWriter, pkg.rels, pkg.Parts()); err != nil {
		physWriter.Close()
		return err
	}
	return physWriter.Close()
}

func (pkg *OpcPackage) SaveToWriter(w io.Writer) error {
	for _, part := range pkg.Parts() {
		part.BeforeMarshal()
	}
	physWriter := NewWriterPhysPkgWriter(w)
	writer := &PackageWriter{}
	if err := writer.Write(physWriter, pkg.rels, pkg.Parts()); err != nil {
		physWriter.Close()
		return err
	}
	return physWriter.Close()
}

func (pkg *OpcPackage) SaveToPath(path string) error {
	for _, part := range pkg.Parts() {
		part.BeforeMarshal()
	}
	physWriter, err := NewPhysPkgWriter(path)
	if err != nil {
		return fmt.Errorf("opc: failed to create package writer: %w", err)
	}

	writer := &PackageWriter{}
	if err := writer.Write(physWriter, pkg.rels, pkg.Parts()); err != nil {
		physWriter.Close()
		return err
	}
	return physWriter.Close()
}

func Unmarshal(pkgReader *PackageReader, pkg *OpcPackage) {
	parts := unmarshalParts(pkgReader, pkg)
	unmarshalRelationships(pkgReader, pkg, parts)
	for _, part := range parts {
		part.AfterUnmarshal()
	}
	pkg.AfterUnmarshal()
}

func unmarshalParts(pkgReader *PackageReader, pkg *OpcPackage) map[PackURI]*Part {
	parts := make(map[PackURI]*Part)
	for _, spart := range pkgReader.IterSparts() {
		part := NewPartFromFactory(spart.partname, spart.contentType, spart.reltype, spart.blob, pkg)
		parts[spart.partname] = part
	}
	return parts
}

func unmarshalRelationships(pkgReader *PackageReader, pkg *OpcPackage, parts map[PackURI]*Part) {
	for _, spart := range pkgReader.IterSparts() {
		source := parts[spart.partname]
		if source == nil {
			continue
		}
		for _, srel := range spart.srels.List() {
			var target interface{}
			if srel.IsExternal() {
				target = srel.TargetRef()
			} else {
				targetPart := parts[srel.TargetPartname()]
				if targetPart == nil {
					continue
				}
				target = targetPart
			}
			source.LoadRel(srel.RelType(), target, srel.RID(), srel.IsExternal())
		}
	}

	for _, srel := range pkgReader.pkgRels.List() {
		var target interface{}
		if srel.IsExternal() {
			target = srel.TargetRef()
		} else {
			targetPart := parts[srel.TargetPartname()]
			if targetPart == nil {
				continue
			}
			target = targetPart
		}
		pkg.LoadRel(srel.RelType(), target, srel.RID(), srel.IsExternal())
	}
}
