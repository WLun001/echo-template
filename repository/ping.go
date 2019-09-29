package repository

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func (r repository) PingDatabase() (string, error) {
	ctx := context.Background()
	err := r.client.Ping(ctx, readpref.Primary())
	if err != nil {
		return "", errors.New("connection error")
	}
	return "connection to database established", nil
}
