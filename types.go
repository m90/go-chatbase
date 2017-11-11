package chatbase

import (
	"encoding/json"
	"fmt"
)

// Message contains info about a platform agnostic chat message
type Message struct {
	APIKey     string `json:"api_key"`
	Type       string `json:"type"`
	UserID     string `json:"user_id"`
	TimeStamp  int    `json:"time_stamp"`
	Platform   string `json:"platform"`
	Message    string `json:"message,omitempty"`
	Intent     string `json:"intent,omitempty"`
	NotHandled bool   `json:"not_handled,omitempty"`
	Feedback   bool   `json:"feedback,omitempty"`
	Version    string `json:"version,omitempty"`
}

// MessageResponse describes an API response to sending a single message
type MessageResponse struct {
	MessageID string `json:"message_id"`
	StatusOK  Status `json:"status"`
}

// Status represents the success of an operation
type Status bool

// UnmarshalJSON handles both int and string values that are being sent
// in responses by chatbase to represent success or failure
func (s *Status) UnmarshalJSON(b []byte) error {
	var i int
	if err := json.Unmarshal(b, &i); err == nil {
		*s = Status(i < 400)
		return nil
	}
	var str string
	if err := json.Unmarshal(b, &str); err == nil {
		*s = Status(str == "success")
		return nil
	}
	return fmt.Errorf("could not unmarshal %s into Status", b)
}

// OK returns the boolean representation of Status
func (s *Status) OK() bool {
	return bool(*s)
}

// Messages is a collection of Message
type Messages []Message

// MessagesResponse describes an API response to sending multiple messages at once
type MessagesResponse struct {
	AllSucceded bool              `json:"all_succeeded"`
	Status      Status            `json:"status"`
	Responses   []MessageResponse `json:"responses"`
}

// Update contains data to be updated about a message that already exists
type Update struct {
	Intent     string `json:"intent,omitempty"`
	NotHandled string `json:"not_handled,omitempty"`
	Feedback   string `json:"feedback,omitempty"`
	Version    string `json:"version,omitempty"`
}

// UpdateResponse describes a service response to an update
type UpdateResponse struct {
	Error   []string `json:"error"`
	Updated []string `json:"updated"`
	Status  Status   `json:"status"`
}

// Event describes an event to be sent to the events API
type Event struct {
	APIKey     string          `json:"api_key"`
	UserID     string          `json:"user_id"`
	Intent     string          `json:"intent"`
	TimeStamp  int             `json:"timestamp_millis,omitempty"`
	Platform   string          `json:"platform,omitempty"`
	Version    string          `json:"version,omitempty"`
	Properties []EventProperty `json:"properties"`
}

// Events is a collection of Event
type Events []Event

// EventProperty is a property that is attached to an event
type EventProperty struct {
	Name         string  `json:"property_name"`
	StringValue  string  `json:"string_value,omitempty"`
	IntegerValue int     `json:"integer_value,omitempty"`
	FloatValue   float64 `json:"float_value,omitempty"`
	BoolValue    bool    `json:"bool_value,omitempty"`
}

// NewEventProperty generates an EventProperty containing the correctly
// typed field for the passed value
func NewEventProperty(name string, value interface{}) (EventProperty, error) {
	p := EventProperty{Name: name}
	if s, ok := value.(string); ok {
		p.StringValue = s
		return p, nil
	}
	if i, ok := value.(int); ok {
		p.IntegerValue = i
		return p, nil
	}
	if f, ok := value.(float64); ok {
		p.FloatValue = f
		return p, nil
	}
	if b, ok := value.(bool); ok {
		p.BoolValue = b
		return p, nil
	}
	return p, fmt.Errorf("could not use %v as event property value", value)
}
