package multi

import (
	"github.com/bcurren/go-hue"
	"time"
)

type rGetLights struct {
	lights []hue.Light
	err error
}

// GetLights() is same as hue.User.GetLights() except all light ids are mapped to
// socket ids.
func (m *MultiAPI) GetLights() ([]hue.Light, error) {
	c := make(chan rGetLights)
	for _, api := range m.apis {
		go gGetLights(c, api)
	}
	return lGetLights(c, len(m.apis))
}

func gGetLights(c chan rGetLights, api hue.API) {
	lights, err := api.GetLights()
	// TODO: Map lights ids to unique api light ids
	c <- rGetLights{lights, err}
}

func lGetLights(c chan rGetLights, nResponses int) ([]hue.Light, error) {
	lErrors := make([]error, 0, 1)
	lLights := make([][]hue.Light, 0, nResponses)
	
	for i := 0; i < nResponses; i++ {
		result := <- c
		if result.err != nil {
			lErrors = append(lErrors, result.err)
		} else {
			lLights = append(lLights, result.lights)
		}
	}
	
	return mergeLights(lLights), mergeErrors(lErrors)
}

// GetNewLights() is same as hue.User.GetNewLights() except all light ids are mapped to
// socket ids.
func (m *MultiAPI) GetNewLights() ([]hue.Light, time.Time, error) {
	return nil, time.Now(), nil
}

// SearchForNewLights() is same as hue.User.SearchForNewLights() except all light ids are mapped to
// socket ids.
func (m *MultiAPI) SearchForNewLights() error {
	return nil
}

// GetLightAttributes() is same as hue.User.GetLightAttributes() except all light ids are mapped to
// socket ids.
func (m *MultiAPI) GetLightAttributes(socketId string) (*hue.LightAttributes, error) {
	return nil, nil
}

// SetLightName() is same as hue.User.SetLightName() except all light ids are mapped to
// socket ids.
func (m *MultiAPI) SetLightName(socketId string, name string) error {
	return nil
}

// SetLightState() is same as hue.User.SetLightState() except all light ids are mapped to
// socket ids.
func (m *MultiAPI) SetLightState(socketId string, state *hue.LightState) error {
	return nil
}

func mergeLights(lLights [][]hue.Light) []hue.Light {
	countOfLights := 0
	for _, lights := range lLights {
		countOfLights += len(lights)
	}
	
	mLights := make([]hue.Light, countOfLights)
	copyTo := 0
	for _, lights := range lLights {
		copy(mLights[copyTo:], lights)
		copyTo += len(lights)
	}
	
	return mLights
}

func mergeErrors(lErrors []error) error {
	return nil
}
