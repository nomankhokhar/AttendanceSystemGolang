package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client // Global client variable

// Initializes the MongoDB client and returns it the main function
func InitDB() (*mongo.Client, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	clientOptions := options.Client().ApplyURI(databaseURL)

	// Initialize the global client variable
	var err error
	client, err = mongo.NewClient(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create MongoDB client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Connect to the MongoDB server
	err = client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Check the connection with database
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	fmt.Println("Connected to database")
	return client, nil
}

func GetDB() *mongo.Database {
	if client == nil {
		log.Fatal("Database not initialized")
	}
	return client.Database("attendance_system")
}
