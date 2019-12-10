package blog_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	. "github.com/nomkhonwaan/myblog/pkg/blog"
	mock_log "github.com/nomkhonwaan/myblog/pkg/log/mock"
	"github.com/nomkhonwaan/myblog/pkg/mongo"
	mock_mongo "github.com/nomkhonwaan/myblog/pkg/mongo/mock"
	"github.com/nomkhonwaan/myblog/pkg/storage"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mgo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func TestPost_MarshalJSON(t *testing.T) {
	// Given
	id := primitive.NewObjectID()
	createdAt := time.Now()
	post := Post{
		ID:          id,
		Title:       "Children of Dune",
		Slug:        "children-of-dune-" + id.Hex(),
		Status:      Draft,
		Markdown:    "Integer tincidunt ante vel ipsum. Praesent blandit lacinia erat. Vestibulum sed magna at nunc commodo placerat. Praesent blandit. Nam nulla. Integer pede justo, lacinia eget, tincidunt eget, tempus vel, pede. Morbi porttitor lorem id ligula. Suspendisse ornare consequat lectus. In est risus, auctor sed, tristique in, tempus sit amet, sem.",
		HTML:        "Nullam sit amet turpis elementum ligula vehicula consequat. Morbi a ipsum. Integer a nibh.",
		PublishedAt: time.Time{},
		AuthorID:    "github|c7834cb0-2b79-4d27-a817-520a6420c11b",
		Categories:  []mongo.DBRef{},
		Tags:        []mongo.DBRef{},
		CreatedAt:   createdAt,
		UpdatedAt:   time.Time{},
	}

	// When
	result, err := json.Marshal(post)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, "{\"id\":\""+id.Hex()+"\",\"title\":\"Children of Dune\",\"slug\":\"children-of-dune-"+id.Hex()+"\",\"status\":\"DRAFT\",\"markdown\":\"Integer tincidunt ante vel ipsum. Praesent blandit lacinia erat. Vestibulum sed magna at nunc commodo placerat. Praesent blandit. Nam nulla. Integer pede justo, lacinia eget, tincidunt eget, tempus vel, pede. Morbi porttitor lorem id ligula. Suspendisse ornare consequat lectus. In est risus, auctor sed, tristique in, tempus sit amet, sem.\",\"html\":\"Nullam sit amet turpis elementum ligula vehicula consequat. Morbi a ipsum. Integer a nibh.\",\"publishedAt\":\"0001-01-01T00:00:00Z\",\"authorId\":\"github|c7834cb0-2b79-4d27-a817-520a6420c11b\",\"engagement\":{\"shareCount\":0},\"createdAt\":\""+createdAt.Format(time.RFC3339Nano)+"\",\"updatedAt\":\"0001-01-01T00:00:00Z\"}", string(result))
}

func TestMongoPostRepository_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	col := mock_mongo.NewMockCollection(ctrl)
	timer := mock_log.NewMockTimer(ctrl)

	t.Run("When insert into the collection successfully", func(t *testing.T) {
		// Given
		ctx := context.Background()
		now := time.Now()
		authorID := "github|303589"

		timer.EXPECT().Now().Return(now)
		col.EXPECT().InsertOne(ctx, gomock.Any()).Return(&mgo.InsertOneResult{}, nil)

		postRepo := NewPostRepository(col, timer)

		// When
		result, err := postRepo.Create(ctx, authorID)

		// Then
		assert.Nil(t, err)
		assert.Equal(t, now, result.CreatedAt)
		assert.Equal(t, fmt.Sprintf("%s", result.ID.Hex()), result.Slug)
		assert.Equal(t, Draft, result.Status)
		assert.Equal(t, authorID, result.AuthorID)
	})

	t.Run("When insert into the collection un-successfully", func(t *testing.T) {
		// Given
		ctx := context.Background()
		now := time.Now()
		authorID := "github|303589"

		timer.EXPECT().Now().Return(now)
		col.EXPECT().InsertOne(ctx, gomock.Any()).Return(&mgo.InsertOneResult{}, errors.New("something went wrong"))

		postRepo := NewPostRepository(col, timer)

		expected := Post{}

		//result When
		result, err := postRepo.Create(ctx, authorID)

		// Then
		assert.EqualError(t, err, "something went wrong")
		assert.Equal(t, expected, result)
	})
}

