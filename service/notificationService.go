package service

import (
	"eventgoapp/kk"
	"eventgoapp/repository"
	"log"
	"time"
)

type NotificationService interface {
	NotifyTodayEvents()
}
type notificationService struct {
	eventRepo repository.EventRepository
	rsvpRepo  repository.RsvpRepository
	producer  *kk.KafkaProducer
}

func NewNotificationService(eventRepo repository.EventRepository, rsvpRepo repository.RsvpRepository, producer *kk.KafkaProducer) NotificationService {
	return &notificationService{
		eventRepo: eventRepo,
		rsvpRepo:  rsvpRepo,
		producer:  producer,
	}
}

// Check for events today and notify attendees
func (ns *notificationService) NotifyTodayEvents() {
	events, err := ns.eventRepo.GetAllEvents()
	if err != nil {
		log.Println("Failed to fetch events:", err)
		return
	}

	today := time.Now().Truncate(24 * time.Hour)

	for _, ev := range events {
		eventDate, _ := time.Parse("2006-01-02", ev.EventDate)

		if eventDate.Equal(today) {
			attendees, err := ns.rsvpRepo.GetAttendeesByEventID(ev.Id)
			if err != nil {
				log.Println("Failed to fetch attendees:", err)
				continue
			}

			for _, attendee := range attendees {

				ns.producer.Publish("events_topic", map[string]interface{}{
					"attendee_id": attendee.Id,
					"event_id":    ev.Id,
					"message":     "Reminder: Your event \"" + ev.Title + "\" is today!",
					"timestamp":   time.Now(),
				})
			}
		}
	}
}
