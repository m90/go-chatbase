package chatbase

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// MessageID normalizes the IDs sent by Chatbase
type MessageID string

// UnmarshalJSON distinguishes ints and strings and normalizes
// both values into a string representation
func (m *MessageID) UnmarshalJSON(b []byte) error {
	var i int
	if err := json.Unmarshal(b, &i); err == nil {
		*m = MessageID(strconv.Itoa(i))
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
