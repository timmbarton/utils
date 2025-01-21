package secs

import (
	"encoding/json"
	"time"
)

type Seconds time.Duration

func (s Seconds) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(s).Seconds())
}

func (s *Seconds) UnmarshalJSON(b []byte) error {
	var seconds int64
	if err := json.Unmarshal(b, &seconds); err != nil {
		return err
	}

	*s = Seconds(seconds * int64(time.Second))

	return nil
}
func (s Seconds) String() string {
	return time.Duration(s).String()
}
