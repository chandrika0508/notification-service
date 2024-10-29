package model

import "time"

type Notification struct {
	Topic   string  `json:"topic"`
	Event   Event   `json:"event"`
	Message Message `json:"message"`
}

type Topic struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type Event struct {
	EventID   string                 `json:"event_id"`
	Timestamp time.Time              `json:"timestamp"`
	Details   map[string]interface{} `json:"details"`
}

type Message struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}
