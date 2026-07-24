package opc

import (
	"fmt"
	"sort"
)

type Relationship struct {
	rID        string
	relType    string
	targetRef  string
	isExternal bool
	targetPart *Part
}

type Relationships struct {
	baseURI string
	rels    map[string]*Relationship
}

func NewRelationships(baseURI string) *Relationships {
	return &Relationships{
		baseURI: baseURI,
		rels:    make(map[string]*Relationship),
	}
}

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

func (rs *Relationships) GetOrAdd(relType string, targetPart *Part) *Relationship {
	rel := rs.findMatching(relType, targetPart, false)
	if rel != nil {
		return rel
	}
	rID := rs.NextRID()
	return rs.AddRelationship(relType, targetPart, rID, false)
}

func (rs *Relationships) GetOrAddExtRel(relType string, targetRef string) string {
	rel := rs.findMatching(relType, targetRef, true)
	if rel != nil {
		return rel.rID
	}
	rID := rs.NextRID()
	rs.AddRelationship(relType, targetRef, rID, true)
	return rID
}

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

func (rs *Relationships) RelatedParts() map[string]*Part {
	result := make(map[string]*Part)
	for rID, rel := range rs.rels {
		if !rel.isExternal {
			result[rID] = rel.targetPart
		}
	}
	return result
}

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

func (rs *Relationships) Len() int {
	return len(rs.rels)
}

func (rs *Relationships) Get(rID string) *Relationship {
	return rs.rels[rID]
}

func (rs *Relationships) Delete(rID string) {
	delete(rs.rels, rID)
}

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

func (rs *Relationships) getRelOfType(relType string) *Relationship {
	var matching []*Relationship
	for _, rel := range rs.rels {
		if rel.relType == relType {
			matching = append(matching, rel)
		}
	}
	if len(matching) == 0 {
		return nil
	}
	if len(matching) > 1 {
		return matching[0]
	}
	return matching[0]
}

func (rs *Relationships) NextRID() string {
	for n := 1; n <= len(rs.rels)+1; n++ {
		candidate := fmt.Sprintf("rId%d", n)
		if _, ok := rs.rels[candidate]; !ok {
			return candidate
		}
	}
	return fmt.Sprintf("rId%d", len(rs.rels)+1)
}

func (rel *Relationship) RID() string {
	return rel.rID
}

func (rel *Relationship) RelType() string {
	return rel.relType
}

func (rel *Relationship) TargetRef() string {
	return rel.targetRef
}

func (rel *Relationship) IsExternal() bool {
	return rel.isExternal
}

func (rel *Relationship) TargetPart() *Part {
	if rel.isExternal {
		panic("opc: target_part is undefined for external relationships")
	}
	return rel.targetPart
}

func (rel *Relationship) SetTargetRef(ref string) {
	rel.targetRef = ref
}
