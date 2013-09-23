package hue

import (
	"testing"
)

func Test_CreateUser(t *testing.T) {
	bridge, stubServer := NewStubBridge("StubBridge", "post/index.json")
	
	user, err := bridge.CreateUser("myDeviceType", "myUsername")
	if err != nil {
		t.Fatal("Error when creating user.", err)
	}
	
	assertEqual(t, "POST", stubServer.method, "method is put")
	assertEqual(t, "/api", stubServer.uri, "request uri")
	assertEqual(t, `{"devicetype":"myDeviceType","username":"myUsername"}`, 
		stubServer.requestJson, "request json")
		
	assertEqual(t, bridge, user.Bridge, "user.Bridge")
	assertEqual(t, "myUsername", user.Username, "user.Username")
}