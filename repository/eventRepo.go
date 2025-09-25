package repository

import (
	models "eventgoapp/model"

	"github.com/jmoiron/sqlx"
)

type EventRepository interface {
	AddEvent(event models.Event) error
	GetAllEvents() ([]models.Event, error)
	GetEventByID(id int) (models.Event, error)
	UpdateEvent(id int, event models.Event) error
	Exists(id int) (bool, error)
	DeleteEvent(id int) error
}

type eventRepo struct {
	db *sqlx.DB
}

func NewEventRepository(db *sqlx.DB) EventRepository {
	return &eventRepo{db: db}
}

func (r *eventRepo) AddEvent(event models.Event) error {
	query := "INSERT INTO events (title, description, event_date) VALUES (:title, :description, :event_date)"
	//_, err := r.db.Exec(query, event.Title, event.Description)
	_, err := r.db.NamedExec(query, event)
	if err != nil {
		return err
	}
	return nil
}

func (r *eventRepo) GetAllEvents() ([]models.Event, error) {
	events := []models.Event{}
	query := "SELECT * FROM EVENTS"
	err := r.db.Select(&events, query)
	return events, err
}

func (r *eventRepo) GetEventByID(id int) (models.Event, error) {
	query := `SELECT * from events where id=$1`
	var event models.Event
	err := r.db.Get(&event, query, id)
	if err != nil {
		return event, err
	}
	return event, nil
}

func (r *eventRepo) UpdateEvent(id int, event models.Event) error {
	query := `UPDATE events SET title=:title, description=:description, event_date=:event_date where id=:id`
	args := map[string]interface{}{
		"id":          id,
		"title":       event.Title,
		"description": event.Description,
		"event_date":        event.EventDate,
	}
	_, err := r.db.NamedExec(query, args)
	if err != nil {
		return err
	}
	return nil
}

func (r *eventRepo) Exists(id int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM events WHERE id=$1)`
	err := r.db.Get(&exists, query, id)
	return exists, err
}

func (r *eventRepo) DeleteEvent(id int) error {
	query := `DELETE FROM events WHERE id = :id`
	args := map[string]interface{}{
		"id": id,
	}
	_, err := r.db.NamedExec(query, args)
	if err != nil {
		return err
	}
	return nil
}
