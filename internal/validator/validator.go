package validator

import (
	"fmt"
	"strings"

	"github.com/johngb/langreg"
)

const (
	// AssertionWarning indicates a potential issue
	AssertionWarning = 0
	// AssertionError indicates an error in the validation
	AssertionError = 1
)

type (
	// Assertion is used to collect validation information
	Assertion struct {
		Type int    // 0 == warning, 1 == error
		Txt  string // description of the problem
		Err  error
	}

	// Validator collects assertions
	Validator struct {
		Name     string
		Issues   []*Assertion
		Errors   int
		Warnings int
	}

	// Validatable is the interface that maust be implemented to support recursive validations of strucs
	Validatable interface {
		Validate(*Validator) *Validator
	}
)

// New initializes and returns a new Validator
func New(name string) *Validator {
	v := Validator{
		Name:   name,
		Issues: make([]*Assertion, 0),
	}
	return &v
}

// Validate starts the chain of validations
func (v *Validator) Validate(src Validatable) *Validator {
	return src.Validate(v)
}

// AssertError add an error assertion
func (v *Validator) AssertError(txt string) {
	v.Issues = append(v.Issues, &Assertion{Type: AssertionError, Txt: txt})
	v.Errors++
}

// AssertWarning add an warning assertion
func (v *Validator) AssertWarning(txt string) {
	v.Issues = append(v.Issues, &Assertion{Type: AssertionWarning, Txt: txt})
	v.Errors++
}

// AssertStringError verifies a string
func (v *Validator) AssertStringError(src, expected string) {
	if len(src) != len(expected) {
		v.AssertError(fmt.Sprintf("Expected '%s', found '%s'", expected, src))
		return
	}

	if src != expected {
		v.AssertError(fmt.Sprintf("Expected '%s', found '%s'", expected, src))
	}
}

// AssertStringExists verifies a string is not empty
func (v *Validator) AssertStringExists(src, name string) {
	if len(src) == 0 {
		v.AssertError(fmt.Sprintf("Expected non empty attribute '%s'", name))
	}
}

// AssertNotNil verifies that an attribute is not nil
func (v *Validator) AssertNotNil(src interface{}, name string) {
	if src == nil {
		v.AssertError(fmt.Sprintf("Expected no nil attribute '%s'", name))
	}
}

// AssertNotEmpty verifies that a map is not empty
func (v *Validator) AssertNotEmpty(src map[string]string, name string) {
	if len(src) == 0 {
		v.AssertError(fmt.Sprintf("Expected none empty map '%s'", name))
	}
}

// AssertNotZero verifies that a map is not empty
func (v *Validator) AssertNotZero(src int, name string) {
	if src == 0 {
		v.AssertError(fmt.Sprintf("Expected no-zero attribute '%s'", name))
	}
}

// AssertISO639 verifies that src complies with ISO 639-1
func (v *Validator) AssertISO639(src string) {
	lang := src
	if !strings.Contains(src, "_") {
		lang = src + "_" + strings.ToUpper(src)
	}
	if !langreg.IsValidLangRegCode(lang) {
		v.AssertError(fmt.Sprintf("Invalid language code '%s'", src))
	}
}

// AssertContains verifies that a map contains key
func (v *Validator) AssertContains(src map[string]string, key, name string) {
	if len(src) == 0 {
		v.AssertError(fmt.Sprintf("Expected none empty map '%s'", name))
		return
	}
	if _, ok := src[key]; !ok {
		v.AssertError(fmt.Sprintf("Expected key '%s' in map '%s'", key, name))
	}
}

// AssertExistsError verifies that a struct exists
func (v *Validator) AssertExistsError(src interface{}, expected string) {
	if src == nil {
		v.AssertError(fmt.Sprintf("Expected '%s', found 'nil'", expected))
	}
}

// IsValid returns true if NError == 0. Warnings are ignored
func (v *Validator) IsValid() bool {
	return v.Errors == 0
}

// IsClean returns true if NError == 0 AND NWarnings == 0
func (v *Validator) IsClean() bool {
	return v.Errors == 0 && v.Warnings == 0
}

// NErrors returns the number of erros
func (v *Validator) NErrors() int {
	return v.Errors
}

// NWarnings returns the number of warnings
func (v *Validator) NWarnings() int {
	return v.Warnings
}

// AsError returns an error if NError > 0, nil otherwise
func (v *Validator) AsError() error {
	if v.Errors == 0 {
		return nil
	}
	return fmt.Errorf(v.Error())
}

// Error returns an error text
func (v *Validator) Error() string {
	//return fmt.Sprintf("validation '%s' has %d errors, %d warnings", v.Name, v.Errors, v.Warnings)
	return v.Report()
}

// Report returns a description of all issues
func (v *Validator) Report() string {
	if v.Errors == 0 {
		return "validation '%s' has zero errors/warnings"
	}
	r := "\n"
	for i, issue := range v.Issues {
		r = r + fmt.Sprintf("Issue %d: %s\n", i+1, issue.Txt)
	}
	return r
}
