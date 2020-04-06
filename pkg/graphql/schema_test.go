package graphql

import (
	"bytes"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	mock_http "github.com/nomkhonwaan/myblog/internal/http/mock"
	"github.com/nomkhonwaan/myblog/pkg/blog"
	mock_blog "github.com/nomkhonwaan/myblog/pkg/blog/mock"
	"github.com/nomkhonwaan/myblog/pkg/facebook"
	"github.com/nomkhonwaan/myblog/pkg/mongo"
	"github.com/nomkhonwaan/myblog/pkg/storage"
	mock_storage "github.com/nomkhonwaan/myblog/pkg/storage/mock"
	"github.com/stretchr/testify/assert"
	"github.com/tkuchiki/faketime"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"
	"time"
)

func TestFindCategoryBySlugFieldFunc(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		repository = mock_blog.NewMockCategoryRepository(ctrl)
	)

	id := primitive.NewObjectID()

	repository.EXPECT().FindByID(gomock.Any(), id).Return(blog.Category{Name: "Test", Slug: "test-" + id.Hex()}, nil)

	// When
	c, err := FindCategoryBySlugFieldFunc(repository).(func(context.Context, struct{ Slug Slug }) (blog.Category, error))(context.Background(), struct{ Slug Slug }{Slug: Slug("test-" + id.Hex())})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, blog.Category{Name: "Test", Slug: "test-" + id.Hex()}, c)
}

func TestFindAllCategoriesFieldFunc(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		repository = mock_blog.NewMockCategoryRepository(ctrl)
	)

	id := primitive.NewObjectID()

	repository.EXPECT().FindAll(gomock.Any()).Return([]blog.Category{{Name: "Test", Slug: "test-" + id.Hex()}}, nil)

	// When
	cats, err := FindAllCategoriesFieldFunc(repository).(func(context.Context) ([]blog.Category, error))(context.Background())

	// Then
	assert.Nil(t, err)
	assert.Equal(t, []blog.Category{{Name: "Test", Slug: "test-" + id.Hex()}}, cats)
}

func TestFindAllCategoriesBelongedToPostFieldFunc(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		repository = mock_blog.NewMockCategoryRepository(ctrl)
	)

	id := primitive.NewObjectID()

	repository.EXPECT().FindAllByIDs(gomock.Any(), []primitive.ObjectID{id}).Return([]blog.Category{{Name: "Test", Slug: "test-" + id.Hex()}}, nil)

	// When
	cats, err := FindAllCategoriesBelongedToPostFieldFunc(repository).(func(context.Context, blog.Post) ([]blog.Category, error))(context.Background(), blog.Post{Categories: []mongo.DBRef{{ID: id}}})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, []blog.Category{{Name: "Test", Slug: "test-" + id.Hex()}}, cats)
}

func TestFindTagBySlugFieldFunc(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		repository = mock_blog.NewMockTagRepository(ctrl)
	)

	id := primitive.NewObjectID()

	repository.EXPECT().FindByID(gomock.Any(), id).Return(blog.Tag{Name: "Test", Slug: "test-" + id.Hex()}, nil)

	// When
	tag, err := FindTagBySlugFieldFunc(repository).(func(context.Context, struct{ Slug Slug }) (blog.Tag, error))(context.Background(), struct{ Slug Slug }{Slug: Slug("test-" + id.Hex())})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, blog.Tag{Name: "Test", Slug: "test-" + id.Hex()}, tag)
}

func TestFindAllTagsFieldFunc(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		repository = mock_blog.NewMockTagRepository(ctrl)
	)

	id := primitive.NewObjectID()

	repository.EXPECT().FindAll(gomock.Any()).Return([]blog.Tag{{Name: "Test", Slug: "test-" + id.Hex()}}, nil)

	// When
	tags, err := FindAllTagsFieldFunc(repository).(func(context.Context) ([]blog.Tag, error))(context.Background())

	// Then
	assert.Nil(t, err)
	assert.Equal(t, []blog.Tag{{Name: "Test", Slug: "test-" + id.Hex()}}, tags)
}

func TestFindAllTagsBelongedToPostFieldFunc(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		repository = mock_blog.NewMockTagRepository(ctrl)
	)

	id := primitive.NewObjectID()

	repository.EXPECT().FindAllByIDs(gomock.Any(), []primitive.ObjectID{id}).Return([]blog.Tag{{Name: "Test", Slug: "test-" + id.Hex()}}, nil)

	// When
	tags, err := FindAllTagsBelongedToPostFieldFunc(repository).(func(context.Context, blog.Post) ([]blog.Tag, error))(context.Background(), blog.Post{Tags: []mongo.DBRef{{ID: id}}})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, []blog.Tag{{Name: "Test", Slug: "test-" + id.Hex()}}, tags)
}

