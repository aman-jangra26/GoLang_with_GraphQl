package main

import (
	"log"
	"net/http"

	"Backend/graphql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/graphql-go/handler"
)

func main() {
	// Set up GraphQL schema
	h := handler.New(&handler.Config{
		Schema:   &graphql.Schema,
		Pretty:   true,
		GraphiQL: true,
	})

	// Set up HTTP handler for GraphQL endpoint
	http.Handle("/graphql", h)

	// Start HTTP server
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
