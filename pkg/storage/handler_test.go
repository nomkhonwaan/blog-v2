package storage_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/nomkhonwaan/myblog/pkg/auth"
	mock_image "github.com/nomkhonwaan/myblog/pkg/image/mock"
	. "github.com/nomkhonwaan/myblog/pkg/storage"
	mock_storage "github.com/nomkhonwaan/myblog/pkg/storage/mock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestDeleteHandlerFunc(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		bucket     = mock_storage.NewMockStorage(ctrl)
		repository = mock_storage.NewMockFileRepository(ctrl)
	)

	newDeleteRequest := func(slug string) *http.Request {
		req := httptest.NewRequest(http.MethodDelete, "/api/v2.1/storage/"+slug+"/delete", nil)
		return req.WithContext(
			context.WithValue(req.Context(), chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{Keys: []string{"slug"}, Values: []string{slug}},
				}))
	}

	t.Run("With successful deleting file", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		id := primitive.NewObjectID()
		slug := "test-" + id.Hex() + ".txt"

		repository.EXPECT().FindByID(gomock.Any(), id).Return(File{Path: filepath.Join("authorizedID", slug)}, nil)
		bucket.EXPECT().Delete(gomock.Any(), filepath.Join("authorizedID", slug)).Return(nil)
		repository.EXPECT().Delete(gomock.Any(), id).Return(nil)

		// When
		DeleteHandlerFunc(bucket, repository).ServeHTTP(w, withAuthorizedID(newDeleteRequest(slug)))

		// Then
		assert.Equal(t, "200 OK", w.Result().Status)
	})

	t.Run("When unable to retrieve authorizedID", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		id := primitive.NewObjectID()
		slug := "test-" + id.Hex() + ".txt"

		// When
		DeleteHandlerFunc(bucket, repository).ServeHTTP(w, newDeleteRequest(slug))

		// Then
		assert.Equal(t, `{"error":{"code":401,"message":"Unauthorized"}}`, w.Body.String())
	})

	t.Run("When unable to find a file", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		id := primitive.NewObjectID()
		slug := "test-" + id.Hex() + ".txt"

		repository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(File{}, errors.New("test unable to find a file"))

		// When
		DeleteHandlerFunc(bucket, repository).ServeHTTP(w, withAuthorizedID(newDeleteRequest(slug)))

		// Then
		assert.Equal(t, `{"error":{"code":404,"message":"test unable to find a file"}}`, w.Body.String())
	})

	t.Run("When unable to delete file from storage", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		id := primitive.NewObjectID()
		slug := "test-" + id.Hex() + ".txt"

		repository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(File{Path: filepath.Join("authorizedID", slug)}, nil)
		bucket.EXPECT().Delete(gomock.Any(), filepath.Join("authorizedID", slug)).Return(errors.New("test unable to delete file from storage"))

		// When
		DeleteHandlerFunc(bucket, repository).ServeHTTP(w, withAuthorizedID(newDeleteRequest(slug)))

		// Then
		assert.Equal(t, `{"error":{"code":500,"message":"test unable to delete file from storage"}}`, w.Body.String())
	})

	t.Run("When unable to delete file from database", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		id := primitive.NewObjectID()
		slug := "test-" + id.Hex() + ".txt"

		repository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(File{Path: filepath.Join("authorizedID", slug)}, nil)
		bucket.EXPECT().Delete(gomock.Any(), filepath.Join("authorizedID", slug)).Return(nil)
		repository.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errors.New("test unable to delete file from database"))

		// When
		DeleteHandlerFunc(bucket, repository).ServeHTTP(w, withAuthorizedID(newDeleteRequest(slug)))

		// Then
		assert.Equal(t, `{"error":{"code":500,"message":"test unable to delete file from database"}}`, w.Body.String())
	})
}

