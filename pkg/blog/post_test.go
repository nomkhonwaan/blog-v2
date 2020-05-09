package blog

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/nomkhonwaan/myblog/pkg/mongo"
	mock_mongo "github.com/nomkhonwaan/myblog/pkg/mongo/mock"
	"github.com/nomkhonwaan/myblog/pkg/storage"
	"github.com/stretchr/testify/assert"
	"github.com/tkuchiki/faketime"
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
		Status:      StatusDraft,
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

	now := time.Date(2020, 3, 29, 18, 57, 0, 0, time.UTC)
	f := faketime.NewFaketimeWithTime(now)
	defer f.Undo()
	f.Do()

	var (
		col = mock_mongo.NewMockCollection(ctrl)
	)

	repo := MongoPostRepository{col: col}

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
		assert.Equal(t, fmt.Sprintf("%s", result.ID.Hex()), result.Slug)
		assert.Equal(t, StatusDraft, result.Status)
		assert.Equal(t, authorID, result.AuthorID)
	})

	t.Run("When unable to create a new record on database", func(t *testing.T) {
		// Given
		authorID := "github|303589"

		col.EXPECT().InsertOne(gomock.Any(), gomock.Any()).Return(nil, errors.New("test unable to create a new record on database"))

		expected := Post{}

		// When
		result, err := repo.Create(context.Background(), authorID)

		// Then
		assert.EqualError(t, err, "test unable to create a new record on database")
		assert.Equal(t, expected, result)
	})
}

func TestMongoPostRepository_FindAll(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		col = mock_mongo.NewMockCollection(ctrl)
		cur = mock_mongo.NewMockCursor(ctrl)
	)

	ctx := context.Background()
	repo := MongoPostRepository{col: col}
	authorizedID := "authorizedID"
	published := StatusPublished
	draft := StatusDraft
	catID := primitive.NewObjectID()
	tagID := primitive.NewObjectID()

	tests := map[string]struct {
		q       PostQuery
		filter  interface{}
		options *options.FindOptions
		err     error
	}{
		"With default query options": {
			q:      NewPostQueryBuilder().Build(),
			filter: bson.M{},
			options: options.Find().
				SetSort(bson.D{
					{"status", 1},
					{"createdAt", -1},
				}).
				SetSkip(0).
				SetLimit(5),
		},
		"With specified offset and limit": {
			q:      NewPostQueryBuilder().WithOffset(10).WithLimit(5).Build(),
			filter: bson.M{},
			options: options.Find().
				SetSort(bson.D{
					{"status", 1},
					{"createdAt", -1},
				}).
				SetSkip(10).
				SetLimit(5),
		},
		"With status draft": {
			q:      NewPostQueryBuilder().WithStatus(draft).Build(),
			filter: bson.M{"status": &draft},
			options: options.Find().
				SetSort(bson.D{{"createdAt", -1}}).
				SetSkip(0).
				SetLimit(5),
		},
		"With status published": {
			q:      NewPostQueryBuilder().WithStatus(published).Build(),
			filter: bson.M{"status": &published},
			options: options.Find().
				SetSort(bson.D{{"publishedAt", -1}}).
				SetSkip(0).
				SetLimit(5),
		},
		"With specific authorID": {
			q:      NewPostQueryBuilder().WithAuthorID(authorizedID).Build(),
			filter: bson.M{"authorId": &authorizedID},
			options: options.Find().
				SetSort(bson.D{
					{"status", 1},
					{"createdAt", -1},
				}).
				SetSkip(0).
				SetLimit(5),
		},
		"With specific category": {
			q:      NewPostQueryBuilder().WithCategory(Category{ID: catID}).Build(),
			filter: bson.M{"categories.$id": catID},
			options: options.Find().
				SetSort(bson.D{
					{"status", 1},
					{"createdAt", -1},
				}).
				SetSkip(0).
				SetLimit(5),
		},
		"With specific tag": {
			q:      NewPostQueryBuilder().WithTag(Tag{ID: tagID}).Build(),
			filter: bson.M{"tags.$id": tagID},
			options: options.Find().
				SetSort(bson.D{
					{"status", 1},
					{"createdAt", -1},
				}).
				SetSkip(0).
				SetLimit(5),
		},
		"When an error has occurred while finding the result": {
			q:      NewPostQueryBuilder().Build(),
			filter: bson.M{},
			options: options.Find().
				SetSort(bson.D{
					{"status", 1},
					{"createdAt", -1},
				}).
				SetSkip(0).
				SetLimit(5),
			err: errors.New("something went wrong"),
		},
	}

	// When
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			col.EXPECT().Find(ctx, gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, filter interface{}, opts *options.FindOptions) (mongo.Cursor, error) {
				assert.EqualValues(t, test.filter, filter)
				assert.EqualValues(t, test.options, opts)

				return cur, test.err
			})

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

	var (
		col          = mock_mongo.NewMockCollection(ctrl)
		singleResult = mock_mongo.NewMockSingleResult(ctrl)
	)

	ctx := context.Background()
	repo := MongoPostRepository{col: col}

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

	now := time.Date(2020, 3, 29, 18, 57, 0, 0, time.UTC)
	f := faketime.NewFaketimeWithTime(now)
	defer f.Undo()
	f.Do()

	var (
		col          = mock_mongo.NewMockCollection(ctrl)
		singleResult = mock_mongo.NewMockSingleResult(ctrl)
	)

	ctx := context.Background()
	repo := MongoPostRepository{col: col}
	publishedAt := time.Now()
	slug := "test-update-post-slug"
	published := StatusPublished
	title := "Test update post title"
	markdown := "Test update post content"
	html := "<p>Test update post content</p>"
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
			q:      NewPostQueryBuilder().WithTitle(title).Build(),
			id:     primitive.NewObjectID(),
			update: bson.M{"$set": bson.M{"title": &title, "updatedAt": now}},
		},
		"When updating post's slug": {
			q:      NewPostQueryBuilder().WithSlug(slug).Build(),
			id:     primitive.NewObjectID(),
			update: bson.M{"$set": bson.M{"slug": &slug, "updatedAt": now}},
		},
		"When updating post's status": {
			q:      NewPostQueryBuilder().WithStatus(published).Build(),
			id:     primitive.NewObjectID(),
			update: bson.M{"$set": bson.M{"status": &published, "updatedAt": now}},
		},
		"When updating post's content": {
			q:      NewPostQueryBuilder().WithMarkdown(markdown).WithHTML(html).Build(),
			id:     primitive.NewObjectID(),
			update: bson.M{"$set": bson.M{"markdown": &markdown, "html": &html, "updatedAt": now}},
		},
		"When updating post's published date-time": {
			q:      NewPostQueryBuilder().WithPublishedAt(publishedAt).Build(),
			id:     primitive.NewObjectID(),
			update: bson.M{"$set": bson.M{"publishedAt": &publishedAt, "updatedAt": now}},
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
}
