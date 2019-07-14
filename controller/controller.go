package controller

import (
	"api/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

// Controller :
type Controller struct {
	repository repository.Repository
	mongo      *mongo.Client
}

// New :
func New(repo repository.Repository, db *mongo.Client) *Controller {
	return &Controller{
		repository: repo,
		mongo:      db,
	}
}
