package chatbase

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestSubmit_Update(t *testing.T) {
	tests := []struct {
		name             string
		handler          http.Handler
		update           Update
		expectError      bool
		expectedResponse *UpdateResponse
	}{
		{
			"success",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				q := r.URL.Query()
				if q.Get("api_key") != "fixture" {
					http.Error(w, "bad api_key", http.StatusBadRequest)
					return
				}
				if q.Get("message_id") != "abc-123" {
					http.Error(w, "bad message_id", http.StatusBadRequest)
					return
				}
				d, _ := ioutil.ReadFile("testdata/update_response_ok.json")
				w.Write(d)
			}),
			Update{
				APIKey:    "fixture",
				MessageID: "abc-123",
			},
			false,
			&UpdateResponse{
				Status:  true,
				Updated: []string{"intent"},
				Error:   []string{},
			},
		},
		{
			"error response",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				d, _ := ioutil.ReadFile("testdata/update_response_error.json")
				w.Write(d)
			}),
			Update{},
			false,
			&UpdateResponse{
				Status:  false,
				Updated: []string{},
				Error:   []string{"intent"},
			},
		},
		{
			"broken response",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				d, _ := ioutil.ReadFile("testdata/update_response_broken.json")
				w.Write(d)
			}),
			Update{},
			true,
			nil,
		},
		{
			"internal server error",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "zalgo", http.StatusInternalServerError)
			}),
			Update{},
			true,
			nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts := httptest.NewServer(test.handler)
			oldEndpoint := updateEndpoint
			updateEndpoint = ts.URL
			defer func() { updateEndpoint = oldEndpoint }()
			res, err := test.update.Submit()
			if test.expectError != (err != nil) {
				t.Errorf("Unexpected error value %v", err)
			}
			if !reflect.DeepEqual(test.expectedResponse, res) {
				t.Errorf("Expected %#v, got %#v", test.expectedResponse, res)
			}
		})
	}
}

func TestUpdate_Setters(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		u := Update{
			APIKey:    "fixture",
			MessageID: "abc123",
		}
		expected := Update{
			APIKey:     "fixture",
			MessageID:  "abc123",
			Intent:     "test-things",
			NotHandled: "true",
			Feedback:   "true",
			Version:    "1.2.34",
		}
		u.SetIntent("test-things").SetNotHandled(true).SetFeedback("true").SetVersion("1.2.34")
		if !reflect.DeepEqual(expected, u) {
			t.Errorf("Expected %#v, got %#v", expected, u)
		}
	})
}
