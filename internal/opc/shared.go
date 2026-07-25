package opc

import "strings"

// CaseInsensitiveDict is a string-to-string map with case-insensitive key
// lookups. Keys are normalised to lowercase for storage and retrieval.
type CaseInsensitiveDict map[string]string

// NewCaseInsensitiveDict creates and returns an empty CaseInsensitiveDict.
func NewCaseInsensitiveDict() CaseInsensitiveDict {
	return make(CaseInsensitiveDict)
}

// Get returns the value for the given key using case-insensitive lookup.
func (d CaseInsensitiveDict) Get(key string) (string, bool) {
	v, ok := d[strings.ToLower(key)]
	return v, ok
}

// Set stores a key-value pair, normalising the key to lowercase.
func (d CaseInsensitiveDict) Set(key, value string) {
	d[strings.ToLower(key)] = value
}

// Has returns true if the key exists in the dict (case-insensitive).
func (d CaseInsensitiveDict) Has(key string) bool {
	_, ok := d[strings.ToLower(key)]
	return ok
}
