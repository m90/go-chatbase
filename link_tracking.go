package chatbase

import "net/url"

var (
	redirectURL = "https://chatbase.com/r"
)

// Link describes a hyperlink that should be tracked using Chatbase
type Link struct {
	APIKey   string
	URL      string
	Platform string
	Version  string
}

// SetVersion adds an optional version parameter
func (l *Link) SetVersion(v string) *Link {
	l.Version = v
	return l
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
