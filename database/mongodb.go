package database

import (
	"context"
	"lxdexplorer-api/config"

	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	log.Printf("Database: Inserted document with ID %v\n", i.InsertedID)

}

func InsertMany(collection string, documents []interface{}) {
	c := connect()
	defer c.Disconnect(context.Background())

	i, err := c.Database("lxd").Collection(collection).InsertMany(context.TODO(), documents)
	if err != nil {
		panic(err)
	}
	log.Printf("Database: Inserted %v documents with IDs %v\n", len(i.InsertedIDs), i.InsertedIDs)

}

func FindOne(collection string, filter interface{}) *mongo.SingleResult {
	c := connect()
	defer c.Disconnect(context.Background())

	return c.Database("lxd").Collection(collection).FindOne(context.Background(), filter)
}

func FindAll(collection string) ([]primitive.M, error) {
	c := connect()
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
	c := connect()
	defer c.Disconnect(context.Background())

	return c.Database("lxd").Collection(collection).ReplaceOne(context.Background(), filter, replacement)
}

func AddTTL(collection string, field string, seconds int32) {
	c := connect()
	defer c.Disconnect(context.Background())

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{field, 1}},
		Options: options.Index().SetExpireAfterSeconds(seconds),
	}

	_, err := c.Database("lxd").Collection(collection).Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		panic(err)
	}
}
