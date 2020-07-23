package dto

import "time"

type DeviceLocation struct {
	Latitude  float64   `json:"latitude" bson:"latitude"`
	Longitude float64   `json:"longitude" bson:"longitude"`
	CreatedAt time.Time `json:"created_at" bson:"created_at,omitempty"`
}
type Device struct {
	ID              string           `json:"id" bson:"_id"`
	DeviceID        string           `json:"device_id" bson:"device_id"`
	CurrentLocation DeviceLocation   `json:"current_location" bson:"current_location"`
	LocationHistory []DeviceLocation `json:"location_history" bson:"location_history"`
}