func TestFindAllLatestPublishedPostsFieldFunc(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		repository = mock_blog.NewMockPostRepository(ctrl)
	)

	id := primitive.NewObjectID()

	repository.EXPECT().FindAll(gomock.Any(), blog.NewPostQueryBuilder().WithStatus(blog.StatusPublished).WithOffset(0).WithLimit(6).Build()).Return([]blog.Post{{Title: "Test", Slug: "test-" + id.Hex()}}, nil)
	// When
	posts, err := FindAllLatestPublishedPostsFieldFunc(repository).(func(context.Context, struct{ Offset, Limit int64 }) ([]blog.Post, error))(context.Background(), struct{ Offset, Limit int64 }{Offset: 0, Limit: 6})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, []blog.Post{{Title: "Test", Slug: "test-" + id.Hex()}}, posts)
}

func TestFindAllLPPBelongedToCategoryFieldFunc(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		repository = mock_blog.NewMockPostRepository(ctrl)
	)

	id := primitive.NewObjectID()
	postID := primitive.NewObjectID()

	repository.EXPECT().FindAll(gomock.Any(), blog.NewPostQueryBuilder().WithCategory(blog.Category{ID: id}).WithStatus(blog.StatusPublished).WithOffset(0).WithLimit(6).Build()).Return([]blog.Post{{Title: "Test", Slug: "test-" + postID.Hex()}}, nil)

	// When
	posts, err := FindAllLPPBelongedToCategoryFieldFunc(repository).(func(context.Context, blog.Category, struct{ Offset, Limit int64 }) ([]blog.Post, error))(context.Background(), blog.Category{ID: id}, struct{ Offset, Limit int64 }{Offset: 0, Limit: 6})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, []blog.Post{{Title: "Test", Slug: "test-" + postID.Hex()}}, posts)
}

func TestFindAllLPPBelongedToTagFieldFunc(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		repository = mock_blog.NewMockPostRepository(ctrl)
	)

	id := primitive.NewObjectID()
	postID := primitive.NewObjectID()

	repository.EXPECT().FindAll(gomock.Any(), blog.NewPostQueryBuilder().WithTag(blog.Tag{ID: id}).WithStatus(blog.StatusPublished).WithOffset(0).WithLimit(6).Build()).Return([]blog.Post{{Title: "Test", Slug: "test-" + postID.Hex()}}, nil)

	// When
	posts, err := FindAllLPPBelongedToTagFieldFunc(repository).(func(context.Context, blog.Tag, struct{ Offset, Limit int64 }) ([]blog.Post, error))(context.Background(), blog.Tag{ID: id}, struct{ Offset, Limit int64 }{Offset: 0, Limit: 6})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, []blog.Post{{Title: "Test", Slug: "test-" + postID.Hex()}}, posts)
}

func TestFindAllMyPostsFieldFunc(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		repository = mock_blog.NewMockPostRepository(ctrl)
	)

	id := primitive.NewObjectID()

	repository.EXPECT().FindAll(gomock.Any(), blog.NewPostQueryBuilder().WithAuthorID("authorizedID").WithOffset(0).WithLimit(6).Build()).Return([]blog.Post{{Title: "Test", Slug: "test-" + id.Hex(), AuthorID: "authorizedID"}}, nil)

	// When
	posts, err := FindAllMyPostsFieldFunc(repository).(func(context.Context, struct{ Offset, Limit int64 }) ([]blog.Post, error))(context.WithValue(context.Background(), AuthorizedID, "authorizedID"), struct{ Offset, Limit int64 }{Offset: 0, Limit: 6})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, []blog.Post{{Title: "Test", Slug: "test-" + id.Hex(), AuthorID: "authorizedID"}}, posts)
}

