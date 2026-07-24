package stypes

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/SamYue1/go-docx/internal/shared"
)

type Simple interface {
	FromXML(string) (any, error)
	ToXML(any) (string, error)
	Validate(any) error
}

type BaseIntType struct{}

func (BaseIntType) Validate(v any) error {
	switch v.(type) {
	case int:
		return nil
	default:
		return fmt.Errorf("expected int, got %T", v)
	}
}

type BaseStringType struct{}

func (BaseStringType) Validate(v any) error {
	switch v.(type) {
	case string:
		return nil
	default:
		return fmt.Errorf("expected string, got %T", v)
	}
}

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
	if v.(bool) {
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
	return strconv.Itoa(v.(int)), nil
}

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
	return v.(string), nil
}

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
	return v.(shared.RGBColor).String(), nil
}

func (STHexColor) Validate(v any) error {
	switch v.(type) {
	case shared.RGBColor:
		return nil
	default:
		return fmt.Errorf("expected RGBColor, got %T", v)
	}
}

type STHpsMeasure struct{}

func (STHpsMeasure) FromXML(s string) (int, error) {
	return strconv.Atoi(s)
}
