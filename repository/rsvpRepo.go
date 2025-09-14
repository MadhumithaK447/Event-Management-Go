package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type RsvpRepository interface {
	RegisterToEvent(eventId int, attendeeId int) error
}

type rsvpRepo struct {
	db *sqlx.DB
	eventRepo    EventRepository
    attendeeRepo AttendeeRepository
}

func NewRsvpRepository(db *sqlx.DB, attendeeRepo AttendeeRepository, eventRepo EventRepository) RsvpRepository {
	return &rsvpRepo{db: db, attendeeRepo: attendeeRepo, eventRepo: eventRepo}
}

func (r *rsvpRepo) RegisterToEvent(eventId int, attendeeId int) error {
	eventExists, err := r.eventRepo.Exists(eventId)
	if err != nil {
		return err
	}
	if !eventExists {
		return fmt.Errorf("event %d does not exist", eventId)
	}

	attendeeExists, err := r.attendeeRepo.AttendeeExists(attendeeId)
	if err != nil {
		return err
	}
	if !attendeeExists {
		return fmt.Errorf("attendee %d does not exist", attendeeId)
	}

	query := "INSERT INTO rsvp (event_id, attendee_id) VALUES (:event_id, :attendee_id)"
	_, err = r.db.NamedExec(query, map[string]interface{}{
		"event_id":    eventId,
		"attendee_id": attendeeId,
	})
	if err != nil {
		return fmt.Errorf("failed to register to event: %w", err)
	}
	return nil
}
