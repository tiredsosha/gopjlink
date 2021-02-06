package pjlink

import (
	"os"
	"testing"
)

var proj *Projector

func TestMain(m *testing.M) {
	proj = NewProjector("10.5.105.30", WithPassword(os.Getenv("PJLINK_PASS")))
	os.Exit(m.Run())
}
