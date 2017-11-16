package chatbase

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
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
			`{"chatbase_fields":{"intent":"test-things","not_handled":true,"feedback":false,"string":"1.4.4"},"foo":"bar","hello":"world","nested":{"false":false,"true":true}}`,
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

func TestFacebookMessage_Submit(t *testing.T) {
	tests := []struct {
		name             string
		handler          http.Handler
		message          FacebookMessage
		expectError      bool
		expectedResponse *MessageResponse
	}{
		{
			"default",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Query().Get("api_key") != "top-secret" {
					http.Error(w, "missing api key", http.StatusBadRequest)
					return
				}
				body, bodyErr := ioutil.ReadAll(r.Body)
				if bodyErr != nil {
					http.Error(w, "bad payload", http.StatusInternalServerError)
					return
				}
				r.Body.Close()
				s := string(body)
				if !strings.Contains(s, `"recipient":"zuck"`) {
					http.Error(w, "missing data from payload", http.StatusBadRequest)
					return
				}
				b, _ := ioutil.ReadFile("testdata/message_response_ok.json")
				w.Write(b)
			}),
			FacebookMessage{
				APIKey: "top-secret",
				Payload: map[string]string{
					"recipient": "zuck",
				},
			},
			false,
			&MessageResponse{
				MessageID: "abc987",
				Status:    true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			oldEndpoint := facebookMessageEndpoint
			defer func() { facebookMessageEndpoint = oldEndpoint }()

			ts := httptest.NewServer(test.handler)
			facebookMessageEndpoint = ts.URL

			res, err := test.message.Submit()
			if test.expectError != (err != nil) {
				t.Errorf("Unexpected error value %v", err)
			}
			if !reflect.DeepEqual(test.expectedResponse, res) {
				t.Errorf("Expected %#v, got %#v", test.expectedResponse, res)
			}
		})
	}
}
func TestFacebookMessages_Submit(t *testing.T) {
	tests := []struct {
		name             string
		handler          http.Handler
		messages         FacebookMessages
		expectError      bool
		expectedResponse *MessagesResponse
	}{
		{
			"default",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Query().Get("api_key") != "top-secret" {
					http.Error(w, "missing api key", http.StatusBadRequest)
					return
				}
				body, bodyErr := ioutil.ReadAll(r.Body)
				if bodyErr != nil {
					http.Error(w, "bad payload", http.StatusInternalServerError)
					return
				}
				r.Body.Close()
				s := string(body)
				if !strings.Contains(s, `"recipient":"zuck"`) {
					http.Error(w, "missing data from payload", http.StatusBadRequest)
					return
				}
				b, _ := ioutil.ReadFile("testdata/messages_response_ok.json")
				w.Write(b)
			}),
			FacebookMessages{
				{
					APIKey: "top-secret",
					Payload: map[string]string{
						"recipient": "zuck",
					},
				},
				{
					APIKey: "top-secret",
					Payload: map[string]string{
						"recipient": "bill",
					},
				},
			},
			false,
			&MessagesResponse{
				AllSucceeded: true,
				Responses: []MessageResponse{
					{MessageID: "123456789", Status: true},
					{MessageID: "987654321", Status: true},
				},
				Status: true,
			},
		},
		{
			"empty",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("OK!"))
			}),
			FacebookMessages{},
			true,
			nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			oldEndpoint := facebookMessagesEndpoint
			defer func() { facebookMessagesEndpoint = oldEndpoint }()

			ts := httptest.NewServer(test.handler)
			facebookMessagesEndpoint = ts.URL

			res, err := test.messages.Submit()
			if test.expectError != (err != nil) {
				t.Errorf("Unexpected error value %v", err)
			}
			if !reflect.DeepEqual(test.expectedResponse, res) {
				t.Errorf("Expected %#v, got %#v", test.expectedResponse, res)
			}
		})
	}
}
