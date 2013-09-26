package hue

import (
	"testing"
)

func Test_UserImplementsApiInterface(t *testing.T) {
	user, _ := NewStubUser("get/username1/config.json", "username1")
	funcThatTakesApiAsParameter(user)
}

func funcThatTakesApiAsParameter(api Api) {
	// noop
}
