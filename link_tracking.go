package chatbase

import (
	"encoding/json"
	"net/url"
)

var (
	redirectURL   = "https://chatbase.com/r"
	clickEndpoint = "https://chatbase.com/api/click"
)

// Link describes a hyperlink that should be tracked using Chatbase
type Link struct {
	APIKey   string `json:"api_key"`
	URL      string `json:"url"`
	Platform string `json:"platform"`
	Version  string `json:"version,omitempty"`
}

// LinkResponse contains the response to posting a click event
type LinkResponse struct {
	Status Status `json:"status"`
	Reason string `json:"reason,omitempty"`
}

// SetVersion adds an optional version parameter
func (l *Link) SetVersion(v string) *Link {
	l.Version = v
	return l
}

// Submit tries to send the click event to Chatbase
func (l *Link) Submit() (*LinkResponse, error) {
	body, bodyErr := apiPost(clickEndpoint, l)
	if bodyErr != nil {
		return nil, bodyErr
	}
	defer body.Close()
	response := LinkResponse{}
	if err := json.NewDecoder(body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (l *Link) String() string {
	u, _ := url.Parse(redirectURL)
	q := u.Query()
	q.Set("api_key", l.APIKey)
	q.Set("url", l.URL)
	q.Set("platform", l.Platform)
	if l.Version != "" {
		q.Set("version", l.Version)
	}
	u.RawQuery = q.Encode()
	return u.String()
}
