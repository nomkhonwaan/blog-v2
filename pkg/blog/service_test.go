package blog_test

import (
	. "github.com/nomkhonwaan/myblog/pkg/blog"
	mock_blog "github.com/nomkhonwaan/myblog/pkg/blog/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewService(t *testing.T) {
	// Given
	categoryRepo := &mock_blog.MockCategoryRepository{}
	postRepo := &mock_blog.MockPostRepository{}
	tagRepo := &mock_blog.MockTagRepository{}

	// When
	service := NewService(categoryRepo, postRepo, tagRepo)

	// Then
	assert.Equal(t, categoryRepo, service.Category())
	assert.Equal(t, postRepo, service.Post())
}
