package main

import (
	"context"
	"fmt"
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

func StoreMetrics(metrics *Metrics) error {
	collection := client.Database("metricsdb").Collection("metrics")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	metricsMap := map[string]interface{}{
		"timestamp":   time.Now(),
		"cpu":         metrics.CPU,
		"ram":         metrics.RAM,
		"queue_size":  metrics.QueueSize,
		"error_count": metrics.ErrorCount,
	}

	fmt.Printf("Guardando métricas no banco: %+v\n", metricsMap)

	_, err := collection.InsertOne(ctx, metricsMap)
	if err != nil {
		return err
	}

	fmt.Println("Métricas guardadas com sucesso no banco de dados.")

	return nil
}