func TestFindPostBySlugFieldFunc(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		repository = mock_blog.NewMockPostRepository(ctrl)
	)

	t.Run("With successful finding my own post", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()

		repository.EXPECT().FindByID(gomock.Any(), id).Return(blog.Post{Title: "Test", Slug: "test-" + id.Hex(), AuthorID: "authorizedID"}, nil)

		// When
		p, err := FindPostBySlugFieldFunc(repository).(func(context.Context, struct{ Slug Slug }) (blog.Post, error))(context.WithValue(context.Background(), AuthorizedID, "authorizedID"), struct{ Slug Slug }{Slug: Slug("test-" + id.Hex())})

		// Then
		assert.Nil(t, err)
		assert.Equal(t, blog.Post{Title: "Test", Slug: "test-" + id.Hex(), AuthorID: "authorizedID"}, p)
	})

	t.Run("When unable to find a post", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()

		repository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(blog.Post{}, errors.New("test unable to find a post"))

		// When
		_, err := FindPostBySlugFieldFunc(repository).(func(context.Context, struct{ Slug Slug }) (blog.Post, error))(context.WithValue(context.Background(), AuthorizedID, "authorizedID"), struct{ Slug Slug }{Slug: Slug("test-" + id.Hex())})

		// Then
		assert.EqualError(t, err, "Not Found")
	})

	t.Run("When finding a published post", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()

		repository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(blog.Post{Title: "Test", Slug: "test-" + id.Hex(), Status: blog.StatusPublished}, nil)

		// When
		p, err := FindPostBySlugFieldFunc(repository).(func(context.Context, struct{ Slug Slug }) (blog.Post, error))(context.Background(), struct{ Slug Slug }{Slug: Slug("test-" + id.Hex())})

		// Then
		assert.Nil(t, err)
		assert.Equal(t, blog.Post{Title: "Test", Slug: "test-" + id.Hex(), Status: blog.StatusPublished}, p)
	})

	t.Run("When finding a non-published post", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()

		repository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(blog.Post{Title: "Test", Slug: "test-" + id.Hex()}, nil)

		// When
		_, err := FindPostBySlugFieldFunc(repository).(func(context.Context, struct{ Slug Slug }) (blog.Post, error))(context.Background(), struct{ Slug Slug }{Slug: Slug("test-" + id.Hex())})

		// Then
		assert.EqualError(t, err, "Forbidden")
	})
}

func TestCreatePostFieldFunc(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		repository = mock_blog.NewMockPostRepository(ctrl)
	)

	id := primitive.NewObjectID()

	repository.EXPECT().Create(gomock.Any(), "authorizedID").Return(blog.Post{Title: "Test", Slug: "test-" + id.Hex(), AuthorID: "authorizedID"}, nil)

	// When
	p, err := CreatePostFieldFunc(repository).(func(context.Context) (blog.Post, error))(context.WithValue(context.Background(), AuthorizedID, "authorizedID"))

	// Then
	assert.Nil(t, err)
	assert.Equal(t, blog.Post{Title: "Test", Slug: "test-" + id.Hex(), AuthorID: "authorizedID"}, p)
}

func TestUpdatePostTitleFieldFunc(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		repository = mock_blog.NewMockPostRepository(ctrl)
	)

	t.Run("With successful updating post title", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()

		repository.EXPECT().FindByID(gomock.Any(), id).Return(blog.Post{Title: "Test", Slug: "test-" + id.Hex(), AuthorID: "authorizedID"}, nil)
		repository.EXPECT().Save(gomock.Any(), id, blog.NewPostQueryBuilder().WithTitle("Test2").WithSlug("test2-"+id.Hex()).Build()).Return(blog.Post{Title: "Test2", Slug: "test2-" + id.Hex(), AuthorID: "authorizedID"}, nil)

		// When
		p, err := UpdatePostTitleFieldFunc(repository).(func(context.Context, struct {
			Slug  Slug
			Title string
		}) (blog.Post, error))(context.WithValue(context.Background(), AuthorizedID, "authorizedID"), struct {
			Slug  Slug
			Title string
		}{
			Slug:  Slug("test-" + id.Hex()),
			Title: "Test2",
		})

		// Then
		assert.Nil(t, err)
		assert.Equal(t, blog.Post{Title: "Test2", Slug: "test2-" + id.Hex(), AuthorID: "authorizedID"}, p)
	})

	t.Run("When unable to find a post", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()

		repository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(blog.Post{}, errors.New("test unable to find a post"))

		// When
		_, err := UpdatePostTitleFieldFunc(repository).(func(context.Context, struct {
			Slug  Slug
			Title string
		}) (blog.Post, error))(context.WithValue(context.Background(), AuthorizedID, "authorizedID"), struct {
			Slug  Slug
			Title string
		}{
			Slug:  Slug("test-" + id.Hex()),
			Title: "Test2",
		})

		// Then
		assert.EqualError(t, err, "Not Found")
	})

	t.Run("When try to update other post title", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()

		repository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(blog.Post{Title: "Test", Slug: "test-" + id.Hex(), AuthorID: "authorizedID"}, nil)

		// When
		_, err := UpdatePostTitleFieldFunc(repository).(func(context.Context, struct {
			Slug  Slug
			Title string
		}) (blog.Post, error))(context.Background(), struct {
			Slug  Slug
			Title string
		}{
			Slug:  Slug("test-" + id.Hex()),
			Title: "Test2",
		})

		// Then
		assert.EqualError(t, err, "Forbidden")
	})
}

