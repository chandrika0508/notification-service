package main

import (
	"database/sql"
	"log"

	"notification-service/api"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/segmentio/kafka-go"
)

func main() {
	dsn := "root:chandu@123@tcp(127.0.0.1:3306)/notification"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()

	kafkaWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "notifications",
	})
	defer kafkaWriter.Close()

	router := gin.Default()
	handler := api.NewHandler(db, kafkaWriter)

	router.POST("/subscribe", handler.Subscribe)
	router.POST("/notifications/send", handler.SendNotification)
	router.POST("/unsubscribe", handler.Unsubscribe)
	router.GET("/subscriptions/:user_id", handler.FetchUserSubscriptions)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
