package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

// ConnectDB initializes the MongoDB connection
func ConnectDB() *mongo.Database {
	uri := os.Getenv("MONGO_LOCAL_URI")
	dbName := os.Getenv("MONGO_DB_NAME")

	clientOptions := options.Client().ApplyURI(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("MongoDB Connection Error: ", err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB Ping Error: ", err)
	}

	log.Println("Connected to MongoDB!")

	DB = client.Database(dbName)

	// Ensure TTL index on blacklisted_tokens.expires_at so expired entries are removed
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		col := DB.Collection("blacklisted_tokens")
		idxModel := mongo.IndexModel{
			Keys:    bson.D{{Key: "expires_at", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(0),
		}
		if _, err := col.Indexes().CreateOne(ctx, idxModel); err != nil {
			log.Println("warning: could not create TTL index for blacklisted_tokens:", err)
		}
	}()
	return DB
}
