package blog_test

import (
	. "github.com/nomkhonwaan/myblog/pkg/blog"
	mock_blog "github.com/nomkhonwaan/myblog/pkg/blog/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_Category(t *testing.T) {
	// Given
	categoryRepo := &mock_blog.MockCategoryRepository{}

	// When
	service := NewService(categoryRepo, nil, nil)

	// Then
	assert.Equal(t, categoryRepo, service.Category())
}

func TestService_Post(t *testing.T) {
	// Given
	postRepo := &mock_blog.MockPostRepository{}

	// When
	service := NewService(nil, postRepo, nil)

	// Then
	assert.Equal(t, postRepo, service.Post())
}

func TestService_Tag(t *testing.T) {
	// Given
	tagRepo := &mock_blog.MockTagRepository{}

	// When
	service := NewService(nil, nil, tagRepo)

	// Then
	assert.Equal(t, tagRepo, service.Tag())
}
