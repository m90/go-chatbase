package chatbase

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApi(t *testing.T) {
	tests := []struct {
		name            string
		handler         http.Handler
		method          string
		payload         interface{}
		expectError     bool
		expectReadError bool
		expectedBody    string
	}{
		{
			"default",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					http.Error(w, "expected POST request", http.StatusMethodNotAllowed)
					return
				}
				w.Write([]byte("OK!"))
			}),
			http.MethodPost,
			map[string]string{},
			false,
			false,
			"OK!",
		},
		{
			"bad payload",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("OK!"))
			}),
			http.MethodPost,
			func() {},
			true,
			false,
			"",
		},
		{
			"server error",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "zalgo", http.StatusInternalServerError)
			}),
			http.MethodPost,
			map[string]string{},
			true,
			false,
			"",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts := httptest.NewServer(test.handler)
			body, err := apiCall(test.method, ts.URL, test.payload)
			if test.expectError != (err != nil) {
				t.Errorf("Unexpected error %v", err)
			}
			if test.expectError {
				return
			}
			defer body.Close()
			read, readErr := ioutil.ReadAll(body)
			if test.expectReadError != (readErr != nil) {
				t.Errorf("Unexpected error %v", readErr)
			}
			if test.expectedBody != string(read) {
				t.Errorf("Expected %v, got %v", test.expectedBody, string(read))
			}
		})
	}
}
