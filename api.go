package chatbase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

var (
	client = http.Client{}
)

// SetAPITransport allows setting a transport for the http.Client that is being used
// for handling Chatbase API calls
func SetAPITransport(t http.RoundTripper) {
	client.Transport = t
}

// SetAPITimeout allows setting a timeout value for the http.Client that is being used
// for handling Chatbase API calls
func SetAPITimeout(t time.Duration) {
	client.Timeout = t
}

func apiCall(method, endpoint string, v interface{}) (io.ReadCloser, error) {
	payload, payloadErr := json.Marshal(v)
	if payloadErr != nil {
		return nil, payloadErr
	}

	req, reqErr := http.NewRequest(method, endpoint, bytes.NewBuffer(payload))
	if reqErr != nil {
		return nil, reqErr
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= http.StatusInternalServerError {
		return nil, fmt.Errorf("request failed with status %v", res.StatusCode)
	}
	return res.Body, nil
}

func apiPost(endpoint string, v interface{}) (io.ReadCloser, error) {
	return apiCall(http.MethodPost, endpoint, v)
}

func apiPut(endpoint string, v interface{}) (io.ReadCloser, error) {
	return apiCall(http.MethodPut, endpoint, v)
}

func newMessageResponse(thunk func() (io.ReadCloser, error)) (*MessageResponse, error) {
	data, err := decodeInto(&MessageResponse{}, thunk)
	if res, ok := data.(*MessageResponse); ok {
		return res, err
	}
	if err == nil {
		err = fmt.Errorf("%v was not of expected type", data)
	}
	return nil, err
}

func newMessagesResponse(thunk func() (io.ReadCloser, error)) (*MessagesResponse, error) {
	data, err := decodeInto(&MessagesResponse{}, thunk)
	if res, ok := data.(*MessagesResponse); ok {
		return res, err
	}
	if err == nil {
		err = fmt.Errorf("%v was not of expected type", data)
	}
	return nil, err
}

func newLinkResponse(thunk func() (io.ReadCloser, error)) (*LinkResponse, error) {
	data, err := decodeInto(&LinkResponse{}, thunk)
	if res, ok := data.(*LinkResponse); ok {
		return res, err
	}
	if err == nil {
		err = fmt.Errorf("%v was not of expected type", data)
	}
	return nil, err
}

func newUpdateResponse(thunk func() (io.ReadCloser, error)) (*UpdateResponse, error) {
	data, err := decodeInto(&UpdateResponse{}, thunk)
	if res, ok := data.(*UpdateResponse); ok {
		return res, err
	}
	if err == nil {
		err = fmt.Errorf("%v was not of expected type", data)
	}
	return nil, err
}

func decodeInto(target interface{}, thunk func() (io.ReadCloser, error)) (interface{}, error) {
	body, err := thunk()
	if err != nil {
		return nil, err
	}
	if body == nil {
		return nil, nil
	}
	defer body.Close()
	if err := json.NewDecoder(body).Decode(target); err != nil {
		return nil, err
	}
	return target, nil
}

func augmentURL(endpoint string, params map[string]string) (string, error) {
	e, endpointErr := url.Parse(endpoint)
	if endpointErr != nil {
		return "", endpointErr
	}
	q := e.Query()
	for key, value := range params {
		q.Set(key, value)
	}
	e.RawQuery = q.Encode()
	return e.String(), nil
}
