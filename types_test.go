package chatbase

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestUnmarshalJSON_Status(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		expected    Status
	}{
		{
			"success int",
			"200",
			false,
			true,
		},
		{
			"error int",
			"400",
			false,
			false,
		},
		{
			"success string",
			`"success"`,
			false,
			true,
		},
		{
			"error string",
			`"failure"`,
			false,
			false,
		},
		{
			"bad value",
			"true",
			true,
			false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var v Status
			err := json.Unmarshal([]byte(test.input), &v)
			if test.expectError != (err != nil) {
				t.Errorf("Expected %v, got %v", test.expectError, err)
			}
			if !test.expectError {
				if v != test.expected {
					t.Errorf("Expected %v, got %v", test.expected, v)
				}
			}
		})
	}
}

func TestOK_Status(t *testing.T) {
	tests := []struct {
		name     string
		status   Status
		expected bool
	}{
		{
			"true",
			true,
			true,
		},
		{
			"false",
			false,
			false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.status.OK() != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, test.status.OK())
			}
		})
	}
}

func TestNewEventProperty(t *testing.T) {
	tests := []struct {
		name        string
		value       interface{}
		expectError bool
		expected    EventProperty
	}{
		{
			"string",
			"foo bar",
			false,
			EventProperty{
				Name:        "string",
				StringValue: "foo bar",
			},
		},
		{
			"int",
			89,
			false,
			EventProperty{
				Name:         "int",
				IntegerValue: 89,
			},
		},
		{
			"float",
			1.2345,
			false,
			EventProperty{
				Name:       "float",
				FloatValue: 1.2345,
			},
		},
		{
			"bool",
			true,
			false,
			EventProperty{
				Name:      "bool",
				BoolValue: true,
			},
		},
		{
			"bad value",
			[]int{2, 9},
			true,
			EventProperty{
				Name: "bad value",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := NewEventProperty(test.name, test.value)
			if test.expectError != (err != nil) {
				t.Errorf("Expected error: %v", test.expectError)
			}
			if !reflect.DeepEqual(test.expected, result) {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestMessage_Setters(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		m := Message{
			APIKey:    "secret",
			Type:      MessageTypeAgent,
			UserID:    "abc-123",
			TimeStamp: 998877,
			Platform:  "test",
		}
		expected := Message{
			APIKey:     "secret",
			Type:       MessageTypeAgent,
			UserID:     "abc-123",
			TimeStamp:  998877,
			Platform:   "test",
			Message:    "Hello world!",
			Intent:     "test-things",
			NotHandled: true,
			Feedback:   true,
			Version:    "1.2.34",
		}
		m.SetMessage("Hello world!").SetIntent("test-things").SetNotHandled(true).SetFeedback(true).SetVersion("1.2.34")
		if !reflect.DeepEqual(expected, m) {
			t.Errorf("Expected %#v, got %#v", expected, m)
		}
	})
}
