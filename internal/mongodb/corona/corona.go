package corona

import (
	"context"
	"pushnotif/dto"
	"pushnotif/pkg/mongodb"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func NewCorona(client mongodb.IMongodb) *Corona {
	return &Corona{
		mongodb: client,
		timeout: time.Duration(10 * time.Second),
	}
}

type Corona struct {
	mongodb mongodb.IMongodb
	timeout time.Duration
}

// collectionName this for collection name
func (d *Corona) collectionName() string {
	return "corona"
}

// db.corona.aggregate([{$geoNear:{near:{type:"Point",coordinates:[106.8373179,-6.2490307]},distanceField: "Distance"}}, { $match: { "Distance": { $lte: 1000 }} } ]).pretty()

func (d *Corona) GetRadius(ctx context.Context, latitude float64, longitude float64, radiusMax float64) (*dto.Corona, error) {

	var queryList []bson.M
	var dtoCorona []dto.Corona
	queryList = append(queryList, bson.M{
		"$geoNear": bson.M{
			"near": bson.M{
				"type":        "Point",
				"coordinates": []float64{longitude, latitude},
			},
			"distanceField": "distance",
		},
	})
	queryList = append(queryList, bson.M{
		"$match": bson.M{
			"distance": bson.M{
				"$lte": radiusMax,
			},
		},
	})
	queryList = append(queryList, bson.M{
		"$limit": 1,
	})
	cs, err := d.mongodb.GetDatabase().Collection(d.collectionName()).Aggregate(context.Background(), queryList)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	err = cs.All(context.Background(), &dtoCorona)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	if len(dtoCorona) <= 0 {
		return nil, nil
	}

	return &dtoCorona[0], nil
}
