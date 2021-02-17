package pjlink

import (
	"os"
	"testing"
)

var proj *Projector

func TestMain(m *testing.M) {
	proj = NewProjector("HBLL-1060-D1.byu.edu", WithPassword(os.Getenv("PJLINK_PASS")))
	os.Exit(m.Run())
}
