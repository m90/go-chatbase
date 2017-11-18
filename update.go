package chatbase

import (
	"io"
	"net/url"
)

var (
	updateEndpoint = "https://chatbase.com/api/message/update"
)

// Update contains data about an existing message that should be updated
type Update struct {
	APIKey     string `json:"-"`
	MessageID  string `json:"-"`
	Intent     string `json:"intent,omitempty"`
	NotHandled string `json:"not_handled,omitempty"`
	Feedback   string `json:"feedback,omitempty"`
	Version    string `json:"version,omitempty"`
}

// SetIntent adds an optional "intent" value to an update
func (u *Update) SetIntent(i string) *Update {
	u.Intent = i
	return u
}

// SetNotHandled adds an optional "not handled" value to an update
func (u *Update) SetNotHandled(n bool) *Update {
	u.NotHandled = ""
	if n {
		u.NotHandled = "true"
	}
	return u
}

// SetFeedback adds an optional "feedback" value to an update
func (u *Update) SetFeedback(f string) *Update {
	u.Feedback = f
	return u
}

// SetVersion adds an optional "version" value to an update
func (u *Update) SetVersion(v string) *Update {
	u.Version = v
	return u
}

// Submit tries to deliver the update to Chatbase
func (u *Update) Submit() (*UpdateResponse, error) {
	return newUpdateResponse(func() (io.ReadCloser, error) {
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
		return body, nil
	})
}

// UpdateResponse describes a Chatbase response to an update submission
type UpdateResponse struct {
	Error   []string `json:"error"`
	Updated []string `json:"updated"`
	Status  Status   `json:"status"`
	Reason  string   `json:"string,omitempty"`
}
