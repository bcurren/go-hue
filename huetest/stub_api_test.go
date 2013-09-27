package huetest

import (
	"github.com/bcurren/go-hue"
	"testing"
)

func Test_UserImplementsAPIInterface(t *testing.T) {
	stub := &StubAPI{}
	funcThatTakesAPIAsParameter(stub)
}

func funcThatTakesAPIAsParameter(api hue.API) {
	// noop
}
