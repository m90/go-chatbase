package chatbase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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
	res, err := http.DefaultClient.Do(req)
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
	data, err := decodeResponse(&MessageResponse{}, thunk)
	if res, ok := data.(*MessageResponse); ok {
		return res, err
	}
	return nil, fmt.Errorf("%v was not of expected type", data)
}

func newMessagesResponse(thunk func() (io.ReadCloser, error)) (*MessagesResponse, error) {
	data, err := decodeResponse(&MessagesResponse{}, thunk)
	if res, ok := data.(*MessagesResponse); ok {
		return res, err
	}
	return nil, fmt.Errorf("%v was not of expected type", data)
}

func newLinkResponse(thunk func() (io.ReadCloser, error)) (*LinkResponse, error) {
	data, err := decodeResponse(&LinkResponse{}, thunk)
	if res, ok := data.(*LinkResponse); ok {
		return res, err
	}
	return nil, fmt.Errorf("%v was not of expected type", data)
}

func newUpdateResponse(thunk func() (io.ReadCloser, error)) (*UpdateResponse, error) {
	data, err := decodeResponse(&UpdateResponse{}, thunk)
	if res, ok := data.(*UpdateResponse); ok {
		return res, err
	}
	return nil, fmt.Errorf("%v was not of expected type", data)
}

func decodeResponse(target interface{}, thunk func() (io.ReadCloser, error)) (interface{}, error) {
	body, err := thunk()
	if err != nil {
		return nil, err
	}
	defer body.Close()
	if err := json.NewDecoder(body).Decode(target); err != nil {
		return nil, err
	}
	return target, nil
}
