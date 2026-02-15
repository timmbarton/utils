package dates

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestDate_String(t *testing.T) {
	tests := []struct {
		name string
		d    Date
		want string
	}{
		{
			name: "обычная дата",
			d:    Date{2023, time.October, 10},
			want: "2023-10-10",
		},
		{
			name: "начало года",
			d:    Date{2024, time.January, 1},
			want: "2024-01-01",
		},
		{
			name: "конец года",
			d:    Date{2022, time.December, 31},
			want: "2022-12-31",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.String(); got != tt.want {
				t.Errorf("Date.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromTime(t *testing.T) {
	loc, _ := time.LoadLocation("UTC")
	SetDefaultLocation(loc)
	defer SetDefaultLocation(time.FixedZone("Europe/Moscow", 60*60*3))

	tests := []struct {
		name string
		t    time.Time
		want Date
	}{
		{
			name: "UTC время",
			t:    time.Date(2023, 10, 10, 15, 0, 0, 0, time.UTC),
			want: Date{2023, time.October, 10},
		},
		{
			name: "другая таймзона (должна нормализоваться в defaultLocation)",
			t:    time.Date(2023, 10, 10, 23, 0, 0, 0, time.FixedZone("Test", -5*3600)),
			want: Date{2023, time.October, 11}, // -5 UTC 23:00 is 04:00 UTC next day
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromTime(tt.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_Unix(t *testing.T) {
	loc := time.UTC
	SetDefaultLocation(loc)
	defer SetDefaultLocation(time.FixedZone("Europe/Moscow", 60*60*3))

	tests := []struct {
		name string
		d    Date
		want int64
	}{
		{
			name: "2023-10-10 UTC",
			d:    Date{2023, time.October, 10},
			want: time.Date(2023, 10, 10, 0, 0, 0, 0, time.UTC).Unix(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Unix(); got != tt.want {
				t.Errorf("Date.Unix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_Time(t *testing.T) {
	loc := time.UTC
	SetDefaultLocation(loc)
	defer SetDefaultLocation(time.FixedZone("Europe/Moscow", 60*60*3))

	tests := []struct {
		name string
		d    Date
		want time.Time
	}{
		{
			name: "2023-10-10",
			d:    Date{2023, time.October, 10},
			want: time.Date(2023, 10, 10, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Time(); !got.Equal(tt.want) {
				t.Errorf("Date.Time() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		d       Date
		want    []byte
		wantErr bool
	}{
		{
			name:    "успешный маршалинг",
			d:       Date{2023, time.October, 10},
			want:    []byte(`"2023-10-10"`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Date.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Date.MarshalJSON() = %s, want %s", string(got), string(tt.want))
			}
		})
	}
}

func TestDate_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    Date
		wantErr bool
	}{
		{
			name:    "правильный формат",
			data:    []byte(`"2023-10-10"`),
			want:    Date{2023, time.October, 10},
			wantErr: false,
		},
		{
			name:    "неправильный формат даты",
			data:    []byte(`"2023/10/10"`),
			wantErr: true,
		},
		{
			name:    "не строка",
			data:    []byte(`123`),
			wantErr: true,
		},
		{
			name:    "невалидный JSON",
			data:    []byte(`{"key":`),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d Date
			err := json.Unmarshal(tt.data, &d)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(d, tt.want) {
				t.Errorf("UnmarshalJSON() = %v, want %v", d, tt.want)
			}
		})
	}
}

func TestDate_MarshalText(t *testing.T) {
	tests := []struct {
		name    string
		d       Date
		want    []byte
		wantErr bool
	}{
		{
			name: "успешный маршалинг",
			d:    Date{2023, time.October, 10},
			want: []byte("2023-10-10"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.MarshalText()
			if (err != nil) != tt.wantErr {
				t.Errorf("Date.MarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Date.MarshalText() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestDate_UnmarshalText(t *testing.T) {
	tests := []struct {
		name    string
		text    []byte
		want    Date
		wantErr bool
	}{
		{
			name:    "правильный формат",
			text:    []byte("2023-10-10"),
			want:    Date{2023, time.October, 10},
			wantErr: false,
		},
		{
			name:    "пустая строка",
			text:    []byte(""),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d Date
			err := d.UnmarshalText(tt.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("Date.UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(d, tt.want) {
				t.Errorf("Date.UnmarshalText() = %v, want %v", d, tt.want)
			}
		})
	}
}

func TestDate_Scan(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name    string
		value   any
		want    Date
		wantErr bool
	}{
		{
			name:    "time.Time",
			value:   now,
			want:    FromTime(now),
			wantErr: false,
		},
		{
			name:    "string",
			value:   "2023-10-10",
			want:    Date{2023, 10, 10},
			wantErr: false,
		},
		{
			name:    "nil",
			value:   nil,
			want:    Date{},
			wantErr: false,
		},
		{
			name:    "unsupported type (int)",
			value:   123,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d Date
			err := d.Scan(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Date.Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(d, tt.want) {
				t.Errorf("Date.Scan() = %v, want %v", d, tt.want)
			}
		})
	}
}

func TestDate_Value(t *testing.T) {
	tests := []struct {
		name    string
		d       Date
		want    any
		wantErr bool
	}{
		{
			name: "валидная дата",
			d:    Date{2023, time.October, 10},
			want: []byte("2023-10-10"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Date.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Date.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDateRange(t *testing.T) {
	tests := []struct {
		name  string
		start Date
		end   Date
		want  []Date
	}{
		{
			name:  "диапазон из 3 дней",
			start: Date{2023, time.October, 1},
			end:   Date{2023, time.October, 3},
			want: []Date{
				{2023, time.October, 1},
				{2023, time.October, 2},
				{2023, time.October, 3},
			},
		},
		{
			name:  "один день",
			start: Date{2023, time.October, 1},
			end:   Date{2023, time.October, 1},
			want: []Date{
				{2023, time.October, 1},
			},
		},
		{
			name:  "пустой диапазон (конец раньше начала)",
			start: Date{2023, time.October, 2},
			end:   Date{2023, time.October, 1},
			want:  []Date{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DateRange(tt.start, tt.end); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DateRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromUnix(t *testing.T) {
	loc := time.UTC
	SetDefaultLocation(loc)
	defer SetDefaultLocation(time.FixedZone("Europe/Moscow", 60*60*3))

	ts := time.Date(2023, 10, 10, 0, 0, 0, 0, time.UTC).Unix()
	want := Date{2023, time.October, 10}

	if got := FromUnix(ts); !reflect.DeepEqual(got, want) {
		t.Errorf("FromUnix() = %v, want %v", got, want)
	}
}

func TestFromUnixMilli(t *testing.T) {
	loc := time.UTC
	SetDefaultLocation(loc)
	defer SetDefaultLocation(time.FixedZone("Europe/Moscow", 60*60*3))

	ts := time.Date(2023, 10, 10, 0, 0, 0, 0, time.UTC).UnixMilli()
	want := Date{2023, time.October, 10}

	if got := FromUnixMilli(ts); !reflect.DeepEqual(got, want) {
		t.Errorf("FromUnixMilli() = %v, want %v", got, want)
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		y, m, d int
		want    Date
		wantErr bool
	}{
		{"валидная дата", 2023, 10, 10, Date{2023, 10, 10}, false},
		{"невалидный день", 2023, 10, 32, Date{}, true},
		{"невалидный месяц", 2023, 13, 10, Date{}, true},
		{"високосный год", 2024, 2, 29, Date{2024, 2, 29}, false},
		{"невисокосный год", 2023, 2, 29, Date{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.y, time.Month(tt.m), tt.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComparison(t *testing.T) {
	d1 := Date{2023, 10, 10}
	d2 := Date{2023, 10, 11}
	d3 := Date{2023, 10, 10}

	if !d1.Equal(d3) {
		t.Errorf("Equal failed")
	}
	if d1.Equal(d2) {
		t.Errorf("Equal failed")
	}
	if !d1.Before(d2) {
		t.Errorf("Before failed")
	}
	if d2.Before(d1) {
		t.Errorf("Before failed")
	}
	if !d2.After(d1) {
		t.Errorf("After failed")
	}
	if d1.After(d2) {
		t.Errorf("After failed")
	}
}

func TestNullJSON(t *testing.T) {
	t.Run("marshal empty", func(t *testing.T) {
		d := Date{}

		got, err := json.Marshal(d)
		if err != nil {
			t.Fatal(err)
		}

		want := "\"0001-01-01\""
		if string(got) != want {
			t.Errorf("Marshal(zero) = %v, want %v", string(got), want)
		}
	})
}
