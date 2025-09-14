package handler

import (
	"eventgoapp/model"
	"eventgoapp/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EventHandler interface {
	AddEvent(ctx *gin.Context)
	GetAllEvents(ctx *gin.Context)
	GetEventByID(ctx *gin.Context)
	UpdateEvent(ctx *gin.Context)
	DeleteEvent(ctx *gin.Context)
}

type eventHandler struct {
	service service.EventService
}

func NewEventHandler(service service.EventService) EventHandler {
	return &eventHandler{service: service}
}

func (h *eventHandler) AddEvent(ctx *gin.Context) {
	var event model.Event

	if err := ctx.ShouldBindJSON(&event); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := event.ValidateEvent(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddEvent(event); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Event created successfully"})
}

func (h *eventHandler) GetAllEvents(ctx *gin.Context) {
	events, err := h.service.GetAllEvents()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, events)
}

func (h *eventHandler) GetEventByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	// strconv.Atoi(...) converts it from string to int
	event, err := h.service.GetEventByID(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, event)
}

func (h *eventHandler) UpdateEvent(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var event model.Event
	event.Id=id
	if err := ctx.ShouldBindJSON(&event); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateEvent(id,event); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Event updated successfully"})
}

func (h *eventHandler) DeleteEvent(ctx *gin.Context) {
    id, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    if err := h.service.DeleteEvent(id); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}

