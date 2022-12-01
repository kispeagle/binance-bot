package repo

import (
	"binance-bot/logger"
	"context"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	l = logger.GetLogger()
)

type MongoRepo struct {
	c      *mongo.Client
	dbName string
}

func NewMongoRepo(dns string, dbName string) *MongoRepo {

	client, err := mongo.NewClient(options.Client().ApplyURI(dns))
	if err != nil {
		l.Error(err.Error())
		return nil
	}

	strTimeout := os.Getenv("mongo_timeout")
	timeout, err := strconv.ParseInt(strTimeout, 10, 64)
	if err != nil {
		l.Error(err.Error())
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		l.Error(err.Error())
		return nil
	}

	return &MongoRepo{
		c:      client,
		dbName: dbName,
	}
}

func GetCollection(m MongoRepo, name string) *mongo.Collection {
	return (*mongo.Collection)(m.c.Database(m.dbName).Collection(name))
}
