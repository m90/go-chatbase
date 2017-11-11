package chatbase

import (
	"reflect"
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
