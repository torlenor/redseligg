package mongobotconfigprovider

import (
	"testing"
)

func TestCreatingNewMongoBackend(t *testing.T) {
	backend, err := NewBackend("mongodb://mongo/ut", "ut")
	if err != nil || backend == nil {
		t.Fatalf("Could not get a new Mongo Backend: %s", err)
	}
}
