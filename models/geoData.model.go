package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type GeoData struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	UserID   string             `json:"user_id" bson:"user_id"`
	FilePath string             `json:"file_path" bson:"file_path"`
	Geometry string             `json:"geometry" bson:"geometry"`
	Title    string             `json:"title" bson:"title"`
	Shapes   string             `json:"shapes" bson:"shapes"`
}
