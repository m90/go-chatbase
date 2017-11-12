package chatbase

import "encoding/json"
import "net/url"
import "net/http"
import "bytes"

var (
	facebookMessageEndpoint = "https://chatbase.com/api/facebook/message_received"
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
	payload, payloadErr := json.Marshal(*f)
	if payloadErr != nil {
		return nil, payloadErr
	}
	u, uErr := url.Parse(facebookMessageEndpoint)
	if uErr != nil {
		return nil, uErr
	}
	u.RawQuery = url.Values{"api_key": []string{f.APIKey}}.Encode()
	res, err := http.Post(u.String(), "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	responseData := MessageResponse{}
	if err := json.NewDecoder(res.Body).Decode(&responseData); err != nil {
		return nil, err
	}
	return &responseData, nil
}
