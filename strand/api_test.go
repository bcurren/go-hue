package strand

import (
	"github.com/bcurren/go-hue"
	"github.com/bcurren/go-hue/huetest"
	"testing"
	"time"
)

func Test_LightStrandImplementsAPIInterface(t *testing.T) {
	lightStrand := NewLightStrand(3, nil)
	funcThatTakesAPIAsParameter(lightStrand)
}

func funcThatTakesAPIAsParameter(api hue.API) {
	// noop
}

func Test_GetLights(t *testing.T) {
	stubAPI := &huetest.StubAPI{}
	stubAPI.GetLightsReturn = []hue.Light{
		hue.Light{Id: "light1"},
		hue.Light{Id: "light2"},
		hue.Light{Id: "light3"},
		hue.Light{Id: "light4"},
	}

	strand := NewLightStrand(3, stubAPI)
	strand.Lights.Set("1", "light1")
	strand.Lights.Set("2", "light2")
	strand.Lights.Set("4", "light4")

	lights, err := strand.GetLights()
	if err != nil {
		t.Fatal("Error returned from GetLights()")
	}

	if len(lights) != 3 {
		t.Error("Should only map lights assigned in strand.", len(lights))
	}
	if lights[0].Id != "1" || lights[1].Id != "2" || lights[2].Id != "4" {
		t.Error("Should have mapped lightIds to socketIds.")
	}
}

func Test_GetNewLights(t *testing.T) {
	stubAPI := &huetest.StubAPI{}
	stubAPI.GetNewLightsReturn = []hue.Light{
		hue.Light{Id: "light1"},
	}
	now := time.Now()
	stubAPI.GetNewLightsReturnTime = now

	strand := NewLightStrand(1, stubAPI)
	strand.Lights.Set("1", "light1")

	lights, lastUpdated, err := strand.GetNewLights()
	if err != nil {
		t.Fatal("Error returned from GetNewLights()")
	}

	if lastUpdated != now {
		t.Error("Didn't pass along the last updated date.")
	}
	if len(lights) != 1 {
		t.Error("Should only map lights assigned in strand.")
	}
	if lights[0].Id != "1" {
		t.Error("Should have mapped lightIds to socketIds.")
	}
}

func Test_GetLightAttributesSuccess(t *testing.T) {
	stubAPI := &huetest.StubAPI{}
	expectedLightAttributes := &hue.LightAttributes{}
	stubAPI.GetLightAttributesReturn = expectedLightAttributes

	strand := NewLightStrand(1, stubAPI)
	strand.Lights.Set("1", "light1")

	actualLightAttributes, err := strand.GetLightAttributes("1")
	if err != nil {
		t.Fatal("Error returned from GetLightAttributes()")
	}

	if actualLightAttributes != expectedLightAttributes {
		t.Error("Didn't pass along the light attributes.")
	}
	if stubAPI.GetLightAttributesParamLightId != "light1" {
		t.Error("Didn't map from socketId to lightId.")
	}
}

func Test_GetLightAttributesError(t *testing.T) {
	stubAPI := &huetest.StubAPI{}

	strand := NewLightStrand(1, stubAPI)
	strand.Lights.Set("1", "light1")

	_, err := strand.GetLightAttributes("invalidId")
	if err == nil {
		t.Fatal("Error for invalid socket id should be returned.")
	}
	apiError, ok := err.(*hue.APIError)
	if !ok {
		t.Fatal("Should return APIError for invalid socket id.")
	}

	if apiError.Errors[0].Type != hue.ResourceNotAvailableErrorType {
		t.Error("Invalid type.")
	}
}

func Test_SetLightNameSuccess(t *testing.T) {
	stubAPI := &huetest.StubAPI{}

	strand := NewLightStrand(1, stubAPI)
	strand.Lights.Set("1", "light1")

	err := strand.SetLightName("1", "MyName")
	if err != nil {
		t.Fatal("Error returned from SetLightName()")
	}

	if stubAPI.SetLightNameParamLightId != "light1" {
		t.Error("Didn't map from socketId to lightId.")
	}
	if stubAPI.SetLightNameParamName != "MyName" {
		t.Error("Didn't pass along name.")
	}
}

func Test_SetLightState(t *testing.T) {
	stubAPI := &huetest.StubAPI{}

	strand := NewLightStrand(1, stubAPI)
	strand.Lights.Set("1", "light1")

	lightState := &hue.LightState{}
	lightState.On = new(bool)
	*lightState.On = true
	err := strand.SetLightState("1", lightState)
	if err != nil {
		t.Fatal("Error returned from SetLightState()")
	}

	if stubAPI.SetLightStateParamLightId != "light1" {
		t.Error("Didn't map from socketId to lightId.")
	}
	if *stubAPI.SetLightStateParamLightState.On != true {
		t.Error("Didn't pass along the state.")
	}
}
