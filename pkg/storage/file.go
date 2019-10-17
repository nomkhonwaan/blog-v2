package storage

import (
	"context"
	"encoding/json"
	"github.com/nomkhonwaan/myblog/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// File is an uploaded file on the storage server
type File struct {
	// Identifier of the file
	ID primitive.ObjectID `bson:"_id" json:"id"`

	// An uploaded file path
	Path string `bson:"path" json:"path"`

	// An original file name
	FileName string `bson:"fileName" json:"fileName"`

	// An optional field #1 for using in some storage server
	OptionalField1 string `bson:"optionalField1" json:"optionalField1,omitempty"`

	// An optional field #2 for using in some storage server
	OptionalField2 string `bson:"optionalField2" json:"optionalField2,omitempty"`

	// An optional field #3 for using in some storage server
	OptionalField3 string `bson:"optionalField3" json:"optionalField3,omitempty"`

	// Date-time that the file was created
	CreatedAt time.Time `bson:"createdAt" json:"createdAt" graphql:"createdAt"`

	// Date-time that the file was updated
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt" graphql:"updatedAt"`
}

// MarshalJSON is a custom JSON marshaling function of file entity
func (f File) MarshalJSON() ([]byte, error) {
	type Alias File
	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    f.ID.Hex(),
		Alias: (*Alias)(&f),
	})
}

// FileRepository is a repository interface of file which defines all file entity related functions
type FileRepository interface {
	// Create or replace (if the file name is exists)
	Create(ctx context.Context, file File) (File, error)
}

// NewFileRepository returns file repository which connects to MongoDB
func NewFileRepository(col mongo.Collection) MongoFileRepository {
	return MongoFileRepository{col}
}

// MongoFileRepository is a MongoDB specified repository for file
type MongoFileRepository struct {
	col mongo.Collection
}

func (repo MongoFileRepository) Create(ctx context.Context, file File) (File, error) {
	if file.ID.IsZero() {
		file.ID = primitive.NewObjectID()
	}
	file.CreatedAt = time.Now()

	doc, _ := bson.Marshal(file)
	_, err := repo.col.InsertOne(ctx, doc)
	if err != nil {
		return File{}, err
	}

	return file, nil
}
