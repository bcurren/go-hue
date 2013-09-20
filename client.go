package hue

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"fmt"
	"io/ioutil"
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
	httpRequest, err := http.NewRequest(method, s.addr + uri, bytes.NewReader(requestBytes))
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
	requestJson string
	responseFile string
}

func (s *stubServer) Do(method string, uri string, requestBytes []byte) ([]byte, error) {
	s.requestJson = string(requestBytes)
	
	
	path := filepath.Join(
		".", 
		"test_responses", 
		strings.ToLower(method), 
		strings.ToLower(strings.Replace(uri, "/api", "", 1)))
	fileBytes, err := ioutil.ReadFile(path + ".json")
	if err != nil {
		return nil, err
	}
	return fileBytes, nil
}

type client struct {
	conn server
}

func NewHttpClient(addr string) *client {
	return &client{conn: &httpServer{addr: addr}}
}

func NewStubClient(responseFile string) *client {
	return &client{conn: &stubServer{responseFile: responseFile}}
}

func (c *client) Get(uri string, responseObj interface{}) (error) {
	return c.Send("GET", uri, nil, responseObj)
}

func (c *client) Send(method string, uri string, requestObj interface{}, responseObj interface{}) (error) {
	// TODO: check if uri starts with /
	
	// Convert object to json
	var requestBytes []byte
	var err error
	if requestObj != nil {
		requestBytes, err = json.Marshal(requestObj)
		if err != nil {
			return err
		}
	}
	
	// Perform http request
	resBytes, err := c.conn.Do(method, uri, requestBytes)
	if err != nil {
		return err
	}
	
	// Parse response json to object
	decoder := json.NewDecoder(bytes.NewReader(resBytes))
	err = decoder.Decode(responseObj)
	if err != nil {
		
		// Parse the error response
		var apiErrorDetails []map[string]*ApiErrorDetail
		decoder = json.NewDecoder(bytes.NewReader(resBytes))
		err = decoder.Decode(&apiErrorDetails)
		if err != nil {
			return err
		}
		
		// Build ApiError structure with slice of errors
		apiError := &ApiError{}
		apiError.Errors = make([]ApiErrorDetail, 0, 1)
		for _, apiErrorDetail := range apiErrorDetails {
			apiError.Errors = append(apiError.Errors, *apiErrorDetail["error"])
		}
		return apiError
	}
	
	return nil
}

type ApiError struct {
	Errors []ApiErrorDetail
}

func (e ApiError) Error() string {
	errors := make([]string, 0, 10)
	for _, error := range e.Errors {
		errors = append(errors, error.Error())
	}
	
	return strings.Join(errors, " ")
}

type ApiErrorDetail struct {
	Type int `json:type`
	Address string `json:address`
	Description string `json:description`
}

func (e ApiErrorDetail) Error() string {
	return fmt.Sprintf("Hue API Error type '%d' with description '%s'.", e.Type, e.Description)
}
