package chatbase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	updateEndpoint = "https://chatbase.com/api/message/update"
)

// Update contains data to be updated about a message that already exists
type Update struct {
	Intent     string `json:"intent,omitempty"`
	NotHandled string `json:"not_handled,omitempty"`
	Feedback   string `json:"feedback,omitempty"`
	Version    string `json:"version,omitempty"`
}

// Submit tries to deliver the update to chatbase
func (u *Update) Submit() (*UpdateResponse, error) {
	payload, payloadErr := json.Marshal(u)
	if payloadErr != nil {
		return nil, payloadErr
	}

	req, reqErr := http.NewRequest(http.MethodPut, updateEndpoint, bytes.NewBuffer(payload))
	if reqErr != nil {
		return nil, reqErr
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	responseData := UpdateResponse{}
	if err := json.NewDecoder(res.Body).Decode(&responseData); err != nil {
		return nil, err
	}

	if !responseData.Status.OK() {
		return nil, fmt.Errorf("failed sending messages with status %v", responseData.Status)
	}
	return &responseData, nil
}

// UpdateResponse describes a service response to an update
type UpdateResponse struct {
	Error   []string `json:"error"`
	Updated []string `json:"updated"`
	Status  Status   `json:"status"`
}
