package hue

import (
	"testing"
)

func NewStubUser(stubFilename string, username string) (*User, *stubServer) {
	stubClient := NewStubClient(stubFilename)
	stubServer, ok := stubClient.conn.(*stubServer)
	if !ok {
		panic("Not using a stub server in tests!")
	}
	bridge := &Bridge{Name: "StubBridge", client: stubClient}

	return &User{Bridge: bridge, Username: username}, stubServer
}

func Test_ApiParseErrorString(t *testing.T) {
	err := NewApiError("string", 1, "user count")
	assertEqual(t, "Parsing error: expected type 'string' but received 'int' for user count.",
		err.Error(), "err message")
}
