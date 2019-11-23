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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestSlug_GetID(t *testing.T) {
	// Given
	id := primitive.NewObjectID()
	slug := Slug("test-featured-image-" + id.Hex() + ".jpg")

	// When
	result, err := slug.GetID()

	// Then
	assert.Nil(t, err)
	assert.Equal(t, id, result)
}

func TestSlug_MustGetID(t *testing.T) {
	t.Run("With valid ObjectID string", func(t *testing.T) {
		// Given
		id := primitive.NewObjectID()
		slug := Slug("test-featured-image-" + id.Hex() + ".jpg")

		// When
		result := slug.MustGetID()

		// Then
		assert.Equal(t, id, result)
	})

	t.Run("With invalid ObjectID string", func(t *testing.T) {
		// Given
		slug := Slug("test-featured-image-invalid-object-id.jpg")

		// When
		result := slug.MustGetID()

		// Then
		assert.IsType(t, primitive.ObjectID{}, result)
		assert.NotEmpty(t, result.(primitive.ObjectID).Hex())
	})
}

func TestHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		downloader = mock_storage.NewMockDownloader(ctrl)
		uploader   = mock_storage.NewMockUploader(ctrl)
		cache      = mock_storage.NewMockCache(ctrl)
		fileRepo   = mock_storage.NewMockFileRepository(ctrl)

		router = mux.NewRouter()
	)

	NewHandler(cache, fileRepo, downloader, uploader).Register(router.PathPrefix("/v1/storage").Subrouter())

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

		uploader.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, path string, _ io.Reader) error {
			id, err := Slug(path).GetID()

			assert.Nil(t, err)
			assert.False(t, id.(primitive.ObjectID).IsZero())

			fileRepo.EXPECT().Create(gomock.Any(), File{
				ID:             id.(primitive.ObjectID),
				Path:           path,
				FileName:       "test.txt",
				Slug:           filepath.Base(path),
				OptionalField1: "CustomizedAmazonS3Client",
			}).Return(File{}, nil)

			return nil
		})

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

		uploader.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test upload file error"))

		// When
		router.ServeHTTP(w, r)

		// Then
		assert.Equal(t, "500 Internal Server Error", w.Result().Status)
	})
}
