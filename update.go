package chatbase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var (
	updateEndpoint = "https://chatbase.com/api/message/update"
)

// Update contains data to be updated about a message that already exists
type Update struct {
	APIKey     string `json:"-"`
	MessageID  string `json:"-"`
	Intent     string `json:"intent,omitempty"`
	NotHandled string `json:"not_handled,omitempty"`
	Feedback   string `json:"feedback,omitempty"`
	Version    string `json:"version,omitempty"`
}

// SetIntent adds an optional intent value
func (u *Update) SetIntent(i string) *Update {
	u.Intent = i
	return u
}

// SetNotHandled adds an optional not handled value
func (u *Update) SetNotHandled(n string) *Update {
	u.NotHandled = n
	return u
}

// SetFeedback adds an optional feedback value
func (u *Update) SetFeedback(f string) *Update {
	u.Feedback = f
	return u
}

// SetVersion adds an optional version value
func (u *Update) SetVersion(v string) *Update {
	u.Version = v
	return u
}

// Submit tries to deliver the update to chatbase
func (u *Update) Submit() (*UpdateResponse, error) {
	payload, payloadErr := json.Marshal(u)
	if payloadErr != nil {
		return nil, payloadErr
	}

	e, endpointErr := url.Parse(updateEndpoint)
	if endpointErr != nil {
		return nil, endpointErr
	}
	q := e.Query()
	q.Set("api_key", u.APIKey)
	q.Set("message_id", u.MessageID)
	e.RawQuery = q.Encode()

	req, reqErr := http.NewRequest(http.MethodPut, e.String(), bytes.NewBuffer(payload))
	if reqErr != nil {
		return nil, reqErr
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("request failed with status code %v", res.StatusCode)
	}

	responseData := UpdateResponse{}
	if err := json.NewDecoder(res.Body).Decode(&responseData); err != nil {
		return nil, err
	}

	if !responseData.Status.OK() {
		return &responseData, fmt.Errorf("failed sending messages with status %v", responseData.Status)
	}
	return &responseData, nil
}

// UpdateResponse describes a service response to an update
type UpdateResponse struct {
	Error   []string `json:"error"`
	Updated []string `json:"updated"`
	Status  Status   `json:"status"`
}