func TestDownloadHandlerFunc(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		cache      = mock_storage.NewMockCache(ctrl)
		bucket     = mock_storage.NewMockStorage(ctrl)
		resizer    = mock_image.NewMockResizer(ctrl)
		repository = mock_storage.NewMockFileRepository(ctrl)
	)

	newDownloadRequest := func(slug string) *http.Request {
		req := httptest.NewRequest(http.MethodGet, "/api/v2.1/storage/"+slug, nil)
		return req.WithContext(
			context.WithValue(req.Context(), chi.RouteCtxKey,
				&chi.Context{
					URLParams: chi.RouteParams{Keys: []string{"slug"}, Values: []string{slug}},
				}))
	}

	t.Run("With successful downloading file", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		id := primitive.NewObjectID()
		slug := "test-" + id.Hex() + ".txt"

		repository.EXPECT().FindByID(gomock.Any(), id).Return(File{Path: filepath.Join("authorizedID", slug)}, nil)
		cache.EXPECT().Exists(filepath.Join("authorizedID", slug)).Return(false)
		cache.EXPECT().Store(gomock.Any(), filepath.Join("authorizedID", slug)).DoAndReturn(func(body io.Reader, path string) error {
			_, _ = ioutil.ReadAll(body)
			return nil
		})
		bucket.EXPECT().Download(gomock.Any(), filepath.Join("authorizedID", slug)).Return(ioutil.NopCloser(bytes.NewReader([]byte("test"))), nil)

		// When
		DownloadHandlerFunc(bucket, cache, resizer, repository).ServeHTTP(w, newDownloadRequest(slug))

		// Then
		assert.Equal(t, "test", w.Body.String())
	})

	t.Run("With successful downloading file from cache storage", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		id := primitive.NewObjectID()
		slug := "test-" + id.Hex() + ".txt"

		repository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(File{Path: filepath.Join("authorizedID", slug)}, nil)
		cache.EXPECT().Exists(gomock.Any()).Return(true)
		cache.EXPECT().Retrieve(gomock.Any()).Return(ioutil.NopCloser(bytes.NewBufferString("test")), nil)

		// When
		DownloadHandlerFunc(bucket, cache, resizer, repository).ServeHTTP(w, newDownloadRequest(slug))

		// Then
		assert.Equal(t, "test", w.Body.String())
	})

	t.Run("With successful retrieving and resizing image file", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		id := primitive.NewObjectID()
		slug := "test-" + id.Hex() + ".png"
		var resizedBody bytes.Buffer
		_ = png.Encode(&resizedBody, image.NewRGBA(image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 100, Y: 100}}))

		repository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(File{Path: filepath.Join("authorizedID", slug)}, nil)
		cache.EXPECT().Exists(filepath.Join("authorizedID", slug)).Return(true)
		cache.EXPECT().Retrieve(gomock.Any()).Return(ioutil.NopCloser(bytes.NewBufferString("test")), nil)
		cache.EXPECT().Exists(gomock.Any()).Return(false)
		resizer.EXPECT().Resize(gomock.Any(), 50, 0).Return(&resizedBody, nil)
		cache.EXPECT().Store(gomock.Any(), gomock.Any()).DoAndReturn(func(body io.Reader, path string) error {
			_, _ = ioutil.ReadAll(body)
			return nil
		})

		// When
		DownloadHandlerFunc(bucket, cache, resizer, repository).ServeHTTP(w, newDownloadRequest(slug+"?width=50"))

		// Then
		assert.Equal(t, "200 OK", w.Result().Status)
	})

	t.Run("With successful retrieving resized image file from cache", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		id := primitive.NewObjectID()
		slug := "test-" + id.Hex() + ".png"
		var resizedBody bytes.Buffer
		_ = png.Encode(&resizedBody, image.NewRGBA(image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 50, Y: 50}}))

		repository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(File{Path: filepath.Join("authorizedID", slug)}, nil)
		cache.EXPECT().Exists(gomock.Any()).Return(true)
		cache.EXPECT().Retrieve(gomock.Any()).Return(ioutil.NopCloser(&resizedBody), nil)

		// When
		DownloadHandlerFunc(bucket, cache, resizer, repository).ServeHTTP(w, newDownloadRequest(slug+"?width=50"))

		// Then
		assert.Equal(t, "200 OK", w.Result().Status)
	})

	t.Run("When unable to find a file", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		id := primitive.NewObjectID()
		slug := "test-" + id.Hex() + ".txt"

		repository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(File{}, errors.New("test unable to find a file"))

		// When
		DownloadHandlerFunc(bucket, cache, resizer, repository).ServeHTTP(w, newDownloadRequest(slug))

		// Then
		assert.Equal(t, `{"error":{"code":404,"message":"test unable to find a file"}}`, w.Body.String())
	})

	t.Run("When unable to retrieve file from cache", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		id := primitive.NewObjectID()
		slug := "test-" + id.Hex() + ".txt"

		repository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(File{Path: filepath.Join("authorizedID", slug)}, nil)
		cache.EXPECT().Exists(gomock.Any()).Return(true)
		cache.EXPECT().Retrieve(gomock.Any()).Return(nil, errors.New("test unable to retrieve file from cache"))
		bucket.EXPECT().Download(gomock.Any(), gomock.Any()).Return(ioutil.NopCloser(bytes.NewBufferString("test")), nil)
		cache.EXPECT().Store(gomock.Any(), gomock.Any()).Return(nil)

		// When
		DownloadHandlerFunc(bucket, cache, resizer, repository).ServeHTTP(w, newDownloadRequest(slug))

		// Then
		assert.Equal(t, "200 OK", w.Result().Status)
	})

	t.Run("When unable to retrieve resized image file from cache", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		id := primitive.NewObjectID()
		slug := "test-" + id.Hex() + ".png"
		var resizedBody bytes.Buffer
		_ = png.Encode(&resizedBody, image.NewRGBA(image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 100, Y: 100}}))

		repository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(File{Path: filepath.Join("authorizedID", slug)}, nil)
		cache.EXPECT().Exists(gomock.Any()).Return(true)
		cache.EXPECT().Retrieve(gomock.Any()).Return(nil, errors.New("test unable to retrieve file from cache"))
		cache.EXPECT().Exists(gomock.Any()).Return(true)
		cache.EXPECT().Retrieve(gomock.Any()).Return(ioutil.NopCloser(&resizedBody), nil)
		resizer.EXPECT().Resize(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
		cache.EXPECT().Store(gomock.Any(), gomock.Any()).Return(nil)

		// When
		DownloadHandlerFunc(bucket, cache, resizer, repository).ServeHTTP(w, newDownloadRequest(slug+"?width=50"))

		// Then
		assert.Equal(t, "200 OK", w.Result().Status)
	})

	t.Run("When unable to resize image", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		id := primitive.NewObjectID()
		slug := "test-" + id.Hex() + ".png"
		var resizedBody bytes.Buffer
		_ = png.Encode(&resizedBody, image.NewRGBA(image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 100, Y: 100}}))

		repository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(File{Path: filepath.Join("authorizedID", slug)}, nil)
		cache.EXPECT().Exists(filepath.Join("authorizedID", slug)).Return(true)
		cache.EXPECT().Retrieve(gomock.Any()).Return(ioutil.NopCloser(&resizedBody), nil)
		cache.EXPECT().Exists(gomock.Any()).Return(false)
		resizer.EXPECT().Resize(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("test unable to resize image file"))

		// When
		DownloadHandlerFunc(bucket, cache, resizer, repository).ServeHTTP(w, newDownloadRequest(slug+"?width=50"))

		// Then
		assert.Equal(t, "test unable to resize image file\n", w.Body.String())
	})

	t.Run("When unable to download file", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		id := primitive.NewObjectID()
		slug := "test-" + id.Hex() + ".txt"

		repository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(File{Path: filepath.Join("authorizedID", slug)}, nil)
		cache.EXPECT().Exists(gomock.Any()).Return(false)
		bucket.EXPECT().Download(gomock.Any(), gomock.Any()).Return(nil, errors.New("test unable to download file"))

		// When
		DownloadHandlerFunc(bucket, cache, resizer, repository).ServeHTTP(w, newDownloadRequest(slug))

		// Then
		assert.Equal(t, "test unable to download file\n", w.Body.String())
	})

	t.Run("When unable to save file to cache", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		id := primitive.NewObjectID()
		slug := "test-" + id.Hex() + ".txt"

		repository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(File{Path: filepath.Join("authorizedID", slug)}, nil)
		cache.EXPECT().Exists(gomock.Any()).Return(false)
		cache.EXPECT().Store(gomock.Any(), gomock.Any()).Return(errors.New("test unable to save file to cache"))
		bucket.EXPECT().Download(gomock.Any(), gomock.Any()).Return(ioutil.NopCloser(bytes.NewBufferString("test")), nil)

		// When
		DownloadHandlerFunc(bucket, cache, resizer, repository).ServeHTTP(w, newDownloadRequest(slug))

		// Then
		assert.Equal(t, "200 OK", w.Result().Status)
	})

	t.Run("When unable to save resized image file to cache", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		id := primitive.NewObjectID()
		slug := "test-" + id.Hex() + ".png"
		var resizedBody bytes.Buffer
		_ = png.Encode(&resizedBody, image.NewRGBA(image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 100, Y: 100}}))

		repository.EXPECT().FindByID(gomock.Any(), id).Return(File{Path: filepath.Join("authorizedID", slug)}, nil)
		cache.EXPECT().Exists(filepath.Join("authorizedID", slug)).Return(true)
		cache.EXPECT().Retrieve(gomock.Any()).Return(ioutil.NopCloser(bytes.NewBufferString("test")), nil)
		cache.EXPECT().Exists(gomock.Any()).Return(false)
		resizer.EXPECT().Resize(gomock.Any(), gomock.Any(), gomock.Any()).Return(ioutil.NopCloser(&resizedBody), nil)
		cache.EXPECT().Store(gomock.Any(), gomock.Any()).DoAndReturn(func(body io.Reader, path string) error {
			_, _ = ioutil.ReadAll(body)
			return errors.New("test unable to save resized image file to cache")
		})

		// When
		DownloadHandlerFunc(bucket, cache, resizer, repository).ServeHTTP(w, newDownloadRequest(slug+"?width=50"))

		// Then
		assert.Equal(t, "200 OK", w.Result().Status)
	})
}

