package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"
)

type mockTimer time.Time

func (timer mockTimer) Now() time.Time {
	return time.Time(timer)
}

type mockOutputer string

func (outputer *mockOutputer) Printf(format string, args ...interface{}) {
	*outputer = mockOutputer(fmt.Sprintf(format, args...))
}

func TestHandler(t *testing.T) {
	// Given
	timer := mockTimer(time.Now())
	outputer := new(mockOutputer)
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "https://api.nomkhonwaan.com/graphql", bytes.NewBufferString("{ categories { name slug } }"))
	r.RequestURI = "https://api.nomkhonwaan.com/graphql"
	r.RemoteAddr = "localhost"

	// When
	NewLoggingInterceptor(timer, outputer).
		Handler(
			http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					data, _ := json.Marshal([]map[string]interface{}{
						{
							"name": "Web Development",
							"slug": "web-development-1",
						},
					})

					_, _ = w.Write(data)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
				},
			),
		).
		ServeHTTP(w, r)

	// Then
	assert.Equal(t, `[{"name":"Web Development","slug":"web-development-1"}]`, w.Body.String())
	assert.Regexp(t, regexp.MustCompile(fmt.Sprintf(`"%s %s" %d %q %q "(\d+|\.)(.*)"`, r.Method, r.RequestURI, http.StatusOK, r.UserAgent(), r.RemoteAddr)), string(*outputer))
}

func TestDefault(t *testing.T) {
	// Given

	// When
	l := Default()

	// Then
	assert.IsType(t, &DefaultTimer{}, l.Timer)
	assert.IsType(t, DefaultOutputer{}, l.Outputer)
}
