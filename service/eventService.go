package service

import (
	models "eventgoapp/model"
	"eventgoapp/repository"
	"fmt"
)

type EventService interface {
	AddEvent(event models.Event) error
	GetAllEvents() ([]models.Event, error)
	GetEventByID(id int) (models.Event, error)
	UpdateEvent(id int, event models.Event) error
	DeleteEvent(id int) error
}

type eventService struct {
	repo repository.EventRepository
}

func NewEventService(repo repository.EventRepository) EventService {
	return &eventService{repo: repo}
}

func (s *eventService) AddEvent(event models.Event) error {
	if err := event.ValidateEvent(); err != nil {
		return err
	}
	return s.repo.AddEvent(event)
}

func (s *eventService) GetAllEvents() ([]models.Event, error) {
	return s.repo.GetAllEvents()
}

func (s *eventService) GetEventByID(id int) (models.Event, error) {
	exists, err := s.repo.Exists(id)
	if err != nil {
		return models.Event{},err
	}
	if !exists {
		return models.Event{},fmt.Errorf("event with id %d does not exist", id)
	}
	return s.repo.GetEventByID(id)
}

func (s *eventService) UpdateEvent(id int, event models.Event) error {
	exists, err := s.repo.Exists(id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("event with id %d does not exist", id)
	}
	return s.repo.UpdateEvent(id, event)
}

func (s *eventService) DeleteEvent(id int) error {
	exists, err := s.repo.Exists(id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("event with id %d does not exist", id)
	}
	return s.repo.DeleteEvent(id)
}
