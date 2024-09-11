package util

import (
	"testing"
	"time"
)

func TestParseDateTime(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    time.Time
	}{
		{
			name:    "ValidDateFormat",
			input:   "2024-01-01 00:00:00",
			wantErr: false,
			want:    time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:    "InvalidDateFormat",
			input:   "This is not a valid date format",
			wantErr: true,
		},
		{
			name:    "EmptyString",
			input:   "",
			wantErr: true,
		},
		{
			name:    "DateFormatWithoutTimeStamp",
			input:   "2024-01-01",
			wantErr: true,
		},
		{
			name:    "TimestampFormatWithoutDate",
			input:   "00:00:00",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDateTime(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDateTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !got.Equal(tt.want) {
				t.Errorf("ParseDateTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWibToUtc(t *testing.T) {
	cases := []struct {
		name     string
		wib      time.Time
		expected time.Time
	}{
		{
			name:     "midnight",
			wib:      time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2023, 12, 31, 17, 0, 0, 0, time.UTC),
		},
		{
			name:     "noon",
			wib:      time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			expected: time.Date(2024, 1, 1, 5, 0, 0, 0, time.UTC),
		},
		{
			name:     "endOfDay",
			wib:      time.Date(2024, 1, 1, 23, 59, 59, 999999999, time.UTC),
			expected: time.Date(2024, 1, 1, 16, 59, 59, 999999999, time.UTC),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			utc := WibToUtc(tc.wib)
			if !utc.Equal(tc.expected) {
				t.Errorf("WibToUtc(%v) = %v; expected %v", tc.wib, utc, tc.expected)
			}
		})
	}
}

func TestTimeToStandardUtc(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "ValidFormat",
			input:   "2000-01-01 12:00:00",
			want:    "2000-01-01T05:00:00+00:00", // assumes WibToUtc reduces 8 hours.
			wantErr: false,
		},
		{
			name:    "InvalidFormat",
			input:   "invalidFormat",
			wantErr: true,
		},
		{
			name:    "LeapYear",
			input:   "2020-03-01 04:59:59",
			want:    "2020-02-29T21:59:59+00:00", // assumes WibToUtc reduces 8 hours.
			wantErr: false,
		},
		{
			name:    "NonLeapYear",
			input:   "2019-02-29T23:59:59",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TimeToStandardUtc(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("TimeToStandardUtc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TimeToStandardUtc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeToStandard(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{"ValidTime", "2022-12-31 23:59:59", "2022-12-31T23:59:59+00:00", false},
		{"InvalidTime", "abcdef", "", true},
		{"EmptyInput", "", "", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := TimeToStandard(test.input)
			if (err != nil) != test.hasError {
				t.Errorf("error mismatch: got %v, want %v", err, test.hasError)
			}
			if err == nil && actual != test.expected {
				t.Errorf("result mismatch: got %s, want %s", actual, test.expected)
			}
		})
	}
}

func TestTimeConvert(t *testing.T) {
	tests := []struct {
		name string
		date string
		utc  bool
		want string
	}{
		{
			name: "ValidTimeInUTC",
			date: "2022-02-10 15:04:05",
			utc:  true,
			want: "2022-02-10T08:04:05+00:00",
		},
		{
			name: "ValidTimeNoUTC",
			date: "2022-02-10 15:04:05",
			utc:  false,
			want: "2022-02-10T15:04:05+00:00",
		},
		{
			name: "EmptyTime",
			date: "",
			utc:  true,
			want: "",
		},
		{
			name: "InvalidTime",
			date: "2022-22-22T99:99:99Z",
			utc:  true,
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TimeConvert(tt.date, tt.utc)
			if got != tt.want {
				t.Errorf("TimeConvert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDateConvert(t *testing.T) {
	var tests = []struct {
		input    string
		isUtc    bool
		expected string
	}{
		{"2022-01-01 00:00:00", true, "2021-12-31"},  // testing with UTC conversion
		{"2022-12-31 23:59:59", true, "2022-12-31"},  // testing with date and time at the edge
		{"2022-01-01 00:00:00", false, "2022-01-01"}, // testing without UTC conversion
		{"bad date format", true, ""},                // testing with bad date format
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			result := DateConvert(tt.input, tt.isUtc)
			if result != tt.expected {
				t.Errorf("DateConvert(%v, %v) = %v; want %v", tt.input, tt.isUtc, result, tt.expected)
			}
		})
	}
}
