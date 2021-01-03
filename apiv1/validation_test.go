package apiv1

import (
	"strings"
	"testing"

	"github.com/johngb/langreg"
)

func TestISO(t *testing.T) {
	src := "de"
	lang := src
	if !strings.Contains(src, "_") {
		lang = src + "_" + strings.ToUpper(src)
	}
	if !langreg.IsValidLangRegCode(lang) {
		t.Errorf("Invalid code '%s'", lang)
	}

}
func TestTemplateShow(t *testing.T) {
	s := DefaultShow("BASE_URL", "NAME", "TITLE", "SUMMARY", "GUID")
	v := s.Validate(NewValidator(ResourceShow))
	if !v.IsClean() {
		t.Errorf(v.AsError().Error())
	}
}

func TestTemplateEpisode(t *testing.T) {
	e := DefaultEpisode("BASE_URL", "NAME", "PARENT_NAME", "GUID", "PARENT_GUID")
	v := e.Validate(NewValidator(ResourceEpisode))
	if !v.IsClean() {
		t.Errorf(v.AsError().Error())
	}
}