func TestMongoPostRepository_FindAll(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cur := mock_mongo.NewMockCursor(ctrl)
	col := mock_mongo.NewMockCollection(ctrl)

	ctx := context.Background()
	repo := NewPostRepository(col, nil)
	catID := primitive.NewObjectID()
	tagID := primitive.NewObjectID()

	tests := map[string]struct {
		q       PostQuery
		filter  interface{}
		options func() *options.FindOptions
		err     error
	}{
		"With default query options": {
			q:      NewPostQueryBuilder().Build(),
			filter: bson.M{},
			options: func() *options.FindOptions {
				return (&options.FindOptions{}).SetSkip(0).SetLimit(5)
			},
		},
		"With specified offset and limit": {
			q:      NewPostQueryBuilder().WithOffset(10).WithLimit(5).Build(),
			filter: bson.M{},
			options: func() *options.FindOptions {
				return (&options.FindOptions{}).SetSkip(10).SetLimit(5)
			},
		},
		"With status draft": {
			q:      NewPostQueryBuilder().WithStatus(Draft).Build(),
			filter: bson.M{"status": Draft},
			options: func() *options.FindOptions {
				return (&options.FindOptions{}).SetSkip(0).SetLimit(5)
			},
		},
		"With status published": {
			q:      NewPostQueryBuilder().WithStatus(Published).Build(),
			filter: bson.M{"status": Published},
			options: func() *options.FindOptions {
				opts := (&options.FindOptions{}).SetSkip(0).SetLimit(5)
				opts.Sort = map[string]interface{}{
					"publishedAt": -1,
				}
				return opts
			},
		},
		"With specific category": {
			q:      NewPostQueryBuilder().WithCategory(Category{ID: catID}).Build(),
			filter: bson.M{"categories.$id": catID},
			options: func() *options.FindOptions {
				return (&options.FindOptions{}).SetSkip(0).SetLimit(5)
			},
		},
		"With specific tag": {
			q:      NewPostQueryBuilder().WithTag(Tag{ID: tagID}).Build(),
			filter: bson.M{"tags.$id": tagID},
			options: func() *options.FindOptions {
				return (&options.FindOptions{}).SetSkip(0).SetLimit(5)
			},
		},
		"When an error has occurred while finding the result": {
			q:      PostQuery{},
			filter: bson.M{},
			options: func() *options.FindOptions {
				return (&options.FindOptions{}).SetSkip(0).SetLimit(0)
			},
			err: errors.New("something went wrong"),
		},
	}

	// When
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			col.EXPECT().Find(ctx, test.filter, test.options()).Return(cur, test.err)

			if test.err == nil {
				cur.EXPECT().Close(ctx).Return(nil)
				cur.EXPECT().Decode(gomock.Any()).Return(nil)

				_, err := repo.FindAll(ctx, test.q)
				assert.Nil(t, err)
			} else {
				_, err := repo.FindAll(ctx, test.q)
				assert.EqualError(t, err, test.err.Error())
			}
		})
	}

	// Then
}

