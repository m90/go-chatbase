package chatbase

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

type badReader int

func (badReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("zalgo")
}
func TestAugmentURL(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		params      map[string]string
		expectError bool
		expected    string
	}{
		{
			"default",
			"https://example.net",
			map[string]string{"foo": "bar", "baz": "qux"},
			false,
			"https://example.net?baz=qux&foo=bar",
		},
		{
			"encoding",
			"https://example.net",
			map[string]string{"foo": "bar baz"},
			false,
			"https://example.net?foo=bar+baz",
		},
		{
			"additional parameters",
			"https://example.net?foo=bar",
			map[string]string{"baz": "qux"},
			false,
			"https://example.net?baz=qux&foo=bar",
		},
		{
			"bad url",
			"%%%%%%%%üü#üü#",
			nil,
			true,
			"",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := augmentURL(test.input, test.params)
			if test.expectError != (err != nil) {
				t.Errorf("Unexpected error %v", err)
			}
			if test.expected != result {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}
func TestDecodeInto(t *testing.T) {
	type testTarget struct {
		Prop string `json:"prop"`
	}
	tests := []struct {
		name           string
		target         interface{}
		thunk          func() (io.ReadCloser, error)
		expectError    bool
		expectedResult interface{}
	}{
		{
			"default",
			&testTarget{},
			func() (io.ReadCloser, error) {
				return ioutil.NopCloser(bytes.NewReader([]byte(`{"prop":"foo"}`))), nil
			},
			false,
			&testTarget{Prop: "foo"},
		},
		{
			"bad thunk",
			&testTarget{},
			func() (io.ReadCloser, error) {
				return nil, errors.New("zalgo")
			},
			true,
			nil,
		},
		{
			"bad reader",
			&testTarget{},
			func() (io.ReadCloser, error) {
				return ioutil.NopCloser(badReader(0)), nil
			},
			true,
			nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := decodeInto(test.target, test.thunk)
			if test.expectError != (err != nil) {
				t.Errorf("Unexpected error %v", err)
			}
			if !reflect.DeepEqual(test.expectedResult, result) {
				t.Errorf("Expected %#v, got %#v", test.expectedResult, result)
			}
		})
	}
}

func TestNewResponses(t *testing.T) {
	okThunk := func() (io.ReadCloser, error) {
		return ioutil.NopCloser(strings.NewReader("{}")), nil
	}
	nilThunk := func() (io.ReadCloser, error) {
		return nil, nil
	}
	{
		response, err := newUpdateResponse(okThunk)
		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}
		if !reflect.DeepEqual(&UpdateResponse{}, response) {
			t.Errorf("Expected zero struct, got %v", response)
		}
	}
	{
		_, err := newUpdateResponse(nilThunk)
		if err == nil {
			t.Errorf("Unexpected error %v", err)
		}
	}
	{
		response, err := newLinkResponse(okThunk)
		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}
		if !reflect.DeepEqual(&LinkResponse{}, response) {
			t.Errorf("Expected zero struct, got %v", response)
		}
	}
	{
		_, err := newLinkResponse(nilThunk)
		if err == nil {
			t.Errorf("Unexpected error %v", err)
		}
	}
	{
		response, err := newMessagesResponse(okThunk)
		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}
		if !reflect.DeepEqual(&MessagesResponse{}, response) {
			t.Errorf("Expected zero struct, got %v", response)
		}
	}
	{
		_, err := newMessagesResponse(nilThunk)
		if err == nil {
			t.Errorf("Unexpected error %v", err)
		}
	}
	{
		response, err := newMessageResponse(okThunk)
		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}
		if !reflect.DeepEqual(&MessageResponse{}, response) {
			t.Errorf("Expected zero struct, got %v", response)
		}
	}
	{
		_, err := newMessageResponse(nilThunk)
		if err == nil {
			t.Errorf("Unexpected error %v", err)
		}
	}
}

func TestApiCall(t *testing.T) {
	defaultClientTimeout := client.Timeout
	defer func() { client.Timeout = defaultClientTimeout }()
	client.Timeout = time.Second

	tests := []struct {
		name         string
		handler      http.Handler
		urlOverride  string
		method       string
		data         interface{}
		expectError  bool
		expectedBody string
	}{
		{
			"bad payload",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("OK!"))
			}),
			"",
			http.MethodPost,
			func() string { return "zalgo" },
			true,
			"",
		},
		{
			"bad url",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("OK!"))
			}),
			"%%%%%%zalgo##üü%&&%%",
			http.MethodPost,
			map[string]string{"hello": "world"},
			true,
			"",
		},
		{
			"request error due to timeout",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(time.Hour)
				w.Write([]byte("OK!"))
			}),
			"",
			http.MethodPost,
			map[string]string{"hello": "world"},
			true,
			"",
		},
		{
			"server error",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}),
			"",
			http.MethodPost,
			map[string]string{"hello": "world"},
			true,
			"",
		},
		{
			"default",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPut {
					http.Error(w, "bad method", http.StatusBadRequest)
					return
				}
				if r.Header.Get("Content-Type") != "application/json" {
					http.Error(w, "bad content type", http.StatusBadRequest)
					return
				}
				defer r.Body.Close()
				b, err := ioutil.ReadAll(r.Body)
				if err != nil {
					http.Error(w, "bad body", http.StatusBadRequest)
					return
				}
				if s := string(b); s != `{"hello":"world"}` {
					http.Error(w, "incorrect body", http.StatusBadRequest)
					return
				}
				w.Write([]byte("OK!"))
			}),
			"",
			http.MethodPut,
			map[string]string{"hello": "world"},
			false,
			"OK!",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts := httptest.NewServer(test.handler)
			endpoint := ts.URL
			if test.urlOverride != "" {
				endpoint = test.urlOverride
			}
			res, err := apiCall(test.method, endpoint, test.data)
			if test.expectError != (err != nil) {
				t.Errorf("Unexpected error %v", err)
			}
			if res != nil {
				defer res.Close()
				body, bodyErr := ioutil.ReadAll(res)
				if bodyErr != nil {
					t.Errorf("Unexpected error %v", bodyErr)
				}
				if s := string(body); test.expectedBody != s {
					t.Errorf("Expected %v, got %v", test.expectedBody, s)
				}
			}
		})
	}
}

func TestApiPost(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "bad method", http.StatusBadRequest)
				return
			}
			w.Write([]byte("OK!"))
		}))
		res, err := apiPost(ts.URL, map[string]string{})
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}
		body, bodyErr := ioutil.ReadAll(res)
		if bodyErr != nil {
			t.Fatalf("Unexpected error %v", bodyErr)
		}
		if s := string(body); s != "OK!" {
			t.Errorf("Unexpected body %v", s)
		}
	})
}
func TestApiPut(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPut {
				http.Error(w, "bad method", http.StatusBadRequest)
				return
			}
			w.Write([]byte("OK!"))
		}))
		res, err := apiPut(ts.URL, map[string]string{})
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}
		body, bodyErr := ioutil.ReadAll(res)
		if bodyErr != nil {
			t.Fatalf("Unexpected error %v", bodyErr)
		}
		if s := string(body); s != "OK!" {
			t.Errorf("Unexpected body %v", s)
		}
	})
}
