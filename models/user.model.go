package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Fullname string             `json:"fullname" bson:"fullname"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password" bson:"password"`
}