func TestMongoPostRepository_FindByID(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	singleResult := mock_mongo.NewMockSingleResult(ctrl)
	col := mock_mongo.NewMockCollection(ctrl)

	ctx := context.Background()
	repo := NewPostRepository(col, nil)

	tests := map[string]struct {
		id  interface{}
		err error
	}{
		"With existing post ID": {
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

func TestMongoPostRepository_Save(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	singleResult := mock_mongo.NewMockSingleResult(ctrl)
	col := mock_mongo.NewMockCollection(ctrl)
	timer := mock_log.NewMockTimer(ctrl)

	ctx := context.Background()
	now := time.Now()
	publishedAt := time.Now()
	repo := NewPostRepository(col, timer)
	catID := primitive.NewObjectID()
	tagID := primitive.NewObjectID()
	featuredImageID := primitive.NewObjectID()
	attachmentID := primitive.NewObjectID()

	tests := map[string]struct {
		q      PostQuery
		id     interface{}
		update interface{}
		err    error
	}{
		"With default query options": {
			q:      NewPostQueryBuilder().Build(),
			id:     primitive.NewObjectID(),
			update: bson.M{"$set": bson.M{"updatedAt": now}},
		},
		"When updating post's title": {
			q:      NewPostQueryBuilder().WithTitle("Test update post title").Build(),
			id:     primitive.NewObjectID(),
			update: bson.M{"$set": bson.M{"title": "Test update post title", "updatedAt": now}},
		},
		"When updating post's slug": {
			q:      NewPostQueryBuilder().WithSlug("test-update-post-slug").Build(),
			id:     primitive.NewObjectID(),
			update: bson.M{"$set": bson.M{"slug": "test-update-post-slug", "updatedAt": now}},
		},
		"When updating post's status": {
			q:      NewPostQueryBuilder().WithStatus(Published).Build(),
			id:     primitive.NewObjectID(),
			update: bson.M{"$set": bson.M{"status": Published, "updatedAt": now}},
		},
		"When updating post's content": {
			q:      NewPostQueryBuilder().WithMarkdown("Test update post content").WithHTML("<p>Test update post content</p>").Build(),
			id:     primitive.NewObjectID(),
			update: bson.M{"$set": bson.M{"markdown": "Test update post content", "html": "<p>Test update post content</p>", "updatedAt": now}},
		},
		"When updating post's published date-time": {
			q:      NewPostQueryBuilder().WithPublishedAt(publishedAt).Build(),
			id:     primitive.NewObjectID(),
			update: bson.M{"$set": bson.M{"publishedAt": publishedAt, "updatedAt": now}},
		},
		"When updating post's categories": {
			q:      NewPostQueryBuilder().WithCategories([]Category{{ID: catID, Name: "Web Development", Slug: "web-development-" + catID.Hex()}}).Build(),
			id:     primitive.NewObjectID(),
			update: bson.M{"$set": bson.M{"categories": bson.A{mongo.DBRef{Ref: "categories", ID: catID}}, "updatedAt": now}},
		},
		"When updating post's tags": {
			q:      NewPostQueryBuilder().WithTags([]Tag{{ID: tagID, Name: "Blog", Slug: "blog-" + tagID.Hex()}}).Build(),
			id:     primitive.NewObjectID(),
			update: bson.M{"$set": bson.M{"tags": bson.A{mongo.DBRef{Ref: "tags", ID: tagID}}, "updatedAt": now}},
		},
		"When updating post's featured image": {
			q:      NewPostQueryBuilder().WithFeaturedImage(storage.File{ID: featuredImageID, Slug: fmt.Sprintf("test-featured-image-%s.jpg", featuredImageID.Hex())}).Build(),
			id:     primitive.NewObjectID(),
			update: bson.M{"$set": bson.M{"featuredImage": mongo.DBRef{Ref: "files", ID: featuredImageID}, "updatedAt": now}},
		},
		"When updating post's attachments": {
			q:      NewPostQueryBuilder().WithAttachments([]storage.File{{ID: attachmentID}}).Build(),
			id:     primitive.NewObjectID(),
			update: bson.M{"$set": bson.M{"attachments": bson.A{mongo.DBRef{Ref: "files", ID: attachmentID}}, "updatedAt": now}},
		},
		"When an error has occurred while updating the post": {
			q:      NewPostQueryBuilder().Build(),
			id:     primitive.NewObjectID(),
			update: bson.M{"$set": bson.M{"updatedAt": now}},
			err:    errors.New("something went wrong"),
		},
	}

	// When
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			timer.EXPECT().Now().Return(now)
			col.EXPECT().UpdateOne(ctx, bson.M{"_id": test.id.(primitive.ObjectID)}, test.update).Return(nil, test.err)

			if test.err == nil {
				col.EXPECT().FindOne(ctx, bson.M{"_id": test.id.(primitive.ObjectID)}).Return(singleResult)
				singleResult.EXPECT().Decode(gomock.Any()).Return(nil)

				_, err := repo.Save(ctx, test.id, test.q)
				assert.Nil(t, err)
			} else {
				_, err := repo.Save(ctx, test.id, test.q)
				assert.EqualError(t, err, test.err.Error())
			}
		})
	}
}
