package strand

import (
	"github.com/bcurren/go-hue"
	"testing"
)

func Test_LightStrandImplementsApiInterface(t *testing.T) {
	lightStrand := NewLightStrand(3, nil)
	funcThatTakesApiAsParameter(lightStrand)
}

func funcThatTakesApiAsParameter(api hue.Api) {
	// noop
}
