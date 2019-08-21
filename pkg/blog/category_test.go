package blog

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http/httptest"
	"testing"
)

func TestJSONMarshalingCategoryEntity(t *testing.T) {
	// Given
	id := primitive.NewObjectID()
	category := Category{
		ID:   id,
		Name: "Web Development",
		Slug: "web-development-" + id.Hex(),
	}
	recorder := httptest.NewRecorder()

	// When
	err := encodeResponse(
		context.Background(),
		recorder,
		category,
	)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, recorder.Body.String(), "{\"id\":\""+id.Hex()+"\",\"name\":\"Web Development\",\"slug\":\"web-development-"+id.Hex()+"\"}\n")
}
