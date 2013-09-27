package hue

import (
	"fmt"
)

type User struct {
	Bridge   *Bridge
	Username string
}

func NewUser(username, bridgeId, addr string) *User {
	return NewUserWithBridge(username, NewBridge(bridgeId, addr))
}

func NewUserWithBridge(username string, bridge *Bridge) *User {
	return &User{Username: username, Bridge: bridge}
}

type APIParseError struct {
	Expected string
	Actual   interface{}
	Context  string
}

func NewAPIParseError(expected string, actual interface{}, context string) error {
	return &APIParseError{Expected: expected, Actual: actual, Context: context}
}

func (e *APIParseError) Error() string {
	return fmt.Sprintf("Parsing error: expected type '%s' but received '%T' for %s.",
		e.Expected, e.Actual, e.Context)
}
