package model

import (
	"errors"
	"strings"
	"time"
)

type Event struct {
	Id          int    `db:"id" json:"id"`
	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`
	EventDate   string `db:"event_date" json:"event_date"`
}

func (e *Event) ValidateEvent() error {
	if e.Title == "" {
		return errors.New("Event name is required")
	}
	if e.Title == " " {
		return errors.New("Event name is Invalid")
	}
	if e.Description == "" {
		return errors.New("Event description is required")
	}
	return nil
}

type DateOnly time.Time

func (d DateOnly) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(d).Format("2006-01-02") + `"`), nil
}

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*d = DateOnly(t)
	return nil
}
