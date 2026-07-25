// Package opc implements the Open Packaging Conventions (OPC) standard for
// reading and writing OOXML packages (.docx, .xlsx, .pptx). It provides the
// core container abstraction: a zip archive holding XML parts connected by
// typed relationships, with a content-type stream. This layer is ported from
// python-docx's opc module.
package opc

import (
	"fmt"
	"io"
	"strings"
)

// OpcPackage is the root of an OPC package. It owns the package-level
// relationships, provides access to all parts, and coordinates open/save
// lifecycle (reading from zip, unmarshalling, marshalling, writing to zip).
type OpcPackage struct {
	rels            *Relationships
	partFactory     func(PackURI, string, string, []byte) *Part
	cachedCoreProps *CoreProperties
}

// NewOpcPackage creates and returns a new empty OpcPackage with no parts and
// a fresh package-level relationships collection rooted at "/".
func NewOpcPackage() *OpcPackage {
	return &OpcPackage{
		rels: NewRelationships(PACKAGE_URI.BaseURI()),
		partFactory: func(partname PackURI, contentType, reltype string, blob []byte) *Part {
			return NewPartFromFactory(partname, contentType, reltype, blob, nil)
		},
	}
}

// AfterUnmarshal is a lifecycle hook called after the package and all its
// parts have been unmarshalled from the physical package. No-op by default.
func (pkg *OpcPackage) AfterUnmarshal() {}

// CoreProperties returns the package-level core properties (author, title,
// created date, etc.), caching them after the first access. If no core
// properties part exists, a default is created in memory without adding it
// to the package.
func (pkg *OpcPackage) CoreProperties() *CoreProperties {
	if pkg.cachedCoreProps != nil {
		return pkg.cachedCoreProps
	}
	cpPart := pkg.CorePropertiesPart()
	if cpPart == nil {
		element := NewDefaultCorePropertiesElement()
		pkg.cachedCoreProps = NewCorePropertiesWithPart(element, nil)
		return pkg.cachedCoreProps
	}
	blob := cpPart.Blob()
	element, err := parseXML(blob)
	if err != nil || element == nil {
		pkg.cachedCoreProps = NewCorePropertiesWithPart(NewDefaultCorePropertiesElement(), cpPart)
		return pkg.cachedCoreProps
	}
	pkg.cachedCoreProps = NewCorePropertiesWithPart(element, cpPart)
	return pkg.cachedCoreProps
}

// CorePropertiesPart returns the Part that holds the core properties XML, or
// nil if no such part is related from the package.
func (pkg *OpcPackage) CorePropertiesPart() *Part {
	cpPart := pkg.PartRelatedBy(RT_CORE_PROPERTIES)
	if cpPart == nil {
		return nil
	}
	return cpPart
}

// createDefaultCorePropertiesPart creates a new core properties part at
// /docProps/core.xml, serialises default core properties into it, relates it
// to the package, and returns the part.
func (pkg *OpcPackage) createDefaultCorePropertiesPart() *Part {
	partname, _ := NewPackURI("/docProps/core.xml")
	element := NewDefaultCorePropertiesElement()
	blob := serializePartXML(element)
	part := NewPart(partname, CT_OPC_CORE_PROPERTIES, blob, pkg)
	pkg.RelateTo(part, RT_CORE_PROPERTIES)
	return part
}

// IterParts returns a slice of every Part reachable from the package-level
// relationships via depth-first traversal.
func (pkg *OpcPackage) IterParts() []*Part {
	return pkg.walkParts()
}

// walkParts performs a depth-first walk of the relationship graph starting
// at the package level, collecting every non-external part. Cycles are
// prevented via the visited set.
func (pkg *OpcPackage) walkParts() []*Part {
	var result []*Part
	visited := make(map[*Part]bool)
	pkg.walkPartsFrom(pkg.rels, visited, &result)
	return result
}

// walkPartsFrom recursively walks the given Relationships, appending every
// internal (non-external) part to result and recursing into each part's own
// relationships. visited prevents duplicate traversal.
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

