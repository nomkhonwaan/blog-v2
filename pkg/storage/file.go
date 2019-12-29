//go:generate mockgen -destination=./mock/file_mock.go github.com/nomkhonwaan/myblog/pkg/storage FileRepository

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
	ID primitive.ObjectID `bson:"_id" json:"id" graphql:"-"`

	// An uploaded file path
	Path string `bson:"path" json:"path" graphql:"path"`

	// An original file name
	FileName string `bson:"fileName" json:"fileName" graphql:"fileName"`

	// Valid URL string composes with file name and ID
	Slug string `bson:"slug" json:"slug" graphql:"slug"`

	// An optional field #1 for using in some storage server
	OptionalField1 string `bson:"optionalField1" json:"optionalField1,omitempty" graphql:"optionalField1"`

	// An optional field #2 for using in some storage server
	OptionalField2 string `bson:"optionalField2" json:"optionalField2,omitempty" graphql:"optionalField2"`

	// An optional field #3 for using in some storage server
	OptionalField3 string `bson:"optionalField3" json:"optionalField3,omitempty" graphql:"optionalField3"`

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
	// Insert a new file record whether exist or not
	Create(ctx context.Context, file File) (File, error)

	// Delete a file record by its ID
	Delete(ctx context.Context, id interface{}) error

	// Return list of files from list of IDs
	FindAllByIDs(ctx context.Context, ids interface{}) ([]File, error)

	// Return a single file from its ID
	FindByID(ctx context.Context, id interface{}) (File, error)
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

func (repo MongoFileRepository) Delete(ctx context.Context, id interface{}) error {
	_, err := repo.col.DeleteOne(ctx, bson.M{"_id": id.(primitive.ObjectID)})
	return err
}

func (repo MongoFileRepository) FindAllByIDs(ctx context.Context, ids interface{}) ([]File, error) {
	cur, err := repo.col.Find(ctx, bson.M{
		"_id": bson.M{
			"$in": ids.([]primitive.ObjectID),
		},
	})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var files []File
	err = cur.Decode(&files)

	return files, err
}

func (repo MongoFileRepository) FindByID(ctx context.Context, id interface{}) (File, error) {
	r := repo.col.FindOne(ctx, bson.M{"_id": id.(primitive.ObjectID)})

	var file File
	err := r.Decode(&file)

	return file, err
}
