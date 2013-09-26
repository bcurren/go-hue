package hue

import (
	"fmt"
)

func (u *User) GetConfiguration() (*Configuration, error) {
	url := fmt.Sprintf("/api/%s/config", u.Username)

	var configuration *Configuration
	err := u.Bridge.client.Get(url, &configuration)
	if err != nil {
		return nil, err
	}

	return configuration, nil
}
