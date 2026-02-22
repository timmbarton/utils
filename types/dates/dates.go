package dates

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type Date struct {
	year  int
	month time.Month
	day   int
}

var (
	defaultLocation = time.FixedZone("Europe/Moscow", 60*60*3)
	locMu           sync.RWMutex
)

//goland:noinspection ALL
func SetDefaultLocation(loc *time.Location) {
	locMu.Lock()
	defer locMu.Unlock()
	defaultLocation = loc
}

func getDefaultLocation() *time.Location {
	locMu.RLock()
	defer locMu.RUnlock()
	return defaultLocation
}

func (d Date) Unix() int64 {
	return time.Date(d.year, d.month, d.day, 0, 0, 0, 0, getDefaultLocation()).Unix()
}
func (d Date) Time() time.Time {
	return time.Date(d.year, d.month, d.day, 0, 0, 0, 0, getDefaultLocation())
}

func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", max(d.year, 1), max(d.month, 1), max(d.day, 1))
}

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Date) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*d = Date{}
		return nil
	}

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
	if len(text) == 0 {
		return fmt.Errorf("empty date string")
	}

	t, err := time.Parse(time.DateOnly, string(text))
	if err != nil {
		return err
	}

	*d = FromTime(t)

	return nil
}

func (d *Date) Scan(value any) error {
	switch v := value.(type) {
	case time.Time:
		*d = FromTime(v)
		return nil
	case string:
		return d.UnmarshalText([]byte(v))
	case []byte:
		return d.UnmarshalText(v)
	case nil:
		return nil
	}

	return fmt.Errorf("unsupported type for (d *Date) Scan(): %T", value)
}
func (d Date) Value() (val driver.Value, err error) {
	return d.MarshalText()
}

// New создает новый объект Date. Если дата невалидна, возвращает ошибку.
func New(year int, month time.Month, day int) (Date, error) {
	d := Date{year: year, month: month, day: day}
	t := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	if t.Year() != year || t.Month() != month || t.Day() != day {
		return Date{}, fmt.Errorf("invalid date: %d-%02d-%02d", year, month, day)
	}
	return d, nil
}

func (d Date) Equal(other Date) bool {
	return d.year == other.year && d.month == other.month && d.day == other.day
}

func (d Date) After(other Date) bool {
	if d.year != other.year {
		return d.year > other.year
	}
	if d.month != other.month {
		return d.month > other.month
	}
	return d.day > other.day
}

func (d Date) Before(other Date) bool {
	if d.year != other.year {
		return d.year < other.year
	}
	if d.month != other.month {
		return d.month < other.month
	}
	return d.day < other.day
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
	t = t.In(getDefaultLocation())
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

//goland:noinspection ALL
func Parse(s string) (Date, error) {
	d := Date{}
	err := d.UnmarshalText([]byte(s))
	return d, err
}
