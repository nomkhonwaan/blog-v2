package graphql

import (
	"net/http"

	"github.com/samsarahq/thunder/graphql"
)

// Handler is a wrapper function to `graphql.HTTPHandler` for avoiding package name conflict
func Handler(schema *graphql.Schema) http.Handler {
	return graphql.Handler(schema)
}
