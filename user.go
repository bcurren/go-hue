package hue

import (
	"fmt"
	"time"
)

const ISO8601 = "2006-01-02T15:04:05"

type User struct {
	Bridge   *Bridge
	Username string
}

type Light struct {
	Id   string
	Name string
}

func (u *User) GetLights() ([]Light, error) {
	url := fmt.Sprintf("/api/%s/lights", u.Username)

	var lightsMap map[string]interface{}
	err := u.Bridge.client.Get(url, &lightsMap)
	if err != nil {
		return nil, err
	}

	lights, err := parseLights(lightsMap)
	if err != nil {
		return nil, err
	}

	return lights, nil
}

func (u *User) GetNewLights() ([]Light, time.Time, error) {
	url := fmt.Sprintf("/api/%s/lights/new", u.Username)

	var lightsMap map[string]interface{}
	err := u.Bridge.client.Get(url, &lightsMap)
	if err != nil {
		return nil, time.Time{}, err
	}

	lastScanString, ok := lightsMap["lastscan"].(string)
	if !ok {
		return nil, time.Time{}, NewApiError("string", lightsMap["lastscan"], "lastscan")
	}
	lastScan, err := time.Parse(ISO8601, lastScanString)
	if err != nil {
		return nil, time.Time{}, err
	}
	delete(lightsMap, "lastscan")

	lights, err := parseLights(lightsMap)
	if err != nil {
		return nil, time.Time{}, err
	}

	return lights, lastScan, nil
}

func parseLights(lightsMap map[string]interface{}) ([]Light, error) {
	lights := make([]Light, 0, 10)

	for lightId, lightInterface := range lightsMap {
		lightMap, ok := lightInterface.(map[string]interface{})
		if !ok {
			return nil, NewApiError("map[string]interface{}", lightInterface, "lights map")
		}
		name, ok := lightMap["name"].(string)
		if !ok {
			return nil, NewApiError("string", lightMap["name"], "lights name")
		}
		lights = append(lights, Light{Id: lightId, Name: name})
	}

	return lights, nil
}

type ApiParseError struct {
	Expected string
	Actual   interface{}
	Context  string
}

func NewApiError(expected string, actual interface{}, context string) error {
	return &ApiParseError{Expected: expected, Actual: actual, Context: context}
}

func (e *ApiParseError) Error() string {
	return fmt.Sprintf("Parsing error: expected type '%s' but received '%T' for %s.",
		e.Expected, e.Actual, e.Context)
}
