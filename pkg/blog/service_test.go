package blog

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewService(t *testing.T) {
	// Given
	categoryRepository := &MockCategoryRepository{}
	postRepository := &MockPostRepository{}

	// When
	service := NewService(categoryRepository, postRepository)

	// Then
	assert.Equal(t, categoryRepository, service.Category())
	assert.Equal(t, postRepository, service.Post())
}
