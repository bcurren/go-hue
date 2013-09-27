package multi

import (
	"github.com/bcurren/go-hue"
	"github.com/bcurren/go-hue/huetest"
	"testing"
)

func Test_MultiBridgeImplementsAPIInterface(t *testing.T) {
	multi := NewMultiAPI()
	funcThatTakesAPIAsParameter(multi)
}

func funcThatTakesAPIAsParameter(api hue.API) {
	// noop
}

func Test_GetLights(t *testing.T) {
	api1 := &huetest.StubAPI{}
	api1.GetLightsReturn = []hue.Light{hue.Light{Id: "1", Name: "Hue Lamp 1"}}
	
	api2 := &huetest.StubAPI{}
	api2.GetLightsReturn = []hue.Light{hue.Light{Id: "1", Name: "Hue Lamp 1"}}
	
	multi := NewMultiAPI()
	multi.AddAPI(api1)
	multi.AddAPI(api2)
	
	lights, err := multi.GetLights()
	if err != nil {
		t.Fatal("Error was returned.", err)
	}
	
	if len(lights) != 2 {
		t.Error("Should merge GetLights and have 2 results.")
	}
}