func TestUpdatePostStatusFieldFunc(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		repository = mock_blog.NewMockPostRepository(ctrl)
	)

	now := time.Date(2020, 4, 6, 9, 42, 0, 0, time.UTC)
	f := faketime.NewFaketimeWithTime(now)
	defer f.Undo()
	f.Do()

	t.Run("With successful updating post status", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()

		repository.EXPECT().FindByID(gomock.Any(), id).Return(blog.Post{Title: "Test", Slug: "test-" + id.Hex(), Status: blog.StatusDraft, AuthorID: "authorizedID"}, nil)
		repository.EXPECT().Save(gomock.Any(), id, blog.NewPostQueryBuilder().WithStatus(blog.StatusPublished).WithPublishedAt(now).Build()).Return(blog.Post{Title: "Test2", Slug: "test2-" + id.Hex(), Status: blog.StatusPublished, AuthorID: "authorizedID", PublishedAt: now}, nil)

		// When
		p, err := UpdatePostStatusFieldFunc(repository).(func(context.Context, struct {
			Slug   Slug
			Status blog.Status
		}) (blog.Post, error))(context.WithValue(context.Background(), AuthorizedID, "authorizedID"), struct {
			Slug   Slug
			Status blog.Status
		}{
			Slug:   Slug("test-" + id.Hex()),
			Status: blog.StatusPublished,
		})

		// Then
		assert.Nil(t, err)
		assert.Equal(t, blog.Post{Title: "Test2", Slug: "test2-" + id.Hex(), Status: blog.StatusPublished, AuthorID: "authorizedID", PublishedAt: now}, p)
	})

	t.Run("When unable to find a post", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()

		repository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(blog.Post{}, errors.New("test unable to find a post"))

		// When
		_, err := UpdatePostStatusFieldFunc(repository).(func(context.Context, struct {
			Slug   Slug
			Status blog.Status
		}) (blog.Post, error))(context.WithValue(context.Background(), AuthorizedID, "authorizedID"), struct {
			Slug   Slug
			Status blog.Status
		}{
			Slug:   Slug("test-" + id.Hex()),
			Status: blog.StatusPublished,
		})

		// Then
		assert.EqualError(t, err, "Not Found")
	})

	t.Run("When updating already published post", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()

		repository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(blog.Post{Title: "Test", Slug: "test-" + id.Hex(), Status: blog.StatusPublished, AuthorID: "authorizedID", PublishedAt: now}, nil)
		repository.EXPECT().Save(gomock.Any(), gomock.Any(), blog.NewPostQueryBuilder().WithStatus(blog.StatusPublished).Build()).Return(blog.Post{Title: "Test2", Slug: "test2-" + id.Hex(), Status: blog.StatusPublished, AuthorID: "authorizedID", PublishedAt: now}, nil)

		// When
		p, err := UpdatePostStatusFieldFunc(repository).(func(context.Context, struct {
			Slug   Slug
			Status blog.Status
		}) (blog.Post, error))(context.WithValue(context.Background(), AuthorizedID, "authorizedID"), struct {
			Slug   Slug
			Status blog.Status
		}{
			Slug:   Slug("test-" + id.Hex()),
			Status: blog.StatusPublished,
		})

		// Then
		assert.Nil(t, err)
		assert.Equal(t, blog.Post{Title: "Test2", Slug: "test2-" + id.Hex(), Status: blog.StatusPublished, AuthorID: "authorizedID", PublishedAt: now}, p)
	})

	t.Run("When try to update other post status", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()

		repository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(blog.Post{Title: "Test", Slug: "test-" + id.Hex(), AuthorID: "authorizedID"}, nil)

		// When
		_, err := UpdatePostStatusFieldFunc(repository).(func(context.Context, struct {
			Slug   Slug
			Status blog.Status
		}) (blog.Post, error))(context.Background(), struct {
			Slug   Slug
			Status blog.Status
		}{
			Slug:   Slug("test-" + id.Hex()),
			Status: blog.StatusPublished,
		})

		// Then
		assert.EqualError(t, err, "Forbidden")
	})
}

