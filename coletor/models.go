package main

type Metric struct {
Timestamp int64  `bson:"timestamp"`
Metrics   string `bson:"metrics"`
}