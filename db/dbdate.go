package db

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// DbDate represents date without time
type DbDate time.Time

// FromTime converts time.Time to DbDate
func (date DbDate) FromTime(t time.Time) DbDate {
	year, month, day := t.Date()
	return DbDate(time.Date(year, month, day, 0, 0, 0, 0, time.UTC))
}

// MarshalJSON marshals DbDate into JSON
func (date DbDate) MarshalJSON() ([]byte, error) {
	string := fmt.Sprintf("\"%s\"", time.Time(date).UTC().Format("2006-01-02"))
	return []byte(string), nil
}

// UnmarshalJSON unmarshalls JSON into DbDate
func (date *DbDate) UnmarshalJSON(data []byte) error {
	if len(data) != 12 {
		return fmt.Errorf("cannot parse %s as YYY-MM-DD", string(data))
	}
	t, err := time.Parse("2006-01-02", string(data[1:11]))
	*date = DbDate(t)
	return err
}

// Time returns underlying time.Time
func (date *DbDate) Time() time.Time {
	return time.Time(*date)
}

// String returns string representation
func (date DbDate) String() string {
	return time.Time(date).UTC().Format("2006-01-02")
}

// Equal checks if two dates are equal
func (date *DbDate) Equal(d DbDate) bool {
	return date.Time().Equal(d.Time())
}

// Value implements driver.Valuer interface
func (date DbDate) Value() (driver.Value, error) {
	return date.Time(), nil
}
