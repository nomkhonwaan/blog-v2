package mongo

import "go.mongodb.org/mongo-driver/bson/primitive"

type DBRef struct {
	Ref string             `bson:"$ref"`
	ID  primitive.ObjectID `bson:"$id"`
}
