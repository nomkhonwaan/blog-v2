package blog

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewService(t *testing.T) {
	// Given
	categoryRepo := &MockCategoryRepository{}
	postRepo := &MockPostRepository{}
	tagRepo := &MockTagRepository{}

	// When
	service := NewService(categoryRepo, postRepo, tagRepo)

	// Then
	assert.Equal(t, categoryRepo, service.Category())
	assert.Equal(t, postRepo, service.Post())
}
