package ffmpeg

import (
	"fmt"
	"strings"
	"time"
)

var currentLocation = time.Now().Location()

const railsTimeLayout = "2006-01-02 15:04:05 MST"

type JSONTime struct {
	time.Time
}

func (jt *JSONTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		jt.Time = time.Time{}
		return nil
	}

	// #731 - returning an error here causes the entire JSON parse to fail for ffprobe.
	jt.Time, _ = ParseDateStringAsTime(s)
	return nil
}

func (jt *JSONTime) MarshalJSON() ([]byte, error) {
	if jt.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", jt.Time.Format(time.RFC3339))), nil
}

func (jt JSONTime) GetTime() time.Time {
	if currentLocation != nil {
		if jt.IsZero() {
			return time.Now().In(currentLocation)
		} else {
			return jt.Time.In(currentLocation)
		}
	} else {
		if jt.IsZero() {
			return time.Now()
		} else {
			return jt.Time
		}
	}
}

func ParseDateStringAsTime(dateString string) (time.Time, error) {
	// https://stackoverflow.com/a/20234207 WTF?

	t, e := time.Parse(time.RFC3339, dateString)
	if e == nil {
		return t, nil
	}

	t, e = time.Parse("2006-01-02", dateString)
	if e == nil {
		return t, nil
	}

	t, e = time.Parse("2006-01-02 15:04:05", dateString)
	if e == nil {
		return t, nil
	}

	t, e = time.Parse(railsTimeLayout, dateString)
	if e == nil {
		return t, nil
	}

	return time.Time{}, fmt.Errorf("ParseDateStringAsTime failed: dateString <%s>", dateString)
}
