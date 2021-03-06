package chatbase

import (
	"encoding/json"
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
