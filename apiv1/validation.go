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
	}

	// Validation collects assertions
	Validation struct {
		Name     string
		Issues   []*Assertion
		Errors   int
		Warnings int
	}
)

// AssertStringError verifies a string
func (v *Validation) AssertStringError(src, expected string) {
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

// AssertExistsError verifies that a struct exists
func (v *Validation) AssertExistsError(src interface{}, expected string) {
	if src == nil {
		issue := &Assertion{
			Type: AssertionError,
			Txt:  fmt.Sprintf("Expected '%s', found 'nil'", expected),
		}
		v.Issues = append(v.Issues, issue)
		v.Errors++
	}
}

// NErrors returns the number of erros
func (v *Validation) NErrors() int {
	return v.Errors
}

// Error returns an error if NError > 0, nil otherwise
func (v *Validation) Error() error {
	if v.Errors == 0 {
		return nil
	}
	return fmt.Errorf("validation: '%s' has %d errors", v.Name, v.Errors)
}

// NWarnings returns the number of warnings
func (v *Validation) NWarnings() int {
	return v.Warnings
}

// NewValidation initializes and returns a new validation
func NewValidation(name string) *Validation {
	v := Validation{
		Name:   name,
		Issues: make([]*Assertion, 0),
	}
	return &v
}

// Validate verifies the integrity of the struct. Aborts on first error.
func (s *Show) Validate() error {
	v := NewValidation("show")

	v.AssertStringError(s.APIVersion, "v1")
	v.AssertStringError(s.Kind, "show")
	v.AssertExistsError(s.Metadata, "metadata")
	v.AssertExistsError(s.Description, "description")
	v.AssertExistsError(s.Image, "image")

	return v.Error()
}

// Validate verifies the integrity of the struct. Aborts on first error.
func (e *Episode) Validate() error {
	v := NewValidation("episode")

	v.AssertStringError(e.APIVersion, "v1")
	v.AssertStringError(e.Kind, "episode")
	v.AssertExistsError(e.Metadata, "metadata")
	v.AssertExistsError(e.Description, "description")
	v.AssertExistsError(e.Image, "image")
	v.AssertExistsError(e.Enclosure, "enclosure")

	return v.Error()
}