func TestUpdatePostContentFieldFunc(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		repository = mock_blog.NewMockPostRepository(ctrl)
	)

	t.Run("With successful updating post content", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()

		repository.EXPECT().FindByID(gomock.Any(), id).Return(blog.Post{Title: "Test", Slug: "test-" + id.Hex(), AuthorID: "authorizedID"}, nil)
		repository.EXPECT().Save(gomock.Any(), id, blog.NewPostQueryBuilder().WithMarkdown("Test").WithHTML("<p>Test</p>\n").Build()).Return(blog.Post{Title: "Test2", Slug: "test2-" + id.Hex(), Markdown: "Test", HTML: "<p>Test</p>\n", AuthorID: "authorizedID"}, nil)

		// When
		p, err := UpdatePostContentFieldFunc(repository).(func(context.Context, struct {
			Slug     Slug
			Markdown string
		}) (blog.Post, error))(context.WithValue(context.Background(), AuthorizedID, "authorizedID"), struct {
			Slug     Slug
			Markdown string
		}{
			Slug:     Slug("test-" + id.Hex()),
			Markdown: "Test",
		})

		// Then
		assert.Nil(t, err)
		assert.Equal(t, blog.Post{Title: "Test2", Slug: "test2-" + id.Hex(), Markdown: "Test", HTML: "<p>Test</p>\n", AuthorID: "authorizedID"}, p)
	})

	t.Run("When unable to find a post", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()

		repository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(blog.Post{}, errors.New("test unable to find a post"))

		// When
		_, err := UpdatePostContentFieldFunc(repository).(func(context.Context, struct {
			Slug     Slug
			Markdown string
		}) (blog.Post, error))(context.WithValue(context.Background(), AuthorizedID, "authorizedID"), struct {
			Slug     Slug
			Markdown string
		}{
			Slug:     Slug("test-" + id.Hex()),
			Markdown: "Test",
		})

		// Then
		assert.EqualError(t, err, "Not Found")
	})

	t.Run("When try to update other post content", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()

		repository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(blog.Post{Title: "Test", Slug: "test-" + id.Hex(), AuthorID: "authorizedID"}, nil)

		// When
		_, err := UpdatePostContentFieldFunc(repository).(func(context.Context, struct {
			Slug     Slug
			Markdown string
		}) (blog.Post, error))(context.Background(), struct {
			Slug     Slug
			Markdown string
		}{
			Slug:     Slug("test-" + id.Hex()),
			Markdown: "Test",
		})

		// Then
		assert.EqualError(t, err, "Forbidden")
	})
}

func TestUpdatePostCategoriesFieldFunc(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		repository = mock_blog.NewMockPostRepository(ctrl)
	)

	t.Run("With successful updating post categories", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()
		catID := primitive.NewObjectID()

		repository.EXPECT().FindByID(gomock.Any(), id).Return(blog.Post{Title: "Test", Slug: "test-" + id.Hex(), AuthorID: "authorizedID"}, nil)
		repository.EXPECT().Save(gomock.Any(), id, blog.NewPostQueryBuilder().WithCategories([]blog.Category{{ID: catID}}).Build()).Return(blog.Post{Title: "Test2", Slug: "test2-" + id.Hex(), AuthorID: "authorizedID", Categories: []mongo.DBRef{{ID: catID}}}, nil)

		// When
		p, err := UpdatePostCategoriesFieldFunc(repository).(func(context.Context, struct {
			Slug          Slug
			CategorySlugs []Slug
		}) (blog.Post, error))(context.WithValue(context.Background(), AuthorizedID, "authorizedID"), struct {
			Slug          Slug
			CategorySlugs []Slug
		}{
			Slug:          Slug("test-" + id.Hex()),
			CategorySlugs: []Slug{Slug("test-" + catID.Hex())},
		})

		// Then
		assert.Nil(t, err)
		assert.Equal(t, blog.Post{Title: "Test2", Slug: "test2-" + id.Hex(), AuthorID: "authorizedID", Categories: []mongo.DBRef{{ID: catID}}}, p)
	})

	t.Run("When unable to find a post", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()
		catID := primitive.NewObjectID()

		repository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(blog.Post{}, errors.New("test unable to find a post"))

		// When
		_, err := UpdatePostCategoriesFieldFunc(repository).(func(context.Context, struct {
			Slug          Slug
			CategorySlugs []Slug
		}) (blog.Post, error))(context.WithValue(context.Background(), AuthorizedID, "authorizedID"), struct {
			Slug          Slug
			CategorySlugs []Slug
		}{
			Slug:          Slug("test-" + id.Hex()),
			CategorySlugs: []Slug{Slug("test-" + catID.Hex())},
		})

		// Then
		assert.EqualError(t, err, "Not Found")
	})

	t.Run("When try to update other post categories", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()
		catID := primitive.NewObjectID()

		repository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(blog.Post{Title: "Test", Slug: "test-" + id.Hex(), AuthorID: "authorizedID"}, nil)

		// When
		_, err := UpdatePostCategoriesFieldFunc(repository).(func(context.Context, struct {
			Slug          Slug
			CategorySlugs []Slug
		}) (blog.Post, error))(context.Background(), struct {
			Slug          Slug
			CategorySlugs []Slug
		}{
			Slug:          Slug("test-" + id.Hex()),
			CategorySlugs: []Slug{Slug("test-" + catID.Hex())},
		})

		// Then
		assert.EqualError(t, err, "Forbidden")
	})
}

