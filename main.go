package main

import (
	"context"
	"eventgoapp/cache"
	"eventgoapp/db"
	"eventgoapp/handler"
	"eventgoapp/kk"
	"eventgoapp/repository"
	"eventgoapp/service"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// ------------------ DB ------------------
	database := db.ConnectDB()
	defer database.Close()

	// ------------------ Kafka ------------------
	producer := kk.NewKafkaWriter("localhost:9092", "events_topic")
	defer producer.Close()

	consumer := kk.NewKafkaConsumer("localhost:9092", "event-consumer-group", "events_topic")
	defer consumer.Close()

	// ------------------ Repositories ------------------
	eventRepo := repository.NewEventRepository(database)
	attendeeRepo := repository.NewAttendeeRepository(database)
	rsvpRepo := repository.NewRsvpRepository(database, attendeeRepo, eventRepo)

	// ------------------ Services ------------------
	eventService := service.NewEventService(eventRepo)
	attendeeService := service.NewAttendeeService(attendeeRepo)
	rsvpService := service.NewRsvpService(rsvpRepo, attendeeRepo, eventRepo)
	notificationService := service.NewNotificationService(eventRepo, rsvpRepo, producer)

	// ------------------ Handlers ------------------
	eventHandler := handler.NewEventHandler(eventService, producer)
	attendeeHandler := handler.NewAttendeeHandler(attendeeService)
	rsvpHandler := handler.NewRsvpHandler(rsvpService)
	notificationHandler := handler.NewNotificationHandler(notificationService, producer, consumer)

	cache.InitRedis()
	fmt.Println("Connected to Redis")

	// ------------------ Gin Router ------------------
	router := gin.Default()
	router.Use(cors.Default())

	router.Static("/static", "./frontend")
	router.GET("/", func(c *gin.Context) {
		c.File("./frontend/index.html")
	})

	// Event endpoints
	router.POST("/events", eventHandler.AddEvent)
	router.GET("/events", eventHandler.GetAllEvents)
	router.GET("/events/:id", eventHandler.GetEventByID)
	router.PATCH("/events/:id", eventHandler.UpdateEvent)
	router.DELETE("/events/:id", eventHandler.DeleteEvent)

	// Attendee endpoints
	router.POST("/attendees", attendeeHandler.AddAttendee)

	// RSVP endpoints
	router.POST("/rsvp/:event_id/:attendee_id", rsvpHandler.RegisterToEvent)

	router.GET("/notifications", notificationHandler.StreamNotifications)
	// ------------------ HTTP Server & Graceful Shutdown ------------------
	srv := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	go func() {
		log.Println("Server is running on port 8081")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server startup failed: %v", err)
		}
	}()

	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for {
			<-ticker.C
			notificationService.NotifyTodayEvents()
		}
	}()

	// Wait for termination signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	if err := consumer.Close(); err != nil {
		log.Fatalf("Failed to close Kafka consumer: %v", err)
	}

	log.Println("Server gracefully stopped")
}
