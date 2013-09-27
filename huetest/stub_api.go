// Package with a Stub Api that conforms to the hue.Api interface for testing
// purposes.
package huetest

import (
	"github.com/bcurren/go-hue"
	"time"
)

// StubApi contains the stub return values for all functions in the hue.Api interface.
//
// Here is an example of how to use this struct:
//
// var stub huetest.StubApi
// stub.GetLightsError = errors.NewError("This is an error")
//
// lights, err := stub.GetLights()
// if err.Error() != "This is an error" {
// 	fmt.Println("Stub is not working")
// }
type StubApi struct {
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

	// Stub GetConfiguration()
	GetConfigurationReturn *hue.Configuration
	GetConfigurationError  error
}

func (s *StubApi) GetLights() ([]hue.Light, error) {
	return s.GetLightsReturn, s.GetLightsError
}

func (s *StubApi) GetNewLights() ([]hue.Light, time.Time, error) {
	return s.GetNewLightsReturn, s.GetNewLightsReturnTime, s.GetNewLightsError
}

func (s *StubApi) SearchForNewLights() error {
	return s.SearchForNewLightsError
}

func (s *StubApi) GetLightAttributes(lightId string) (*hue.LightAttributes, error) {
	s.GetLightAttributesParamLightId = lightId
	return s.GetLightAttributesReturn, s.GetLightAttributesError
}

func (s *StubApi) SetLightName(lightId string, name string) error {
	s.SetLightNameParamLightId = lightId
	s.SetLightNameParamName = name
	return s.SetLightNameError
}

func (s *StubApi) SetLightState(lightId string, state hue.LightState) error {
	s.SetLightStateParamLightId = lightId
	s.SetLightStateParamLightState = state
	return s.SetLightStateError
}

func (s *StubApi) GetConfiguration() (*hue.Configuration, error) {
	return s.GetConfigurationReturn, s.GetConfigurationError
}
