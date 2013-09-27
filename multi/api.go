package multi

import (
	"github.com/bcurren/go-hue"
	"time"
)

// GetLights() is same as hue.User.GetLights() except all light ids are mapped to
// socket ids.
func (m *MultiAPI) GetLights() ([]hue.Light, error) {
	listOfErrors := make([]error, 0, 1)
	listOfLights := make([][]hue.Light, 0, len(m.apis))
	
	for _, api := range m.apis {
		lights, err := api.GetLights()
		if err != nil {
			listOfErrors = append(listOfErrors, err)
		} else {
			listOfLights = append(listOfLights, lights)
		}
	}
	
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
	
	return mergedLights, nil
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
