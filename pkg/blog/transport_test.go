package blog

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestEncodeResponse(t *testing.T) {
	// Given
	recorder := httptest.NewRecorder()

	// When
	err := encodeResponse(
		context.Background(),
		recorder,
		map[string]interface{}{
			"hello": "world!",
		},
	)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, recorder.Header().Get("Content-Type"), "application/json; charset=utf-8")
	assert.Equal(t, recorder.Body.String(), "{\"hello\":\"world!\"}\n")
}
