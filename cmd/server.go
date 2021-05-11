package cmd

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
	"github.com/nebisin/gograph/graph/generated"
	"github.com/nebisin/gograph/graph/resolver"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	defaultPort       = "8080"
	defaultComplexity = 200 // TODO: Test the limit and define later
)

func Run() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	initEnv()

	client, database := initDatabase(ctx)
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(client, ctx)

	initServer(database)
}

func initServer(database *mongo.Database) {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	config := initServerConfig(database)

	schema := generated.NewExecutableSchema(config)
	srv := handler.NewDefaultServer(schema)
	srv.Use(extension.FixedComplexityLimit(defaultComplexity))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("🚀 connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func initServerConfig(database *mongo.Database) generated.Config {
	config := generated.Config{
		Resolvers: &resolver.Resolver{DB: database},
	}

	countComplexity := func(childComplexity int, limit *int, _ *int) int {
		count := *limit
		return count * childComplexity
	}

	config.Complexity.User.Tweets = countComplexity
	config.Complexity.Query.ListTweet = countComplexity

	return config
}

func initDatabase(ctx context.Context) (*mongo.Client, *mongo.Database) {
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")
	dbURI := fmt.Sprintf("mongodb://%s:%s@%s:%s", dbUsername, dbPassword, dbHost, dbPort)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Fatal("cannot connect the mongodb: ", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("ping to the database is failed: ", err)
	}

	database := client.Database("social")

	return client, database
}

func initEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
