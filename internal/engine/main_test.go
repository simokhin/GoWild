package engine

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	AllInit()
	os.Exit(m.Run())
}
