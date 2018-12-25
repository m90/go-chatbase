package chatbase

import (
	"context"
	"encoding/json"
	"fmt"
)

var (
	eventEndpoint  = "https://api.chatbase.com/apis/v1/events/insert"
	eventsEndpoint = "https://api.chatbase.com/apis/v1/events/insert_batch"
)

// Event contains data about an event
type Event struct {
	APIKey     string          `json:"api_key"`
	UserID     string          `json:"user_id"`
	Intent     string          `json:"intent"`
	TimeStamp  int             `json:"timestamp_millis,omitempty"`
	Platform   string          `json:"platform,omitempty"`
	Version    string          `json:"version,omitempty"`
	Properties []EventProperty `json:"properties"`
}

// SetTimeStamp adds an optional "timestamp" value to the event
func (e *Event) SetTimeStamp(t int) *Event {
	e.TimeStamp = t
	return e
}

// SetPlatform adds an optional "platform" value to the event
func (e *Event) SetPlatform(p string) *Event {
	e.Platform = p
	return e
}

// SetVersion adds an optional "version" value to the event
func (e *Event) SetVersion(v string) *Event {
	e.Version = v
	return e
}

// AddProperty adds a new property to the event using the given name and value.
// The passed value needs to be one of `int`, `string`, `bool` or `float64`
func (e *Event) AddProperty(name string, v interface{}) error {
	prop, err := NewEventProperty(name, v)
	if err != nil {
		return err
	}
	e.Properties = append(e.Properties, prop)
	return nil
}

// Submit tries to deliver the event to Chatbase
func (e *Event) Submit() error {
	_, err := apiPost(eventEndpoint, e)
	return err
}

// SubmitWithContext tries to deliver the event to Chatbase
// while considering the given context's deadline
func (e *Event) SubmitWithContext(ctx context.Context) error {
	return withContext(ctx, e.Submit)
}

// Events is a collection of Event
type Events []Event

// MarshalJSON ensure the collection is correctly wrapped
// into an object and added the api_key value
func (e Events) MarshalJSON() ([]byte, error) {
	var apiKey string
	if len(e) > 0 {
		apiKey = e[0].APIKey
	}
	return json.Marshal(map[string]interface{}{
		"api_key": apiKey,
		"events":  []Event(e),
	})
}

// Submit tries to deliver the set of events to Chatbase
func (e *Events) Submit() error {
	_, err := apiPost(eventsEndpoint, e)
	return err
}

// SubmitWithContext tries to deliver the set of events to Chatbase
// while considering the context's deadline
func (e *Events) SubmitWithContext(ctx context.Context) error {
	return withContext(ctx, e.Submit)
}

// Append adds events to the the collection. The collection should not
// contain events using different API keys
func (e *Events) Append(addition ...*Event) *Events {
	for _, a := range addition {
		*e = append(*e, *a)
	}
	return e
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
// typed field for the passed value. The passed value needs to be one of
// `int`, `string`, `bool` or `float64`
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
