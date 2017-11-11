package chatbase

import (
	"reflect"
	"testing"
)

func TestMessage_Setters(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		m := Message{
			APIKey:    "secret",
			Type:      MessageTypeAgent,
			UserID:    "abc-123",
			TimeStamp: 998877,
			Platform:  "test",
		}
		expected := Message{
			APIKey:     "secret",
			Type:       MessageTypeAgent,
			UserID:     "abc-123",
			TimeStamp:  998877,
			Platform:   "test",
			Message:    "Hello world!",
			Intent:     "test-things",
			NotHandled: true,
			Feedback:   true,
			Version:    "1.2.34",
		}
		m.SetMessage("Hello world!").SetIntent("test-things").SetNotHandled(true).SetFeedback(true).SetVersion("1.2.34")
		if !reflect.DeepEqual(expected, m) {
			t.Errorf("Expected %#v, got %#v", expected, m)
		}
	})
}
