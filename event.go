package chatbase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	eventEndpoint  = "https://api.chatbase.com/apis/v1/events/insert"
	eventsEndpoint = "https://api.chatbase.com/apis/v1/events/insert_batch"
)

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

// SetTimeStamp adds an optional timestamp value
func (e *Event) SetTimeStamp(t int) *Event {
	e.TimeStamp = t
	return e
}

// SetPlatform adds an optional platform value
func (e *Event) SetPlatform(p string) *Event {
	e.Platform = p
	return e
}

// SetVersion adds an optional version value
func (e *Event) SetVersion(v string) *Event {
	e.Version = v
	return e
}

// AddProperty adds a new property to the event using the given name and value
func (e *Event) AddProperty(name string, value interface{}) *Event {
	prop, _ := NewEventProperty(name, value)
	e.Properties = append(e.Properties, prop)
	return e
}

func postEvent(endpoint string, v interface{}) error {
	payload, payloadErr := json.Marshal(v)
	if payloadErr != nil {
		return payloadErr
	}

	res, err := http.Post(endpoint, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("request failed with status %v", res.StatusCode)
	}
	return nil
}

// Submit tries to deliver the set of events to chatbase
func (e *Event) Submit() error {
	return postEvent(eventEndpoint, e)
}

// Events is a collection of Event
type Events []Event

// Submit tries to deliver the set of events to chatbase
func (e *Events) Submit() error {
	return postEvent(eventsEndpoint, e)
}

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
