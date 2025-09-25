package model

import (
	"errors"
	"time"
)

type Event struct {
	Id          int       `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Date        time.Time `db:"date" json:"date"`
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
