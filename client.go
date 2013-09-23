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

func (c *client) Get(uri string, responseObj interface{}) error {
	resultBytes, err := c.conn.Do("GET", uri, nil)
	if err != nil {
		return err
	}

	err = decode(resultBytes, responseObj)
	if err != nil {
		_, errors, err := decodeApiResult(resultBytes)
		if err != nil {
			return err
		}
		return &ApiError{errors}
	}

	return nil
}

func (c *client) Post(uri string, requestObj interface{}) ([]ApiSuccessDetail, error) {
	return c.Send("POST", uri, requestObj)
}

func (c *client) Put(uri string, requestObj interface{}) ([]ApiSuccessDetail, error) {
	return c.Send("PUT", uri, requestObj)
}

func (c *client) Send(method string, uri string, requestObj interface{}) ([]ApiSuccessDetail, error) {
	requestBytes, err := encode(requestObj)
	if err != nil {
		return nil, err
	}

	resultBytes, err := c.conn.Do(method, uri, requestBytes)
	if err != nil {
		return nil, err
	}

	successes, errors, err := decodeApiResult(resultBytes)
	if err != nil {
		return nil, err
	}
	if len(errors) != 0 {
		return successes, &ApiError{errors}
	}

	return successes, nil
}

func decodeApiResult(resultBytes []byte) ([]ApiSuccessDetail, []ApiErrorDetail, error) {
	// Decode api result
	var responseObj []*ApiResult
	err := decode(resultBytes, &responseObj)
	if err != nil {
		return nil, nil, err
	}

	// Create slice of successes and errors
	successes := make([]ApiSuccessDetail, 0, 5)
	errors := make([]ApiErrorDetail, 0, 5)
	for _, result := range responseObj {
		if result.Success != nil {
			successes = append(successes, *result.Success)
		}
		if result.Error != nil {
			errors = append(errors, *result.Error)
		}
	}

	// Set slice to nil if empty
	if len(successes) == 0 {
		successes = nil
	}
	if len(errors) == 0 {
		errors = nil
	}

	return successes, errors, nil
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

type ApiResult struct {
	Success *ApiSuccessDetail `json:"success"`
	Error   *ApiErrorDetail   `json:"error"`
}

type ApiSuccessDetail map[string]interface{}
