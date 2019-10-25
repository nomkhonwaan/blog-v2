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

func TestMarshalTagJSON(t *testing.T) {
	// Given
	id := primitive.NewObjectID()
	tag := Tag{
		ID:   id,
		Name: "GraphQL",
		Slug: "graphql-" + id.Hex(),
	}

	// When
	result, err := json.Marshal(tag)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, "{\"id\":\""+id.Hex()+"\",\"name\":\"GraphQL\",\"slug\":\"graphql-"+id.Hex()+"\"}", string(result))

}
func TestMongoTagRepository_FindAllByIDs(t *testing.T) {
	t.Run("With successful finding all tags by list of IDs", func(t *testing.T) {
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

		repo := NewTagRepository(col)

		// When
		_, err := repo.FindAllByIDs(ctx, ids)

		// Then
		assert.Nil(t, err)
	})

	t.Run("When unable to find all tags by list of IDs", func(t *testing.T) {
		// Given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		col := mock_mongo.NewMockCollection(ctrl)
		ctx := context.Background()
		ids := []primitive.ObjectID{primitive.NewObjectID()}

		col.EXPECT().Find(gomock.Any(), gomock.Any()).Return(nil, errors.New("test unable to find all tags by list of IDs"))

		repo := NewTagRepository(col)

		// When
		_, err := repo.FindAllByIDs(ctx, ids)

		// Then
		assert.EqualError(t, err, "test unable to find all tags by list of IDs")
	})
}
