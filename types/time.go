package types

import (
	"time"
)

type Time time.Time

const (
	timeFormat = "2006-01-02 15:04:05"
	dateFormat = "2006-01-02"
)

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	var now time.Time
	dataLen := len(data)

	if dataLen == len(dateFormat) + 2 {
		now, err = time.ParseInLocation(`"`+dateFormat+`"`, string(data), time.Local)
	} else if dataLen == len(timeFormat) + 2 {
		now, err = time.ParseInLocation(`"`+timeFormat+`"`, string(data), time.Local)
	} else {
		now, err = time.Parse(time.RFC3339, string(data))
	}

	*t = Time(now)
	return
}

func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timeFormat) + 2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, timeFormat)
	b = append(b, '"')
	return b, nil
}

func (t Time) String() string {
	return time.Time(t).Format(timeFormat)
}

func (t *Time) ToDate(plusDays int) *time.Time  {
	date := time.Time(*t)
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	if plusDays != 0 {
		date = date.AddDate(0, 0, plusDays)
	}

	return &date
}