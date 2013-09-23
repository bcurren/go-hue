package hue

import (
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"
)

func NewStubUser(stubFilename string, username string) (*User, *stubServer) {
	bridge, stubServer := NewStubBridge(stubFilename)
	return &User{Bridge: bridge, Username: username}, stubServer
}

func NewStubBridge(stubFilename string) (*Bridge, *stubServer) {
	stubClient := NewStubClient(stubFilename)
	stubServer, ok := stubClient.conn.(*stubServer)
	if !ok {
		panic("Not using a stub server in tests!")
	}
	return &Bridge{client: stubClient}, stubServer
}

func NewStubClient(responseFile string) *client {
	return &client{conn: &stubServer{responseFile: responseFile}}
}

type stubServer struct {
	requestJson  string
	uri          string
	method       string
	responseFile string
}

func (s *stubServer) Do(method string, uri string, requestBytes []byte) ([]byte, error) {
	s.requestJson = string(requestBytes)
	s.uri = uri
	s.method = method

	path := filepath.Join(".", "test_responses", s.responseFile)

	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return fileBytes, nil
}

func assertEqual(t *testing.T, expected interface{}, actual interface{}, errorMessage string) {
	if reflect.TypeOf(expected) != reflect.TypeOf(actual) {
		t.Errorf("Received 'expected' of type %T and 'actual' of type %T. %q", expected, actual, errorMessage)
		return
	}
	if expected != actual {
		t.Errorf("%q is not equal to %q. %q", expected, actual, errorMessage)
	}
}

func assertNotNil(t *testing.T, obj interface{}, errorMessage string) {
	if obj == nil {
		t.Errorf("%q should not be nil. %q", obj, errorMessage)
	}
}
