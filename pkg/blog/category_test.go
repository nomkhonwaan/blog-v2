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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		cur = mock_mongo.NewMockCursor(ctrl)
		col = mock_mongo.NewMockCollection(ctrl)
	)

	t.Run("With successful finding all categories", func(t *testing.T) {
		// Given
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

	t.Run("When an error has occurred while finding all categories", func(t *testing.T) {
		// Given
		ctx := context.Background()

		col.EXPECT().Find(ctx, bson.D{}).Return(nil, errors.New("test find all categories error"))

		repo := NewCategoryRepository(col)

		// When
		_, err := repo.FindAll(ctx)

		// Then
		assert.EqualError(t, err, "test find all categories error")
	})
}

func TestMongoCategoryRepository_FindAllByIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		cur = mock_mongo.NewMockCursor(ctrl)
		col = mock_mongo.NewMockCollection(ctrl)
	)

	t.Run("With successful finding all categories by list of IDs", func(t *testing.T) {
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

		repo := NewCategoryRepository(col)

		// When
		_, err := repo.FindAllByIDs(ctx, ids)

		// Then
		assert.Nil(t, err)
	})

	t.Run("When unable to find all categories by list of IDs", func(t *testing.T) {
		// Given
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

func TestMongoCategoryRepository_FindByID(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	singleResult := mock_mongo.NewMockSingleResult(ctrl)
	col := mock_mongo.NewMockCollection(ctrl)

	ctx := context.Background()
	repo := NewCategoryRepository(col)

	tests := map[string]struct {
		id  interface{}
		err error
	}{
		"With existing tag ID": {
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
