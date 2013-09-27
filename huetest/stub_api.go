// Package with a Stub API that conforms to the hue.API interface for testing
// purposes.
package huetest

import (
	"github.com/bcurren/go-hue"
	"time"
)

// StubAPI contains the stub return values for all functions in the hue.API interface.
//
// Here is an example of how to use this struct:
//
// var stub huetest.StubAPI
// stub.GetLightsError = errors.NewError("This is an error")
//
// lights, err := stub.GetLights()
// if err.Error() != "This is an error" {
// 	fmt.Println("Stub is not working")
// }
type StubAPI struct {
	// Stub GetLights()
	GetLightsReturn []hue.Light
	GetLightsError  error

	// Stub GetNewLights()
	GetNewLightsReturn     []hue.Light
	GetNewLightsReturnTime time.Time
	GetNewLightsError      error

	// Stub SearchForNewLights()
	SearchForNewLightsError error

	// Stub GetLightAttributes()
	GetLightAttributesReturn       *hue.LightAttributes
	GetLightAttributesError        error
	GetLightAttributesParamLightId string

	// Stub SetLightName()
	SetLightNameError        error
	SetLightNameParamLightId string
	SetLightNameParamName    string

	// Stub SetLightState()
	SetLightStateError           error
	SetLightStateParamLightId    string
	SetLightStateParamLightState hue.LightState
}

func (s *StubAPI) GetLights() ([]hue.Light, error) {
	return s.GetLightsReturn, s.GetLightsError
}

func (s *StubAPI) GetNewLights() ([]hue.Light, time.Time, error) {
	return s.GetNewLightsReturn, s.GetNewLightsReturnTime, s.GetNewLightsError
}

func (s *StubAPI) SearchForNewLights() error {
	return s.SearchForNewLightsError
}

func (s *StubAPI) GetLightAttributes(lightId string) (*hue.LightAttributes, error) {
	s.GetLightAttributesParamLightId = lightId
	return s.GetLightAttributesReturn, s.GetLightAttributesError
}

func (s *StubAPI) SetLightName(lightId string, name string) error {
	s.SetLightNameParamLightId = lightId
	s.SetLightNameParamName = name
	return s.SetLightNameError
}

func (s *StubAPI) SetLightState(lightId string, state hue.LightState) error {
	s.SetLightStateParamLightId = lightId
	s.SetLightStateParamLightState = state
	return s.SetLightStateError
}
