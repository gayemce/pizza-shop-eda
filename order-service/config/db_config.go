package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gayemce/pizza-shop-eda/order-service/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DBClient *mongo.Client
	MONGO_DB_NAME = "pizza-shop-eda"
)

func init() {
	InitializeDB()
}

func InitializeDB() (*mongo.Client, error) {
	var err error
	logger.Log("initilizing database once more")
	if DBClient == nil {
		DBClient, err = initDatabase()
		if err != nil {
			log.Fatal("Failed to intialize database: %v", err)
		}
	}
	return DBClient, err
}

func initDatabase() (*mongo.Client, error) {
	dbURL := GetEnvProperty("database_url")
	if dbURL == "" {
		return nil, fmt.Errorf("database_url is not set in the environment variable")
	}
	clientOptions := options.Client().ApplyURI(dbURL).
		SetMaxPoolSize(600).
		SetMinPoolSize(50).
		SetMaxConnIdleTime(30 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect with mongodb: %v", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("mongodb is unreachable: %v", err)
	}

	log.Printf("Connected to MongoDB: %s\n", dbURL)
	return client, nil
}

func GetDatabaseCollection(dbName *string, collectionName string) *mongo.Collection {
	if dbName == nil {
		dbName = &MONGO_DB_NAME
	}
	client, err := InitializeDB()
	if err != nil || client == nil {
		log.Fatalf("MongoDB client initialization failed: %v", err)
	}
	return client.Database(*dbName).Collection(collectionName)
}

func GetMongoClient() *mongo.Client {
	client, err := InitializeDB()
	if err != nil || client == nil {
		log.Fatalf("MongoDB client initialization failed: %v", err)
	}
	return client
}