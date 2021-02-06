package app

import (
	"context"

	"github.com/nomkhonwaan/myblog/pkg/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

func provideMongoDB(uri, dbName string) (mongo.Database, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if dbName == "" {
		connString, err := connstring.Parse(uri)
		if err != nil {
			return nil, err
		}
		dbName = connString.Database
	}

	return client.Database(dbName), nil
}
