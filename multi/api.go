package multi

import (
	"github.com/bcurren/go-hue"
	"time"
)

type sGetLights struct {
	lights []hue.Light
	apiError error
}

// GetLights() is same as hue.User.GetLights() except all light ids are mapped to
// socket ids.
func (m *MultiAPI) GetLights() ([]hue.Light, error) {
	c := make(chan sGetLights)
	for _, api := range m.apis {
		go gGetLights(c, api)
	}
	return rGetLights(c, len(m.apis))
}

func gGetLights(c chan sGetLights, api hue.API) {
	lights, err := api.GetLights()
	c <- sGetLights{lights, err}
}

func rGetLights(c chan sGetLights, numRespones int) ([]hue.Light, error) {
	listOfErrors := make([]error, 0, 1)
	listOfLights := make([][]hue.Light, 0, numRespones)
	for i := 0; i < numRespones; i++ {
		result := <- c
		if result.apiError != nil {
			listOfErrors = append(listOfErrors, result.apiError)
		} else {
			listOfLights = append(listOfLights, result.lights)
		}
	}
	
	return mergeLightSlice(listOfLights), nil
}

func mergeLightSlice(listOfLights [][]hue.Light) []hue.Light {
	countOfLights := 0
	for _, lights := range listOfLights {
		countOfLights += len(lights)
	}
	
	mergedLights := make([]hue.Light, countOfLights)
	copyTo := 0
	for _, lights := range listOfLights {
		copy(mergedLights[copyTo:], lights)
		copyTo += len(lights)
	}
	
	return mergedLights
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
