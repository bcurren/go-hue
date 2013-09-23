package hue

import (
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
	assertEqual(t, "/api/username1/lights", stubServer.uri, "request uri")

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
	assertEqual(t, "/api/username1/lights/new", stubServer.uri, "request uri")

	assertEqual(t, 1, len(lights), "len(lights)")

	assertEqual(t, "7", lights[0].Id, "lights[0].Id")
	assertEqual(t, "Hue Lamp 7", lights[0].Name, "lights[0].Name")
}

func Test_SearchForNewLights(t *testing.T) {
	user, stubServer := NewStubUser("post/username1/lights.json", "username1")

	err := user.SearchForNewLights()
	if err != nil {
		t.Fatal(err)
	}

	assertEqual(t, "POST", stubServer.method, "method is post")
	assertEqual(t, "/api/username1/lights", stubServer.uri, "request uri")
}

func Test_GetLightAttributes(t *testing.T) {
	user, stubServer := NewStubUser("get/username1/lights/light1.json", "username1")

	lightAttributes, err := user.GetLightAttributes("light1")
	if err != nil {
		t.Fatal(err)
	}

	// Request parameters
	assertEqual(t, "GET", stubServer.method, "method is get")
	assertEqual(t, "/api/username1/lights/light1", stubServer.uri, "request uri")

	// Light attributes
	assertEqual(t, "LC 1", lightAttributes.Name, "Name")
	assertEqual(t, "Living Colors", lightAttributes.Type, "Type")
	assertEqual(t, "LC0015", lightAttributes.ModelId, "ModelId")
	assertEqual(t, "1.0.3", lightAttributes.SoftwareVersion, "SoftwareVersion")
	pointSymbol := lightAttributes.PointSymbol.(map[string]interface{})
	assertEqual(t, "none", pointSymbol["1"].(string), "pointSymbol['1']")
	assertEqual(t, "none", pointSymbol["2"].(string), "pointSymbol['2']")
	assertEqual(t, "none", pointSymbol["3"].(string), "pointSymbol['3']")
	assertEqual(t, "none", pointSymbol["4"].(string), "pointSymbol['4']")
	assertEqual(t, "none", pointSymbol["5"].(string), "pointSymbol['5']")
	assertEqual(t, "none", pointSymbol["6"].(string), "pointSymbol['6']")
	assertEqual(t, "none", pointSymbol["7"].(string), "pointSymbol['7']")
	assertEqual(t, "none", pointSymbol["8"].(string), "pointSymbol['8']")

	// Light state
	lightState := lightAttributes.State
	assertEqual(t, true, *lightState.On, "On")
	assertEqual(t, uint8(200), *lightState.Brightness, "Brightness")
	assertEqual(t, uint16(50000), *lightState.Hue, "Hue")
	assertEqual(t, uint8(200), *lightState.Saturation, "Saturation")
	assertEqual(t, 0.5, lightState.Xy[0], "Xy")
	assertEqual(t, 0.25, lightState.Xy[1], "Xy")
	assertEqual(t, uint16(500), *lightState.ColorTemp, "ColorTemp")
	assertEqual(t, "none", lightState.Alert, "Alert")
	assertEqual(t, "none", lightState.Effect, "Effect")
	assertEqual(t, "hs", lightState.ColorMode, "ColorMode")
	assertEqual(t, true, lightState.Reachable, "Reachable")
}

func Test_ApiParseErrorString(t *testing.T) {
	err := NewApiError("string", 1, "user count")
	assertEqual(t, "Parsing error: expected type 'string' but received 'int' for user count.",
		err.Error(), "err message")
}
