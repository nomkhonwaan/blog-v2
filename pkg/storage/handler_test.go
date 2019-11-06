package storage_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/nomkhonwaan/myblog/pkg/auth"
	. "github.com/nomkhonwaan/myblog/pkg/storage"
	mock_storage "github.com/nomkhonwaan/myblog/pkg/storage/mock"
	"github.com/stretchr/testify/assert"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	downloader := mock_storage.NewMockDownloader(ctrl)
	uploader := mock_storage.NewMockUploader(ctrl)
	router := mux.NewRouter()

	Register(router.PathPrefix("/v1/storage").Subrouter(), downloader, uploader)

	newFileUploadRequest := func(fileName string, body io.Reader) *http.Request {
		buf := &bytes.Buffer{}
		wtr := multipart.NewWriter(buf)
		defer wtr.Close()

		w, _ := wtr.CreateFormFile("file", fileName)
		_, _ = io.Copy(w, body)

		r := httptest.NewRequest(http.MethodPost, "/v1/storage/upload", buf)
		r.Header.Set("Content-Type", wtr.FormDataContentType())

		return r
	}

	withAuthorizedID := func(r *http.Request) *http.Request {
		return r.WithContext(context.WithValue(context.Background(), auth.UserProperty, &jwt.Token{
			Claims: jwt.MapClaims{
				"sub": "authorizedID",
			},
		}))
	}

	t.Run("With successful uploading file", func(t *testing.T) {
		// Given
		body := bytes.NewBufferString("test")
		fileName := "test.txt"

		w := httptest.NewRecorder()
		r := withAuthorizedID(newFileUploadRequest(fileName, body))

		uploader.EXPECT().Upload(gomock.Any(), "authorizedID/"+fileName, gomock.Any()).Return(File{}, nil)

		// When
		router.ServeHTTP(w, r)

		// Then
		assert.Equal(t, "200 OK", w.Result().Status)
	})

	t.Run("When unable to retrieve authorized ID from", func(t *testing.T) {
		// Given
		body := bytes.NewBufferString("test")
		fileName := "test.txt"

		w := httptest.NewRecorder()
		r := newFileUploadRequest(fileName, body)

		// When
		router.ServeHTTP(w, r)

		// Then
		var result map[string]interface{}
		_ = json.NewDecoder(w.Body).Decode(&result)

		assert.Equal(t, "401 Unauthorized", w.Result().Status)
		assert.EqualValues(t, 401, result["error"].(map[string]interface{})["code"])
		assert.Equal(t, "Unauthorized", result["error"].(map[string]interface{})["message"])
	})

	t.Run("When unable to read form file", func(t *testing.T) {
		// Given
		body := bytes.NewBufferString("")

		w := httptest.NewRecorder()
		r := withAuthorizedID(newFileUploadRequest("", body))

		// When
		router.ServeHTTP(w, r)

		// Then
		assert.Equal(t, "500 Internal Server Error", w.Result().Status)
	})

	t.Run("When unable to upload file to the storage server", func(t *testing.T) {
		// Given
		body := bytes.NewBufferString("test")
		fileName := "test.txt"

		w := httptest.NewRecorder()
		r := withAuthorizedID(newFileUploadRequest(fileName, body))

		uploader.EXPECT().Upload(gomock.Any(), "authorizedID/"+fileName, gomock.Any()).Return(File{}, errors.New("test upload file error"))

		// When
		router.ServeHTTP(w, r)

		// Then
		assert.Equal(t, "500 Internal Server Error", w.Result().Status)
	})
}
