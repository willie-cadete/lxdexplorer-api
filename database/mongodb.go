package database

import (
	"context"
	"lxdexplorer-api/config"

	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var conf, _ = config.LoadConfig()

func connect() *mongo.Client {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(conf.MongoDB.URI))
	if err != nil {
		panic(err)
	}

	ping(client)

	return client
}

func ping(c *mongo.Client) {
	err := c.Ping(context.Background(), nil)
	if err != nil {
		panic(err)
	}
}

func InsertOne(collection string, document interface{}) {
	c := connect()
	defer c.Disconnect(context.Background())

	i, err := c.Database("lxd").Collection(collection).InsertOne(context.TODO(), document)
	if err != nil {
		panic(err)
	}
	log.Printf("Inserted document with ID %v\n", i.InsertedID)

}

func InsertMany(collection string, documents []interface{}) {
	c := connect()
	defer c.Disconnect(context.Background())

	i, err := c.Database("lxd").Collection(collection).InsertMany(context.TODO(), documents)
	if err != nil {
		panic(err)
	}
	log.Printf("Inserted %v documents with IDs %v\n", len(i.InsertedIDs), i.InsertedIDs)

}
