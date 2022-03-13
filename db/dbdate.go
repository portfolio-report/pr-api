package db

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type DbDate time.Time

func (date DbDate) FromTime(t time.Time) DbDate {
	year, month, day := t.Date()
	return DbDate(time.Date(year, month, day, 0, 0, 0, 0, time.UTC))
}

func (date DbDate) MarshalJSON() ([]byte, error) {
	string := fmt.Sprintf("\"%s\"", time.Time(date).UTC().Format("2006-01-02"))
	return []byte(string), nil
}

func (date *DbDate) UnmarshalJSON(data []byte) error {
	if len(data) != 12 {
		return fmt.Errorf("cannot parse %s as YYY-MM-DD", string(data))
	}
	t, err := time.Parse("2006-01-02", string(data[1:11]))
	*date = DbDate(t)
	return err
}

func (date *DbDate) Time() time.Time {
	return time.Time(*date)
}

func (date *DbDate) TimePtr() *time.Time {
	t := time.Time(*date)
	return &t
}

func (date DbDate) String() string {
	return time.Time(date).UTC().Format("2006-01-02")
}

func (date *DbDate) Equal(d DbDate) bool {
	return date.Time().Equal(d.Time())
}

func (date DbDate) Value() (driver.Value, error) {
	return date.Time(), nil
}
