package discussion

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/nomkhonwaan/myblog/pkg/mongo"
	mock_mongo "github.com/nomkhonwaan/myblog/pkg/mongo/mock"
	"github.com/stretchr/testify/assert"
	"github.com/tkuchiki/faketime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mgo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func TestComment_MarshalJSON(t *testing.T) {
	// Given
	now := time.Now()
	id := primitive.NewObjectID()
	parentID := primitive.NewObjectID()
	comment := Comment{
		ID:        id,
		Parent:    mongo.DBRef{Ref: "comments", ID: parentID},
		AuthorID:  "github|c7834cb0-2b79-4d27-a817-520a6420c11b",
		Text:      "Just a sample comment",
		CreatedAt: now,
	}

	// When
	result, err := json.Marshal(comment)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, "{\"id\":\""+id.Hex()+"\",\"authorId\":\"github|c7834cb0-2b79-4d27-a817-520a6420c11b\",\"text\":\"Just a sample comment\",\"createdAt\":\""+now.Format(time.RFC3339Nano)+"\",\"updatedAt\":\"0001-01-01T00:00:00Z\"}", string(result))
}

func TestMongoCommentRepository_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now := time.Date(2020, 4, 26, 12, 31, 0, 0, time.UTC)
	f := faketime.NewFaketimeWithTime(now)
	defer f.Undo()
	f.Do()

	var (
		col = mock_mongo.NewMockCollection(ctrl)
	)

	repo := MongoCommentRepository{col: col}

	t.Run("With successful creating a new record", func(t *testing.T) {
		// Given
		ctx := context.Background()
		authorID := "github|303589"

		col.EXPECT().InsertOne(ctx, gomock.Any()).Return(&mgo.InsertOneResult{}, nil)

		// When
		result, err := repo.Create(ctx, authorID)

		// Then
		assert.Nil(t, err)
		assert.Equal(t, now, result.CreatedAt)
		assert.Equal(t, authorID, result.AuthorID)
	})

	t.Run("When unable to create a new record on database", func(t *testing.T) {
		// Given
		authorID := "github|303589"

		col.EXPECT().InsertOne(gomock.Any(), gomock.Any()).Return(nil, errors.New("test unable to create a new record on database"))

		expected := Comment{}

		// When
		result, err := repo.Create(context.Background(), authorID)

		// Then
		assert.EqualError(t, err, "test unable to create a new record on database")
		assert.Equal(t, expected, result)
	})
}

func TestMongoCommentRepository_FindAllByIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		col = mock_mongo.NewMockCollection(ctrl)
		cur = mock_mongo.NewMockCursor(ctrl)
	)

	repo := MongoCommentRepository{col: col}

	t.Run("With successful finding all children by IDs", func(t *testing.T) {
		// Given
		ctx := context.Background()
		ids := []primitive.ObjectID{primitive.NewObjectID()}
		filter := bson.M{"_id": bson.M{"$in": ids}}
		opts := options.Find().SetSort(bson.D{{"createdAt", -1}})

		col.EXPECT().Find(ctx, filter, opts).Return(cur, nil)
		cur.EXPECT().Close(ctx).Return(nil)
		cur.EXPECT().Decode(gomock.Any()).Return(nil)

		// When
		_, err := repo.FindAllByIDs(ctx, ids)

		// Then
		assert.Nil(t, err)
	})

	t.Run("When an error has occurred while finding all children by IDs", func(t *testing.T) {
		// Given
		ids := []primitive.ObjectID{primitive.NewObjectID()}

		col.EXPECT().Find(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("test unable to find all children by list of IDs"))

		// When
		_, err := repo.FindAllByIDs(context.Background(), ids)

		// Then
		assert.EqualError(t, err, "test unable to find all children by list of IDs")
	})
}

func TestMongoCommentRepository_FindByID(t *testing.T) {

}

func TestMongoCommentRepository_Save(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now := time.Date(2020, 4, 26, 12, 44, 0, 0, time.UTC)
	f := faketime.NewFaketimeWithTime(now)
	defer f.Undo()
	f.Do()

	var (
		col          = mock_mongo.NewMockCollection(ctrl)
		singleResult = mock_mongo.NewMockSingleResult(ctrl)
	)

	ctx := context.Background()
	repo := MongoCommentRepository{col: col}
	text := "Test update comment content"
	parent := Comment{ID: primitive.NewObjectID()}
	child := Comment{ID: primitive.NewObjectID()}

	tests := map[string]struct {
		q      CommentQuery
		id     interface{}
		update interface{}
		err    error
	}{
		"With default query options": {
			q:      NewCommentQueryBuilder().Build(),
			id:     primitive.NewObjectID(),
			update: bson.M{"$set": bson.M{"updatedAt": now}},
		},
		"When updating comment's content": {
			q:      NewCommentQueryBuilder().WithText(text).Build(),
			id:     primitive.NewObjectID(),
			update: bson.M{"$set": bson.M{"text": &text, "updatedAt": now}},
		},
		"When updating comment's parent": {
			q:      NewCommentQueryBuilder().WithParent(parent).Build(),
			id:     primitive.NewObjectID(),
			update: bson.M{"$set": bson.M{"parent": mongo.DBRef{Ref: "comments", ID: parent.ID}, "updatedAt": now}},
		},
		"When updating comment's children": {
			q:      NewCommentQueryBuilder().WithChildren([]Comment{child}).Build(),
			id:     primitive.NewObjectID(),
			update: bson.M{"$set": bson.M{"children": bson.A{mongo.DBRef{Ref: "comments", ID: child.ID}}, "updatedAt": now}},
		},
		"When an error has occurred while updating the comment": {
			q:      NewCommentQueryBuilder().Build(),
			id:     primitive.NewObjectID(),
			update: bson.M{"$set": bson.M{"updatedAt": now}},
			err:    errors.New("something went wrong"),
		},
	}

	// When
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			col.EXPECT().UpdateOne(ctx, bson.M{"_id": test.id.(primitive.ObjectID)}, test.update).Return(nil, test.err)

			if test.err == nil {
				col.EXPECT().FindOne(ctx, bson.M{"_id": test.id.(primitive.ObjectID)}).Return(singleResult)
				singleResult.EXPECT().Decode(gomock.Any()).Return(nil)

				_, err := repo.Save(ctx, test.id, test.q)
				assert.Nil(t, err)
			} else {
				_, err := repo.Save(ctx, test.id, test.q)
				assert.Equal(t, err, test.err)
			}
		})
	}

	// Then
}
