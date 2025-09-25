package handler

import (
	"encoding/json"
	"eventgoapp/kk"
	"eventgoapp/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NotificationHandler interface {
	StreamNotifications(ctx *gin.Context)
}

type notificationHandler struct {
	service  service.NotificationService
	producer *kk.KafkaProducer
	consumer *kk.KafkaConsumer
}

func NewNotificationHandler(service service.NotificationService, producer *kk.KafkaProducer, consumer *kk.KafkaConsumer) NotificationHandler {
	return &notificationHandler{service: service, producer: producer, consumer: consumer}
}

func (h *notificationHandler) StreamNotifications(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")

	flusher, ok := ctx.Writer.(http.Flusher)
	if !ok {
		ctx.String(http.StatusInternalServerError, "Streaming unsupported!")
		return
	}
	reader := h.consumer.Reader()
	go func() {
		for {
			msg, err := reader.ReadMessage(ctx) // blocking call
			if err != nil {
				log.Println("Kafka read error:", err)
				break
			}
			
			var parsed map[string]interface{}
			if err := json.Unmarshal(msg.Value, &parsed); err == nil {
				data, _ := json.Marshal(parsed)
				_, _ = ctx.Writer.Write([]byte("data: " + string(data) + "\n\n"))
			} else {
				// If it's plain text
				_, _ = ctx.Writer.Write([]byte("data: " + string(msg.Value) + "\n\n"))
			}

			flusher.Flush()
		}
	}()
	<-ctx.Done()

}
