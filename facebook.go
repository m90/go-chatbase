package chatbase

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

var (
	facebookMessageEndpoint  = "https://chatbase.com/api/facebook/message_received"
	facebookMessagesEndpoint = "https://chatbase.com/api/facebook/message_received_batch"
)

// FacebookFields contains metadata about a native Facebook message
type FacebookFields struct {
	Intent     string `json:"intent"`
	NotHandled bool   `json:"not_handled"`
	Feedback   bool   `json:"feedback"`
	Version    string `json:"string"`
}

// FacebookMessage is a single native Facebook message
type FacebookMessage struct {
	Fields  *FacebookFields
	Payload interface{}
	APIKey  string
}

// MarshalJSON ensures the message is encoded in the way that
// Chatbase expects
func (f FacebookMessage) MarshalJSON() ([]byte, error) {
	intermediate, intermediateErr := json.Marshal(f.Payload)
	if intermediateErr != nil {
		return nil, intermediateErr
	}
	m := make(map[string]interface{})
	if err := json.Unmarshal(intermediate, &m); err != nil {
		return nil, err
	}
	if f.Fields != nil {
		m["chatbase_fields"] = f.Fields
	}
	return json.Marshal(m)
}

// SetIntent adds an optional intent value to the message
func (f *FacebookMessage) SetIntent(i string) *FacebookMessage {
	if f.Fields == nil {
		f.Fields = &FacebookFields{}
	}
	f.Fields.Intent = i
	return f
}

// SetNotHandled adds an optional not handled value to the message
func (f *FacebookMessage) SetNotHandled(n bool) *FacebookMessage {
	if f.Fields == nil {
		f.Fields = &FacebookFields{}
	}
	f.Fields.NotHandled = n
	return f
}

// SetFeedback adds an optional feedback value to the message
func (f *FacebookMessage) SetFeedback(n bool) *FacebookMessage {
	if f.Fields == nil {
		f.Fields = &FacebookFields{}
	}
	f.Fields.Feedback = n
	return f
}

// SetVersion adds an optional version value to the message
func (f *FacebookMessage) SetVersion(v string) *FacebookMessage {
	if f.Fields == nil {
		f.Fields = &FacebookFields{}
	}
	f.Fields.Version = v
	return f
}

// Submit tries to deliver a single Facebook message to chatbase
func (f *FacebookMessage) Submit() (*MessageResponse, error) {
	body, err := postFacebook(facebookMessageEndpoint, f.APIKey, f)
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

// FacebookMessages is a collection of Facecbook Message
type FacebookMessages []FacebookMessage

// Append adds the additional message to the collection
func (f *FacebookMessages) Append(addition *FacebookMessage) *FacebookMessages {
	*f = append(*f, *addition)
	return f
}

// Submit tries to deliver the set of messages to chatbase
func (f *FacebookMessages) Submit() (*MessagesResponse, error) {
	if len(*f) == 0 {
		return nil, errors.New("cannot submit empty collection")
	}
	apiKey := (*f)[0].APIKey
	body, err := postFacebook(facebookMessagesEndpoint, apiKey, f)
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

func postFacebook(endpoint, apiKey string, v interface{}) (io.ReadCloser, error) {
	payload, payloadErr := json.Marshal(v)
	if payloadErr != nil {
		return nil, payloadErr
	}
	u, uErr := url.Parse(endpoint)
	if uErr != nil {
		return nil, uErr
	}
	u.RawQuery = url.Values{"api_key": []string{apiKey}}.Encode()
	res, err := http.Post(u.String(), "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}
