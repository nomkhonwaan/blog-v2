package slug_test

import (
	. "github.com/nomkhonwaan/myblog/pkg/slug"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMake(t *testing.T) {
	// Given
	tests := map[string]struct {
		s        string
		expected string
	}{
		"With alphanumeric string (alphabet and digit)": {
			s:        "Test post title with number 1 or number 2",
			expected: "test-post-title-with-number-1-or-number-2",
		},
		"With non-unicode character in the string": {
			s:        "√ We're 100% genuine software",
			expected: "we-re-100-genuine-software",
		},
		"With Thai character in the string": {
			s:        "ทางที่ดี คือ ทางลาดยาง",
			expected: "ทางที่ดี-คือ-ทางลาดยาง",
		},
	}

	// When
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, Make(test.s))
		})
	}

	// Then
}
