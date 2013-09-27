# go-hue/huetest

A stub implementation of the hue.API interface for testing purposes. This package
allows you to easily stub out calls to the server.

Please see (godoc)[http://godoc.org/github.com/bcurren/go-hue/huetest] for detailed API
description.

## Example

See the examples directory for running code using this package. Below is an inline exerpt
from that directory.

```Go
package example

import (
	"github.com/bcurren/go-hue"
	"github.com/bcurren/go-hue/huetest"
	"testing"
)

func Test_GetLightAttributesSuccess(t *testing.T) {
	stubAPI := &huetest.StubAPI{}
	stubAPI.GetLightAttributesReturn = &hue.LightAttributes{}
	stubAPI.GetLightAttributesError = nil

	attrs, err := stubAPI.GetLightAttributes("light1")
	if err != stubAPI.GetLightAttributesError {
		t.Fatal("err is GetLightAttributesError so it must be nil")
	}
	if attrs != stubAPI.GetLightAttributesReturn {
		t.Fatal("attrs is GetLightAttributesReturn so it must be equal")
	}

	if "light1" != stubAPI.GetLightAttributesParamLightId {
		t.Fatal("GetLightAttributesParamLightId is set to light id parameter so it must be equal")
	}
}
```
