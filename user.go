package hue

import (
	"fmt"
)

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
