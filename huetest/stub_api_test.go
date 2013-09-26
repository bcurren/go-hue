package huetest

import (
	"testing"
	"github.com/bcurren/go-hue"
)

func Test_UserImplementsApiInterface(t *testing.T) {
	var stub StubApi
	funcThatTakesApiAsParameter(stub)
}

func funcThatTakesApiAsParameter(api hue.Api) {
	// noop
}
