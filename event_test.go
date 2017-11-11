package chatbase

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

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

func TestEvent_Setters(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		e := Event{
			APIKey: "foo-bar-baz",
			UserID: "abc-123",
			Intent: "test-things",
		}
		expected := Event{
			APIKey:    "foo-bar-baz",
			UserID:    "abc-123",
			Intent:    "test-things",
			TimeStamp: 667722,
			Platform:  "fantasy-chat",
			Version:   "1.2.45",
		}
		e.SetTimeStamp(667722).SetPlatform("fantasy-chat").SetVersion("1.2.45")
		if !reflect.DeepEqual(expected, e) {
			t.Errorf("Expected %v, got %v", expected, e)
		}
	})
}

func TestEvent_SetProperty(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		e := Event{
			APIKey: "foo-bar-baz",
			UserID: "abc-123",
			Intent: "test-things",
		}
		expected := Event{
			APIKey: "foo-bar-baz",
			UserID: "abc-123",
			Intent: "test-things",
			Properties: []EventProperty{
				{Name: "one", StringValue: "one"},
				{Name: "two", IntegerValue: 2},
				{Name: "three", FloatValue: 3.333},
			},
		}
		if err := e.AddProperty("one", "one"); err != nil {
			t.Errorf("Unexpected error %v", err)
		}
		if err := e.AddProperty("two", 2); err != nil {
			t.Errorf("Unexpected error %v", err)
		}
		if err := e.AddProperty("three", 3.333); err != nil {
			t.Errorf("Unexpected error %v", err)
		}
		if !reflect.DeepEqual(expected, e) {
			t.Errorf("Expected %v, got %v", expected, e)
		}
	})
	t.Run("error", func(t *testing.T) {
		e := Event{
			APIKey: "foo-bar-baz",
			UserID: "abc-123",
			Intent: "test-things",
		}
		if err := e.AddProperty("nope", []int{99}); err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

func TestSubmit_Event(t *testing.T) {
	tests := []struct {
		name        string
		handler     http.Handler
		event       Event
		expectError bool
	}{
		{
			"default",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					http.Error(w, "expected POST request", http.StatusMethodNotAllowed)
					return
				}
				body, bodyErr := ioutil.ReadAll(r.Body)
				if bodyErr != nil {
					http.Error(w, bodyErr.Error(), http.StatusInternalServerError)
					return
				}
				s := string(body)
				if !strings.Contains(s, "foo-bar-baz") {
					t.Error("expected body to contain correct payload")
				}
				if !strings.Contains(s, "abc123") {
					t.Error("expected body to contain correct payload")
				}
				if !strings.Contains(s, "test-things") {
					t.Error("expected body to contain correct payload")
				}
				w.Write([]byte("OK!"))
			}),
			Event{
				APIKey: "foo-bar-baz",
				UserID: "abc123",
				Intent: "test-things",
			},
			false,
		},
		{
			"server error",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "zalgo!!", http.StatusInternalServerError)
			}),
			Event{
				APIKey: "foo-bar-baz",
				UserID: "abc123",
				Intent: "test-things",
			},
			true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			oldEventEndpoint := eventEndpoint
			defer func() { eventEndpoint = oldEventEndpoint }()

			ts := httptest.NewServer(test.handler)
			eventEndpoint = ts.URL

			err := test.event.Submit()
			if test.expectError != (err != nil) {
				t.Errorf("Unexpected error result %v", err)
			}
		})
	}

}
func TestSubmit_Events(t *testing.T) {
	tests := []struct {
		name        string
		handler     http.Handler
		events      Events
		expectError bool
	}{
		{
			"default",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					http.Error(w, "expected POST request", http.StatusMethodNotAllowed)
					return
				}
				body, bodyErr := ioutil.ReadAll(r.Body)
				if bodyErr != nil {
					http.Error(w, bodyErr.Error(), http.StatusInternalServerError)
					return
				}
				s := string(body)
				if !strings.Contains(s, "foo-bar-baz") {
					t.Error("expected body to contain correct payload")
				}
				if !strings.Contains(s, "bar-bar-bar") {
					t.Error("expected body to contain correct payload")
				}
				if !strings.Contains(s, "abc123") {
					t.Error("expected body to contain correct payload")
				}
				if !strings.Contains(s, "123abc") {
					t.Error("expected body to contain correct payload")
				}
				if !strings.Contains(s, "test-things") {
					t.Error("expected body to contain correct payload")
				}
				if !strings.Contains(s, "things-to-test") {
					t.Error("expected body to contain correct payload")
				}
				w.Write([]byte("OK!"))
			}),
			Events{
				{
					APIKey: "foo-bar-baz",
					UserID: "abc123",
					Intent: "test-things",
				},
				{
					APIKey: "bar-bar-bar",
					UserID: "123abc",
					Intent: "things-to-test",
				},
			},
			false,
		},
		{
			"server error",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "zalgo!!", http.StatusInternalServerError)
			}),
			Events{{
				APIKey: "foo-bar-baz",
				UserID: "abc123",
				Intent: "test-things",
			}},
			true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			oldEventsEndpoint := eventsEndpoint
			defer func() { eventsEndpoint = oldEventsEndpoint }()

			ts := httptest.NewServer(test.handler)
			eventsEndpoint = ts.URL

			err := test.events.Submit()
			if test.expectError != (err != nil) {
				t.Errorf("Unexpected error result %v", err)
			}
		})
	}
}

func TestAppend_Events(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		evs := Events{}
		evs.Append(&Event{APIKey: "foo-bar", UserID: "bar-foo"})
		expected := Events{{APIKey: "foo-bar", UserID: "bar-foo"}}
		if !reflect.DeepEqual(expected, evs) {
			t.Errorf("Expected %#v, got %#v", expected, evs)
		}
	})
}