func TestUpdatePostTagsFieldFunc(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		repository = mock_blog.NewMockPostRepository(ctrl)
	)

	t.Run("With successful updating post tags", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()
		tagID := primitive.NewObjectID()

		repository.EXPECT().FindByID(gomock.Any(), id).Return(blog.Post{Title: "Test", Slug: "test-" + id.Hex(), AuthorID: "authorizedID"}, nil)
		repository.EXPECT().Save(gomock.Any(), id, blog.NewPostQueryBuilder().WithTags([]blog.Tag{{ID: tagID}}).Build()).Return(blog.Post{Title: "Test2", Slug: "test2-" + id.Hex(), AuthorID: "authorizedID", Tags: []mongo.DBRef{{ID: tagID}}}, nil)

		// When
		p, err := UpdatePostTagsFieldFunc(repository).(func(context.Context, struct {
			Slug     Slug
			TagSlugs []Slug
		}) (blog.Post, error))(context.WithValue(context.Background(), AuthorizedID, "authorizedID"), struct {
			Slug     Slug
			TagSlugs []Slug
		}{
			Slug:     Slug("test-" + id.Hex()),
			TagSlugs: []Slug{Slug("test-" + tagID.Hex())},
		})

		// Then
		assert.Nil(t, err)
		assert.Equal(t, blog.Post{Title: "Test2", Slug: "test2-" + id.Hex(), AuthorID: "authorizedID", Tags: []mongo.DBRef{{ID: tagID}}}, p)
	})

	t.Run("When unable to find a post", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()
		catID := primitive.NewObjectID()

		repository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(blog.Post{}, errors.New("test unable to find a post"))

		// When
		_, err := UpdatePostTagsFieldFunc(repository).(func(context.Context, struct {
			Slug     Slug
			TagSlugs []Slug
		}) (blog.Post, error))(context.WithValue(context.Background(), AuthorizedID, "authorizedID"), struct {
			Slug     Slug
			TagSlugs []Slug
		}{
			Slug:     Slug("test-" + id.Hex()),
			TagSlugs: []Slug{Slug("test-" + catID.Hex())},
		})

		// Then
		assert.EqualError(t, err, "Not Found")
	})

	t.Run("When try to update other post tags", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()
		catID := primitive.NewObjectID()

		repository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(blog.Post{Title: "Test", Slug: "test-" + id.Hex(), AuthorID: "authorizedID"}, nil)

		// When
		_, err := UpdatePostTagsFieldFunc(repository).(func(context.Context, struct {
			Slug     Slug
			TagSlugs []Slug
		}) (blog.Post, error))(context.Background(), struct {
			Slug     Slug
			TagSlugs []Slug
		}{
			Slug:     Slug("test-" + id.Hex()),
			TagSlugs: []Slug{Slug("test-" + catID.Hex())},
		})

		// Then
		assert.EqualError(t, err, "Forbidden")
	})
}

