package opc

import (
	"fmt"
	"sort"
)

// Relationship represents a single OPC relationship: a directed, typed link
// from a source part (or the package itself) to a target part or external
// resource, identified by a relationship ID (e.g. "rId1").
type Relationship struct {
	rID        string
	relType    string
	targetRef  string
	isExternal bool
	targetPart *Part
}

// Relationships is a collection of Relationship objects owned by a single
// source (a Part or the package root). It provides lookup, addition,
// deletion, and XML serialisation of the relationship items.
type Relationships struct {
	baseURI string
	rels    map[string]*Relationship
}

// NewRelationships creates an empty Relationships collection with the given
// base URI (used to compute relative target references).
func NewRelationships(baseURI string) *Relationships {
	return &Relationships{
		baseURI: baseURI,
		rels:    make(map[string]*Relationship),
	}
}

// AddRelationship creates a new Relationship with the given parameters and
// adds it to the collection. For internal relationships, the target must be
// a *Part and the targetRef is computed as a relative path from baseURI.
// For external relationships, the target must be a string (the external URI).
func (rs *Relationships) AddRelationship(relType string, target interface{}, rID string, isExternal bool) *Relationship {
	var rel *Relationship
	if isExternal {
		targetStr, ok := target.(string)
		if !ok {
			panic("opc: external relationship target must be a string")
		}
		rel = &Relationship{
			rID:        rID,
			relType:    relType,
			targetRef:  targetStr,
			isExternal: true,
		}
	} else {
		targetPart, ok := target.(*Part)
		if !ok {
			panic("opc: internal relationship target must be a *Part")
		}
		rel = &Relationship{
			rID:        rID,
			relType:    relType,
			isExternal: false,
			targetPart: targetPart,
			targetRef:  targetPart.Partname().RelativeRef(rs.baseURI),
		}
	}
	rs.rels[rID] = rel
	return rel
}

// GetOrAdd returns the existing relationship matching the given type and
// target part, or creates a new one with the next available RID if none
// matches.
func (rs *Relationships) GetOrAdd(relType string, targetPart *Part) *Relationship {
	rel := rs.findMatching(relType, targetPart, false)
	if rel != nil {
		return rel
	}
	rID := rs.NextRID()
	return rs.AddRelationship(relType, targetPart, rID, false)
}

// GetOrAddExtRel returns the RID of an existing external relationship
// matching the given type and target ref, or creates a new one and returns
// its RID.
func (rs *Relationships) GetOrAddExtRel(relType string, targetRef string) string {
	rel := rs.findMatching(relType, targetRef, true)
	if rel != nil {
		return rel.rID
	}
	rID := rs.NextRID()
	rs.AddRelationship(relType, targetRef, rID, true)
	return rID
}

// PartWithReltype returns the first internal target part whose relationship
// has the given type, or nil if none matches.
func (rs *Relationships) PartWithReltype(relType string) *Part {
	rel := rs.getRelOfType(relType)
	if rel == nil {
		return nil
	}
	if rel.isExternal {
		return nil
	}
	return rel.targetPart
}

// RelatedParts returns a map of RID to *Part for every internal (non-external)
// relationship in the collection.
func (rs *Relationships) RelatedParts() map[string]*Part {
	result := make(map[string]*Part)
	for rID, rel := range rs.rels {
		if !rel.isExternal {
			result[rID] = rel.targetPart
		}
	}
	return result
}

// XML serialises the relationships collection into an OPC relationships XML
// document. Items are sorted by RID for deterministic output.
func (rs *Relationships) XML() []byte {
	relsEl := NewRelationshipsElement()

	rIDs := make([]string, 0, len(rs.rels))
	for rID := range rs.rels {
		rIDs = append(rIDs, rID)
	}
	sort.Strings(rIDs)

	for _, rID := range rIDs {
		rel := rs.rels[rID]
		targetMode := RTM_INTERNAL
		if rel.isExternal {
			targetMode = RTM_EXTERNAL
		}
		child := NewRelationshipElement(rID, rel.relType, rel.targetRef, targetMode)
		relsEl.AddChild(child)
	}

	return serializePartXML(relsEl)
}

// Len returns the number of relationships in the collection.
func (rs *Relationships) Len() int {
	return len(rs.rels)
}

// Get returns the relationship with the given RID, or nil if not found.
func (rs *Relationships) Get(rID string) *Relationship {
	return rs.rels[rID]
}

// Delete removes the relationship identified by the given RID from the
// collection.
func (rs *Relationships) Delete(rID string) {
	delete(rs.rels, rID)
}

// findMatching searches the collection for a relationship matching the given
// type, target, and external flag. For internal relationships it compares
// the *Part pointer; for external it compares the target string.
func (rs *Relationships) findMatching(relType string, target interface{}, isExternal bool) *Relationship {
	for _, rel := range rs.rels {
		if rel.relType != relType {
			continue
		}
		if rel.isExternal != isExternal {
			continue
		}
		if isExternal {
			targetStr, ok := target.(string)
			if ok && rel.targetRef == targetStr {
				return rel
			}
		} else {
			targetPart, ok := target.(*Part)
			if ok && rel.targetPart == targetPart {
				return rel
			}
		}
	}
	return nil
}

// getRelOfType returns the first relationship with the given type, or nil if
// none exists. If multiple match, the first one (iteration order) is returned.
func (rs *Relationships) getRelOfType(relType string) *Relationship {
	for _, rel := range rs.rels {
		if rel.relType == relType {
			return rel
		}
	}
	return nil
}

// NextRID returns the next available relationship ID in the sequence
// "rId1", "rId2", ..., skipping any already in use.
func (rs *Relationships) NextRID() string {
	for n := 1; n <= len(rs.rels)+1; n++ {
		candidate := fmt.Sprintf("rId%d", n)
		if _, ok := rs.rels[candidate]; !ok {
			return candidate
		}
	}
	return fmt.Sprintf("rId%d", len(rs.rels)+1)
}

// RID returns the relationship ID (e.g. "rId1").
func (rel *Relationship) RID() string {
	return rel.rID
}

// RelType returns the relationship type URI.
func (rel *Relationship) RelType() string {
	return rel.relType
}

// TargetRef returns the target reference string (a relative path for internal
// relationships, or an absolute URI for external ones).
func (rel *Relationship) TargetRef() string {
	return rel.targetRef
}

// IsExternal returns true if this relationship points to an external resource
// outside the package.
func (rel *Relationship) IsExternal() bool {
	return rel.isExternal
}

// TargetPart returns the target *Part for internal relationships. Panics if
// called on an external relationship.
func (rel *Relationship) TargetPart() *Part {
	if rel.isExternal {
		panic("opc: target_part is undefined for external relationships")
	}
	return rel.targetPart
}

// SetTargetRef updates the target reference string for this relationship.
func (rel *Relationship) SetTargetRef(ref string) {
	rel.targetRef = ref
}
