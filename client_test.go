package hue

import (
	"testing"
	"fmt"
)

func Test_jsonConn(t *testing.T) {
	c := NewHttpClient("http://10.0.1.2:80")
	
	var lights map[string]map[string]string
	err := c.Get("/api/f8946c33ae3512f1abeefbb23bf5ca7/lights", &lights)
	if err != nil {
		t.Error(err)
	}
	
	m := lights
	for k, v := range m {
		fmt.Printf("%s -> %s\n", k, v["name"])
	}
}
