package repository

import (
	"context"
	"encoding/json"
	"eventgoapp/cache"
	models "eventgoapp/model"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type RsvpRepository interface {
	RegisterToEvent(eventId int, attendeeId int) error
	GetAttendeesByEventID(eventID int) ([]models.Attendee, error)
}

type rsvpRepo struct {
	db           *sqlx.DB
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

	cache.Rdb.Del(cache.Ctx, fmt.Sprintf("event:%d:attendees", eventId))

	return nil
}

func (r *rsvpRepo) GetAttendeesByEventID(eventID int) ([]models.Attendee, error) {
	cacheKey := fmt.Sprintf("event:%d:attendees", eventID)

	val, err := cache.Rdb.Get(cache.Ctx, cacheKey).Result()
	if err == nil {
		var attendees []models.Attendee
		if json.Unmarshal([]byte(val), &attendees) == nil {
			fmt.Println("✅ Returning from Redis cache")
			return attendees, nil
		}
	}
	query := `
	SELECT a.id, a.name, a.phone_no, a.email
	FROM attendees a
	INNER JOIN rsvp r ON r.attendee_id = a.id
	WHERE r.event_id = $1
	`
	attendees := []models.Attendee{}
	//query := `SELECT * FROM rsvp WHERE r.event_id = $1`
	// err := r.db.Select(&attendees, query, eventID)
	// return attendees, err

	err = r.db.Select(&attendees, query, eventID)
	if err != nil {
		return nil, err
	}

	// Save in Redis
	data, _ := json.Marshal(attendees)
	cache.Rdb.Set(context.Background(), cacheKey, data, time.Minute*10)

	fmt.Println("📦 Saved attendees in Redis cache")
	return attendees, nil
}
