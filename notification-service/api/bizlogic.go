package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"notification-service/dataservice"
	"notification-service/model"

	"github.com/segmentio/kafka-go"
)

func SubscribeUser(db *sql.DB, subscription model.Subscription) error {
	if subscription.UserID == "" {
		return errors.New("user_id is required")
	}

	if len(subscription.Topics) == 0 {
		return errors.New("at least one topic must be provided for subscription")
	}

	err := dataservice.AddSubscription(db, subscription)
	if err != nil {
		return err
	}

	return nil
}

func SendNotification(producer *kafka.Writer, notification model.Notification) error {
	if notification.Topic == "" {
		return errors.New("notification topic is required")
	}

	notificationJSON, err := json.Marshal(notification)
	if err != nil {
		return errors.New("failed to marshal notification")
	}

	msg := kafka.Message{
		Value: notificationJSON,
		Topic: notification.Topic,
	}

	if err := producer.WriteMessages(context.Background(), msg); err != nil {
		return errors.New("failed to send notification to Kafka")
	}

	return nil
}

func FetchUserSubscriptions(db *sql.DB, userID string) ([]model.Subscription, error) {
	if userID == "" {
		return nil, errors.New("user_id is required")
	}

	subscriptions, err := dataservice.GetUserSubscriptions(db, userID)
	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}
