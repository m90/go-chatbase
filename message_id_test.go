package chatbase

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalJSON_MessageID(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		expected    MessageID
	}{
		{
			"int",
			"200",
			false,
			"200",
		},
		{
			"string",
			`"200"`,
			false,
			"200",
		},
		{
			"bad value",
			"true",
			true,
			"",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var v MessageID
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

func TestString_MessageID(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		m := MessageID("foobar")
		if m.String() != "foobar" {
			t.Errorf("Expected foobar, got %v", m.String())
		}
	})
}
