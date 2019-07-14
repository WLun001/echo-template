package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	CollectionUser = "User"
)

type Model struct {
	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
}

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Email    string             `bson:"email"`
	Username string             `bson:"username"`
	Security struct {
		PasswordHash string `bson:"passwordHash"`
		RefreshToken string `bson:"refreshToken"`
	} `bson:"security"`
	Model `bson:",inline"`
}

type Random struct {
	Name         string `json:"name" validate:"required"`
	RandomNumber int64  `json:"randomNumber" validate:"required,numeric"`
}

type Response struct {
	Message string `json:"message"`
}
