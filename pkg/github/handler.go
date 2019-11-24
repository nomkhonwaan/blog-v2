package github

import (
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strings"
)

type Handler struct {
	transport http.RoundTripper
}

func NewHandler(transport http.RoundTripper) Handler {
	return Handler{
		transport: transport,
	}
}

func (h Handler) Register(r *mux.Router) {
	r.Path("/gist").HandlerFunc(h.proxyGist).Methods(http.MethodGet)
}

func (h Handler) proxyGist(w http.ResponseWriter, r *http.Request) {
	src := r.URL.Query().Get("src")
	if src == "" {
		http.Error(w, "`src` is required for downloading Gist content", http.StatusUnprocessableEntity)
		return
	}

	req, _ := http.NewRequest(http.MethodGet, "https://gist.github.com/"+strings.Replace(src[24:], ".js", ".json", 1), nil)
	res, err := h.transport.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	length, _ := io.Copy(w, res.Body)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", length))
}
