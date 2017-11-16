// +build integration

package chatbase_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	chatbase "github.com/m90/go-chatbase"
)

var apiKey string

const (
	userID   = "abc-123"
	platform = "integration-test"
)

func readFixture(p string) (map[string]interface{}, error) {
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(b, &payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func TestMain(m *testing.M) {
	if apiKey = os.Getenv("CHATBASE_API_KEY"); apiKey == "" {
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func TestMessages(t *testing.T) {
	t.Run("single", func(t *testing.T) {
		client := chatbase.New(apiKey)
		message := client.UserMessage(userID, platform)
		message.SetFeedback(true)
		message.SetIntent("always-on-time")
		msgRes, msgErr := message.Submit()
		if msgErr != nil {
			t.Fatalf("Unexpected error %v", msgErr)
		}
		if !msgRes.Status.OK() {
			t.Fatalf("Unexpected status %v with reason %v", msgRes.Status, msgRes.Reason)
		}
		update := client.Update(msgRes.MessageID.String())
		update.SetVersion("1.2.3")
		update.SetNotHandled("true")
		updateRes, updateErr := update.Submit()
		if updateErr != nil {
			t.Fatalf("Unexpected error %v", updateErr)
		}
		if !updateRes.Status.OK() {
			t.Fatalf("Unexpected status %v with reason %v", updateRes.Status, updateRes.Reason)
		}
	})
	t.Run("multiple", func(t *testing.T) {
		client := chatbase.New(apiKey)
		messages := chatbase.Messages{}
		messages.Append(
			client.UserMessage(userID, platform).SetMessage("Hello Bot!"),
			client.AgentMessage(userID, platform).SetMessage("Hello User!"),
		)
		msgRes, msgErr := messages.Submit()
		if msgErr != nil {
			t.Fatalf("Unexpected error %v", msgErr)
		}
		if !msgRes.Status.OK() {
			t.Fatalf("Unexpected status %v with reason %v", msgRes.Status, msgRes.Reason)
		}
		update := client.Update(msgRes.Responses[0].MessageID.String())
		update.SetIntent("slightly-too-late")
		updateRes, updateErr := update.Submit()
		if updateErr != nil {
			t.Fatalf("Unexpected error %v", updateErr)
		}
		if !updateRes.Status.OK() {
			t.Fatalf("Unexpected status %v with reason %v", updateRes, updateRes.Reason)
		}
	})
}

func TestEvents(t *testing.T) {
	t.Run("single", func(t *testing.T) {
		client := chatbase.New(apiKey)
		ev := client.Event(userID, "send-an-event")
		ev.SetPlatform(platform).AddProperty("is-this-a-test", true)
		if err := ev.Submit(); err != nil {
			t.Fatalf("Unexpected error %v", err)
		}
	})
	t.Run("multiple", func(t *testing.T) {
		client := chatbase.New(apiKey)
		events := chatbase.Events{}
		for i := 1; i < 4; i++ {
			ev := client.Event(userID, "send-multiple-events")
			if err := ev.AddProperty("number", i); err != nil {
				t.Fatalf("Unexpected error %v", err)
			}
			events.Append(ev)
		}
		if err := events.Submit(); err != nil {
			t.Fatalf("Unexpected error %v", err)
		}
	})
}

func TestFacebookMessages(t *testing.T) {
	t.Run("single", func(t *testing.T) {
		client := chatbase.New(apiKey)
		payload, err := readFixture("testdata/facebook_single_payload.json")
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}
		fbMessage := client.FacebookMessage(payload)
		fbMessage.SetIntent("test-facebook")
		response, responseErr := fbMessage.Submit()
		if responseErr != nil {
			t.Fatalf("Unexpected error %v", responseErr)
		}
		if !response.Status.OK() {
			t.Fatalf("Unexpected status %v with reason %v", response.Status, response.Reason)
		}
		update := client.Update(response.MessageID.String())
		update.SetVersion("1.2.4")
		updateRes, updateErr := update.Submit()
		if updateErr != nil {
			t.Fatalf("Unexpected error %v", updateErr)
		}
		if !updateRes.Status.OK() {
			t.Fatalf("Unexpected status %v with reason %v", updateRes.Status, updateRes.Reason)
		}
	})
	t.Run("multiple", func(t *testing.T) {
		client := chatbase.New(apiKey)
		payload, err := readFixture("testdata/facebook_single_payload.json")
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}
		fbMessages := chatbase.FacebookMessages{}
		for i := 0; i < 4; i++ {
			fbMessage := client.FacebookMessage(payload)
			fbMessage.SetIntent(fmt.Sprintf("number-%d", i))
			fbMessages.Append(fbMessage)
		}
		response, responseErr := fbMessages.Submit()
		if responseErr != nil {
			t.Fatalf("Unexpected error %v", responseErr)
		}
		if !response.Status.OK() {
			t.Fatalf("Unexpected status %v with reason %v", response.Status, response.Reason)
		}
		update := client.Update(response.Responses[0].MessageID.String())
		update.SetVersion("9.99.1")
		updateRes, updateErr := update.Submit()
		if updateErr != nil {
			t.Fatalf("Unexpected error %v", err)
		}
		if !updateRes.Status.OK() {
			t.Fatalf("Unexpected status %v with reason %v", updateRes.Status, updateRes.Reason)
		}
	})

	t.Run("multiple request response objects", func(t *testing.T) {
		t.SkipNow() // as of 2017-11-16 chatbase does return `all_succeeded: true` with reason: "All messages errored."

		requestPayload, requestErr := readFixture("testdata/facebook_single_request.json")
		if requestErr != nil {
			t.Fatalf("Unexpected error %v", requestErr)
		}
		responsePayload, responseErr := readFixture("testdata/facebook_single_response.json")
		if responseErr != nil {
			t.Fatalf("Unexpected error %v", responseErr)
		}

		client := chatbase.New(apiKey)
		fbMessages := chatbase.FacebookRequestResponses{}
		for i := 0; i < 1; i++ {
			fbMessage := client.FacebookRequestResponse(requestPayload, responsePayload)
			fbMessage.SetIntent(fmt.Sprintf("number-%d", i))
			fbMessages.Append(fbMessage)
		}
		response, responseErr := fbMessages.Submit()
		if responseErr != nil {
			t.Fatalf("Unexpected error %v", responseErr)
		}
		fmt.Printf("response %#v\n", response)
		if !response.Status.OK() {
			t.Fatalf("Unexpected status %v with reason %v", response.Status, response.Reason)
		}
	})
}

func TestLink(t *testing.T) {
	t.Run("single", func(t *testing.T) {
		client := chatbase.New(apiKey)
		click := client.Link("https://www.chatbase.com", "integration-test")
		click.SetVersion("9.8.7")
		res, err := click.Submit()
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}
		if !res.Status.OK() {
			t.Errorf("Unexpected status %v with reason %v", res.Status, res.Reason)
		}
	})
}
