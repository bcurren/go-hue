package hue

import (
	"testing"
)

func Test_GetWhenSuccess(t *testing.T) {
	c := NewStubClient("")

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

func Test_GetWhenError(t *testing.T) {
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

func assertEqual(t *testing.T, expected interface{}, actual interface{}, errorMessage string) {
	if expected != actual {
		t.Errorf("%q is not equal to %q. %q", expected, actual, errorMessage)
	}
}

func assertNotNil(t *testing.T, obj interface{}, errorMessage string) {
	if obj == nil {
		t.Errorf("%q should not be nil. %q", obj, errorMessage)
	}
}