func TestUploadHandlerFunc(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		bucket     = mock_storage.NewMockStorage(ctrl)
		repository = mock_storage.NewMockFileRepository(ctrl)
	)

	newUploadRequest := func(fileName string, body io.Reader) *http.Request {
		var buf bytes.Buffer
		multipartForm := multipart.NewWriter(&buf)
		defer multipartForm.Close()

		w, _ := multipartForm.CreateFormFile("file", fileName)
		_, _ = io.Copy(w, body)

		req := httptest.NewRequest(http.MethodGet, "/api/v2.1/storage/upload", &buf)
		req.Header.Set("Content-Type", multipartForm.FormDataContentType())

		return req
	}

	t.Run("With successful uploading file", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()

		bucket.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, _ io.Reader, path string) error {
			id, err := Slug(path).GetID()

			assert.Nil(t, err)
			assert.False(t, id.(primitive.ObjectID).IsZero())

			f := File{
				ID:       id.(primitive.ObjectID),
				Path:     path,
				FileName: "test.txt",
				Slug:     filepath.Base(path),
			}
			repository.EXPECT().Create(gomock.Any(), f).Return(f, nil)

			return nil
		})

		// When
		UploadHandlerFunc(bucket, repository).ServeHTTP(w, withAuthorizedID(newUploadRequest("test.txt", bytes.NewBufferString("test"))))

		// Then
		var f File
		err := json.Unmarshal(w.Body.Bytes(), &f)
		assert.Nil(t, err)
		assert.Equal(t, "test.txt", f.FileName)
		assert.Equal(t, filepath.Join("authorizedID", "test-"+f.ID.Hex()+".txt"), f.Path)
	})

	t.Run("When unable to retrieve authorizedID", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()

		// When
		UploadHandlerFunc(bucket, repository).ServeHTTP(w, newUploadRequest("test.txt", bytes.NewBufferString("test")))

		// Then
		assert.Equal(t, `{"error":{"code":401,"message":"Unauthorized"}}`, w.Body.String())
	})

	t.Run("When unable to read body", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()

		// When
		UploadHandlerFunc(bucket, repository).ServeHTTP(w, withAuthorizedID(newUploadRequest("", bytes.NewBufferString("test"))))

		// Then
		assert.Equal(t, `{"error":{"code":500,"message":"http: no such file"}}`, w.Body.String())
	})

	t.Run("When unable to upload file to storage", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()

		bucket.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test unable to upload file to storage"))

		// When
		UploadHandlerFunc(bucket, repository).ServeHTTP(w, withAuthorizedID(newUploadRequest("test.txt", bytes.NewBufferString("test"))))

		// Then
		assert.Equal(t, `{"error":{"code":500,"message":"test unable to upload file to storage"}}`, w.Body.String())
	})

	t.Run("When unable to create new record on database", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()

		bucket.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		repository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(File{}, errors.New("test unable to create new record on database"))

		// When
		UploadHandlerFunc(bucket, repository).ServeHTTP(w, withAuthorizedID(newUploadRequest("test.txt", bytes.NewBufferString("test"))))

		// Then
		assert.Equal(t, `{"error":{"code":500,"message":"test unable to create new record on database"}}`, w.Body.String())
	})
}

func withAuthorizedID(r *http.Request) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), auth.UserProperty, &jwt.Token{
		Claims: jwt.MapClaims{
			"sub": "authorizedID",
		},
	}))
}
