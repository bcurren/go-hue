package hue

import (
	"bytes"
	"encoding/json"
)

type client struct {
	conn server
}

func NewHTTPClient(addr string) *client {
	return &client{conn: &httpServer{addr: addr}}
}

func (c *client) Get(uri string, responseObj interface{}) error {
	resultBytes, err := c.conn.Do("GET", uri, nil)
	if err != nil {
		return err
	}

	err = decode(resultBytes, responseObj)
	if err != nil {
		_, errors, err := decodeAPIResult(resultBytes)
		if err != nil {
			return err
		}
		return &APIError{errors}
	}

	return nil
}

func (c *client) Post(uri string, requestObj interface{}) ([]APISuccessDetail, error) {
	return c.Send("POST", uri, requestObj)
}

func (c *client) Put(uri string, requestObj interface{}) ([]APISuccessDetail, error) {
	return c.Send("PUT", uri, requestObj)
}

func (c *client) Send(method string, uri string, requestObj interface{}) ([]APISuccessDetail, error) {
	requestBytes, err := encode(requestObj)
	if err != nil {
		return nil, err
	}

	resultBytes, err := c.conn.Do(method, uri, requestBytes)
	if err != nil {
		return nil, err
	}

	successes, errors, err := decodeAPIResult(resultBytes)
	if err != nil {
		return nil, err
	}
	if len(errors) != 0 {
		return successes, &APIError{errors}
	}

	return successes, nil
}

func decodeAPIResult(resultBytes []byte) ([]APISuccessDetail, []APIErrorDetail, error) {
	// Decode api result
	var responseObj []*APIResult
	err := decode(resultBytes, &responseObj)
	if err != nil {
		return nil, nil, err
	}

	// Create slice of successes and errors
	successes := make([]APISuccessDetail, 0, 5)
	errors := make([]APIErrorDetail, 0, 5)
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

func decodeAPIError(resultBytes []byte) error {
	// Parse the error response
	var apiErrorDetails []map[string]*APIErrorDetail
	err := decode(resultBytes, &apiErrorDetails)
	if err != nil {
		return err
	}

	// Build APIError structure with slice of errors
	apiError := &APIError{}
	apiError.Errors = make([]APIErrorDetail, 0, 1)
	for _, apiErrorDetail := range apiErrorDetails {
		apiError.Errors = append(apiError.Errors, *apiErrorDetail["error"])
	}

	return apiError
}

type APIResult struct {
	Success *APISuccessDetail `json:"success"`
	Error   *APIErrorDetail   `json:"error"`
}

type APISuccessDetail map[string]interface{}
