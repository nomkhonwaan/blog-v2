package mongo

import "go.mongodb.org/mongo-driver/bson/primitive"

// DBRef is a MongoDB DBRef type
type DBRef struct {
	// A reference collection
	Ref string `bson:"$ref"`

	// A reference identifier
	ID primitive.ObjectID `bson:"$id"`
}
