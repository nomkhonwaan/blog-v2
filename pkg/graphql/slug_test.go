package graphql

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestSlug_GetID(t *testing.T) {
	// Given
	id := primitive.NewObjectID()
	slug := Slug(id.Hex())

	// When
	result, err := slug.GetID()

	// Then
	assert.Nil(t, err)
	assert.Equal(t, id, result)
}

func TestSlug_MustGetID(t *testing.T) {
	t.Run("With valid ObjectID string", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()
		slug := Slug(id.Hex())

		// When
		result := slug.MustGetID()

		// Then
		assert.Equal(t, id, result)
	})

	t.Run("With invalid ObjectID string", func(t *testing.T) {
		// Given
		slug := Slug("invalid-object-id")

		// When
		result := slug.MustGetID()

		// Then
		assert.IsType(t, primitive.ObjectID{}, result)
		assert.NotEmpty(t, result.(primitive.ObjectID).Hex())
	})
}
