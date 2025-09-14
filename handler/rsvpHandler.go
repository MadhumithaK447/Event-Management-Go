package handler

import (
	"eventgoapp/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RsvpHandler interface {
	RegisterToEvent(ctx *gin.Context)
}

type rsvpHandler struct {
	service service.RsvpService
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
