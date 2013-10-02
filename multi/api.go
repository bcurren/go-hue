package multi

import (
	"fmt"
	"github.com/bcurren/go-hue"
	"strconv"
	"strings"
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
	for _, api := range m.APIs {
		go m.gGetLights(c, api)
	}
	return m.lGetLights(c, len(m.APIs))
}

func (m *MultiAPI) gGetLights(c chan rGetLights, api hue.API) {
	lights, err := api.GetLights()
	c <- rGetLights{m.mapLightIds(api, lights), err}
}

func (m *MultiAPI) lGetLights(c chan rGetLights, nResponses int) ([]hue.Light, error) {
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
	for _, api := range m.APIs {
		go m.gGetNewLights(c, api)
	}
	return m.lGetNewLights(c, len(m.APIs))
}

func (m *MultiAPI) gGetNewLights(c chan rGetNewLights, api hue.API) {
	lights, lastScan, err := api.GetNewLights()
	c <- rGetNewLights{m.mapLightIds(api, lights), lastScan, err}
}

func (m *MultiAPI) lGetNewLights(c chan rGetNewLights, nResponses int) ([]hue.Light, time.Time, error) {
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
	for _, api := range m.APIs {
		go m.gSearchForNewLights(c, api)
	}
	return m.lSearchForNewLights(c, len(m.APIs))
}

func (m *MultiAPI) gSearchForNewLights(c chan rSearchForNewLights, api hue.API) {
	err := api.SearchForNewLights()
	c <- rSearchForNewLights{err}
}

func (m *MultiAPI) lSearchForNewLights(c chan rSearchForNewLights, nResponses int) error {
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
	api, newLightId, err := m.findAPIAndLightId(lightId, fmt.Sprintf("/lights/%s", lightId))
	if err != nil {
		return nil, err
	}
	return api.GetLightAttributes(newLightId)
}

// SetLightName() is same as hue.User.SetLightName() except all light ids are mapped to
// socket ids.
func (m *MultiAPI) SetLightName(lightId string, name string) error {
	api, newLightId, err := m.findAPIAndLightId(lightId, fmt.Sprintf("/lights/%s", lightId))
	if err != nil {
		return err
	}
	return api.SetLightName(newLightId, name)
}

// SetLightState() is same as hue.User.SetLightState() except all light ids are mapped to
// socket ids.
func (m *MultiAPI) SetLightState(lightId string, state *hue.LightState) error {
	api, newLightId, err := m.findAPIAndLightId(lightId, fmt.Sprintf("/lights/%s/state", lightId))
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

func (m *MultiAPI) mapLightIds(api hue.API, lights []hue.Light) []hue.Light {
	for i, light := range lights {
		lights[i].Id = m.mapAPIAndLightId(api, light.Id)
	}
	return lights
}

func (m *MultiAPI) apiToId(api hue.API) string {
	for i, checkApi := range m.APIs {
		if api == checkApi {
			return strconv.Itoa(i)
		}
	}
	return ""
}

func (m *MultiAPI) idToApi(id string) hue.API {
	idAsInt, err := strconv.Atoi(id)
	if err != nil {
		return nil
	}
	if idAsInt < 0 || idAsInt >= len(m.APIs) {
		return nil
	}

	return m.APIs[idAsInt]
}

func (m *MultiAPI) mapAPIAndLightId(api hue.API, lightId string) string {
	return fmt.Sprintf("%s#%s", m.apiToId(api), lightId)
}

func (m *MultiAPI) findAPIAndLightId(lightId, address string) (hue.API, string, error) {
	splitLightId := strings.Split(lightId, "#")
	api := m.idToApi(splitLightId[0])
	if api == nil {
		return nil, "", createResourceNotAvailableAPIError(lightId, address)
	}

	return api, splitLightId[1], nil
}

func createResourceNotAvailableAPIError(resourceId, address string) error {
	apiError := &hue.APIError{}
	apiError.Errors = make([]hue.APIErrorDetail, 1, 1)
	apiError.Errors[0].Type = hue.ResourceNotAvailableErrorType
	apiError.Errors[0].Address = address
	apiError.Errors[0].Description = fmt.Sprintf("resource, %s, not available", resourceId)

	return apiError
}
