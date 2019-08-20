package graphql

import "net/http"

// Graphiql provides GraphQL IDE handler
func Graphiql(gzipAsset []byte) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Type", "text/html")
		w.Write(gzipAsset)
	}
}