// LoadRel adds a relationship to the package-level relationships with the
// given type, target, relationship ID, and external flag, then returns it.
func (pkg *OpcPackage) LoadRel(relType string, target interface{}, rID string, isExternal bool) *Relationship {
	return pkg.rels.AddRelationship(relType, target, rID, isExternal)
}

// MainDocumentPart returns the part related to the package by the office
// document relationship type, or nil if none exists.
func (pkg *OpcPackage) MainDocumentPart() *Part {
	return pkg.PartRelatedBy(RT_OFFICE_DOCUMENT)
}

// NextPartname returns the next available pack URI by applying the given
// fmt template (e.g. "/word/chapter%d.xml") with the lowest positive
// integer not already used by an existing part in the package.
func (pkg *OpcPackage) NextPartname(template string) PackURI {
	if !strings.Contains(template, "%d") {
		pu, err := NewPackURI(template)
		if err == nil {
			return pu
		}
		pu, _ = NewPackURI(template + "1")
		return pu
	}
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

// PartRelatedBy returns the first internal part related to this package by
// the given relationship type, or nil if no such relationship exists.
func (pkg *OpcPackage) PartRelatedBy(relType string) *Part {
	rel := pkg.rels.getRelOfType(relType)
	if rel == nil || rel.isExternal {
		return nil
	}
	return rel.targetPart
}

// Parts returns all parts reachable from the package-level relationships
// via depth-first traversal (same as IterParts).
func (pkg *OpcPackage) Parts() []*Part {
	return pkg.walkParts()
}

// PartsCount returns the number of parts reachable from the package-level
// relationships.
func (pkg *OpcPackage) PartsCount() int {
	return len(pkg.walkParts())
}

// RelateTo creates or retrieves a relationship from the package to the
// given part with the given type, and returns the relationship ID.
func (pkg *OpcPackage) RelateTo(part *Part, relType string) string {
	return pkg.rels.GetOrAdd(relType, part).RID()
}

// Rels returns the package-level Relationships collection.
func (pkg *OpcPackage) Rels() *Relationships {
	return pkg.rels
}

// SetRels replaces the package-level Relationships collection.
func (pkg *OpcPackage) SetRels(rels *Relationships) {
	pkg.rels = rels
}

// Open reads an OPC package from an io.ReaderAt (e.g. a bytes.Reader or
// open file) with the given size. It returns the fully unmarshalled
// OpcPackage or an error.
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

// OpenFromPath opens and unmarshals an OPC package from the given file path.
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

// Save marshals all parts (calling BeforeMarshal on each) and writes the
// complete OPC package as a zip archive to the given io.Writer. The writer
// must also implement Name() string (e.g. an os.File) to derive the output
// path for the zip writer.
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

// SaveToWriter marshals the package and writes the zip archive to any
// io.Writer (without requiring a file path).
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

// SaveToPath marshals the package and writes the zip archive to the given
// file path.
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

// Unmarshal populates an OpcPackage from a PackageReader. It deserialises
// all parts, wires up relationships between them, and calls AfterUnmarshal
// on every part and the package itself.
func Unmarshal(pkgReader *PackageReader, pkg *OpcPackage) {
	parts := unmarshalParts(pkgReader, pkg)
	unmarshalRelationships(pkgReader, pkg, parts)
	for _, part := range parts {
		part.AfterUnmarshal()
	}
	pkg.AfterUnmarshal()
}

// unmarshalParts converts every SerializedPart from the PackageReader into a
// Part (via the part factory), registering them by PackURI in a map.
func unmarshalParts(pkgReader *PackageReader, pkg *OpcPackage) map[PackURI]*Part {
	parts := make(map[PackURI]*Part)
	for _, spart := range pkgReader.IterSparts() {
		part := NewPartFromFactory(spart.partname, spart.contentType, spart.reltype, spart.blob, pkg)
		parts[spart.partname] = part
	}
	return parts
}

// unmarshalRelationships wires up all relationships from serialised data
// into the in-memory part graph. For each serialised part and the package
// itself, it loads every serialised relationship, resolves the target part
// from the parts map (or stores the external target string), and calls
// LoadRel on the source.
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
