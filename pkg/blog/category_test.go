package blog_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	. "github.com/nomkhonwaan/myblog/pkg/blog"
	mock_mongo "github.com/nomkhonwaan/myblog/pkg/mongo/mock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestCategory_MarshalJSON(t *testing.T) {
	// Given
	id := primitive.NewObjectID()
	cat := Category{
		ID:   id,
		Name: "Test",
		Slug: "test-" + id.Hex(),
	}

	// When
	result, err := json.Marshal(cat)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, "{\"id\":\""+id.Hex()+"\",\"name\":\"Test\",\"slug\":\"test-"+id.Hex()+"\"}", string(result))
}

func TestMongoCategoryRepository_FindAll(t *testing.T) {
	t.Run("With successful finding all categories", func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cur := mock_mongo.NewMockCursor(ctrl)
		col := mock_mongo.NewMockCollection(ctrl)
		ctx := context.Background()

		col.EXPECT().Find(ctx, bson.D{}).Return(cur, nil)
		cur.EXPECT().Close(ctx).Return(nil)
		cur.EXPECT().Decode(gomock.Any()).Return(nil)

		repo := NewCategoryRepository(col)

		// When
		_, err := repo.FindAll(ctx)

		// Then
		assert.Nil(t, err)
	})

	t.Run("With an error has occurred while finding all categories", func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		col := mock_mongo.NewMockCollection(ctrl)
		ctx := context.Background()

		col.EXPECT().Find(ctx, bson.D{}).Return(nil, errors.New("something went wrong"))

		repo := NewCategoryRepository(col)

		// When
		_, err := repo.FindAll(ctx)

		// Then
		assert.EqualError(t, err, "something went wrong")
	})
}

func TestMongoCategoryRepository_FindAllByIDs(t *testing.T) {
	t.Run("With successful finding all categories by list of IDs", func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cur := mock_mongo.NewMockCursor(ctrl)
		col := mock_mongo.NewMockCollection(ctrl)
		ctx := context.Background()
		ids := []primitive.ObjectID{primitive.NewObjectID()}

		col.EXPECT().Find(ctx, bson.M{
			"_id": bson.M{
				"$in": ids,
			},
		}).Return(cur, nil)
		cur.EXPECT().Close(ctx).Return(nil)
		cur.EXPECT().Decode(gomock.Any()).Return(nil)

		repo := NewCategoryRepository(col)

		// When
		_, err := repo.FindAllByIDs(ctx, ids)

		// Then
		assert.Nil(t, err)
	})

	t.Run("When unable to find all categories by list of IDs", func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		col := mock_mongo.NewMockCollection(ctrl)
		ctx := context.Background()
		ids := []primitive.ObjectID{primitive.NewObjectID()}

		col.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, errors.New("test unable to find all categories by list of IDs"))

		repo := NewCategoryRepository(col)

		// When
		_, err := repo.FindAllByIDs(ctx, ids)

		// Then
		assert.EqualError(t, err, "test unable to find all categories by list of IDs")
	})
}
