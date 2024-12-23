package main

import (
"context"
"log"
"time"

"go.mongodb.org/mongo-driver/mongo"
"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func init() {
var err error
clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
client, err = mongo.Connect(context.TODO(), clientOptions)
if err != nil {
	log.Fatal(err)
}

err = client.Ping(context.TODO(), nil)
if err != nil {
	log.Fatal(err)
}
}

func StoreMetrics(metrics string) error {
collection := client.Database("metricsdb").Collection("metrics")

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

_, err := collection.InsertOne(ctx, map[string]interface{}{
	"timestamp": time.Now(),
	"metrics":   metrics,
})
return err
}