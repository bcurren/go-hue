package hue

import (
// "github.com/bcurren/go-udpn"
// "time"
)

type Bridge struct {
	client *client
}

func NewBridge(addr string) *Bridge {
	client := NewHttpClient(addr)
	return &Bridge{client: client}
}

// TODO: Need to add some functionality to udpn library
// func FindBridges() ([]*Bridge, error) {
// 	_, err := udpn.Search("upnp:rootdevice", 3*time.Second)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return nil, nil
// }

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
