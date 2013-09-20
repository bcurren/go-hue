package hue

import (
	"testing"
	"fmt"
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
	assertEqual(t, fmt.Sprintf("/api/%s/lights", user.Username) , stubServer.uri, "uri is correct")
	
	assertEqual(t, 2, len(lights), "len(lights)")
	
	assertEqual(t, "1", lights[0].Id, "lights[0].Id")
	assertEqual(t, "Bedroom", lights[0].Name, "lights[0].Name")
	
	assertEqual(t, "2", lights[1].Id, "lights[1].Id")
	assertEqual(t, "Kitchen", lights[1].Name, "lights[1].Name")
}