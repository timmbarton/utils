package dates

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Date time.Time

func (d Date) Unix() int64 {
	return time.Time(d).Unix()
}

func (d Date) String() string {
	return time.Time(d).Format(time.DateOnly)
}

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}
func (d *Date) UnmarshalJSON(data []byte) error {
	dateStr := ""

	err := json.Unmarshal(data, &dateStr)
	if err != nil {
		return err
	}

	return d.UnmarshalText([]byte(dateStr))
}

func (d Date) MarshalText() (text []byte, err error) {
	return []byte(d.String()), nil
}
func (d *Date) UnmarshalText(text []byte) error {
	timeTime, err := time.Parse(time.DateOnly, string(text))
	if err != nil {
		return err
	}

	*d = Date(timeTime)

	return nil
}

func (d *Date) Scan(value any) error {
	if timeTime, ok := value.(time.Time); ok {
		*d = Date(timeTime)

		return nil
	} else {
		return errors.New(fmt.Sprintf("unsupported type for (d *Date) Scan()"))
	}
}
func (d Date) Value() (val driver.Value, err error) {
	return d.MarshalText()
}

//goland:noinspection ALL
func DateRange(start, end Date) []Date {
	dateRange := make([]Date, 0, (end.Unix()-start.Unix())/86400+1)
	for d := time.Time(start); d.After(time.Time(end)) == false; d = d.AddDate(0, 0, 1) {
		dateRange = append(dateRange, Date(d))
	}

	return dateRange
}

const day = 24 * time.Hour

//goland:noinspection ALL
func FromUnix(unixSeconds int64) Date {
	return FromTime(time.Unix(unixSeconds, 0))
}

//goland:noinspection ALL
func FromUnixMilli(unixMilliSeconds int64) Date {
	return FromTime(time.UnixMilli(unixMilliSeconds))
}

func FromTime(t time.Time) Date {
	return Date(t.Truncate(day))
}

//goland:noinspection ALL
func Today() Date {
	return FromTime(time.Now())
}

//goland:noinspection ALL
func Yesterday() Date {
	return FromTime(time.Now().Add(-day))
}

//goland:noinspection ALL
func Tomorrow() Date {
	return FromTime(time.Now().Add(day))
}
