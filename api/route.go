package api

import (
	"fmt"
	"pushnotif/internal/app/device"
	"pushnotif/internal/mongodb/corona"
	"pushnotif/pkg/redis"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func NewRoute(port int, device *device.Device, redis redis.Adapter, corona *corona.Corona) *Route {
	return &Route{
		httpPort:    port,
		routeEngine: gin.Default(),
		device:      device,
		redis:       redis,
		corona:      corona,
	}
}

type Route struct {
	httpPort    int
	routeEngine *gin.Engine
	device      *device.Device
	redis       redis.Adapter
	corona      *corona.Corona
}

func (rr *Route) Run() {
	rr.registerController()
	err := rr.routeEngine.Run(fmt.Sprintf(":%d", rr.httpPort))
	if err != nil {
		log.Error(err.Error())
	}
}

func (rr *Route) registerController() {
	rr.routeEngine.GET("/ping", func(c *gin.Context) {
		log.Info(rr.corona.GetRadius(c.Request.Context(), -6.2490307, 106.8373179, 1))
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	rr.routeEngine.POST("/device/location", rr.deviceLocationController)
}
