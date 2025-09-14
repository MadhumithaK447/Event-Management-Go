package handler

import (
	"eventgoapp/model"
	"eventgoapp/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AttendeeHandler interface {
	AddAttendee(ctx *gin.Context)
}

type attendeeHandler struct {
	service service.AttendeeService
}

func NewAttendeeHandler(service service.AttendeeService) AttendeeHandler {
	return &attendeeHandler{service: service}
}

func (h *attendeeHandler) AddAttendee(ctx *gin.Context) {
	var attendee model.Attendee
	if err := ctx.ShouldBindJSON(&attendee); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddAttendee(attendee); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "Attendee Added Successfully"})
}