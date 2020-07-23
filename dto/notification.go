package dto

import "time"

type Notification struct {
	ID       string    `json:"id" bson:"_id"`
	DeviceID string    `json:"device_id"`
	Channel  string    `json:"channel" bson:"channel"`
	Name     string    `json:"name" bson:"name"`
	Message  string    `json:"message" bson:"message"`
	UpdateAt time.Time `json:"update_at" bson:"update_at"`
	Unique   string    `json:"unique" bson:"unique"`
}
