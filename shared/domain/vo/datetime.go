package vo

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// DateTime represents a timestamp that works with both PostgreSQL and SQLite
type DateTime time.Time

// Now returns the current time as DateTime
func Now() DateTime {
	return DateTime(time.Now())
}

// ParseDateTime parses a string into DateTime
func ParseDateTime(s string) (DateTime, error) {
	if s == "" {
		return DateTime{}, nil
	}
	// Clean SQLite monotonic clock suffix (e.g., "m=+8.005836901")
	if idx := strings.Index(s, " m="); idx != -1 {
		s = s[:idx]
	}
	// Clean duplicate timezone: " -0500 -05" -> " -0500"
	if idx := strings.LastIndex(s, " -05"); idx > 10 && strings.HasSuffix(s, "-05") {
		s = s[:idx]
	}
	// Clean extra timezone parts (e.g., " -0500 -05" at end)
	// Use regex to find and fix the timezone duplication
	re := regexp.MustCompile(`(\d{2}:\d{2})(.*)-\d{2}$`)
	if matches := re.FindStringSubmatch(s); len(matches) > 2 {
		s = strings.TrimSuffix(s, matches[2])
	}

	// Try parsing with flexible decimal seconds using regex replacement
	// Match any number of decimal digits
	formats := []string{
		"2006-01-02 15:04:05.000000000 -0700",
		"2006-01-02 15:04:05.000000 -0700",
		"2006-01-02 15:04:05.00000 -0700",
		"2006-01-02 15:04:05.0000 -0700",
		"2006-01-02 15:04:05.000 -0700",
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05.999999999Z07:00",
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02 15:04:05.999999",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			return DateTime(t), nil
		}
	}

	// Last resort: try to parse just the date part and ignore time
	if len(s) >= 10 {
		datePart := s[:10]
		if t, err := time.Parse("2006-01-02", datePart); err == nil {
			return DateTime(t), nil
		}
	}

	return DateTime{}, fmt.Errorf("invalid datetime format: %s", s)
}

func (d DateTime) String() string {
	return time.Time(d).Format(time.RFC3339)
}

// Time returns the underlying time.Time
func (d DateTime) Time() time.Time {
	return time.Time(d)
}

// IsZero checks if the DateTime is zero
func (d DateTime) IsZero() bool {
	return time.Time(d).IsZero()
}

// MarshalJSON returns the JSON encoding of DateTime
func (d DateTime) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(time.Time(d).Format(time.RFC3339))
}

// UnmarshalJSON parses JSON bytes into a valid DateTime
func (d *DateTime) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == "" || s == "null" {
		*d = DateTime{}
		return nil
	}
	dt, err := ParseDateTime(s)
	if err != nil {
		return err
	}
	*d = dt
	return nil
}

// Value implements the driver.Valuer interface
func (d DateTime) Value() (driver.Value, error) {
	if d.IsZero() {
		return nil, nil
	}
	return time.Time(d), nil
}

// Scan implements the sql.Scanner interface
func (d *DateTime) Scan(value interface{}) error {
	if value == nil {
		*d = DateTime{}
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		*d = DateTime(v)
	case string:
		dt, err := ParseDateTime(v)
		if err != nil {
			return err
		}
		*d = dt
	case []byte:
		dt, err := ParseDateTime(string(v))
		if err != nil {
			return err
		}
		*d = dt
	default:
		return fmt.Errorf("cannot scan type %T into DateTime", value)
	}
	return nil
}