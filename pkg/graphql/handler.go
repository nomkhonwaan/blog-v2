package graphql

import (
	"net/http"

	"github.com/nomkhonwaan/myblog/pkg/data"
	"github.com/samsarahq/thunder/graphql"
)

// Playground provides GraphQL IDE from static asset file /data/graphql-playground.html
func Playground(w http.ResponseWriter, _ *http.Request) {
	d, _ := data.GzipAsset("data/graphql-playground.html")

	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Type", "text/html")
	w.Write(d)
}

// HTTPHandler is a wrapper function to `graphql.HTTPHandler` for avoiding package name conflict
func HTTPHandler(schema *graphql.Schema, middlewares ...graphql.MiddlewareFunc) http.Handler {
	return graphql.HTTPHandler(schema, middlewares...)
}
