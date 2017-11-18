package chatbase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func api(method, endpoint string, v interface{}) (io.ReadCloser, error) {
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
	return api(http.MethodPost, endpoint, v)
}

func apiPut(endpoint string, v interface{}) (io.ReadCloser, error) {
	return api(http.MethodPut, endpoint, v)
}

func newMessageResponse(thunk func() (io.ReadCloser, error)) (*MessageResponse, error) {
	body, err := thunk()
	if err != nil {
		return nil, err
	}
	defer body.Close()
	responseData := MessageResponse{}
	if err := json.NewDecoder(body).Decode(&responseData); err != nil {
		return nil, err
	}
	return &responseData, nil
}

func newMessagesResponse(thunk func() (io.ReadCloser, error)) (*MessagesResponse, error) {
	body, err := thunk()
	if err != nil {
		return nil, err
	}
	defer body.Close()
	responseData := MessagesResponse{}
	if err := json.NewDecoder(body).Decode(&responseData); err != nil {
		return nil, err
	}
	return &responseData, nil
}
