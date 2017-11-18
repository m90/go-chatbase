package chatbase

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// MessageID is a named string type that normalizes the IDs sent by Chatbase
type MessageID string

// UnmarshalJSON distinguishes ints and strings and normalizes
// both values into a string representation
func (m *MessageID) UnmarshalJSON(b []byte) error {
	var i int64
	if err := json.Unmarshal(b, &i); err == nil {
		*m = MessageID(strconv.FormatInt(i, 10))
		return nil
	}
	var str string
	if err := json.Unmarshal(b, &str); err == nil {
		*m = MessageID(str)
		return nil
	}
	return fmt.Errorf("could not unmarshal %s into MessageID", b)
}

func (m *MessageID) String() string {
	return string(*m)
}

// Int64 returns the ID's int64 representation
func (m *MessageID) Int64() (int64, error) {
	return strconv.ParseInt(m.String(), 10, 0)
}
