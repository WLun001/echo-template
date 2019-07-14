package repository

import "go.mongodb.org/mongo-driver/mongo"

type Repository interface {

}

type repository struct {
	client *mongo.Client
	db     *mongo.Database
}

// New :
func New(dbName string, mongo *mongo.Client) Repository {
	return &repository{
		client: mongo,
		db:     mongo.Database(dbName),
	}
}
