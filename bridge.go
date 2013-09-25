package hue

import (
	"github.com/bcurren/go-ssdp"
	"time"
)

const HueModelUrl = "http://www.meethue.com"

type Bridge struct {
	UniqueId string
	client   *client
}

func NewBridge(uniqueId, addr string) *Bridge {
	client := NewHttpClient(addr)
	return &Bridge{UniqueId: uniqueId, client: client}
}

func FindBridges() ([]*Bridge, error) {
	devices, err := ssdp.SearchForDevices("upnp:rootdevice", 3*time.Second)
	if err != nil {
		return nil, err
	}

	hueDevices := reduceToHueDevices(devices)
	bridges := convertHueDevicesToBridges(hueDevices)

	return bridges, nil
}

func reduceToHueDevices(devices []ssdp.Device) []ssdp.Device {
	hueDevices := make([]ssdp.Device, 0, len(devices))

	for _, device := range devices {
		if device.ModelUrl == HueModelUrl {
			hueDevices = append(hueDevices, device)
		}
	}

	return hueDevices
}

func convertHueDevicesToBridges(devices []ssdp.Device) []*Bridge {
	bridges := make([]*Bridge, 0, len(devices))
	for _, device := range devices {
		bridges = append(bridges, NewBridge(device.Udn, device.UrlBase))
	}

	return bridges
}

func (b *Bridge) CreateUser(deviceType, username string) (*User, error) {
	url := "/api"

	requestObj := map[string]string{
		"devicetype": deviceType,
		"username":   username,
	}
	_, err := b.client.Post(url, &requestObj)
	if err != nil {
		return nil, err
	}

	return &User{Bridge: b, Username: username}, nil
}

func (b *Bridge) IsValidUser(username string) (bool, error) {
	testUser := NewUserWithBridge(username, b)

	// Get Configuration to determine if valid user
	_, err := testUser.GetLights()
	if err != nil {
		if apiError, ok := err.(*ApiError); ok {
			for _, apiErrorDetail := range apiError.Errors {
				if apiErrorDetail.Type == UnauthorizedUserErrorType {
					return false, nil
				}
			}
		}

		return false, err
	}

	return true, nil
}
