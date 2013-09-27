package hue

import (
	"testing"
)

func Test_UserImplementsAPIInterface(t *testing.T) {
	user, _ := NewStubUser("get/username1/config.json", "username1")
	funcThatTakesAPIAsParameter(user)
}

func funcThatTakesAPIAsParameter(api API) {
	// noop
}
