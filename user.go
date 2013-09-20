package hue

import (
	"fmt"
	"time"
	"errors"
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

	var lightsMap map[string]map[string]string
	err := u.Bridge.client.Get(url, &lightsMap)
	if err != nil {
		return nil, err
	}

	lights := make([]Light, 0, 10)
	for lightId, lightMap := range lightsMap {
		lights = append(lights, Light{Id: lightId, Name: lightMap["name"]})
	}

	return lights, nil
}

func (u *User) GetNewLights() ([]Light, time.Time, error) {
	url := fmt.Sprintf("/api/%s/lights/new", u.Username)

	var newLightsResponse map[string]interface{}
	err := u.Bridge.client.Get(url, &newLightsResponse)
	if err != nil {
		return nil, time.Time{}, err
	}

	lastScanString, ok := newLightsResponse["lastscan"].(string)
	if !ok {
		return nil, time.Time{}, errors.New("Error parsing lastscan")
	}
	lastScan, err := time.Parse(ISO8601, lastScanString)
	if err != nil {
		return nil, time.Time{}, err
	}
	delete(newLightsResponse, "lastscan")
	
	lights, err := parseLights(newLightsResponse)
	if err != nil {
		return nil, time.Time{}, err
	}

	return lights, lastScan, nil
}

func parseLights(newLightsResponse map[string]interface{}) ([]Light, error) {
	lights := make([]Light, 0, 10)
	
	for lightId, lightInterface := range newLightsResponse {
		lightMap, ok := lightInterface.(map[string]interface{})
		if !ok {
			// return nil, newParseError(reflect.TypeOf(map[string]interface{}), lightInterface, "lights map")
			return nil, newParseError()
		}
		name, ok := lightMap["name"].(string)
		if !ok {
			// return nil, newParseError(string, lightMap["name"], "light name")
			return nil, newParseError()
		}
		lights = append(lights, Light{Id: lightId, Name: name})
	}
	
	return lights, nil
}

// func newParseError(expected reflect.Type, actual interface{}, context string) error {
func newParseError() error {
	return errors.New("Error parsing api response.")
}
