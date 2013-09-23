package hue

type Bridge struct {
	client *client
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
