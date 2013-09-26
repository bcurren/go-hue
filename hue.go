package hue

import (
	"time"
)

type Api interface {
	GetLights() ([]Light, error)
	GetNewLights() ([]Light, time.Time, error)
	SearchForNewLights() error
	GetLightAttributes(lightId string) (*LightAttributes, error)
	SetLightName(lightId string, name string) error
	SetLightState(lightId string, state LightState) error
	GetConfiguration() (*Configuration, error)

	// TODO Configuration Api Methods:
	// func (u *User) DeleteUser(*User) error
	// func (u *User) UpdateConfiguration(*Configuration) error
	// func (u *User) GetDataStore() (DataStore, error)
}

type Light struct {
	Id   string
	Name string
}

type LightState struct {
	On             *bool     `json:"on,omitempty"`
	Brightness     *uint8    `json:"bri,omitempty"`
	Hue            *uint16   `json:"hue,omitempty"`
	Saturation     *uint8    `json:"sat,omitempty"`
	Xy             []float64 `json:"xy,omitempty"`
	ColorTemp      *uint16   `json:"ct,omitempty"`
	Alert          string    `json:"alert,omitempty"`
	Effect         string    `json:"effect,omitempty"`
	TransitionTime *uint16   `json:"transitiontime,omitempty"` /* write only */
	ColorMode      string    `json:"colormode,omitempty"`      /* read only */
	Reachable      bool      `json:"reachable,omitempty"`      /* read only */
}

type LightAttributes struct {
	Name            string      `json:"name"`
	State           *LightState `json:"state"`
	Type            string      `json:"type"`
	ModelId         string      `json:"modelid"`
	SoftwareVersion string      `json:"swversion"`
	PointSymbol     interface{} `json:"pointsymbol"`
}

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
	Utc             string                       `json:"utc,omitempty"`
	Whitelist       map[string]map[string]string `json:"whitelist,omitempty"`
	SoftwareVersion string                       `json:"swversion,omitempty"`
	Mac             string                       `json:"mac,omitempty"`
}

type SoftwareUpdateInfo struct {
	UpdateState *uint  `json:"updatestate,omitempty"`
	Url         string `json:"url,omitempty"`
	Text        string `json:"text,omitempty"`
	Notify      *bool  `json:"notify,omitempty"`
}
