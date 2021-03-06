package chatbase

import (
	"testing"
)

func TestLink(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		c := New("foo-bar-baz")
		l := c.Link("http://www.example.net/article", "fantasy").SetVersion("1.2.3")
		expected := "https://chatbase.com/r?api_key=foo-bar-baz&platform=fantasy&url=http%3A%2F%2Fwww.example.net%2Farticle&version=1.2.3"
		href, err := l.Encode()
		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}
		if expected != href {
			t.Errorf("Expected %v, got %v", expected, href)
		}
	})
}
