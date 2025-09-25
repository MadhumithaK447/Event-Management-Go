package handler

import (
	"eventgoapp/kk"
	"eventgoapp/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RsvpHandler interface {
	RegisterToEvent(ctx *gin.Context)
	GetAttendeesByEventID(ctx *gin.Context)
	SubscribeNotifications(ctx *gin.Context)
}

type rsvpHandler struct {
	service  service.RsvpService
	Producer *kk.KafkaProducer
}

func NewRsvpHandler(service service.RsvpService) RsvpHandler {
	return &rsvpHandler{service: service}
}

func (h *rsvpHandler) RegisterToEvent(ctx *gin.Context) {
	eventId, err := strconv.Atoi(ctx.Param("event_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	attendeeId, err := strconv.Atoi(ctx.Param("attendee_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid attendee ID"})
		return
	}

	err = h.service.RegisterToEvent(eventId, attendeeId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully registered to event"})
}

func (h *rsvpHandler) GetAttendeesByEventID(ctx *gin.Context) {
	eventID, err := strconv.Atoi(ctx.Param("event_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID"})
		return
	}

	attendees, err := h.service.GetAttendeesByEventID(eventID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, attendees)
}

func (h *rsvpHandler) SubscribeNotifications(ctx *gin.Context) {
	// Set headers for SSE
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")

	// Flush immediately
	ctx.Writer.Flush()

	// Add this client to producer (or a global clients slice)
	h.Producer.AddClient(ctx.Writer)

	// Keep connection open
	notify := ctx.Writer.CloseNotify()
	<-notify

	// Remove client when disconnected
	h.Producer.RemoveClient(ctx.Writer)
}
