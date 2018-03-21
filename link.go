package chatbase

import (
	"context"
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

// SubmitWithContext tries to send the link to Chatbase
// while considering the given context's deadline
func (l *Link) SubmitWithContext(ctx context.Context) (*LinkResponse, error) {
	data, err := resultWithContext(ctx, func() (interface{}, error) {
		return l.Submit()
	})
	if err != nil {
		return nil, err
	}
	if res, ok := data.(*LinkResponse); ok {
		return res, nil
	}
	return nil, errBadData
}

// Encode turns the link object into a URL
func (l *Link) Encode() (string, error) {
	params := map[string]string{
		"api_key":  l.APIKey,
		"url":      l.URL,
		"platform": l.Platform,
	}
	if l.Version != "" {
		params["version"] = l.Version
	}
	return augmentURL(redirectURL, params)
}
