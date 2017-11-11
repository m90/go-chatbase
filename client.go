package chatbase

// Client wraps a chatbase API Key and can be used to
// collect and submit messages and events
type Client string

// NewClient returns a new Client using the given chatbase API Key
func NewClient(apiKey string) *Client {
	c := Client(apiKey)
	return &c
}

func (c *Client) String() string {
	return string(*c)
}

// NewMessage returns a new Message using the client's key and
// the current time as a timestamp
func (c *Client) NewMessage(typ, userID, platform string) *Message {
	return &Message{
		APIKey:    c.String(),
		Type:      typ,
		UserID:    userID,
		TimeStamp: TimeStamp(),
		Platform:  platform,
	}
}

// NewUserMessage is a convenience method for creating a user created message
func (c *Client) NewUserMessage(userID, platform string) *Message {
	return c.NewMessage(MessageTypeUser, userID, platform)
}

// NewAgentMessage is a convenience method for creating an agent created message
func (c *Client) NewAgentMessage(userID, platform string) *Message {
	return c.NewMessage(MessageTypeAgent, userID, platform)
}

// NewEvent creates a new Event using the client's chatbase API Key
func (c *Client) NewEvent(userID, intent string) *Event {
	return &Event{
		APIKey: c.String(),
		UserID: userID,
		Intent: intent,
	}
}
