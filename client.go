package hue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type client struct {
	conn server
}

func NewHttpClient(addr string) *client {
	return &client{conn: &httpServer{addr: addr}}
}

func NewStubClient(responseFile string) *client {
	return &client{conn: &stubServer{responseFile: responseFile}}
}

func (c *client) Get(uri string, responseObj interface{}) error {
	return c.Send("GET", uri, nil, responseObj)
}

func (c *client) Post(uri string, requestObj interface{}, responseObj interface{}) error {
	return c.Send("POST", uri, requestObj, responseObj)
}

func (c *client) Put(uri string, requestObj interface{}, responseObj interface{}) error {
	return c.Send("PUT", uri, requestObj, responseObj)
}

func (c *client) Send(method string, uri string, requestObj interface{}, responseObj interface{}) error {
	requestBytes, err := encode(requestObj)
	if err != nil {
		return err
	}

	resultBytes, err := c.conn.Do(method, uri, requestBytes)
	if err != nil {
		return err
	}

	// If response object is nil, use default response to ensure response is not an error
	if responseObj == nil {
		var defaultResponse []map[string]map[string]string
		responseObj = &defaultResponse
	}

	err = decode(resultBytes, responseObj)
	if err != nil {
		return decodeApiError(resultBytes)
	}

	return nil
}

func encode(requestObj interface{}) ([]byte, error) {
	var requestBytes []byte
	var err error

	if requestObj != nil {
		requestBytes, err = json.Marshal(requestObj)
		if err != nil {
			return nil, err
		}
	}

	return requestBytes, nil
}

func decode(resultBytes []byte, responseObj interface{}) error {
	decoder := json.NewDecoder(bytes.NewReader(resultBytes))
	return decoder.Decode(responseObj)
}

func decodeApiError(resultBytes []byte) error {
	// Parse the error response
	var apiErrorDetails []map[string]*ApiErrorDetail
	err := decode(resultBytes, &apiErrorDetails)
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
	Type        int    `json:"type"`
	Address     string `json:"address"`
	Description string `json:"description"`
}

func (e ApiErrorDetail) Error() string {
	return fmt.Sprintf("Hue API Error type '%d' with description '%s'.", e.Type, e.Description)
}
