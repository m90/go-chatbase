// +build integration

package chatbase_test

import (
	"fmt"
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
		)
		messages.Append(
			client.AgentMessage(userID, platform).SetMessage("Hello User!"),
		)
		msgRes, msgErr := messages.Submit()
		if msgErr != nil {
			t.Errorf("Unexpected error %v", msgErr)
		}
		if !msgRes.Status.OK() {
			t.Errorf("Unexpected status %v", msgRes.Status)
		}
		fmt.Printf("%#v\n", msgRes)
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
