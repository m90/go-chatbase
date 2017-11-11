package chatbase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// types of messages used to identify the origin of
// a message in Chatbase
const (
	MessageTypeUser  = "user"
	MessageTypeAgent = "agent"
)

var (
	messagesEndpoint = "https://chatbase.com/api/messages"
)

// Message contains info about a generic chat message
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
func (m *Message) SetMessage(msg string) *Message {
	m.Message = msg
	return m
}

// SetIntent adds an optional intent value
func (m *Message) SetIntent(i string) *Message {
	m.Intent = i
	return m
}

// SetNotHandled adds an optional not handled flag
func (m *Message) SetNotHandled(n bool) *Message {
	m.NotHandled = n
	return m
}

// SetFeedback adds an optional not feedback flag
func (m *Message) SetFeedback(f bool) *Message {
	m.Feedback = f
	return m
}

// SetVersion adds an optional version value
func (m *Message) SetVersion(v string) *Message {
	m.Version = v
	return m
}

// SetTimeStamp overrides the message's timestamp value
func (m *Message) SetTimeStamp(t int) *Message {
	m.TimeStamp = t
	return m
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

// Submit tries to deliver the message to Chatbase
func (m *Message) Submit() (*MessageResponse, error) {
	body, err := postMessage(m)
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

// MessageResponse describes a Chatbase response to sending a single message
// or is contained in a set of responses when performing batch operations
type MessageResponse struct {
	MessageID MessageID `json:"message_id"`
	Status    Status    `json:"status"`
}

// Messages is a collection of Message
type Messages []Message

// Submit tries to deliver the set of messages to Chatbase
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

// Append adds a message to the the collection
func (m *Messages) Append(addition *Message) *Messages {
	*m = append(*m, *addition)
	return m
}

// MessagesResponse describes a Chatbase response to sending multiple messages at once
type MessagesResponse struct {
	AllSucceded bool              `json:"all_succeeded"`
	Status      Status            `json:"status"`
	Responses   []MessageResponse `json:"responses"`
}
