package blog

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStatuses(t *testing.T) {
	// Given
	tests := map[string]struct {
		status   Status
		expected string
	}{
		"Published is PUBLISHED":          {status: Published, expected: "PUBLISHED"},
		"Draft is DRAFT":                  {status: Draft, expected: "DRAFT"},
		"PendingReview is PENDING_REVIEW": {status: PendingReview, expected: "PENDING_REVIEW"},
	}

	// When
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, string(test.status), test.expected)
		})
	}

	// Then
}
