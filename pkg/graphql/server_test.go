package graphql

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/nomkhonwaan/myblog/pkg/blog"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestMakeFieldFuncCategories(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	id := primitive.NewObjectID()
	ctx := context.Background()
	expected := []blog.Category{
		{
			ID:   id,
			Name: "Web Development",
			Slug: "web-development-" + id.Hex(),
		},
	}
	categoryRepository := blog.NewMockCategoryRepository(ctrl)
	service := blog.NewMockService(ctrl)

	service.EXPECT().Category().Return(categoryRepository)
	categoryRepository.EXPECT().FindAll(ctx).Return(expected, nil)

	server := NewServer(service)

	// When
	categories, err := server.makeFieldFuncCategories(ctx)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, expected, categories)
}
