package storage_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	mock_mongo "github.com/nomkhonwaan/myblog/pkg/mongo/mock"
	. "github.com/nomkhonwaan/myblog/pkg/storage"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mgo "go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
)

func TestFile_MarshalJSON(t *testing.T) {
	// Given
	id := primitive.NewObjectID()
	createdAt := time.Now()
	file := File{
		ID:             id,
		Path:           "/path/to/the/file.txt",
		FileName:       "file.txt",
		Slug:           fmt.Sprintf("file-%s.txt", id.Hex()),
		OptionalField1: "",
		OptionalField2: "",
		OptionalField3: "",
		CreatedAt:      createdAt,
		UpdatedAt:      time.Time{},
	}

	// When
	result, err := json.Marshal(file)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, "{\"id\":\""+id.Hex()+"\",\"path\":\"/path/to/the/file.txt\",\"fileName\":\"file.txt\",\"slug\":\"file-"+id.Hex()+".txt\",\"createdAt\":\""+createdAt.Format(time.RFC3339Nano)+"\",\"updatedAt\":\"0001-01-01T00:00:00Z\"}", string(result))
}

func TestMongoFileRepository_Create(t *testing.T) {
	t.Run("When insert into the collection successfully", func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		col := mock_mongo.NewMockCollection(ctrl)
		ctx := context.Background()
		id := primitive.NewObjectID()
		path := "/path/to/the/file.txt"
		fileName := "file.txt"

		file := File{
			ID:       id,
			Path:     path,
			FileName: fileName,
		}

		col.EXPECT().InsertOne(ctx, gomock.Any()).Return(&mgo.InsertOneResult{}, nil)

		fileRepo := NewFileRepository(col)

		// When
		result, err := fileRepo.Create(ctx, file)

		// Then
		assert.Nil(t, err)
		assert.Equal(t, id, result.ID)
		assert.Equal(t, path, result.Path)
		assert.Equal(t, fileName, result.FileName)
		assert.True(t, time.Since(result.CreatedAt) < time.Minute)
	})

	t.Run("With empty ID field", func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		col := mock_mongo.NewMockCollection(ctrl)
		ctx := context.Background()
		path := "/path/to/the/file.txt"
		fileName := "file.txt"

		file := File{
			Path:     path,
			FileName: fileName,
		}

		col.EXPECT().InsertOne(ctx, gomock.Any()).Return(&mgo.InsertOneResult{}, nil)

		fileRepo := NewFileRepository(col)

		// When
		result, err := fileRepo.Create(ctx, file)

		// Then
		assert.Nil(t, err)
		assert.NotEmpty(t, result.ID)
	})

	t.Run("When insert into the collection un-successfully", func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		col := mock_mongo.NewMockCollection(ctrl)
		ctx := context.Background()
		path := "/path/to/the/file.txt"
		fileName := "file.txt"

		file := File{
			Path:     path,
			FileName: fileName,
		}

		col.EXPECT().InsertOne(ctx, gomock.Any()).Return(&mgo.InsertOneResult{}, errors.New("something went wrong"))

		fileRepo := NewFileRepository(col)

		expected := File{}

		// When
		result, err := fileRepo.Create(ctx, file)

		// Then
		assert.EqualError(t, err, "something went wrong")
		assert.Equal(t, expected, result)
	})
}
