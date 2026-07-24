package enums

import "fmt"

type XmlEnum interface {
	~int
	Value() int
	XMLValue() string
}

type xmlEntry[T any] struct {
	value   int
	xmlVal  string
	member  T
}

func FromXML[T ~int](xmlVal string, entries []xmlEntry[T]) (T, error) {
	for _, e := range entries {
		if e.xmlVal == xmlVal {
			return e.member, nil
		}
	}
	var zero T
	return zero, fmt.Errorf("no XML mapping for '%s'", xmlVal)
}

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
