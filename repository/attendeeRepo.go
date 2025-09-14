package repository

import (
	"eventgoapp/model"

	"github.com/jmoiron/sqlx"
)

type AttendeeRepository interface {
	AddAttendee(attendee model.Attendee) error
	AttendeeExists(id int) (bool, error)
}

type attendeeRepo struct {
	db *sqlx.DB
}

func NewAttendeeRepository(db *sqlx.DB) AttendeeRepository {
	return &attendeeRepo{db: db}
}

func (r *attendeeRepo) AddAttendee(attendee model.Attendee) error {
	query := "insert into attendees (name,phone_no,email) values (:name,:phone_no,:email)"
	_, err := r.db.NamedExec(query,attendee)
	if err != nil {
		return err
	}
	return nil
}

func (r *attendeeRepo) AttendeeExists(id int) (bool, error) {
    var exists bool
    query := `SELECT EXISTS (SELECT 1 FROM attendees WHERE id=$1)`
    err := r.db.Get(&exists, query, id)
    return exists, err
}
