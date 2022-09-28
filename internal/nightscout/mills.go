package nightscout

import (
	"encoding/json"
	"time"
)

type Mills struct {
	time.Time
}

func (m *Mills) UnmarshalJSON(bytes []byte) error {
	var mills int64
	if err := json.Unmarshal(bytes, &mills); err != nil {
		return err
	}
	m.Time = time.Unix(int64(mills/1000), 0)
	return nil
}

func (m *Mills) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Unix())
}
