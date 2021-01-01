package apiv1

import (
	"testing"
)

func TestTemplateShow(t *testing.T) {
	s := DefaultShow("BASE_URL", "NAME", "TITLE", "SUMMARY", "GUID")
	v := s.Validate(NewValidator("show"))
	if !v.IsClean() {
		t.Errorf(v.AsError().Error())
	}
}

func TestTemplateEpisode(t *testing.T) {
	e := DefaultEpisode("BASE_URL", "NAME", "PARENT_NAME", "GUID", "PARENT_GUID")
	v := e.Validate(NewValidator("episode"))
	if !v.IsClean() {
		t.Errorf(v.AsError().Error())
	}
}
