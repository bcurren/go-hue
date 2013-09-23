package hue

import (
	"testing"
)

func Test_ApiParseErrorString(t *testing.T) {
	err := NewApiError("string", 1, "user count")
	assertEqual(t, "Parsing error: expected type 'string' but received 'int' for user count.",
		err.Error(), "err message")
}
