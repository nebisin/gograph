package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"testing"
	"time"
)

var testRepository *Repository

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")
	dbURI := fmt.Sprintf("mongodb://%s:%s@%s:%s", dbUsername, dbPassword, dbHost, dbPort)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Fatal("cannot connect the mongodb: ", err)
	}

	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(client, ctx)

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("ping to the database is failed: ", err)
	}

	database := client.Database("social_test")

	testRepository = NewRepository(database)

	os.Exit(m.Run())
}