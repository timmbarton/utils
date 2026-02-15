package dates

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Date struct {
	year  int
	month time.Month
	day   int
}

var defaultLocation = time.FixedZone("Europe/Moscow", 60*60*3)

//goland:noinspection ALL
func SetDefaultLocation(loc *time.Location) {
	defaultLocation = loc
}

func (d Date) Unix() int64 {
	return time.Date(d.year, d.month, d.day, 0, 0, 0, 0, defaultLocation).Unix()
}
func (d Date) Time() time.Time {
	return time.Date(d.year, d.month, d.day, 0, 0, 0, 0, defaultLocation)
}

func (d Date) String() string {
	return fmt.Sprintf("%d-%02d-%02d", d.year, d.month, d.day)
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
	t, err := time.Parse(time.DateOnly, string(text))
	if err != nil {
		return err
	}

	*d = FromTime(t)

	return nil
}

func (d *Date) Scan(value any) error {
	if t, ok := value.(time.Time); ok {
		*d = FromTime(t)
		return nil
	}

	return fmt.Errorf("unsupported type for (d *Date) Scan(): %T", value)
}
func (d Date) Value() (val driver.Value, err error) {
	return d.MarshalText()
}

//goland:noinspection ALL
func DateRange(start, end Date) []Date {
	dateRange := make([]Date, 0)

	for d := start.Time(); d.After(end.Time()) == false; d = d.AddDate(0, 0, 1) {
		dateRange = append(dateRange, FromTime(d))
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
	t = t.In(defaultLocation)
	return Date{
		year:  t.Year(),
		month: t.Month(),
		day:   t.Day(),
	}
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
