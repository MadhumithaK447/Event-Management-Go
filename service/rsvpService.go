package service

import (
	models "eventgoapp/model"
	"eventgoapp/repository"
)

type RsvpService interface {
	RegisterToEvent(eventId int, attendeeId int) error
	GetAttendeesByEventID(eventID int) ([]models.Attendee, error)
}

type rsvpService struct {
	repo         repository.RsvpRepository
	attendeeRepo repository.AttendeeRepository
	eventRepo    repository.EventRepository
}

func NewRsvpService(repo repository.RsvpRepository, attendeeRepo repository.AttendeeRepository, eventRepo repository.EventRepository) RsvpService {
	return &rsvpService{repo: repo, attendeeRepo: attendeeRepo, eventRepo: eventRepo}
}

func (s *rsvpService) RegisterToEvent(eventId int, attendeeId int) error {
	return s.repo.RegisterToEvent(eventId, attendeeId)
}

func (s *rsvpService) GetAttendeesByEventID(eventID int) ([]models.Attendee, error) {
	return s.repo.GetAttendeesByEventID(eventID)
}
