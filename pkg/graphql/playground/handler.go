package playground

import "net/http"

// Handler provides Graphiql for playing on GraphQL server
func Handler(data []byte) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write(data)
	})
}
