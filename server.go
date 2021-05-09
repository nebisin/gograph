package main

import (
	"context"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/nebisin/gograph/db"
	"github.com/nebisin/gograph/graph"
	"github.com/nebisin/gograph/graph/generated"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"os"
	"time"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client := db.InitClient(ctx)
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(client, ctx)

	database := client.Database("social")

	config := generated.Config{
		Resolvers: &graph.Resolver{DB: database},
	}

	countComplexity := func(childComplexity int, limit *int, _ *int) int {
		count := *limit
		return count * childComplexity
	}

	config.Complexity.User.Tweets = countComplexity
	config.Complexity.Query.ListTweet = countComplexity

	schema := generated.NewExecutableSchema(
		config,
	)

	srv := handler.NewDefaultServer(schema)
	srv.Use(extension.FixedComplexityLimit(200)) // TODO: Test the limit and define later

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
