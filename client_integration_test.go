// +build integration

package chatbase_test

import (
	"os"
	"testing"

	chatbase "github.com/m90/go-chatbase"
)

var apiKey string

const (
	userID   = "abc-123"
	platform = "integration-test"
)

func TestMain(m *testing.M) {
	if apiKey = os.Getenv("CHATBASE_API_KEY"); apiKey == "" {
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func TestMessages(t *testing.T) {
	t.Run("single", func(t *testing.T) {
		client := chatbase.NewClient(apiKey)
		message := client.AgentMessage(userID, platform)
		message.SetIntent("always-on-time")
		msgRes, msgErr := message.Submit()
		if msgErr != nil {
			t.Errorf("Unexpected error %v", msgErr)
		}
		if !msgRes.Status.OK() {
			t.Errorf("Unexpected status %v", msgRes.Status)
		}
		update := client.Update(msgRes.MessageID.String())
		update.SetVersion("1.2.3")
		updateRes, updateErr := update.Submit()
		if updateErr != nil {
			t.Errorf("Unexpected error %v", updateErr)
		}
		if !updateRes.Status.OK() {
			t.Errorf("Unexpected status %v", updateRes.Status)
		}
	})
	t.Run("multiple", func(t *testing.T) {
		client := chatbase.NewClient(apiKey)
		messages := chatbase.Messages{}
		messages.Append(
			client.UserMessage(userID, platform).SetMessage("Hello Bot!"),
			client.AgentMessage(userID, platform).SetMessage("Hello User!"),
		)
		msgRes, msgErr := messages.Submit()
		if msgErr != nil {
			t.Errorf("Unexpected error %v", msgErr)
		}
		if !msgRes.Status.OK() {
			t.Errorf("Unexpected status %v", msgRes.Status)
		}
		update := client.Update(msgRes.Responses[0].MessageID.String())
		update.SetIntent("slightly-too-late")
		updateRes, updateErr := update.Submit()
		if updateErr != nil {
			t.Errorf("Unexpected error %v", updateErr)
		}
		if !updateRes.Status.OK() {
			t.Errorf("Unexpected status %v", updateRes)
		}
	})
}

func TestEvents(t *testing.T) {
	t.Run("single", func(t *testing.T) {
		client := chatbase.NewClient(apiKey)
		ev := client.Event(userID, "send-an-event")
		ev.SetPlatform(platform).AddProperty("is-this-a-test", true)
		if err := ev.Submit(); err != nil {
			t.Errorf("Unexpected error %v", err)
		}
	})
	t.Run("multiple", func(t *testing.T) {
		client := chatbase.NewClient(apiKey)
		events := chatbase.Events{}
		ev1 := client.Event(userID, "send-multiple-events")
		ev1.AddProperty("number", 1)
		ev2 := client.Event(userID, "send-multiple-events")
		ev2.AddProperty("number", 2)
		ev3 := client.Event(userID, "send-multiple-events")
		ev3.AddProperty("number", 3)
		events.Append(ev1, ev2, ev3)
		if err := events.Submit(); err != nil {
			t.Errorf("Unexpected error %v", err)
		}
	})
}
