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
	SetLightStateParamLightState *hue.LightState
}

// Stub version of GetLights. This function just returns GetLightsReturn and
// GetLightsError in StubAPI struct.
func (s *StubAPI) GetLights() ([]hue.Light, error) {
	return s.GetLightsReturn, s.GetLightsError
}

// Stub version of GetNewLights. This function just returns GetNewLightsReturn,
// GetNewLightsReturnTime and GetNewLightsError in StubAPI struct.
func (s *StubAPI) GetNewLights() ([]hue.Light, time.Time, error) {
	return s.GetNewLightsReturn, s.GetNewLightsReturnTime, s.GetNewLightsError
}

// Stub version of SearchForNewLights. This function just returns SearchForNewLightsError
// in StubAPI struct.
func (s *StubAPI) SearchForNewLights() error {
	return s.SearchForNewLightsError
}

// Stub version of GetLightAttributes. This function just returns GetLightAttributesReturn and
// GetLightAttributesError in StubAPI struct. It will also set GetLightAttributesParamLightId in
// StubAPI to the value of the lightId parameter.
func (s *StubAPI) GetLightAttributes(lightId string) (*hue.LightAttributes, error) {
	s.GetLightAttributesParamLightId = lightId
	return s.GetLightAttributesReturn, s.GetLightAttributesError
}

// Stub version of SetLightName. This function just returns SetLightNameError in StubAPI struct.
// It will also set SetLightNameParamLightId to the value of lightId parameter and
// SetLightNameParamName to the value of the name parameter.
func (s *StubAPI) SetLightName(lightId string, name string) error {
	s.SetLightNameParamLightId = lightId
	s.SetLightNameParamName = name
	return s.SetLightNameError
}

// Stub version of SetLightState. This function just returns SetLightStateError in StubAPI struct.
// It will also set SetLightStateParamLightId to the value of lightId parameter and
// SetLightStateParamLightState to the value of the state parameter.
func (s *StubAPI) SetLightState(lightId string, state *hue.LightState) error {
	s.SetLightStateParamLightId = lightId
	s.SetLightStateParamLightState = state
	return s.SetLightStateError
}
