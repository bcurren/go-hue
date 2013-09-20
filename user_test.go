package hue

import (
	"fmt"
	"testing"
	"time"
)

func NewStubUser(stubFilename string, username string) (*User, *stubServer) {
	stubClient := NewStubClient(stubFilename)
	stubServer, ok := stubClient.conn.(*stubServer)
	if !ok {
		panic("Not using a stub server in tests!")
	}
	bridge := &Bridge{Name: "StubBridge", client: stubClient}

	return &User{Bridge: bridge, Username: username}, stubServer
}

func Test_GetLights(t *testing.T) {
	user, stubServer := NewStubUser("get/username1/lights.json", "username1")

	lights, err := user.GetLights()
	if err != nil {
		t.Fatal(err)
	}

	assertEqual(t, "GET", stubServer.method, "method is get")
	assertEqual(t, fmt.Sprintf("/api/%s/lights", user.Username), stubServer.uri, "uri is correct")

	assertEqual(t, 2, len(lights), "len(lights)")

	assertEqual(t, "1", lights[0].Id, "lights[0].Id")
	assertEqual(t, "Bedroom", lights[0].Name, "lights[0].Name")

	assertEqual(t, "2", lights[1].Id, "lights[1].Id")
	assertEqual(t, "Kitchen", lights[1].Name, "lights[1].Name")
}

func Test_GetNewLights(t *testing.T) {
	user, stubServer := NewStubUser("get/username1/lights/new.json", "username1")

	lights, lastScan, err := user.GetNewLights()
	if err != nil {
		t.Fatal(err)
	}

	expectatedLastScan, err := time.Parse(ISO8601, "2012-10-29T12:00:00")
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, expectatedLastScan, lastScan, "lastScan")

	assertEqual(t, "GET", stubServer.method, "method is get")
	assertEqual(t, fmt.Sprintf("/api/%s/lights/new", user.Username), stubServer.uri, "uri is correct")

	assertEqual(t, 1, len(lights), "len(lights)")

	assertEqual(t, "7", lights[0].Id, "lights[0].Id")
	assertEqual(t, "Hue Lamp 7", lights[0].Name, "lights[0].Name")
}

func Test_ApiParseErrorString(t *testing.T) {
	err := NewApiError("string", 1, "user count")
	assertEqual(t, "Parsing error: expected type 'string' but received 'int' for user count.",
		err.Error(), "err message")
}
