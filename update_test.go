package chatbase

import (
	"reflect"
	"testing"
)

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
