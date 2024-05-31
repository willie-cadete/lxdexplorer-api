package database

import (
	"context"
	"lxdexplorer-api/config"
	"time"

	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var conf, _ = config.LoadConfig()

func connect() (*mongo.Client, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(conf.MongoDB.URI).SetTimeout(time.Second*5)) // Somehow is 10 second timeout
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return client, err
}

func Ping() error {
	c, err := connect()
	if err != nil {
		return err
	}

	defer c.Disconnect(context.Background())

	return nil
}

func InsertOne(collection string, document interface{}) {
	c, _ := connect()
	defer c.Disconnect(context.Background())

	i, err := c.Database("lxd").Collection(collection).InsertOne(context.TODO(), document)
	if err != nil {
		panic(err)
	}
	log.Printf("Database: Inserted document with ID %v\n", i.InsertedID)

}

func InsertMany(collection string, documents []interface{}) {
	c, _ := connect()
	defer c.Disconnect(context.Background())

	i, err := c.Database("lxd").Collection(collection).InsertMany(context.TODO(), documents)
	if err != nil {
		panic(err)
	}
	log.Printf("Database: Inserted %v documents with IDs %v\n", len(i.InsertedIDs), i.InsertedIDs)

}

func FindOne(collection string, filter interface{}) *mongo.SingleResult {
	c, _ := connect()
	defer c.Disconnect(context.Background())

	return c.Database("lxd").Collection(collection).FindOne(context.Background(), filter)
}

func FindAll(collection string) ([]primitive.M, error) {
	c, err := connect()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer c.Disconnect(context.Background())

	cur, err := c.Database("lxd").Collection(collection).Find(context.Background(), bson.D{})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	var results []primitive.M
	for cur.Next(context.Background()) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil

}

func ReplaceOne(collection string, filter interface{}, replacement interface{}) (*mongo.UpdateResult, error) {
	c, _ := connect()
	defer c.Disconnect(context.Background())

	return c.Database("lxd").Collection(collection).ReplaceOne(context.Background(), filter, replacement)
}

func AddTTL(collection string, field string, seconds int32) {
	c, _ := connect()
	defer c.Disconnect(context.Background())

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: field, Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(seconds),
	}

	indexView := c.Database("lxd").Collection(collection).Indexes()

	// Create a new index
	var err error
	_, err = indexView.CreateOne(context.Background(), indexModel)
	if err != nil {
		// Drop the existing index
		log.Println(err)
		_, err := indexView.DropOne(context.Background(), string(field+"_1"))
		if err != nil {
			// Handle error
			log.Println(err)
		}
		// log.Printf("Database: Dropped existing TTL index on %s\n", field)
		// Create a new index
		_, err = indexView.CreateOne(context.Background(), indexModel)
		if err != nil {
			// Handle error
			log.Println(err)
		}
	}

}
