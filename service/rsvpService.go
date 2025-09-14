package service

import "eventgoapp/repository"

type RsvpService interface {
	RegisterToEvent(eventId int, attendeeId int) error
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
