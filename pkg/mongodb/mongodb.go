package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

const (
	ErrorDataNotFound = "data not found"
)

type IMongodb interface {
	GetConnection() *mongo.Client
	GetDatabase() *mongo.Database
}

// Config this for config mongodb
type Config struct {
	AppName      string
	Hosts        []string
	Timeout      time.Duration
	Auth         string
	Username     string
	Password     string
	DatabaseName string
}

// NewMongodb this for new connection
func NewMongodb(cfg Config) (*Mongodb, error) {
	var err error
	mongodb := new(Mongodb)
	mongodb.dbName = cfg.DatabaseName
	mongodb.client, err = mongodb.newConnection(cfg)
	if err != nil {
		return nil, err
	}
	return mongodb, nil
}

type Mongodb struct {
	client *mongo.Client
	dbName string
}

func (m *Mongodb) newConnection(cfg Config) (*mongo.Client, error) {

	opts := &options.ClientOptions{
		Hosts:   cfg.Hosts,
		AppName: &cfg.AppName,
		WriteConcern: writeconcern.New(writeconcern.WMajority(), writeconcern.J(true),
			writeconcern.WTimeout(10*time.Second)),
		ReadPreference: readpref.SecondaryPreferred(),
	}
	if cfg.Username != "" && cfg.Password != "" {
		opts.Auth = &options.Credential{
			Username:   cfg.Username,
			Password:   cfg.Password,
			AuthSource: cfg.DatabaseName,
		}
	}
	client, err := mongo.NewClient(opts)

	opts.SetMaxPoolSize(5)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("error while ping mongodb client %s", err.Error())
	}

	return client, err

}

func (m *Mongodb) GetConnection() *mongo.Client {
	return m.client
}

func (m *Mongodb) GetDatabase() *mongo.Database {
	return m.client.Database(m.dbName)
}
