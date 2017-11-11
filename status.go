package chatbase

import (
	"encoding/json"
	"fmt"
)

// Status describes if an operation was successful
type Status bool

// UnmarshalJSON normalizes the int and string values that are being used
// by Chatbase to represent success or failure. Int values will map to HTTP
// status codes, otherwise `"success"` is considered true
func (s *Status) UnmarshalJSON(b []byte) error {
	var i int
	if err := json.Unmarshal(b, &i); err == nil {
		*s = Status(i < 400)
		return nil
	}
	var str string
	if err := json.Unmarshal(b, &str); err == nil {
		*s = Status(str == "success")
		return nil
	}
	return fmt.Errorf("could not unmarshal %s into Status", b)
}

// OK returns the boolean representation of Status
func (s *Status) OK() bool {
	return bool(*s)
}
