// Package enums provides generic utilities for mapping between integer enum values and their XML string representations.
package enums

import "fmt"

// XmlEnum is an interface for enum types that have an integer value and an XML string representation.
type XmlEnum interface {
	~int
	Value() int
	XMLValue() string
}

// xmlEntry pairs an integer value with its XML string and member of the enum type.
type xmlEntry[T any] struct {
	value  int
	xmlVal string
	member T
}

// FromXML looks up an enum member by its XML string representation. Returns an error if no match is found.
func FromXML[T ~int](xmlVal string, entries []xmlEntry[T]) (T, error) {
	for _, e := range entries {
		if e.xmlVal == xmlVal {
			return e.member, nil
		}
	}
	var zero T
	return zero, fmt.Errorf("no XML mapping for '%s'", xmlVal)
}

// ToXML returns the XML string representation of an enum member. Returns an error if no mapping exists.
func ToXML[T ~int](v T, entries []xmlEntry[T]) (string, error) {
	for _, e := range entries {
		if e.member == v {
			if e.xmlVal == "" {
				return "", fmt.Errorf("member %d has no XML representation", int(v))
			}
			return e.xmlVal, nil
		}
	}
	return "", fmt.Errorf("value %d is not a valid member", int(v))
}
