package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"pushnotif/dto"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type paramDeviceLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (rr *Route) deviceLocationController(ctx *gin.Context) {
	var response = struct {
		Message string `json:"message"`
	}{}
	var payload paramDeviceLocation
	param := dto.Device{}
	param.DeviceID = ctx.Request.Header.Get("Device-ID")
	bData, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Warn("invalid validation ", err.Error())
		response.Message = "invalid json format"
		ctx.JSON(400, response)
		return
	}
	err = json.Unmarshal(bData, &payload)
	if err != nil {
		log.Warn("error decode payload ", err.Error())
		response.Message = "invalid json format"
		ctx.JSON(400, response)
		return
	}
	param.CurrentLocation.Longitude = payload.Longitude
	param.CurrentLocation.Latitude = payload.Latitude
	param.LocationHistory = append(param.LocationHistory, param.CurrentLocation)

	if param.CurrentLocation.Latitude == 0 || param.CurrentLocation.Longitude == 0 || param.DeviceID == "" {
		log.WithFields(log.Fields{
			"data":    param,
			"payload": string(bData),
		}).Warn("invalid validation ")
		response.Message = "invalid data"
		ctx.JSON(400, response)
		return
	}
	bBody, _ := json.Marshal(param)
	_, err = rr.redis.Publish(context.Background(), "device_location", string(bBody))
	if err != nil {
		log.Error(err.Error())
		response.Message = "error system"
		ctx.JSON(500, response)
		return
	}
	response.Message = "success"
	ctx.JSON(200, response)
}
