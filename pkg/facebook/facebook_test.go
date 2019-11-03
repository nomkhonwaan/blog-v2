package facebook_test

import (
	. "github.com/nomkhonwaan/myblog/pkg/facebook"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsFacebookCrawlerRequest(t *testing.T) {
	// Given
	tests := map[string]struct {
		data     string
		expected bool
	}{
		"With the first Facebook's user-agent string": {
			data:     "facebookexternalhit/1.1 (+http://www.facebook.com/externalhit_uatext.php)",
			expected: true,
		},
		"With the second Facebook's user-agent string": {
			data:     "facebookexternalhit/1.1",
			expected: true,
		},
		"With non-Facebook's user-agent string": {
			data:     "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.87 Safari/537.36",
			expected: false,
		},
	}

	// When
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, IsFacebookCrawlerRequest(test.data))
		})
	}

	// Then
}

func TestCrawlerMiddleware_Handler(t *testing.T) {
	//ctrl := gomock.NewController(t)
	//defer ctrl.Finish()
	//
	//var (
	//	postRepo = mock_blog.NewMockPostRepository(ctrl)
	//
	//	mw          = NewCrawlerMiddleware(postRepo)
	//	nextHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	})
	//	newSinglePageRequest = func(id interface{}) *http.Request {
	//		return httptest.NewRequest(http.MethodGet, "/2019/11/03/test-post-"+id.(string), nil)
	//	}
	//)
	//
	//t.Run("With Facebook's user-agent string on a single page", func(t *testing.T) {
	//	// Given
	//	id := primitive.NewObjectID()
	//	w := httptest.NewRecorder()
	//
	//	// When
	//	mw.Handler(nextHandler).ServeHTTP(w, newSinglePageRequest(id.Hex()))
	//
	//	// Then
	//})
}
