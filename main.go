package main

import (
	"api/repository"
	"api/routes"
	"context"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

func main() {

	connectionString := os.Getenv("DATABASE_CONNECTION")
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)

	defer cancel()
	if err := mongoClient.Connect(ctx); err != nil {
		panic(err)
	}

	repo := repository.New(os.Getenv("DATABASE_NAME"), mongoClient)
	//hosts := make(map[string]*echo.Echo)
	apiHost := routes.APIRoutes(repo, mongoClient)
	//hosts["localhost:8080"] = apiHost

	e := echo.New()
	e.Any("/*", func(c echo.Context) (err error) {
		req := c.Request()
		res := c.Response()
		host := apiHost

		if host == nil {
			err = echo.ErrNotFound
			return
		}

		host.ServeHTTP(res, req)
		return
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	e.Logger.Fatal(e.Start(":" + port))
}
