package graphiql

import (
	"github.com/samsarahq/thunder/graphql/graphiql"
	"net/http"
)

// Handler is a wrapped function to the `graphiql.Handler` for avoiding package name conflict
func Handler(prefix string) http.Handler {
	return http.StripPrefix(prefix, graphiql.Handler())
}
