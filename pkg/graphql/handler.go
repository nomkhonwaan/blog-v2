package graphql

import (
	"net/http"

	"github.com/samsarahq/thunder/graphql"
)

// Handler is a wrapped function to the original graphql.HTTPHandler for avoiding package name conflict
func Handler(schema *graphql.Schema, middlewares ...graphql.MiddlewareFunc) http.Handler {
	return graphql.HTTPHandler(schema, middlewares...)
}

// ServeGraphiqlHandlerFunc provides a GraphQL Playground page
func ServeGraphiqlHandlerFunc(tmpl []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write(tmpl)
	}
}
