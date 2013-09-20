package hue

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

type server interface {
	Do(method string, uri string, requestBytes []byte) ([]byte, error)
}

type httpServer struct {
	addr string
}

func (s *httpServer) Do(method string, uri string, requestBytes []byte) ([]byte, error) {
	if requestBytes == nil {
		requestBytes = make([]byte, 0, 0)
	}
	httpRequest, err := http.NewRequest(method, s.addr+uri, bytes.NewReader(requestBytes))
	if err != nil {
		return nil, err
	}

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	// Get json body as a string
	bodyBuffer := new(bytes.Buffer)
	defer httpResponse.Body.Close()
	_, err = bodyBuffer.ReadFrom(httpResponse.Body)
	if err != nil {
		return nil, err
	}

	return bodyBuffer.Bytes(), nil
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
