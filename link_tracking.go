package chatbase

import (
	"io"
)

var (
	redirectURL   = "https://chatbase.com/r"
	clickEndpoint = "https://chatbase.com/api/click"
)

// Link describes a hyperlink to be tracked using Chatbase
type Link struct {
	APIKey   string `json:"api_key"`
	URL      string `json:"url"`
	Platform string `json:"platform"`
	Version  string `json:"version,omitempty"`
}

// LinkResponse contains the response to submitting a link
type LinkResponse struct {
	Status Status `json:"status"`
	Reason string `json:"reason,omitempty"`
}

// SetVersion adds an optional "version" parameter to a link
func (l *Link) SetVersion(v string) *Link {
	l.Version = v
	return l
}

// Submit tries to send the link to Chatbase
func (l *Link) Submit() (*LinkResponse, error) {
	return newLinkResponse(func() (io.ReadCloser, error) {
		return apiPost(clickEndpoint, l)
	})
}

func (l *Link) String() string {
	params := map[string]string{
		"api_key":  l.APIKey,
		"url":      l.URL,
		"platform": l.Platform,
	}
	if l.Version != "" {
		params["version"] = l.Version
	}
	u, _ := augmentURL(redirectURL, params)
	return u
}
