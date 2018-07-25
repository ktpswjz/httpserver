package types

import (
	"time"
)

type Time time.Time

const (
	timeFormart = "2006-01-02 15:04:05"
	dateFormart = "2006-01-02"
)

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	format := timeFormart
	if len(data) == len(dateFormart) + 2 {
		format = dateFormart
	}

	now, err := time.ParseInLocation(`"`+format+`"`, string(data), time.Local)
	*t = Time(now)
	return
}

func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timeFormart)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, timeFormart)
	b = append(b, '"')
	return b, nil
}

func (t Time) String() string {
	return time.Time(t).Format(timeFormart)
}
