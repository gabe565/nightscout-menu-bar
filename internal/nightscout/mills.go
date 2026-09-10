package nightscout

import (
	"encoding/json"
	"strings"
	"time"
)

type Mills struct {
	time.Time
}

func (m *Mills) UnmarshalJSON(bytes []byte) error {
	// Some Nightscout instances return mills as a float
	var mills float64
	if err := json.Unmarshal(bytes, &mills); err != nil {
		return err
	}
	m.Time = time.UnixMicro(int64(mills * 1000))
	return nil
}

func (m *Mills) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.UnixMilli())
}

func (m *Mills) Relative(round bool) string {
	if m.Unix() == 0 {
		return ""
	}

	// Drop resolution to minutes
	duration := time.Since(m.Time)
	if round {
		duration = duration.Round(time.Minute)
	} else {
		duration = duration.Truncate(time.Minute)
	}

	str := duration.String()
	str = strings.TrimSuffix(str, "0s")
	if str == "" {
		str = "0m"
	}
	return str
}
