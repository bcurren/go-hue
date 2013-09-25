package hue

import (
	"github.com/bcurren/go-ssdp"
	"testing"
)

func Test_NewBridge(t *testing.T) {
	bridge := NewBridge("uuid:456-455", "192.168.0.1")
	httpServer, ok := bridge.client.conn.(*httpServer)
	if !ok {
		t.Fatal("Didn't create an httpServer properly.")
	}

	assertEqual(t, "192.168.0.1", httpServer.addr, "httpServer.addr")
}

func Test_reduceToHueDevices(t *testing.T) {
	devices := make([]ssdp.Device, 2, 2)
	devices[0].ModelUrl = HueModelUrl
	devices[1].ModelUrl = "http://someotherdevice.com"

	hueDevices := reduceToHueDevices(devices)

	assertEqual(t, 1, len(hueDevices), "len(hueDevices)")
	assertEqual(t, "http://www.meethue.com", hueDevices[0].ModelUrl, "ModelUrl")
}

func Test_convertHueDevicesToBridges(t *testing.T) {
	devices := make([]ssdp.Device, 1, 1)
	devices[0].UrlBase = "http://192.168.1.10:80/"
	devices[0].Udn = "uuid:8678-9078"

	bridges := convertHueDevicesToBridges(devices)

	assertEqual(t, 1, len(bridges), "len(bridges)")
	assertEqual(t, "uuid:8678-9078", bridges[0].UniqueId, "bridge.UniqueId")

	httpServer, ok := bridges[0].client.conn.(*httpServer)
	if !ok {
		t.Fatal("Bridge doesn't have httpServer.")
	}
	assertEqual(t, "http://192.168.1.10:80/", httpServer.addr, "httpServer.addr")
}

func Test_CreateUser(t *testing.T) {
	bridge, stubServer := NewStubBridge("post/index.json")

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
