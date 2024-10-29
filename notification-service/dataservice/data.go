package dataservice

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"notification-service/model"
	"strings"
)

func AddSubscription(db *sql.DB, sub model.Subscription) error {
	query := "INSERT INTO subscriptions (user_id, topic, channels) VALUES (?, ?, ?)"

	channelsJSON, err := json.Marshal(sub.Channels)
	if err != nil {
		return err
	}

	for _, topic := range sub.Topics {
		_, err = db.Exec(query, sub.UserID, topic, channelsJSON)
		if err != nil {
			return err
		}
	}

	return nil
}

func RemoveSubscription(db *sql.DB, req model.UnsubscribeRequest) error {
	query := "DELETE FROM subscriptions WHERE user_id = ? AND topic IN ("

	placeholders := make([]string, len(req.Topics))
	for i := range req.Topics {
		placeholders[i] = "?"
	}

	query += strings.Join(placeholders, ",") + ")"

	args := []interface{}{req.UserID}
	for _, topic := range req.Topics {
		args = append(args, topic)
	}

	result, err := db.Exec(query, args...)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("no subscription found to remove")
	}

	return nil
}

func GetUserSubscriptions(db *sql.DB, userID string) ([]model.Subscription, error) {
	query := "SELECT topic, channels FROM subscriptions WHERE user_id = ?"
	log.Printf("Executing query: %s with user_id: %s", query, userID)

	rows, err := db.Query(query, userID)
	if err != nil {
		log.Printf("Query execution error: %v", err)
		return nil, err
	}
	defer rows.Close()

	subscriptionMap := make(map[string]*model.Subscription)

	for rows.Next() {
		var topic string
		var channelsJSON []byte
		if err := rows.Scan(&topic, &channelsJSON); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}

		sub, exists := subscriptionMap[userID]
		if !exists {
			sub = &model.Subscription{
				UserID: userID,
				Topics: []string{},
			}
			if err := json.Unmarshal(channelsJSON, &sub.Channels); err != nil {
				log.Printf("Error unmarshalling channels JSON: %v", err)
				return nil, err
			}
			subscriptionMap[userID] = sub
		}

		sub.Topics = append(sub.Topics, topic)
		log.Printf("Added topic %s for user %s", topic, userID)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Row iteration error: %v", err)
		return nil, err
	}

	subscriptions := make([]model.Subscription, 0, len(subscriptionMap))
	for _, sub := range subscriptionMap {
		subscriptions = append(subscriptions, *sub)
	}

	log.Printf("Retrieved subscriptions: %v", subscriptions)
	return subscriptions, nil
}
