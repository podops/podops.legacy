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
)

// Validate verifies the integrity of the struct. Aborts on first error.
func (s *Show) Validate() error {
	var issues []*Assertion

	issues = assertString(s.APIVersion, "v1", issues)
	issues = assertString(s.Kind, "show", issues)
	issues = assertExists(s.Metadata, "metadata", issues)
	issues = assertExists(s.Description, "description", issues)
	issues = assertExists(s.Image, "image", issues)

	// FIXME: return an error or something

	return nil
}

// Validate verifies the integrity of the struct. Aborts on first error.
func (e *Episode) Validate() error {
	var issues []*Assertion

	issues = assertString(e.APIVersion, "v1", issues)
	issues = assertString(e.Kind, "episode", issues)
	issues = assertExists(e.Metadata, "metadata", issues)
	issues = assertExists(e.Description, "description", issues)
	issues = assertExists(e.Image, "image", issues)
	issues = assertExists(e.Enclosure, "enclosure", issues)

	// FIXME: return an error or something

	return nil
}

func assertString(src, expected string, issues []*Assertion) []*Assertion {
	if len(src) != len(expected) {
		issue := &Assertion{
			Type: AssertionError,
			Txt:  fmt.Sprintf("Expected '%s', found '%s'", src, expected),
		}
		return append(issues, issue)
	}

	if src != expected {
		issue := &Assertion{
			Type: AssertionError,
			Txt:  fmt.Sprintf("Expected '%s', found '%s'", src, expected),
		}
		return append(issues, issue)
	}
	return issues
}

func assertExists(src interface{}, expected string, issues []*Assertion) []*Assertion {
	if src == nil {
		issue := &Assertion{
			Type: AssertionError,
			Txt:  fmt.Sprintf("Expected '%s', found 'nil'", expected),
		}
		return append(issues, issue)
	}
	return issues
}

/*
Show struct {
		APIVersion  string          `json:"apiVersion" yaml:"apiVersion" binding:"required"`   // REQUIRED default: v1.0
		Kind        string          `json:"kind" yaml:"kind" binding:"required"`               // REQUIRED default: show
		Metadata    Metadata        `json:"metadata" yaml:"metadata" binding:"required"`       // REQUIRED
		Description ShowDescription `json:"description" yaml:"description" binding:"required"` // REQUIRED
		Image       Resource        `json:"image" yaml:"image" binding:"required"`             // REQUIRED 'channel.itunes.image'
	}
*/
