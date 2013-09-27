package huetest

import (
	"github.com/bcurren/go-hue"
	"testing"
)

func Test_UserImplementsApiInterface(t *testing.T) {
	stub := &StubApi{}
	funcThatTakesApiAsParameter(stub)
}

func funcThatTakesApiAsParameter(api hue.Api) {
	// noop
}
