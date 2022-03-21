package model

import (
	"database/sql/driver"
	"fmt"
	"io"
	"time"
)

// Date represents date without time
type Date time.Time

// FromTime converts from time.Time
func (date Date) FromTime(t time.Time) Date {
	year, month, day := t.Date()
	return Date(time.Date(year, month, day, 0, 0, 0, 0, time.UTC))
}

// MarshalJSON marshals into JSON
func (date Date) MarshalJSON() ([]byte, error) {
	string := fmt.Sprintf("\"%s\"", time.Time(date).UTC().Format("2006-01-02"))
	return []byte(string), nil
}

// UnmarshalJSON unmarshalls JSON
func (date *Date) UnmarshalJSON(data []byte) error {
	if len(data) != 12 {
		return fmt.Errorf("cannot parse %s as YYYY-MM-DD", string(data))
	}
	t, err := time.Parse("2006-01-02", string(data[1:11]))
	*date = Date(t)
	return err
}

// Time returns underlying time.Time
func (date *Date) Time() time.Time {
	return time.Time(*date)
}

// String returns string representation
func (date Date) String() string {
	return time.Time(date).UTC().Format("2006-01-02")
}

// Equal checks if two dates are equal
func (date *Date) Equal(d Date) bool {
	return date.Time().Equal(d.Time())
}

// Value implements driver.Valuer interface
func (date Date) Value() (driver.Value, error) {
	return date.Time(), nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (date Date) MarshalGQL(w io.Writer) {
	w.Write([]byte(`"` + date.String() + `"`))
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (date *Date) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("must be a string")
	}
	if len(str) != 12 {
		return fmt.Errorf("cannot parse %s as YYYY-MM-DD", str)
	}
	t, err := time.Parse("2006-01-02", string(str[1:11]))
	*date = Date(t)
	return err
}
