package blog

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/nomkhonwaan/myblog/pkg/mongo"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestMongoCategoryRepository_FindAll(t *testing.T) {
	t.Run("With successful finding all categories", func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cur := mongo.NewMockCursor(ctrl)
		col := mongo.NewMockCollection(ctrl)
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

		col := mongo.NewMockCollection(ctrl)
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

		cur := mongo.NewMockCursor(ctrl)
		col := mongo.NewMockCollection(ctrl)
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
}
