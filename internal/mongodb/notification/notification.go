package notification

import (
	"context"
	"pushnotif/dto"
	"pushnotif/pkg/mongodb"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

func NewNotification(client mongodb.IMongodb) *Notification {
	return &Notification{
		mongodb: client,
		timeout: time.Duration(10 * time.Second),
	}
}

type Notification struct {
	mongodb mongodb.IMongodb
	timeout time.Duration
}

// collectionName this for collection name
func (d *Notification) collectionName() string {
	return "notification"
}

func (d *Notification) FindDeviceID(ctx context.Context, deviceID string) (*dto.Notification, error) {

	var result dto.Notification

	err := d.mongodb.GetDatabase().Collection(d.collectionName()).FindOne(context.Background(), bson.M{
		"device_id": deviceID,
	}).Decode(&result)

	if err == mongoDriver.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (d *Notification) Upsert(ctx context.Context, param dto.Notification) error {
	queryParam := bson.M{}
	updateParam := bson.M{
		"device_id": param.DeviceID,
		"name":      param.Name,
		"channel":   param.Channel,
		"message":   param.Message,
		"update_at": param.UpdateAt,
		"unique":    param.Unique,
	}

	queryParam["$set"] = updateParam
	upsert := true
	_, err := d.mongodb.GetDatabase().Collection(d.collectionName()).UpdateOne(context.Background(), bson.M{"device_id": param.DeviceID}, queryParam, &options.UpdateOptions{
		Upsert: &upsert,
	})
	if err != nil {
		return err
	}

	return nil
}
