package main

import (
	"context"
	"encoding/json"
	"net"
	"os"
	"pushnotif/api"
	"pushnotif/dto"
	"pushnotif/internal/app/device"
	dbCorona "pushnotif/internal/mongodb/corona"
	dbDevice "pushnotif/internal/mongodb/device"
	dbNotification "pushnotif/internal/mongodb/notification"
	"pushnotif/pkg/config"
	pkgMongo "pushnotif/pkg/mongodb"
	pkgRedis "pushnotif/pkg/redis"
	"reflect"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/appleboy/go-fcm"

	goRedis "github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

var (
	_config             config.Config
	_pkgMongo           *pkgMongo.Mongodb
	_pkgredis           pkgRedis.Adapter
	_dbCorona           *dbCorona.Corona
	_dbDevice           *dbDevice.Device
	_dbNotification     *dbNotification.Notification
	deviceLocationData1 chan string
	deviceLocationData2 chan string
	_fcmClient          *fcm.Client
	_device             *device.Device
)

func initLog() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
	log.SetReportCaller(true)

}
func initConfig() {
	_config = config.Load("./config/config.yaml")

}
func initDepedency() {
	var err error
	_pkgMongo, err = pkgMongo.NewMongodb(pkgMongo.Config{
		AppName:      _config.App.Name,
		Hosts:        _config.Mongo.Host,
		Username:     _config.Mongo.Username,
		Password:     _config.Mongo.Password,
		DatabaseName: _config.Mongo.DatabaseName,
	})
	if err != nil {
		log.Panic("initiate mongodb ", err.Error())
	}
	_pkgredis = pkgRedis.NewRedis(_config.App.Name, pkgRedis.Config{
		Addrs: strings.Split(_config.Redis.Host, ","),
		DB:    _config.Redis.Db,
	})
	_dbCorona = dbCorona.NewCorona(_pkgMongo)
	_dbDevice = dbDevice.NewDevice(_pkgMongo)
	_dbNotification = dbNotification.NewNotification(_pkgMongo)

	// Create a FCM client to send the message.
	_fcmClient, err = fcm.NewClient(_config.FCM.ApiKey)
	if err != nil {
		log.Fatal(err)
	}

	_device = device.NewDevice(_fcmClient, _dbDevice, _dbCorona, _dbNotification)

}
func serve() {
	api.NewRoute(_config.Serve.Port, _device, _pkgredis, _dbCorona).Run()
}

func init() {
	deviceLocationData1 = make(chan string)
	deviceLocationData2 = make(chan string)
	initConfig()
	initLog()
	initDepedency()
}

func fetchDataFromDeviceLocation() {
	pubsub := _pkgredis.Subscribe(context.Background(), "device_location")

	for {
		msg, err := pubsub.Receive()
		if err != nil {
			if reflect.TypeOf(err) == reflect.TypeOf(&net.OpError{}) && reflect.TypeOf(err.(*net.OpError).Err).String() == "*net.timeoutError" {
				// Timeout, ignore
				continue
			}
			// Actual error
			log.Print("Error in ReceiveTimeout()", err)
		}

		switch m := msg.(type) {
		case *goRedis.Subscription:
			log.Printf("Subscription Message: %v to channel '%v'. %v total subscriptions.", m.Kind, m.Channel, m.Count)
			continue
		case *goRedis.Message:

			deviceLocationData1 <- m.Payload
			deviceLocationData2 <- m.Payload

		}

	}
}

func deviceLocation() {
	for {
		select {
		case msg := <-deviceLocationData1:
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Error("Recovered in f ", r, string(debug.Stack()))
					}
				}()

				param := dto.Device{}
				err := json.Unmarshal([]byte(msg), &param)
				if err != nil {
					log.Error("error parsing payload ", err.Error())
				}

				err = _device.Upsert(context.Background(), param)
				if err != nil {
					log.Error("error parsing payload ", err.Error())
				}
			}()
		}
	}
}
func notification() {
	for {
		select {
		case msg := <-deviceLocationData2:
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Error("Recovered in f ", r, string(debug.Stack()))
					}
				}()

				param := dto.Device{}
				err := json.Unmarshal([]byte(msg), &param)
				if err != nil {
					log.Error("error parsing payload ", err.Error())
				}

				err = _device.SendNotification(context.Background(), param)
				if err != nil {
					log.Error("error parsing payload ", err.Error())
				}
			}()
		}
	}
}
func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		serve()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		fetchDataFromDeviceLocation()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		deviceLocation()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		notification()
	}()
	wg.Wait()

}
