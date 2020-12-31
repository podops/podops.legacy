package apiv1

import (
	"testing"
)

func TestTeplateShow(t *testing.T) {
	show := DefaultShow("BASE_URL", "NAME", "TITLE", "SUMMARY", "GUID")

	err := show.Validate()
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestTeplateEpisode(t *testing.T) {
	show := DefaultEpisode("BASE_URL", "NAME", "PARENT_NAME", "GUID", "PARENT_GUID")
	err := show.Validate()
	if err != nil {
		t.Errorf(err.Error())
	}
}
