package chatbase

// Client wraps a Chatbase API Key and can be used to
// generate messages, events and link
type Client string

// New returns a new Client using the given Chatbase API Key
func New(apiKey string) *Client {
	c := Client(apiKey)
	return &c
}

func (c *Client) String() string {
	return string(*c)
}

// Message returns a new Message using the client's key and
// the current time as its "TimeStamp" value
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
		MessageID: MessageID(messageID),
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

// Link returns a trackable link to the given URL
func (c *Client) Link(url, platform string) *Link {
	return &Link{
		APIKey:   c.String(),
		URL:      url,
		Platform: platform,
	}
}
