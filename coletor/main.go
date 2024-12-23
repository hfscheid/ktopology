package main

import (
"context"
"log"
"os"
"time"

"go.mongodb.org/mongo-driver/mongo"
"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func init() {
var err error
mongoURI := os.Getenv("MONGO_URI")
if mongoURI == "" {
	mongoURI = "mongodb://localhost:27017"
}

clientOptions := options.Client().ApplyURI(mongoURI)
client, err = mongo.Connect(context.TODO(), clientOptions)
if err != nil {
	log.Fatal(err)
}

err = client.Ping(context.TODO(), nil)
if err != nil {
	log.Fatal(err)
}
}