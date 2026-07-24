package opc

import "strings"

type CaseInsensitiveDict map[string]string

func NewCaseInsensitiveDict() CaseInsensitiveDict {
	return make(CaseInsensitiveDict)
}

func (d CaseInsensitiveDict) Get(key string) (string, bool) {
	v, ok := d[strings.ToLower(key)]
	return v, ok
}

func (d CaseInsensitiveDict) Set(key, value string) {
	d[strings.ToLower(key)] = value
}

func (d CaseInsensitiveDict) Has(key string) bool {
	_, ok := d[strings.ToLower(key)]
	return ok
}
