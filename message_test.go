package chatbase

import (
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
