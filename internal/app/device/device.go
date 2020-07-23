package device

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"pushnotif/dto"
	dbCorona "pushnotif/internal/mongodb/corona"
	dbDevice "pushnotif/internal/mongodb/device"
	dbNotification "pushnotif/internal/mongodb/notification"
	"time"

	"github.com/appleboy/go-fcm"
	log "github.com/sirupsen/logrus"
)

func NewDevice(fcm *fcm.Client, dbDevice *dbDevice.Device, dbCorona *dbCorona.Corona, dbNotification *dbNotification.Notification) *Device {
	return &Device{fcm: fcm, dbDevice: dbDevice, dbCorona: dbCorona, dbNotification: dbNotification}
}

type Device struct {
	dbDevice       *dbDevice.Device
	dbCorona       *dbCorona.Corona
	dbNotification *dbNotification.Notification
	fcm            *fcm.Client
}

func (d *Device) Upsert(ctx context.Context, paramDevice dto.Device) error {
	device, err := d.dbDevice.FindDeviceID(ctx, paramDevice.DeviceID)
	if err != nil {
		return err
	}
	// insert device
	if device == nil {
		err = d.dbDevice.StoreDevice(ctx, paramDevice)
		if err != nil {
			return err
		}
		return nil
	}
	// update device
	err = d.dbDevice.Update(ctx, paramDevice)
	if err != nil {
		return err
	}
	return nil
}

func (d *Device) SendNotification(ctx context.Context, paramDevice dto.Device) error {
	// radius 1km
	coronaData, err := d.dbCorona.GetRadius(ctx, paramDevice.CurrentLocation.Latitude, paramDevice.CurrentLocation.Longitude, 1000)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if coronaData == nil {
		log.Warn("data not found")
		return nil
	}
	notif, err := d.dbNotification.FindDeviceID(ctx, paramDevice.DeviceID)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	h := sha512.New()
	h.Write([]byte(fmt.Sprintf("wilayah:%s city:%s country_id: %s device_id:%s", coronaData.Region.UrbanVillage, coronaData.Region.City, coronaData.Region.CountryID, paramDevice.DeviceID)))
	unique := base64.StdEncoding.EncodeToString([]byte(h.Sum(nil)))
	if notif != nil {
		if notif.Unique == unique {
			log.Warn("your still the same place")
			return nil
		}
	}

	pushNotivicationMessage := map[string]interface{}{
		"status":      coronaData.Zona,
		"location":    coronaData.Region.UrbanVillage,
		"active_case": coronaData.Confirmed,
	}
	bpushNotivicationMessagr, _ := json.Marshal(pushNotivicationMessage)
	notificationParam := dto.Notification{
		DeviceID: paramDevice.DeviceID,
		Channel:  "push_notification",
		Name:     "virus corona",
		Message:  string(bpushNotivicationMessagr),
		UpdateAt: time.Now(),
		Unique:   unique,
	}
	err = d.dbNotification.Upsert(ctx, notificationParam)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	msg := &fcm.Message{
		To:   paramDevice.DeviceID,
		Data: pushNotivicationMessage,
		Notification: &fcm.Notification{
			Title: "Covid Alert",
			Body:  string(bpushNotivicationMessagr),
		},
	}

	// Send the message and receive the response without retries.
	_, err = d.fcm.Send(msg)
	if err != nil {
		log.Error(err)
	}
	return nil
}
