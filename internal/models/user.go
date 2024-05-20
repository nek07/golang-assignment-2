package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email       string             `bson:"email" json:"email"`
	Username    string             `bson:"username" json:"username"`
	Password    string             `bson:"password" json:"password"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	AccessToken string             `bson:"access_token" json:"access_token"`
	Role        string             `bson:"role" json:"role"`
}
