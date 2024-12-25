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

	fmt.Printf("Guardando métricas:\n")
	fmt.Printf("  Timestamp: %s\n", metricsMap["timestamp"].(time.Time).Format(time.RFC3339))
	fmt.Printf("  CPU: %.2f\n", metricsMap["cpu"].(float64))
	fmt.Printf("  RAM: %.2f MB\n", metricsMap["ram"].(float64))
	fmt.Printf("  Queue Size: %d\n", metricsMap["queue_size"].(int))
	fmt.Printf("  Error Count: %d\n", metricsMap["error_count"].(int))

	_, err := collection.InsertOne(ctx, metricsMap)
	if err != nil {
		return err
	}

	fmt.Println("Métricas guardadas com sucesso no banco de dados.")

	return nil
}
