package hue

import (
	"testing"
)

func Test_cleanURL(t *testing.T) {
	url := cleanURL("http://10.0.1.1//api/newdeveloper")
	assertEqual(t, "http://10.0.1.1/api/newdeveloper", url, "cleanURL")
}
