package playground

import "net/http"

// HandlerFunc provides GraphQL Playground handler
func HandlerFunc(gzipAsset []byte) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write(gzipAsset)
	}
}
