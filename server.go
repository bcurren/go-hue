package hue

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"
	"strings"
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

	url := cleanURL(fmt.Sprintf("%s/%s", s.addr, uri))

	httpRequest, err := http.NewRequest(method, url, bytes.NewReader(requestBytes))
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

func cleanURL(url string) string {
	doubleSlash := regexp.MustCompile(`/+`)
	url = doubleSlash.ReplaceAllString(url, "/")
	url = strings.Replace(url, "http:/", "http://", -1)
	url = strings.Replace(url, "https:/", "https://", -1)
	return url
}
