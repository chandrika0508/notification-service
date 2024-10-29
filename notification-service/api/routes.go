package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

func RegisterRoutes(db *sql.DB, producer *kafka.Writer) *gin.Engine {
	r := gin.Default()
	handler := NewHandler(db, producer)

	r.Use(gin.Logger())

	r.POST("/subscribe", handler.Subscribe)
	r.POST("/notifications/send", handler.SendNotification)
	r.POST("/unsubscribe", handler.Unsubscribe)
	r.GET("/subscriptions/:user_id", handler.FetchUserSubscriptions)

	return r
}
