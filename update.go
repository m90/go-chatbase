package chatbase

import (
	"encoding/json"
	"fmt"
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

// SetIntent adds an optional intent value to update
func (u *Update) SetIntent(i string) *Update {
	u.Intent = i
	return u
}

// SetNotHandled adds an optional not handled value to update
func (u *Update) SetNotHandled(n bool) *Update {
	u.NotHandled = ""
	if n {
		u.NotHandled = "true"
	}
	return u
}

// SetFeedback adds an optional feedback value to update
func (u *Update) SetFeedback(f string) *Update {
	u.Feedback = f
	return u
}

// SetVersion adds an optional version value to update
func (u *Update) SetVersion(v string) *Update {
	u.Version = v
	return u
}

// Submit tries to deliver the update to Chatbase
func (u *Update) Submit() (*UpdateResponse, error) {
	e, endpointErr := url.Parse(updateEndpoint)
	if endpointErr != nil {
		return nil, endpointErr
	}
	q := e.Query()
	q.Set("api_key", u.APIKey)
	q.Set("message_id", u.MessageID)
	e.RawQuery = q.Encode()

	body, bodyErr := apiPut(e.String(), u)
	if bodyErr != nil {
		return nil, bodyErr
	}
	defer body.Close()

	responseData := UpdateResponse{}
	if err := json.NewDecoder(body).Decode(&responseData); err != nil {
		return nil, err
	}

	if !responseData.Status.OK() {
		return &responseData, fmt.Errorf("failed sending messages with status %v", responseData.Status)
	}
	return &responseData, nil
}

// UpdateResponse describes a Chatbase response to an update
type UpdateResponse struct {
	Error   []string `json:"error"`
	Updated []string `json:"updated"`
	Status  Status   `json:"status"`
	Reason  string   `json:"string"`
}
