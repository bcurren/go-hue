package example

import (
	"github.com/bcurren/go-hue"
	"github.com/bcurren/go-hue/multi"
	"testing"
)

func Test_Blah(t *testing.T) {
	multi := multi.NewMultiAPI()

	bridges, err := hue.FindBridges()
	if err != nil {
		panic(err)
	}

	for _, bridge := range bridges {
		user := hue.NewUserWithBridge("test", bridge)
		multi.AddAPI(user)
	}

	_, err = multi.GetLights()
	if err != nil {
		panic(err)
	}
}
