package hue

import (
	"testing"
)

func Test_cleanUrl(t *testing.T) {
	url := cleanUrl("http://10.0.1.1//api/newdeveloper")
	assertEqual(t, "http://10.0.1.1/api/newdeveloper", url, "cleanUrl")
}
