package log

import (
	"net/http"
	"time"
)

// LoggingInterceptor uses for logging an HTTP request with response status code and processing time
type LoggingInterceptor struct {
	now    func() time.Time
	printf func(format string, args ...interface{})
}

func (interceptor LoggingInterceptor) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := interceptor.now()
		cw := &CustomResponseWriter{ResponseWriter: w}

		h.ServeHTTP(cw, r)

		interceptor.printf("\"%s %s\" %d %q %q \"%v\"", r.Method, r.RequestURI, cw.statusCode, r.UserAgent(), r.RemoteAddr, time.Since(start))
	})
}

// CustomResponseWriter embeds http.ResponseWriter interface for overriding method WriteHeader
// which will save the response status code for logging purpose
type CustomResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *CustomResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// NewLoggingInterceptor returns the LoggingInterceptor object
func NewLoggingInterceptor(now func() time.Time, printf func(format string, args ...interface{})) LoggingInterceptor {
	return LoggingInterceptor{now, printf}
}
