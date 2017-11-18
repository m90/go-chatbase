package chatbase

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestFacebookMessage_Setters(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		payload := map[string]string{
			"hello": "world",
		}
		m := FacebookMessage{
			Payload: payload,
		}
		expected := FacebookMessage{
			Payload: payload,
			Fields: &FacebookFields{
				Intent:     "test-things",
				NotHandled: true,
				Feedback:   true,
				Version:    "1.3.1",
			},
		}
		m.SetIntent("test-things").SetNotHandled(true).SetFeedback(true).SetVersion("1.3.1")
		if !reflect.DeepEqual(expected, m) {
			t.Errorf("Expected %#v, got %#v", expected, m)
		}
	})
	t.Run("nil fields", func(t *testing.T) {
		fixture := FacebookMessage{}

		m := fixture
		m.SetIntent("test")
		if m.Fields.Intent != "test" {
			t.Errorf("Unexpected fields %v", m.Fields)
		}

		m = fixture
		m.SetNotHandled(true)
		if m.Fields.NotHandled != true {
			t.Errorf("Unexpected fields %v", m.Fields)
		}

		m = fixture
		m.SetFeedback(true)
		if m.Fields.Feedback != true {
			t.Errorf("Unexpected fields %v", m.Fields)
		}

		m = fixture
		m.SetVersion("1.2.3")
		if m.Fields.Version != "1.2.3" {
			t.Errorf("Unexpected fields %v", m.Fields)
		}
	})
}
func TestFacebookRequestResponse_Setters(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		m := FacebookRequestResponse{
			Request:  "hello",
			Response: "goodbye",
		}
		expected := FacebookRequestResponse{
			Request:  "hello",
			Response: "goodbye",
			Fields: &FacebookFields{
				Intent:     "test-things",
				NotHandled: true,
				Feedback:   true,
				Version:    "1.3.1",
			},
		}
		m.SetIntent("test-things").SetNotHandled(true).SetFeedback(true).SetVersion("1.3.1")
		if !reflect.DeepEqual(expected, m) {
			t.Errorf("Expected %#v, got %#v", expected, m)
		}
	})
	t.Run("nil fields", func(t *testing.T) {
		fixture := FacebookRequestResponse{}

		m := fixture
		m.SetIntent("test")
		if m.Fields.Intent != "test" {
			t.Errorf("Unexpected fields %v", m.Fields)
		}

		m = fixture
		m.SetNotHandled(true)
		if m.Fields.NotHandled != true {
			t.Errorf("Unexpected fields %v", m.Fields)
		}

		m = fixture
		m.SetFeedback(true)
		if m.Fields.Feedback != true {
			t.Errorf("Unexpected fields %v", m.Fields)
		}

		m = fixture
		m.SetVersion("1.2.3")
		if m.Fields.Version != "1.2.3" {
			t.Errorf("Unexpected fields %v", m.Fields)
		}
	})
}

func TestFacebookMessage_MarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		input       FacebookMessage
		expected    string
		expectError bool
	}{
		{
			"default",
			FacebookMessage{
				Payload: map[string]interface{}{
					"hello": "world",
					"foo":   "bar",
					"nested": map[string]bool{
						"true":  true,
						"false": false,
					},
				},
				Fields: &FacebookFields{
					Intent:     "test-things",
					NotHandled: true,
					Version:    "1.4.4",
				},
			},
			`{"chatbase_fields":{"intent":"test-things","not_handled":true,"version":"1.4.4"},"foo":"bar","hello":"world","nested":{"false":false,"true":true}}`,
			false,
		},
		{
			"pass through",
			FacebookMessage{
				Payload: map[string]string{
					"hello": "world",
					"foo":   "bar",
				},
			},
			`{"foo":"bar","hello":"world"}`,
			false,
		},
		{
			"bad payload",
			FacebookMessage{
				Payload: "bingo bongo",
			},
			"",
			true,
		},
		{
			"invalid payload",
			FacebookMessage{
				Payload: func() {},
			},
			"",
			true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b, err := json.Marshal(test.input)
			if test.expectError != (err != nil) {
				t.Errorf("Unexpected error %v", err)
			}
			if string(b) != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, string(b))
			}
		})
	}
}

func TestMarshalJSON_FacebookRequestResponses(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		f := FacebookRequestResponses{}
		b, err := json.Marshal(f)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}
		if s := string(b); s != `{"messages":[]}` {
			t.Errorf("Unexpected result %v", s)
		}
	})
}

func TestMarshalJSON_FacebookMessages(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		f := FacebookMessages{}
		b, err := json.Marshal(f)
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}
		if s := string(b); s != `{"messages":[]}` {
			t.Errorf("Unexpected result %v", s)
		}
	})
}

func TestAppend_FacebookMessages(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		m := FacebookMessages{}
		m.Append(&FacebookMessage{}, &FacebookMessage{})
		if l := len(m); l != 2 {
			t.Errorf("Unexpected length %v", l)
		}
	})
}
func TestAppend_FacebookRequestResponses(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		m := FacebookRequestResponses{}
		m.Append(&FacebookRequestResponse{}, &FacebookRequestResponse{})
		if l := len(m); l != 2 {
			t.Errorf("Unexpected length %v", l)
		}
	})
}
