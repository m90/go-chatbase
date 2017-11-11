package chatbase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Types of messages used to identify the origin of
// a message in chatbase
const (
	MessageTypeUser  = "user"
	MessageTypeAgent = "agent"
)

var (
	messagesEndpoint = "https://chatbase.com/api/messages"
)

// Message contains info about a platform agnostic chat message
type Message struct {
	APIKey     string `json:"api_key"`
	Type       string `json:"type"`
	UserID     string `json:"user_id"`
	TimeStamp  int    `json:"time_stamp"`
	Platform   string `json:"platform"`
	Message    string `json:"message,omitempty"`
	Intent     string `json:"intent,omitempty"`
	NotHandled bool   `json:"not_handled,omitempty"`
	Feedback   bool   `json:"feedback,omitempty"`
	Version    string `json:"version,omitempty"`
}

// SetMessage adds an optional message value
func (msg *Message) SetMessage(m string) *Message {
	msg.Message = m
	return msg
}

// SetIntent adds an optional intent value
func (msg *Message) SetIntent(i string) *Message {
	msg.Intent = i
	return msg
}

// SetNotHandled adds an optional not handled flag
func (msg *Message) SetNotHandled(n bool) *Message {
	msg.NotHandled = n
	return msg
}

// SetFeedback adds an optional not feedback flag
func (msg *Message) SetFeedback(f bool) *Message {
	msg.Feedback = f
	return msg
}

// SetVersion adds an optional version value
func (msg *Message) SetVersion(v string) *Message {
	msg.Version = v
	return msg
}

func postMessage(v interface{}) (io.ReadCloser, error) {
	payload, payloadErr := json.Marshal(v)
	if payloadErr != nil {
		return nil, payloadErr
	}
	res, err := http.Post(messagesEndpoint, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= http.StatusInternalServerError {
		return nil, fmt.Errorf("request failed with status %v", res.StatusCode)
	}
	return res.Body, nil
}

// Submit tries to deliver the set of messages to chatbase
func (msg *Message) Submit() (*MessageResponse, error) {
	body, err := postMessage(msg)
	if err != nil {
		return nil, err
	}
	defer body.Close()
	responseData := MessageResponse{}
	if err := json.NewDecoder(body).Decode(&responseData); err != nil {
		return nil, err
	}
	return &responseData, nil
}

// MessageResponse describes an API response to sending a single message
type MessageResponse struct {
	MessageID MessageID `json:"message_id"`
	Status    Status    `json:"status"`
}

// Messages is a collection of Message
type Messages []Message

// Submit tries to deliver the set of messages to chatbase
func (m *Messages) Submit() (*MessagesResponse, error) {
	body, err := postMessage(m)
	if err != nil {
		return nil, err
	}
	defer body.Close()
	responseData := MessagesResponse{}
	if err := json.NewDecoder(body).Decode(&responseData); err != nil {
		return nil, err
	}
	return &responseData, nil
}

// MessagesResponse describes an API response to sending multiple messages at once
type MessagesResponse struct {
	AllSucceded bool              `json:"all_succeeded"`
	Status      Status            `json:"status"`
	Responses   []MessageResponse `json:"responses"`
}
