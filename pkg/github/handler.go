package github

import (
	"bytes"
	"errors"
	"github.com/nomkhonwaan/myblog/pkg/storage"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// GetGistHandlerFunc handles GitHub Gist downloading request
func GetGistHandlerFunc(cache storage.Cache, transport http.RoundTripper) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		src := r.URL.Query().Get("src")
		if src == "" {
			http.Error(w, errors.New("src value is empty").Error(), http.StatusBadRequest)
			return
		}

		cacheFileName := url.QueryEscape(src) + ".json"

		if cache.Exists(cacheFileName) {
			body, err := cache.Retrieve(cacheFileName)
			if err == nil {
				defer body.Close()
				length, _ := io.Copy(w, body)
				w.Header().Set("Content-Length", strconv.Itoa(int(length)))
				return
			}
			logrus.Errorf("unable to retrieve Gist file %q from cache", cacheFileName)
		}

		var (
			c    = &http.Client{Transport: transport}
			u, _ = url.Parse(src)
		)

		u.Host = "gist.github.com"
		u.Path = strings.Replace(u.Path, ".js", ".json", 1)

		req, _ := http.NewRequest(http.MethodGet, u.String(), nil)
		res, err := c.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()

		data, _ := ioutil.ReadAll(res.Body)
		err = cache.Store(bytes.NewReader(data), cacheFileName)
		if err != nil {
			logrus.Errorf("unable to store Gist file %q to cache", cacheFileName)
		}

		w.Header().Set("Content-Length", strconv.Itoa(len(data)))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(data)
	}
}
