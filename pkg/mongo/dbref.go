package mongo

import "go.mongodb.org/mongo-driver/bson/primitive"

// DBRef represents a MongoDB DBRef type
type DBRef struct {
	// Reference collection name
	Ref string `bson:"$ref"`

	// Reference identifier
	ID primitive.ObjectID `bson:"$id"`
}
