package hue

import (
	"fmt"
)

type Configuration struct {
	// Read and Write
	Name         string  `json:"name,omitempty"`
	ProxyAddress string  `json:"proxyaddress,omitempty"`
	ProxyPort    *uint16 `json:"proxyport,omitempty"`
	IpAddress    string  `json:"ipaddress,omitempty"`
	Netmask      string  `json:"netmask,omitempty"`
	Gateway      string  `json:"gateway,omitempty"`
	Dhcp         *bool   `json:"dhcp,omitempty"`

	// These can be updated? Strange
	SoftwareUpdate interface{} `json:"swupdate,omitempty"`
	LinkButton     *bool       `json:"linkbutton,omitempty"`
	PortalServices *bool       `json:"portalservices,omitempty"`

	// Read only
	Utc             string      `json:"utc,omitempty"`
	Whitelist       interface{} `json:"whitelist,omitempty"`
	SoftwareVersion string      `json:"swversion,omitempty"`
	Mac             string      `json:"mac,omitempty"`
}

func (u *User) GetConfiguration() (*Configuration, error) {
	url := fmt.Sprintf("/api/%s/config", u.Username)

	var configuration *Configuration
	err := u.Bridge.client.Get(url, &configuration)
	if err != nil {
		return nil, err
	}

	return configuration, nil
}

// TODO Configuration Api Methods:
// func (u *User) DeleteUser(*User) error
// func (u *User) UpdateConfiguration(*Configuration) error
// func (u *User) GetDataStore() (DataStore, error)
//

// type DataStore struct {
// 	Lights []Light
// 	Groups []Group
// 	Schedules []Schedule
// 	Config Configuration
// }
