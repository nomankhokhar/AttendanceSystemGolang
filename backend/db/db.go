package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func InitDB() *mongo.Client {
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

    client, err := mongo.NewClient(clientOptions)
    if err != nil {
        log.Fatal("Failed to create new MongoDB client:", err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    err = client.Connect(ctx)
    if err != nil {
        log.Fatal("Failed to connect to MongoDB:", err)
    }

    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal("Failed to ping MongoDB:", err)
    }

    fmt.Println("Connected to MongoDB!")

    MongoClient = client
    return client
}

func GetCollection(client *mongo.Client, dbName string, collectionName string) *mongo.Collection {
    return client.Database(dbName).Collection(collectionName)
}
