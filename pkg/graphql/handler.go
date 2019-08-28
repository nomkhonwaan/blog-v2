package graphql

import (
	"net/http"

	"github.com/samsarahq/thunder/graphql"
)

// Handler is a wrapper function to `graphql.HTTPHandler` for avoiding package name conflict
func Handler(schema *graphql.Schema, middlewares ...graphql.MiddlewareFunc) http.Handler {
	return graphql.HTTPHandler(schema, middlewares...)
}
