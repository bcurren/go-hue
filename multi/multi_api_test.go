package multi

import (
	"testing"
	// "github.com/bcurren/go-hue"
	"github.com/bcurren/go-hue/huetest"
)

func Test_AddAPI(t *testing.T) {
	multi := NewMultiAPI()

	if len(multi.APIs) != 0 {
		t.Error("Should init to an empty slice.")
	}

	api1 := &huetest.StubAPI{}
	multi.AddAPI(api1)
	if len(multi.APIs) != 1 {
		t.Error("Should have one api.")
	}

	api2 := &huetest.StubAPI{}
	multi.AddAPI(api2)
	if len(multi.APIs) != 2 {
		t.Error("Should have two apis.")
	}
}
