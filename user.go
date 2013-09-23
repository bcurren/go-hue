package hue

import (
	"fmt"
)

type User struct {
	Bridge   *Bridge
	Username string
}

type ApiParseError struct {
	Expected string
	Actual   interface{}
	Context  string
}

func NewApiError(expected string, actual interface{}, context string) error {
	return &ApiParseError{Expected: expected, Actual: actual, Context: context}
}

func (e *ApiParseError) Error() string {
	return fmt.Sprintf("Parsing error: expected type '%s' but received '%T' for %s.",
		e.Expected, e.Actual, e.Context)
}
