package service

import (
	"eventgoapp/model"
	"eventgoapp/repository"
)

type AttendeeService interface {
	AddAttendee (attendee model.Attendee) error
}

type attendeeService struct {
	repo repository.AttendeeRepository
}

func NewAttendeeService(repo repository.AttendeeRepository) AttendeeService{
	return &attendeeService{repo: repo}
}

func (s *attendeeService) AddAttendee(attendee model.Attendee) error {
	return s.repo.AddAttendee(attendee)
}