func TestUpdatePostFeaturedImageFieldFunc(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		repository = mock_blog.NewMockPostRepository(ctrl)
	)

	t.Run("With successful updating featured image", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()
		fileID := primitive.NewObjectID()

		repository.EXPECT().FindByID(gomock.Any(), id).Return(blog.Post{Title: "Test", Slug: "test-" + id.Hex(), AuthorID: "authorizedID"}, nil)
		repository.EXPECT().Save(gomock.Any(), id, blog.NewPostQueryBuilder().WithFeaturedImage(storage.File{ID: fileID}).Build()).Return(blog.Post{Title: "Test2", Slug: "test2-" + id.Hex(), AuthorID: "authorizedID", FeaturedImage: mongo.DBRef{ID: fileID}}, nil)

		// When
		p, err := UpdatePostFeaturedImageFieldFunc(repository).(func(context.Context, struct {
			Slug              Slug
			FeaturedImageSlug storage.Slug `graphql:",optional"`
		}) (blog.Post, error))(context.WithValue(context.Background(), AuthorizedID, "authorizedID"), struct {
			Slug              Slug
			FeaturedImageSlug storage.Slug `graphql:",optional"`
		}{
			Slug:              Slug("test-" + id.Hex()),
			FeaturedImageSlug: storage.Slug("test-" + fileID.Hex() + ".png"),
		})

		// Then
		assert.Nil(t, err)
		assert.Equal(t, blog.Post{Title: "Test2", Slug: "test2-" + id.Hex(), AuthorID: "authorizedID", FeaturedImage: mongo.DBRef{ID: fileID}}, p)
	})

	t.Run("When unable to find a post", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()
		fileID := primitive.NewObjectID()

		repository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(blog.Post{}, errors.New("test unable to  find a post"))

		// When
		_, err := UpdatePostFeaturedImageFieldFunc(repository).(func(context.Context, struct {
			Slug              Slug
			FeaturedImageSlug storage.Slug `graphql:",optional"`
		}) (blog.Post, error))(context.WithValue(context.Background(), AuthorizedID, "authorizedID"), struct {
			Slug              Slug
			FeaturedImageSlug storage.Slug `graphql:",optional"`
		}{
			Slug:              Slug("test-" + id.Hex()),
			FeaturedImageSlug: storage.Slug("test-" + fileID.Hex() + ".png"),
		})

		// Then
		assert.EqualError(t, err, "Not Found")
	})

	t.Run("With empty featured image slug", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()

		repository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(blog.Post{Title: "Test", Slug: "test-" + id.Hex(), AuthorID: "authorizedID"}, nil)
		repository.EXPECT().Save(gomock.Any(), gomock.Any(), blog.NewPostQueryBuilder().WithFeaturedImage(storage.File{}).Build()).Return(blog.Post{Title: "Test2", Slug: "test2-" + id.Hex(), AuthorID: "authorizedID", FeaturedImage: mongo.DBRef{}}, nil)

		// When
		p, err := UpdatePostFeaturedImageFieldFunc(repository).(func(context.Context, struct {
			Slug              Slug
			FeaturedImageSlug storage.Slug `graphql:",optional"`
		}) (blog.Post, error))(context.WithValue(context.Background(), AuthorizedID, "authorizedID"), struct {
			Slug              Slug
			FeaturedImageSlug storage.Slug `graphql:",optional"`
		}{
			Slug: Slug("test-" + id.Hex()),
		})

		// Then
		assert.Nil(t, err)
		assert.Equal(t, blog.Post{Title: "Test2", Slug: "test2-" + id.Hex(), AuthorID: "authorizedID", FeaturedImage: mongo.DBRef{}}, p)
	})

	t.Run("When try to update other post featured image", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()
		fileID := primitive.NewObjectID()

		repository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(blog.Post{Title: "Test", Slug: "test-" + id.Hex(), AuthorID: "authorizedID"}, nil)

		// When
		_, err := UpdatePostFeaturedImageFieldFunc(repository).(func(context.Context, struct {
			Slug              Slug
			FeaturedImageSlug storage.Slug `graphql:",optional"`
		}) (blog.Post, error))(context.Background(), struct {
			Slug              Slug
			FeaturedImageSlug storage.Slug `graphql:",optional"`
		}{
			Slug:              Slug("test-" + id.Hex()),
			FeaturedImageSlug: storage.Slug("test-" + fileID.Hex() + ".png"),
		})

		// Then
		assert.EqualError(t, err, "Forbidden")
	})
}

func TestUpdatePostAttachmentsFieldFunc(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		repository = mock_blog.NewMockPostRepository(ctrl)
	)

	t.Run("With successful updating post attachments", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()
		fileID := primitive.NewObjectID()

		repository.EXPECT().FindByID(gomock.Any(), id).Return(blog.Post{Title: "Test", Slug: "test-" + id.Hex(), AuthorID: "authorizedID"}, nil)
		repository.EXPECT().Save(gomock.Any(), id, blog.NewPostQueryBuilder().WithAttachments([]storage.File{{ID: fileID}}).Build()).Return(blog.Post{Title: "Test2", Slug: "test2-" + id.Hex(), AuthorID: "authorizedID", Attachments: []mongo.DBRef{{ID: fileID}}}, nil)

		// When
		p, err := UpdatePostAttachmentsFieldFunc(repository).(func(context.Context, struct {
			Slug            Slug
			AttachmentSlugs []storage.Slug
		}) (blog.Post, error))(context.WithValue(context.Background(), AuthorizedID, "authorizedID"), struct {
			Slug            Slug
			AttachmentSlugs []storage.Slug
		}{
			Slug:            Slug("test-" + id.Hex()),
			AttachmentSlugs: []storage.Slug{storage.Slug("test-" + fileID.Hex() + ".png")},
		})

		// Then
		assert.Nil(t, err)
		assert.Equal(t, blog.Post{Title: "Test2", Slug: "test2-" + id.Hex(), AuthorID: "authorizedID", Attachments: []mongo.DBRef{{ID: fileID}}}, p)
	})

	t.Run("When unable to find a post", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()
		fileID := primitive.NewObjectID()

		repository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(blog.Post{}, errors.New("test unable to  find a post"))

		// When
		_, err := UpdatePostAttachmentsFieldFunc(repository).(func(context.Context, struct {
			Slug            Slug
			AttachmentSlugs []storage.Slug
		}) (blog.Post, error))(context.WithValue(context.Background(), AuthorizedID, "authorizedID"), struct {
			Slug            Slug
			AttachmentSlugs []storage.Slug
		}{
			Slug:            Slug("test-" + id.Hex()),
			AttachmentSlugs: []storage.Slug{storage.Slug("test-" + fileID.Hex() + ".png")},
		})

		// Then
		assert.EqualError(t, err, "Not Found")
	})

	t.Run("When try to update other post attachments", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()
		fileID := primitive.NewObjectID()

		repository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(blog.Post{Title: "Test", Slug: "test-" + id.Hex(), AuthorID: "authorizedID"}, nil)

		// When
		_, err := UpdatePostAttachmentsFieldFunc(repository).(func(context.Context, struct {
			Slug            Slug
			AttachmentSlugs []storage.Slug
		}) (blog.Post, error))(context.Background(), struct {
			Slug            Slug
			AttachmentSlugs []storage.Slug
		}{
			Slug:            Slug("test-" + id.Hex()),
			AttachmentSlugs: []storage.Slug{storage.Slug("test-" + fileID.Hex() + ".png")},
		})

		// Then
		assert.EqualError(t, err, "Forbidden")
	})
}

