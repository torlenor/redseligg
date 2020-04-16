package utils

import (
	"testing"
)

func TestVersion(t *testing.T) {
	expectedVersion := "1.2.3"
	expectedCompTime := "12345"

	actualVersion := Version().Get()
	if actualVersion != "" {
		t.Fatalf("Default version was not empty but %s", actualVersion)
	}
	Version().Set(expectedVersion)
	actualVersion = Version().Get()
	if actualVersion != expectedVersion {
		t.Fatalf("Version was not %s but %s", expectedVersion, actualVersion)
	}

	actualCompTime := Version().GetCompTime()
	if actualCompTime != "" {
		t.Fatalf("Default compTime was not empty but %s", actualVersion)
	}
	Version().SetCompTime(expectedCompTime)
	actualCompTime = Version().GetCompTime()
	if actualCompTime != expectedCompTime {
		t.Fatalf("CompTime was not %s but %s", expectedCompTime, actualCompTime)
	}

}
