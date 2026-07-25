// Package stypes provides simple-type converters for OOXML schema types
// (ST_OnOff, ST_DecimalNumber, ST_String, ST_HexColor, ST_HpsMeasure).
// Each type implements the Simple interface for XML-to-Go and Go-to-XML
// conversion with validation, mirroring the python-docx oxml.simpletypes layer.
package stypes

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/SamYue1/go-docx/internal/shared"
)

// Simple is the interface for OOXML simple-type converters. Each implementation
// handles XML string serialization, deserialization, and Go value validation.
type Simple interface {
	FromXML(string) (any, error)
	ToXML(any) (string, error)
	Validate(any) error
}

// BaseIntType provides a Validate method that checks whether a value is an int.
type BaseIntType struct{}

func (BaseIntType) Validate(v any) error {
	switch v.(type) {
	case int:
		return nil
	default:
		return fmt.Errorf("expected int, got %T", v)
	}
}

// BaseStringType provides a Validate method that checks whether a value is a string.
type BaseStringType struct{}

func (BaseStringType) Validate(v any) error {
	switch v.(type) {
	case string:
		return nil
	default:
		return fmt.Errorf("expected string, got %T", v)
	}
}

// STOnOff converts OOXML ST_OnOff values ("true"/"false"/"1"/"0"/"on"/"off")
// to and from Go bool.
type STOnOff struct{}

func (STOnOff) FromXML(s string) (bool, error) {
	switch strings.ToLower(s) {
	case "true", "1", "on":
		return true, nil
	case "false", "0", "off":
		return false, nil
	default:
		return false, fmt.Errorf("invalid ST_OnOff value: %q", s)
	}
}

func (STOnOff) ToXML(v any) (string, error) {
	if err := (STOnOff{}).Validate(v); err != nil {
		return "", err
	}
	b, ok := v.(bool)
	if !ok {
		return "", fmt.Errorf("STOnOff: expected bool, got %T", v)
	}
	if b {
		return "true", nil
	}
	return "false", nil
}

func (STOnOff) Validate(v any) error {
	_, ok := v.(bool)
	if !ok {
		return fmt.Errorf("expected bool, got %T", v)
	}
	return nil
}

// STDecimalNumber converts OOXML ST_DecimalNumber (integer) values to and from Go int.
type STDecimalNumber struct {
	BaseIntType
}

func (STDecimalNumber) FromXML(s string) (int, error) {
	return strconv.Atoi(s)
}

func (STDecimalNumber) ToXML(v any) (string, error) {
	if err := (STDecimalNumber{}).Validate(v); err != nil {
		return "", err
	}
	n, ok := v.(int)
	if !ok {
		return "", fmt.Errorf("STDecimalNumber: expected int, got %T", v)
	}
	return strconv.Itoa(n), nil
}

// STString is a pass-through converter for OOXML ST_String values.
type STString struct {
	BaseStringType
}

func (STString) FromXML(s string) (string, error) {
	return s, nil
}

func (STString) ToXML(v any) (string, error) {
	if err := (STString{}).Validate(v); err != nil {
		return "", err
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("STString: expected string, got %T", v)
	}
	return s, nil
}

// STHexColor converts OOXML ST_HexColor values (6-digit hex, or "auto") to
// and from shared.RGBColor. A value of "auto" maps to nil.
type STHexColor struct{}

func (STHexColor) FromXML(s string) (any, error) {
	if strings.ToLower(s) == "auto" {
		return nil, nil
	}
	if len(s) != 6 {
		return nil, fmt.Errorf("invalid ST_HexColor value %q: expected 6 hex digits", s)
	}
	return shared.RGBColorFromString(s)
}

func (STHexColor) ToXML(v any) (string, error) {
	if err := (STHexColor{}).Validate(v); err != nil {
		return "", err
	}
	c, ok := v.(shared.RGBColor)
	if !ok {
		return "", fmt.Errorf("STHexColor: expected RGBColor, got %T", v)
	}
	return c.String(), nil
}

func (STHexColor) Validate(v any) error {
	switch v.(type) {
	case shared.RGBColor:
		return nil
	default:
		return fmt.Errorf("expected RGBColor, got %T", v)
	}
}

// STHpsMeasure converts OOXML ST_HpsMeasure (half-point measurement) values
// to and from Go int.
type STHpsMeasure struct{}

func (STHpsMeasure) FromXML(s string) (int, error) {
	return strconv.Atoi(s)
}
