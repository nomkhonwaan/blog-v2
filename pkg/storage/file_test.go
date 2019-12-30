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
	"go.mongodb.org/mongo-driver/bson"
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

func TestMongoFileRepository_Delete(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		col = mock_mongo.NewMockCollection(ctrl)
	)

	ctx := context.Background()
	repo := NewFileRepository(col)
	id := primitive.NewObjectID()

	col.EXPECT().DeleteOne(ctx, bson.M{"_id": id}).Return(nil, nil)

	// When
	err := repo.Delete(ctx, id)

	// Then
	assert.Nil(t, err)
}

func TestMongoFileRepository_FindAllByIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		col = mock_mongo.NewMockCollection(ctrl)
		cur = mock_mongo.NewMockCursor(ctrl)
	)

	repo := NewFileRepository(col)

	t.Run("With successful finding all files by list of IDs", func(t *testing.T) {
		// Given
		ctx := context.Background()
		ids := []primitive.ObjectID{primitive.NewObjectID()}

		col.EXPECT().Find(ctx, bson.M{
			"_id": bson.M{
				"$in": ids,
			},
		}).Return(cur, nil)
		cur.EXPECT().Close(ctx).Return(nil)
		cur.EXPECT().Decode(gomock.Any()).Return(nil)

		// When
		_, err := repo.FindAllByIDs(ctx, ids)

		// Then
		assert.Nil(t, err)
	})

	t.Run("With empty list of IDs", func(t *testing.T) {
		// Given
		ctx := context.Background()
		var ids []primitive.ObjectID

		// When
		_, err := repo.FindAllByIDs(ctx, ids)

		// Then
		assert.Nil(t, err)
	})

	t.Run("When unable to find all files by list of IDs", func(t *testing.T) {
		// Given
		ctx := context.Background()
		ids := []primitive.ObjectID{primitive.NewObjectID()}

		col.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, errors.New("test unable to find all files by list of IDs"))

		// When
		_, err := repo.FindAllByIDs(ctx, ids)

		// Then
		assert.EqualError(t, err, "test unable to find all files by list of IDs")
	})
}

func TestMongoFileRepository_FindByID(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	singleResult := mock_mongo.NewMockSingleResult(ctrl)
	col := mock_mongo.NewMockCollection(ctrl)

	ctx := context.Background()
	repo := NewFileRepository(col)

	tests := map[string]struct {
		id  interface{}
		err error
	}{
		"With existing file ID": {
			id: primitive.NewObjectID(),
		},
		"When an error has occurred while finding the result": {
			id:  primitive.NewObjectID(),
			err: errors.New("test find by ID error"),
		},
	}

	// When
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			col.EXPECT().FindOne(ctx, bson.M{"_id": test.id.(primitive.ObjectID)}, gomock.Any()).Return(singleResult)
			singleResult.EXPECT().Decode(gomock.Any()).Return(test.err)

			if test.err == nil {
				_, err := repo.FindByID(ctx, test.id)
				assert.Nil(t, err)
			} else {
				_, err := repo.FindByID(ctx, test.id)
				assert.EqualError(t, err, test.err.Error())
			}
		})
	}

	// Then
}
