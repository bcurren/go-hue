package multi

import (
	"fmt"
	"github.com/bcurren/go-hue"
	"time"
)

type rGetLights struct {
	lights []hue.Light
	err    error
}

type rGetNewLights struct {
	lights   []hue.Light
	lastScan time.Time
	err      error
}

type rSearchForNewLights struct {
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
	c <- rGetLights{mapLightIds(api, lights), err}
}

func lGetLights(c chan rGetLights, nResponses int) ([]hue.Light, error) {
	lErrors := make([]error, 0, 1)
	lLights := make([][]hue.Light, 0, nResponses)

	for i := 0; i < nResponses; i++ {
		result := <-c
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
	c := make(chan rGetNewLights)
	for _, api := range m.apis {
		go gGetNewLights(c, api)
	}
	return lGetNewLights(c, len(m.apis))
}

func gGetNewLights(c chan rGetNewLights, api hue.API) {
	lights, lastScan, err := api.GetNewLights()
	c <- rGetNewLights{mapLightIds(api, lights), lastScan, err}
}

func lGetNewLights(c chan rGetNewLights, nResponses int) ([]hue.Light, time.Time, error) {
	lErrors := make([]error, 0, 1)
	lLights := make([][]hue.Light, 0, nResponses)
	lLastScan := make([]time.Time, 0, nResponses)

	for i := 0; i < nResponses; i++ {
		result := <-c
		if result.err != nil {
			lErrors = append(lErrors, result.err)
		}
		if result.lights != nil {
			lLights = append(lLights, result.lights)
		}
		lLastScan = append(lLastScan, result.lastScan)
	}

	return mergeLights(lLights), mergeTime(lLastScan), mergeErrors(lErrors)
}

// SearchForNewLights() is same as hue.User.SearchForNewLights() except all light ids are mapped to
// socket ids.
func (m *MultiAPI) SearchForNewLights() error {
	c := make(chan rSearchForNewLights)
	for _, api := range m.apis {
		go gSearchForNewLights(c, api)
	}
	return lSearchForNewLights(c, len(m.apis))
}

func gSearchForNewLights(c chan rSearchForNewLights, api hue.API) {
	err := api.SearchForNewLights()
	c <- rSearchForNewLights{err}
}

func lSearchForNewLights(c chan rSearchForNewLights, nResponses int) error {
	lErrors := make([]error, 0, 1)

	for i := 0; i < nResponses; i++ {
		result := <-c
		if result.err != nil {
			lErrors = append(lErrors, result.err)
		}
	}

	return mergeErrors(lErrors)
}

// GetLightAttributes() is same as hue.User.GetLightAttributes() except all light ids are mapped to
// socket ids.
func (m *MultiAPI) GetLightAttributes(lightId string) (*hue.LightAttributes, error) {
	api, newLightId, err := m.findAPIAndLightId(lightId)
	if err != nil {
		return nil, err
	}
	return api.GetLightAttributes(newLightId)
}

// SetLightName() is same as hue.User.SetLightName() except all light ids are mapped to
// socket ids.
func (m *MultiAPI) SetLightName(lightId string, name string) error {
	api, newLightId, err := m.findAPIAndLightId(lightId)
	if err != nil {
		return err
	}
	return api.SetLightName(newLightId, name)
}

// SetLightState() is same as hue.User.SetLightState() except all light ids are mapped to
// socket ids.
func (m *MultiAPI) SetLightState(lightId string, state *hue.LightState) error {
	api, newLightId, err := m.findAPIAndLightId(lightId)
	if err != nil {
		return err
	}
	return api.SetLightState(newLightId, state)
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
	if lErrors == nil || len(lErrors) <= 0 {
		return nil
	}

	apiErrorDetails := make([]hue.APIErrorDetail, 0, len(lErrors))

	// Collect all error details and exit if non API error found
	// TODO: Collect all non API Errors for an array of errors return type
	for _, err := range lErrors {
		if apiError, ok := err.(*hue.APIError); ok {
			for _, apiErrorDetail := range apiError.Errors {
				apiErrorDetails = append(apiErrorDetails, apiErrorDetail)
			}
		} else {
			return err
		}
	}
	return &hue.APIError{apiErrorDetails}
}

func mergeTime(lTime []time.Time) time.Time {
	return lTime[0]
}

func mapLightIds(api hue.API, lLights []hue.Light) []hue.Light {
	return lLights
}

func (m *MultiAPI) findAPIAndLightId(lightId string) (hue.API, string, error) {
	return m.apis[0], lightId, nil
}

func createResourceNotAvailableAPIError(resourceId, address string) error {
	apiError := &hue.APIError{}
	apiError.Errors = make([]hue.APIErrorDetail, 1, 1)
	apiError.Errors[0].Type = hue.ResourceNotAvailableErrorType
	apiError.Errors[0].Address = address
	apiError.Errors[0].Description = fmt.Sprintf("resource, %s, not available", resourceId)

	return apiError
}
