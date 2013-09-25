package hue

import (
	"testing"
)

func Test_NewUser(t *testing.T) {
	user := NewUser("myUsername", "uuid:454-454", "10.4.5.2")
	httpServer, ok := user.Bridge.client.conn.(*httpServer)
	if !ok {
		t.Fatal("User doesn't have an httpServer set up.")
	}

	assertEqual(t, "myUsername", user.Username, "user.Username")
	assertEqual(t, "10.4.5.2", httpServer.addr, "httpServer.addr")
}

func Test_ApiParseErrorString(t *testing.T) {
	err := NewApiError("string", 1, "user count")
	assertEqual(t, "Parsing error: expected type 'string' but received 'int' for user count.",
		err.Error(), "err message")
}