func TestFindFeaturedImageBelongedToPostFieldFunc(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		repository = mock_storage.NewMockFileRepository(ctrl)
	)

	id := primitive.NewObjectID()

	repository.EXPECT().FindByID(gomock.Any(), id).Return(storage.File{Path: filepath.Join("authorizedID", "test.png"), FileName: "test.png", Slug: "test-" + id.Hex() + ".png"}, nil)

	// When
	f := FindFeaturedImageBelongedToPostFieldFunc(repository).(func(context.Context, blog.Post) storage.File)(context.Background(), blog.Post{FeaturedImage: mongo.DBRef{ID: id}})

	// Then
	assert.Equal(t, storage.File{Path: filepath.Join("authorizedID", "test.png"), FileName: "test.png", Slug: "test-" + id.Hex() + ".png"}, f)
}

func TestFindAllAttachmentsBelongedToPostFieldFunc(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		repository = mock_storage.NewMockFileRepository(ctrl)
	)

	id := primitive.NewObjectID()

	repository.EXPECT().FindAllByIDs(gomock.Any(), []primitive.ObjectID{id}).Return([]storage.File{{Path: filepath.Join("authorizedID", "test.png"), FileName: "test.png", Slug: "test-" + id.Hex() + ".png"}}, nil)

	// When
	files, err := FindAllAttachmentsBelongedToPostFieldFunc(repository).(func(context.Context, blog.Post) ([]storage.File, error))(context.Background(), blog.Post{Attachments: []mongo.DBRef{{ID: id}}})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, []storage.File{{Path: filepath.Join("authorizedID", "test.png"), FileName: "test.png", Slug: "test-" + id.Hex() + ".png"}}, files)
}

func TestGetURLNodeShareCountFieldFunc(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		transport = mock_http.NewMockRoundTripper(ctrl)
	)

	now := time.Date(2020, 4, 6, 9, 42, 0, 0, time.UTC)
	f := faketime.NewFaketimeWithTime(now)
	defer f.Undo()
	f.Do()

	c := facebook.NewClient("", transport)

	t.Run("With successful getting URLNode", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()

		transport.EXPECT().RoundTrip(gomock.Any()).DoAndReturn(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString(`{"engagement":{"share_count":1}}`)),
			}, nil
		})

		// When
		engagement := GetURLNodeShareCountFieldFunc("http://localhost", c).(func(context.Context, blog.Post) (engagement blog.Engagement))(context.Background(), blog.Post{Slug: "test-" + id.Hex(), PublishedAt: now})

		// Then
		assert.Equal(t, blog.Engagement{ShareCount: 1}, engagement)
	})

	t.Run("When unable to getting URLNode", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()

		transport.EXPECT().RoundTrip(gomock.Any()).Return(nil, errors.New("test unable to getting URLNode"))

		// When
		engagement := GetURLNodeShareCountFieldFunc("http://localhost", c).(func(context.Context, blog.Post) (engagement blog.Engagement))(context.Background(), blog.Post{Slug: "test-" + id.Hex(), PublishedAt: now})

		// Then
		assert.Equal(t, blog.Engagement{}, engagement)
	})
}
