package hue

import (
	"fmt"
)

type Configuration struct {
	// Read and Write
	Name           string              `json:"name,omitempty"`
	ProxyAddress   string              `json:"proxyaddress,omitempty"`
	ProxyPort      *uint16             `json:"proxyport,omitempty"`
	IpAddress      string              `json:"ipaddress,omitempty"`
	Netmask        string              `json:"netmask,omitempty"`
	Gateway        string              `json:"gateway,omitempty"`
	Dhcp           *bool               `json:"dhcp,omitempty"`
	PortalServices *bool               `json:"portalservices,omitempty"`
	LinkButton     *bool               `json:"linkbutton,omitempty"`
	SoftwareUpdate *SoftwareUpdateInfo `json:"swupdate,omitempty"`

	// Read only
	Utc             string      `json:"utc,omitempty"`
	Whitelist       interface{} `json:"whitelist,omitempty"`
	SoftwareVersion string      `json:"swversion,omitempty"`
	Mac             string      `json:"mac,omitempty"`
}

type SoftwareUpdateInfo struct {
	UpdateState *uint  `json:"updatestate,omitempty"`
	Url         string `json:"url,omitempty"`
	Text        string `json:"text,omitempty"`
	Notify      *bool  `json:"notify,omitempty"`
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
// Add struct for and Whitelist
// Provide strong type for times, addresses, etc
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
