package chatbase

// Client wraps a Chatbase API Key and can be used to
// collect and submit messages and events
type Client string

// NewClient returns a new Client using the given Chatbase API Key
func NewClient(apiKey string) *Client {
	c := Client(apiKey)
	return &c
}

func (c *Client) String() string {
	return string(*c)
}

// Message returns a new Message using the client's key and
// the current time as a timestamp
func (c *Client) Message(typ MessageType, userID, platform string) *Message {
	return &Message{
		APIKey:    c.String(),
		Type:      typ,
		UserID:    userID,
		TimeStamp: TimeStamp(),
		Platform:  platform,
	}
}

// UserMessage is a convenience method for creating a user created message
func (c *Client) UserMessage(userID, platform string) *Message {
	return c.Message(UserType, userID, platform)
}

// AgentMessage is a convenience method for creating an agent created message
func (c *Client) AgentMessage(userID, platform string) *Message {
	return c.Message(AgentType, userID, platform)
}

// Event creates a new Event using the client's API Key
func (c *Client) Event(userID, intent string) *Event {
	return &Event{
		APIKey: c.String(),
		UserID: userID,
		Intent: intent,
	}
}

// Update creates a new Update using the client's API key
func (c *Client) Update(messageID string) *Update {
	return &Update{
		APIKey:    c.String(),
		MessageID: messageID,
	}
}

// FacebookMessage creates a new native Facebook message
func (c *Client) FacebookMessage(payload interface{}) *FacebookMessage {
	return &FacebookMessage{
		Payload: payload,
		APIKey:  c.String(),
	}
}

// FacebookRequestResponse creates a new wrapper around a request and response
func (c *Client) FacebookRequestResponse(request, response interface{}) *FacebookRequestResponse {
	return &FacebookRequestResponse{
		APIKey:   c.String(),
		Request:  request,
		Response: response,
	}
}
