package apiv1

import "fmt"

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

// Validate starts the chain of validations
func (v *Validator) Validate(src Validatable) *Validator {
	return src.Validate(v)
}

// AssertStringError verifies a string
func (v *Validator) AssertStringError(src, expected string) {
	if len(src) != len(expected) {
		issue := &Assertion{
			Type: AssertionError,
			Txt:  fmt.Sprintf("Expected '%s', found '%s'", expected, src),
		}
		v.Issues = append(v.Issues, issue)
		v.Errors++
		return
	}

	if src != expected {
		issue := &Assertion{
			Type: AssertionError,
			Txt:  fmt.Sprintf("Expected '%s', found '%s'", expected, src),
		}
		v.Issues = append(v.Issues, issue)
		v.Errors++
	}
}

// AssertStringExists verifies a string is not empty
func (v *Validator) AssertStringExists(src, name string) {
	if len(src) == 0 {
		issue := &Assertion{
			Type: AssertionError,
			Txt:  fmt.Sprintf("Expected non empty attribute '%s'", name),
		}
		v.Issues = append(v.Issues, issue)
		v.Errors++
		return
	}
}

// AssertNotNil verifies that an attribute is not nil
func (v *Validator) AssertNotNil(src interface{}, name string) {
	if src == nil {
		issue := &Assertion{
			Type: AssertionError,
			Txt:  fmt.Sprintf("Expected no nil attribute '%s'", name),
		}
		v.Issues = append(v.Issues, issue)
		v.Errors++
		return
	}
}

// AssertNotEmpty verifies that a map is not empty
func (v *Validator) AssertNotEmpty(src map[string]string, name string) {
	if len(src) == 0 {
		issue := &Assertion{
			Type: AssertionError,
			Txt:  fmt.Sprintf("Expected none empty map '%s'", name),
		}
		v.Issues = append(v.Issues, issue)
		v.Errors++
		return
	}
}

// AssertNotZero verifies that a map is not empty
func (v *Validator) AssertNotZero(src int, name string) {
	if src == 0 {
		issue := &Assertion{
			Type: AssertionError,
			Txt:  fmt.Sprintf("Expected no-zero attribute '%s'", name),
		}
		v.Issues = append(v.Issues, issue)
		v.Errors++
		return
	}
}

// AssertContains verifies that a map contains key
func (v *Validator) AssertContains(src map[string]string, key, name string) {
	if len(src) == 0 {
		issue := &Assertion{
			Type: AssertionError,
			Txt:  fmt.Sprintf("Expected none empty map '%s'", name),
		}
		v.Issues = append(v.Issues, issue)
		v.Errors++
		return
	}
	if _, ok := src[key]; !ok {
		issue := &Assertion{
			Type: AssertionError,
			Txt:  fmt.Sprintf("Expected key '%s' in map '%s'", key, name),
		}
		v.Issues = append(v.Issues, issue)
		v.Errors++
		return
	}
}

// AssertExistsError verifies that a struct exists
func (v *Validator) AssertExistsError(src interface{}, expected string) {
	if src == nil {
		issue := &Assertion{
			Type: AssertionError,
			Txt:  fmt.Sprintf("Expected '%s', found 'nil'", expected),
		}
		v.Issues = append(v.Issues, issue)
		v.Errors++
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

// NewValidator initializes and returns a new Validator
func NewValidator(name string) *Validator {
	v := Validator{
		Name:   name,
		Issues: make([]*Assertion, 0),
	}
	return &v
}
