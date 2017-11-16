package chatbase

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestMessage_Setters(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		m := Message{
			APIKey:    "secret",
			Type:      AgentType,
			UserID:    "abc-123",
			TimeStamp: 998877,
			Platform:  "test",
		}
		expected := Message{
			APIKey:     "secret",
			Type:       AgentType,
			UserID:     "abc-123",
			TimeStamp:  111444,
			Platform:   "test",
			Message:    "Hello world!",
			Intent:     "test-things",
			NotHandled: true,
			Feedback:   true,
			Version:    "1.2.34",
		}
		m.SetMessage("Hello world!").SetIntent("test-things").SetNotHandled(true).SetFeedback(true).SetVersion("1.2.34").SetTimeStamp(111444)
		if !reflect.DeepEqual(expected, m) {
			t.Errorf("Expected %#v, got %#v", expected, m)
		}
	})
}

func TestSubmit_Message(t *testing.T) {
	tests := []struct {
		name             string
		handler          http.Handler
		message          Message
		expectError      bool
		expectedResponse *MessageResponse
	}{
		{
			"default",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				b, _ := ioutil.ReadFile("testdata/message_response_ok.json")
				w.Write(b)
			}),
			Message{},
			false,
			&MessageResponse{
				MessageID: "abc987",
				Status:    true,
			},
		},
		{
			"error",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				b, _ := ioutil.ReadFile("testdata/message_response_error.json")
				w.Write(b)
			}),
			Message{},
			false,
			&MessageResponse{
				Status: false,
			},
		},
		{
			"broken payload",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				b, _ := ioutil.ReadFile("testdata/message_response_broken.json")
				w.Write(b)
			}),
			Message{},
			true,
			nil,
		},
		{
			"internal server error",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "zalgo!", http.StatusInternalServerError)
			}),
			Message{},
			true,
			nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			oldMessageEndpoint := messageEndpoint
			defer func() { messageEndpoint = oldMessageEndpoint }()

			ts := httptest.NewServer(test.handler)
			messageEndpoint = ts.URL

			result, err := test.message.Submit()
			if test.expectError != (err != nil) {
				t.Errorf("Unexpected error value %v", err)
			}
			if !reflect.DeepEqual(test.expectedResponse, result) {
				t.Errorf("Expected %#v, got %#v", test.expectedResponse, result)
			}
		})
	}
}
func TestSubmit_Messages(t *testing.T) {
	tests := []struct {
		name             string
		handler          http.Handler
		message          Messages
		expectError      bool
		expectedResponse *MessagesResponse
	}{
		{
			"default",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				b, _ := ioutil.ReadFile("testdata/messages_response_ok.json")
				w.Write(b)
			}),
			Messages{{}},
			false,
			&MessagesResponse{
				AllSucceeded: true,
				Status:       true,
				Responses: []MessageResponse{
					{
						MessageID: "123456789",
						Status:    true,
					},
					{
						MessageID: "987654321",
						Status:    true,
					},
				},
			},
		},
		{
			"error",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				b, _ := ioutil.ReadFile("testdata/messages_response_error.json")
				w.Write(b)
			}),
			Messages{{}},
			false,
			&MessagesResponse{
				AllSucceeded: false,
				Status:       false,
				Responses:    []MessageResponse{},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			oldMessagesEndpoint := messagesEndpoint
			defer func() { messagesEndpoint = oldMessagesEndpoint }()

			ts := httptest.NewServer(test.handler)
			messagesEndpoint = ts.URL

			result, err := test.message.Submit()
			if test.expectError != (err != nil) {
				t.Errorf("Unexpected error value %v", err)
			}
			if !reflect.DeepEqual(test.expectedResponse, result) {
				t.Errorf("Expected %#v, got %#v", test.expectedResponse, result)
			}
		})
	}
}

func TestMessage_Append(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		msgs := Messages{}
		msgs.Append(&Message{APIKey: "foo-bar", UserID: "bar-foo"})
		expected := Messages{{APIKey: "foo-bar", UserID: "bar-foo"}}
		if !reflect.DeepEqual(expected, msgs) {
			t.Errorf("Expected %#v, got %#v", expected, msgs)
		}
	})
}
