package model

type Subscription struct {
	UserID   string               `json:"user_id"`
	Topics   []string             `json:"topics"`
	Channels NotificationChannels `json:"channels"`
}

type NotificationChannels struct {
	Email             string `json:"email,omitempty"`
	SMS               string `json:"sms,omitempty"`
	PushNotifications bool   `json:"push_notifications"`
}

type UnsubscribeRequest struct {
	UserID string   `json:"user_id"`
	Topics []string `json:"topics"`
}
