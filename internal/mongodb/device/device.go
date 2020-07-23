package device

import (
	"context"
	"pushnotif/dto"
	"pushnotif/pkg/mongodb"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

func NewDevice(client mongodb.IMongodb) *Device {
	return &Device{
		mongodb: client,
		timeout: time.Duration(10 * time.Second),
	}
}

type Device struct {
	mongodb mongodb.IMongodb
	timeout time.Duration
}

// collectionName this for collection name
func (d *Device) collectionName() string {
	return "device"
}

func (d *Device) StoreDevice(ctx context.Context, param dto.Device) error {
	param.CurrentLocation.CreatedAt = time.Now()
	_, err := d.mongodb.GetDatabase().Collection(d.collectionName()).InsertOne(context.Background(), bson.M{
		"device_id":        param.DeviceID,
		"current_location": param.CurrentLocation,
		"location_history": []dto.DeviceLocation{param.CurrentLocation},
	})
	if err != nil {
		return err
	}
	return nil
}

func (d *Device) FindDeviceID(ctx context.Context, deviceID string) (*dto.Device, error) {

	var result dto.Device

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

func (d *Device) Update(ctx context.Context, param dto.Device) error {
	queryParam := bson.M{}
	updateParam := bson.M{}
	if param.CurrentLocation.Latitude != 0 || param.CurrentLocation.Longitude != 0 {
		param.CurrentLocation.CreatedAt = time.Now()
		updateParam["current_location"] = param.CurrentLocation
		queryParam["$push"] = bson.M{"location_history": param.CurrentLocation}
	}

	queryParam["$set"] = updateParam
	upsert := false
	_, err := d.mongodb.GetDatabase().Collection(d.collectionName()).UpdateOne(context.Background(), bson.M{"device_id": param.DeviceID}, queryParam, &options.UpdateOptions{
		Upsert: &upsert,
	})
	if err != nil {
		return err
	}

	return nil
}
