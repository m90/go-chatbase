package chatbase

import (
	"context"
	"io"
)

var (
	updateEndpoint = "https://chatbase.com/api/message/update"
)

// Update contains data about an existing message that should be updated
type Update struct {
	APIKey     string    `json:"-"`
	MessageID  MessageID `json:"-"`
	Intent     string    `json:"intent,omitempty"`
	NotHandled string    `json:"not_handled,omitempty"`
	Feedback   string    `json:"feedback,omitempty"`
	Version    string    `json:"version,omitempty"`
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
		ep, epErr := augmentURL(updateEndpoint, map[string]string{
			"api_key":    u.APIKey,
			"message_id": u.MessageID.String(),
		})
		if epErr != nil {
			return nil, epErr
		}
		body, bodyErr := apiPut(ep, u)
		if bodyErr != nil {
			return nil, bodyErr
		}
		return body, nil
	})
}

// SubmitWithContext tries to deliver the update to Chatbase while
// considering the given context's deadline
func (u *Update) SubmitWithContext(ctx context.Context) (*UpdateResponse, error) {
	data, err := resultWithContext(ctx, func() (interface{}, error) {
		return u.Submit()
	})
	if err != nil {
		return nil, err
	}
	if res, ok := data.(*UpdateResponse); ok {
		return res, nil
	}
	return nil, errBadData
}

// UpdateResponse describes a Chatbase response to an update submission
type UpdateResponse struct {
	Error   []string `json:"error"`
	Updated []string `json:"updated"`
	Status  Status   `json:"status"`
	Reason  string   `json:"string,omitempty"`
}
