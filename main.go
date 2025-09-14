package main

import (
	"eventgoapp/db"
	"eventgoapp/handler"
	"eventgoapp/repository"
	"eventgoapp/service"

	"github.com/gin-gonic/gin"
)

func main() {
	database := db.ConnectDB()
	defer database.Close()

	var eventRepository repository.EventRepository = repository.NewEventRepository(database)
	var eventService service.EventService = service.NewEventService(eventRepository)
	var eventHandler handler.EventHandler = handler.NewEventHandler(eventService)

	attendeeRepo := repository.NewAttendeeRepository(database)
	attendeeService := service.NewAttendeeService(attendeeRepo)
	attendeeHandler := handler.NewAttendeeHandler(attendeeService)

	rsvpRepo := repository.NewRsvpRepository(database, attendeeRepo, eventRepository)
	rsvpService := service.NewRsvpService(rsvpRepo, attendeeRepo, eventRepository)
	rsvpHandler := handler.NewRsvpHandler(rsvpService)
	server := gin.Default()

	server.POST("/events", eventHandler.AddEvent)
	server.GET("/events", eventHandler.GetAllEvents)
	server.GET("/events/:id", eventHandler.GetEventByID)
	server.PATCH("/events/:id", eventHandler.UpdateEvent)
	server.DELETE("/events/:id", eventHandler.DeleteEvent)

	server.POST("/attendees", attendeeHandler.AddAttendee)
	
	server.POST("/rsvp/:event_id/:attendee_id", rsvpHandler.RegisterToEvent)

	server.Run(":8080")
}
