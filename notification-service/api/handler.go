package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"notification-service/dataservice"
	"notification-service/model"
	"notification-service/queue"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

type Handler struct {
	db       *sql.DB
	producer *kafka.Writer
}

func NewHandler(db *sql.DB, producer *kafka.Writer) *Handler {
	return &Handler{db: db, producer: producer}
}

func (h *Handler) Subscribe(c *gin.Context) {
	var subscription model.Subscription
	if err := c.BindJSON(&subscription); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	if err := dataservice.AddSubscription(h.db, subscription); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to subscribe"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Subscribed successfully"})
}

func (h *Handler) SendNotification(c *gin.Context) {
	var request model.Notification

	// Bind JSON request to the Notification struct
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	kafkaConfig := queue.KafkaConfig{
		Broker: "localhost:9092",
		Topic:  request.Topic,
	}

	producer := queue.NewProducer(kafkaConfig)
	defer producer.Close()

	notificationJSON, err := json.Marshal(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal notification"})
		return
	}

	if err := producer.SendMessage(c, notificationJSON); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send notification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification sent successfully"})
}

func (h *Handler) Unsubscribe(c *gin.Context) {
	var request model.UnsubscribeRequest

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := dataservice.RemoveSubscription(h.db, request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unsubscribe"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unsubscribed successfully"})
}

func (h *Handler) FetchUserSubscriptions(c *gin.Context) {
	userID := c.Param("user_id")
	log.Printf("Received user_id from request: %s", userID)

	subscriptions, err := dataservice.GetUserSubscriptions(h.db, userID)
	if err != nil {
		log.Printf("Error fetching subscriptions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subscriptions"})
		return
	}

	log.Printf("Fetched subscriptions: %v", subscriptions)
	c.JSON(http.StatusOK, gin.H{"user_id": userID, "subscriptions": subscriptions})
}
