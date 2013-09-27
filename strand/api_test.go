package strand

import (
	"github.com/bcurren/go-hue"
	"testing"
)

func Test_LightStrandImplementsAPIInterface(t *testing.T) {
	lightStrand := NewLightStrand(3, nil)
	funcThatTakesAPIAsParameter(lightStrand)
}

func funcThatTakesAPIAsParameter(api hue.API) {
	// noop
}
