package graphql

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/nomkhonwaan/myblog/pkg/blog"
	mock_blog "github.com/nomkhonwaan/myblog/pkg/blog/mock"
	"github.com/nomkhonwaan/myblog/pkg/mongo"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestBuildCategorySchema(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		repository = mock_blog.NewMockCategoryRepository(ctrl)
	)

	t.Run("FindCategoryBySlugFieldFunc", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()

		repository.EXPECT().FindByID(gomock.Any(), id).Return(blog.Category{Name: "Test", Slug: "test-" + id.Hex()}, nil)

		// When
		cat, err := FindCategoryBySlugFieldFunc(repository).(func(context.Context, struct{ Slug Slug }) (blog.Category, error))(context.Background(), struct{ Slug Slug }{Slug: Slug("test-" + id.Hex())})

		// Then
		assert.Nil(t, err)
		assert.Equal(t, blog.Category{Name: "Test", Slug: "test-" + id.Hex()}, cat)
	})

	t.Run("FindAllCategoriesFieldFunc", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()

		repository.EXPECT().FindAll(gomock.Any()).Return([]blog.Category{{Name: "Test", Slug: "test-" + id.Hex()}}, nil)

		// When
		cats, err := FindAllCategoriesFieldFunc(repository).(func(context.Context) ([]blog.Category, error))(context.Background())

		// Then
		assert.Nil(t, err)
		assert.Equal(t, []blog.Category{{Name: "Test", Slug: "test-" + id.Hex()}}, cats)
	})

	t.Run("FindAllCategoriesBelongedToPostFieldFunc", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()

		repository.EXPECT().FindAllByIDs(gomock.Any(), []primitive.ObjectID{id}).Return([]blog.Category{{Name: "Test", Slug: "test-" + id.Hex()}}, nil)

		// When
		cats, err := FindAllCategoriesBelongedToPostFieldFunc(repository).(func(context.Context, blog.Post) ([]blog.Category, error))(context.Background(), blog.Post{Categories: []mongo.DBRef{{ID: id}}})

		// Then
		assert.Nil(t, err)
		assert.Equal(t, []blog.Category{{Name: "Test", Slug: "test-" + id.Hex()}}, cats)
	})
}
