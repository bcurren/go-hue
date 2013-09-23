package hue

import (
	"testing"
)

func Test_NewHttpClient(t *testing.T) {
	client := NewHttpClient("192.168.10.2")
	httpServer, ok := client.conn.(*httpServer)
	if !ok {
		t.Fatal("Client doesn't have an httpServer")
	}

	assertEqual(t, "192.168.10.2", httpServer.addr, "httpServer.addr")
}

func Test_GetWithEmptyRequestBody(t *testing.T) {
	c := NewStubClient("get/username1/lights.json")

	var lights map[string]map[string]string
	err := c.Get("/api/username1/lights", &lights)
	if err != nil {
		t.Fatal(err)
	}

	assertEqual(t, 2, len(lights), "Num lights returned.")

	assertNotNil(t, lights["1"], "lights[1]")
	assertEqual(t, "Bedroom", lights["1"]["name"], `lights["1"]["name"]`)

	assertNotNil(t, lights["2"], "lights[2]")
	assertEqual(t, "Kitchen", lights["2"]["name"], `lights["2"]["name"]`)
}

func Test_GetWithResponseError(t *testing.T) {
	c := NewStubClient("errors/unauthorized_user.json")

	var lights map[string]map[string]string
	err := c.Get("/api/username1/lights", &lights)
	apiError, ok := err.(*ApiError)
	if !ok {
		t.Fatal("Should return an unauthorized user error.")
	}

	assertEqual(t, 1, len(apiError.Errors), "Num errors returned.")

	assertEqual(t, 1, apiError.Errors[0].Type, "error.Type")
	assertEqual(t, "/lights", apiError.Errors[0].Address, "error.Address")
	assertEqual(t, "unauthorized user", apiError.Errors[0].Description, "error.Description")
}

func Test_SendNonGetAllSuccessResponse(t *testing.T) {
	c := NewStubClient("post/username1/lights.json")

	successes, err := c.Send("POST", "/api/username1/lights", nil)
	if err != nil {
		t.Fatal("Should be successful.")
	}

	assertEqual(t, 1, len(successes), "len(successes)")
	assertEqual(t, "Searching for new devices", successes[0]["/lights"], "val of /lights")
}

func Test_SendNonGetAllErrorResponse(t *testing.T) {
	c := NewStubClient("errors/unauthorized_user.json")

	successes, err := c.Send("POST", "/api/username1/lights", nil)
	apiError, ok := err.(*ApiError)
	if !ok {
		t.Fatal("Error should be ApiError.")
	}

	errors := apiError.Errors
	assertEqual(t, 1, len(errors), "len(errors)")
	assertEqual(t, "/lights", errors[0].Address, "errors[0].Address")
	assertEqual(t, "unauthorized user", errors[0].Description, "errors[0].Description")

	if successes != nil {
		t.Error("Success should be nil when 0 are returned")
	}
}

func Test_SendNonGetMixedSuccessAndErrorResponse(t *testing.T) {
	c := NewStubClient("errors/mixed_errors.json")

	successes, err := c.Send("POST", "/api/username1/lights", nil)
	apiError, ok := err.(*ApiError)
	if !ok {
		t.Fatal("Error should be ApiError.")
	}

	assertEqual(t, 1, len(successes), "len(successes)")
	assertEqual(t, true, successes[0]["/lights/light1/state/on"], "val of light on")

	errors := apiError.Errors
	assertEqual(t, 1, len(errors), "len(errors)")
	assertEqual(t, "/fake", errors[0].Address, "errors[0].Address")
	assertEqual(t, "link button not pressed", errors[0].Description, "errors[0].Description")
}
