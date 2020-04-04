package storage_test

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/nomkhonwaan/myblog/pkg/auth"
	. "github.com/nomkhonwaan/myblog/pkg/storage"
	mock_storage "github.com/nomkhonwaan/myblog/pkg/storage/mock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		req := httptest.NewRequest(http.MethodDelete, "/api/v2.1/storage/delete/"+slug, nil)
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

}

func TestUploadHandlerFunc(t *testing.T) {

}

func withAuthorizedID(r *http.Request) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), auth.UserProperty, &jwt.Token{
		Claims: jwt.MapClaims{
			"sub": "authorizedID",
		},
	}))
}

//
//import (
//	"bytes"
//	"context"
//	"errors"
//	"github.com/dgrijalva/jwt-go"
//	"github.com/golang/mock/gomock"
//	"github.com/gorilla/mux"
//	"github.com/nomkhonwaan/myblog/pkg/auth"
//	mock_image "github.com/nomkhonwaan/myblog/pkg/image/mock"
//	. "github.com/nomkhonwaan/myblog/pkg/storage"
//	mock_storage "github.com/nomkhonwaan/myblog/pkg/storage/mock"
//	"github.com/stretchr/testify/assert"
//	"go.mongodb.org/mongo-driver/bson/primitive"
//	"image"
//	"image/png"
//	"io"
//	"mime/multipart"
//	"net/http"
//	"net/http/httptest"
//	"path/filepath"
//	"testing"
//)
//
//func TestHandler_Delete(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	var (
//		resizer        = mock_image.NewMockResizer(ctrl)
//		cacheService   = mock_storage.NewMockCache(ctrl)
//		storageService = mock_storage.NewMockStorage(ctrl)
//		file           = mock_storage.NewMockFileRepository(ctrl)
//	)
//
//	r := mux.NewRouter()
//	NewHandler(cacheService, storageService, file, resizer).Register(r.PathPrefix("/v1/storage").Subrouter())
//
//	newDeleteRequest := func(slug string) *http.Request {
//		return httptest.NewRequest(http.MethodDelete, "/v1/storage/delete/"+slug, nil)
//	}
//
//	withAuthorizedID := func(r *http.Request) *http.Request {
//		return r.WithContext(context.WithValue(context.Background(), auth.UserProperty, &jwt.Token{
//			Claims: jwt.MapClaims{
//				"sub": "authorizedID",
//			},
//		}))
//	}
//
//	t.Run("With successful deleting file", func(t *testing.T) {
//		// Given
//		w := httptest.NewRecorder()
//		id := primitive.NewObjectID()
//		slug := "test-" + id.Hex() + ".txt"
//
//		file.EXPECT().FindByID(gomock.Any(), id).Return(File{Path: "authorizedID/" + slug}, nil)
//		storageService.EXPECT().Delete(gomock.Any(), "authorizedID/"+slug).Return(nil)
//		file.EXPECT().Delete(gomock.Any(), id).Return(nil)
//
//		// When
//		r.ServeHTTP(w, withAuthorizedID(newDeleteRequest(slug)))
//
//		// Then
//		assert.Equal(t, "200 OK", w.Result().Status)
//	})
//
//	t.Run("When unable to retrieve authorized ID from", func(t *testing.T) {
//		// Given
//		w := httptest.NewRecorder()
//		id := primitive.NewObjectID()
//		slug := "test-" + id.Hex() + ".txt"
//
//		// When
//		r.ServeHTTP(w, newDeleteRequest(slug))
//
//		// Then
//		assert.Equal(t, "401 Unauthorized", w.Result().Status)
//	})
//
//	t.Run("When unable to find file on database", func(t *testing.T) {
//		// Given
//		w := httptest.NewRecorder()
//		id := primitive.NewObjectID()
//		slug := "test-" + id.Hex() + ".txt"
//
//		file.EXPECT().FindByID(gomock.Any(), id).Return(File{}, errors.New("test unable to find file on database"))
//
//		// When
//		r.ServeHTTP(w, withAuthorizedID(newDeleteRequest(slug)))
//
//		// Then
//		assert.Equal(t, "404 Not Found", w.Result().Status)
//	})
//
//	t.Run("When unable to delete file", func(t *testing.T) {
//		// Given
//		w := httptest.NewRecorder()
//		id := primitive.NewObjectID()
//		slug := "test-" + id.Hex() + ".txt"
//
//		file.EXPECT().FindByID(gomock.Any(), id).Return(File{Path: "authorizedID/" + slug}, nil)
//		storageService.EXPECT().Delete(gomock.Any(), "authorizedID/"+slug).Return(errors.New("test unable to delete file"))
//
//		// When
//		r.ServeHTTP(w, withAuthorizedID(newDeleteRequest(slug)))
//
//		// Then
//		assert.Equal(t, "500 Internal Server Error", w.Result().Status)
//	})
//
//	t.Run("When unable to delete file from database", func(t *testing.T) {
//		// Given
//		w := httptest.NewRecorder()
//		id := primitive.NewObjectID()
//		slug := "test-" + id.Hex() + ".txt"
//
//		file.EXPECT().FindByID(gomock.Any(), id).Return(File{Path: "authorizedID/" + slug}, nil)
//		storageService.EXPECT().Delete(gomock.Any(), "authorizedID/"+slug).Return(nil)
//		file.EXPECT().Delete(gomock.Any(), id).Return(errors.New("test unable to delete file from database"))
//
//		// When
//		r.ServeHTTP(w, withAuthorizedID(newDeleteRequest(slug)))
//
//		// Then
//		assert.Equal(t, "500 Internal Server Error", w.Result().Status)
//	})
//}
//
//func TestHandler_Download(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	var (
//		resizer        = mock_image.NewMockResizer(ctrl)
//		cacheService   = mock_storage.NewMockCache(ctrl)
//		storageService = mock_storage.NewMockStorage(ctrl)
//		file           = mock_storage.NewMockFileRepository(ctrl)
//	)
//
//	r := mux.NewRouter()
//	NewHandler(cacheService, storageService, file, resizer).Register(r.PathPrefix("/v1/storage").Subrouter())
//
//	newDownloadRequest := func(slug string) *http.Request {
//		return httptest.NewRequest(http.MethodGet, "/v1/storage/"+slug, nil)
//	}
//
//	withAuthorizedID := func(r *http.Request) *http.Request {
//		return r.WithContext(context.WithValue(context.Background(), auth.UserProperty, &jwt.Token{
//			Claims: jwt.MapClaims{
//				"sub": "authorizedID",
//			},
//		}))
//	}
//
//	t.Run("With successful downloading file", func(t *testing.T) {
//		// Given
//		w := httptest.NewRecorder()
//		id := primitive.NewObjectID()
//		slug := "test-" + id.Hex() + ".txt"
//		body := bytes.NewBufferString("test")
//
//		file.EXPECT().FindByID(gomock.Any(), id).Return(File{Path: "authorizedID/" + slug}, nil)
//		cacheService.EXPECT().Exist("authorizedID/" + slug).Return(false)
//		cacheService.EXPECT().Store(gomock.Any(), "authorizedID/"+slug).Return(nil)
//		storageService.EXPECT().Download(gomock.Any(), "authorizedID/"+slug).Return(body, nil)
//
//		// When
//		r.ServeHTTP(w, withAuthorizedID(newDownloadRequest(slug)))
//
//		// Then
//		assert.Equal(t, "200 OK", w.Result().Status)
//	})
//
//	t.Run("With successful downloading file from cache storage", func(t *testing.T) {
//		// Given
//		w := httptest.NewRecorder()
//		id := primitive.NewObjectID()
//		slug := "test-" + id.Hex() + ".txt"
//		body := bytes.NewBufferString("test")
//
//		file.EXPECT().FindByID(gomock.Any(), id).Return(File{Path: "authorizedID/" + slug}, nil)
//		cacheService.EXPECT().Exist("authorizedID/" + slug).Return(true)
//		cacheService.EXPECT().Retrieve("authorizedID/"+slug).Return(body, nil)
//
//		// When
//		r.ServeHTTP(w, withAuthorizedID(newDownloadRequest(slug)))
//
//		// Then
//		assert.Equal(t, "200 OK", w.Result().Status)
//	})
//
//	t.Run("With successful retrieving and resizing image file", func(t *testing.T) {
//		// Given
//		w := httptest.NewRecorder()
//		id := primitive.NewObjectID()
//		slug := "test-" + id.Hex() + ".png"
//		var body bytes.Buffer
//		_ = png.Encode(&body, image.NewRGBA(image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 100, Y: 100}}))
//
//		file.EXPECT().FindByID(gomock.Any(), id).Return(File{Path: "authorizedID/" + slug}, nil)
//		cacheService.EXPECT().Exist("authorizedID/" + slug).Return(true)
//		cacheService.EXPECT().Retrieve("authorizedID/"+slug).Return(&body, nil)
//		cacheService.EXPECT().Exist("authorizedID/test-" + id.Hex() + "-50-0.png").Return(false)
//		resizer.EXPECT().Resize(gomock.Any(), 50, 0).Return(nil, nil)
//		cacheService.EXPECT().Store(gomock.Any(), "authorizedID/test-"+id.Hex()+"-50-0.png").Return(nil)
//
//		// When
//		r.ServeHTTP(w, newDownloadRequest(slug+"?width=50"))
//
//		// Then
//		assert.Equal(t, "200 OK", w.Result().Status)
//	})
//
//	t.Run("With successful retrieving resized image file from cache storage", func(t *testing.T) {
//		// Given
//		w := httptest.NewRecorder()
//		id := primitive.NewObjectID()
//		slug := "test-" + id.Hex() + ".png"
//		var body bytes.Buffer
//		_ = png.Encode(&body, image.NewRGBA(image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 50, Y: 50}}))
//
//		file.EXPECT().FindByID(gomock.Any(), id).Return(File{Path: "authorizedID/" + slug}, nil)
//		cacheService.EXPECT().Exist("authorizedID/test-" + id.Hex() + "-50-0.png").Return(true)
//		cacheService.EXPECT().Retrieve("authorizedID/test-"+id.Hex()+"-50-0.png").Return(&body, nil)
//
//		// When
//		r.ServeHTTP(w, newDownloadRequest(slug+"?width=50"))
//
//		// Then
//		assert.Equal(t, "200 OK", w.Result().Status)
//	})
//
//	t.Run("When unable to find file on database", func(t *testing.T) {
//		// Given
//		w := httptest.NewRecorder()
//		id := primitive.NewObjectID()
//		slug := "test-" + id.Hex() + ".txt"
//
//		file.EXPECT().FindByID(gomock.Any(), id).Return(File{}, errors.New("test unable to find file on database"))
//
//		// When
//		r.ServeHTTP(w, withAuthorizedID(newDownloadRequest(slug)))
//
//		// Then
//		assert.Equal(t, "404 Not Found", w.Result().Status)
//	})
//
//	t.Run("When unable to retrieve file from cache storage even file exist", func(t *testing.T) {
//		// Given
//		w := httptest.NewRecorder()
//		id := primitive.NewObjectID()
//		slug := "test-" + id.Hex() + ".txt"
//		body := bytes.NewBufferString("test")
//
//		file.EXPECT().FindByID(gomock.Any(), id).Return(File{Path: "authorizedID/" + slug}, nil)
//		cacheService.EXPECT().Exist("authorizedID/" + slug).Return(true)
//		cacheService.EXPECT().Retrieve("authorizedID/"+slug).Return(nil, errors.New("test unable to retrieve file from cache storage"))
//		storageService.EXPECT().Download(gomock.Any(), "authorizedID/"+slug).Return(body, nil)
//		cacheService.EXPECT().Store(gomock.Any(), "authorizedID/"+slug).Return(nil)
//
//		// When
//		r.ServeHTTP(w, withAuthorizedID(newDownloadRequest(slug)))
//
//		// Then
//		assert.Equal(t, "200 OK", w.Result().Status)
//	})
//
//	t.Run("When unable to retrieve resized image file from cache storage even file exist", func(t *testing.T) {
//		// Given
//		w := httptest.NewRecorder()
//		id := primitive.NewObjectID()
//		slug := "test-" + id.Hex() + ".png"
//		var body bytes.Buffer
//		_ = png.Encode(&body, image.NewRGBA(image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 100, Y: 100}}))
//
//		file.EXPECT().FindByID(gomock.Any(), id).Return(File{Path: "authorizedID/" + slug}, nil)
//		cacheService.EXPECT().Exist("authorizedID/test-" + id.Hex() + "-50-0.png").Return(true)
//		cacheService.EXPECT().Retrieve("authorizedID/test-"+id.Hex()+"-50-0.png").Return(nil, errors.New("test unable to retrieve file from cache storage"))
//		cacheService.EXPECT().Exist("authorizedID/" + slug).Return(true)
//		cacheService.EXPECT().Retrieve("authorizedID/"+slug).Return(&body, nil)
//		resizer.EXPECT().Resize(gomock.Any(), 50, 0).Return(nil, nil)
//		cacheService.EXPECT().Store(gomock.Any(), "authorizedID/test-"+id.Hex()+"-50-0.png").Return(nil)
//
//		// When
//		r.ServeHTTP(w, newDownloadRequest(slug+"?width=50"))
//
//		// Then
//		assert.Equal(t, "200 OK", w.Result().Status)
//	})
//
//	t.Run("When unable to resize image file", func(t *testing.T) {
//		// Given
//		w := httptest.NewRecorder()
//		id := primitive.NewObjectID()
//		slug := "test-" + id.Hex() + ".png"
//		var body bytes.Buffer
//		_ = png.Encode(&body, image.NewRGBA(image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 100, Y: 100}}))
//
//		file.EXPECT().FindByID(gomock.Any(), id).Return(File{Path: "authorizedID/" + slug}, nil)
//		cacheService.EXPECT().Exist("authorizedID/" + slug).Return(true)
//		cacheService.EXPECT().Retrieve("authorizedID/"+slug).Return(&body, nil)
//		cacheService.EXPECT().Exist("authorizedID/test-" + id.Hex() + "-50-0.png").Return(false)
//		resizer.EXPECT().Resize(gomock.Any(), 50, 0).Return(nil, errors.New("test unable to resize image file"))
//		cacheService.EXPECT().Store(gomock.Any(), "authorizedID/test-"+id.Hex()+"-50-0.png").Return(nil)
//
//		// When
//		r.ServeHTTP(w, newDownloadRequest(slug+"?width=50"))
//
//		// Then
//		assert.Equal(t, "200 OK", w.Result().Status)
//	})
//
//	t.Run("When unable to download file", func(t *testing.T) {
//		// Given
//		w := httptest.NewRecorder()
//		id := primitive.NewObjectID()
//		slug := "test-" + id.Hex() + ".txt"
//
//		file.EXPECT().FindByID(gomock.Any(), id).Return(File{Path: "authorizedID/" + slug}, nil)
//		cacheService.EXPECT().Exist("authorizedID/" + slug).Return(false)
//		storageService.EXPECT().Download(gomock.Any(), "authorizedID/"+slug).Return(nil, errors.New("test unable to downloading file"))
//
//		// When
//		r.ServeHTTP(w, withAuthorizedID(newDownloadRequest(slug)))
//
//		// Then
//		assert.Equal(t, "404 Not Found", w.Result().Status)
//	})
//
//	t.Run("When unable to save file to cache storage", func(t *testing.T) {
//		// Given
//		w := httptest.NewRecorder()
//		id := primitive.NewObjectID()
//		slug := "test-" + id.Hex() + ".txt"
//		body := bytes.NewBufferString("test")
//
//		file.EXPECT().FindByID(gomock.Any(), id).Return(File{Path: "authorizedID/" + slug}, nil)
//		cacheService.EXPECT().Exist("authorizedID/" + slug).Return(false)
//		cacheService.EXPECT().Store(gomock.Any(), "authorizedID/"+slug).Return(errors.New("test unable to save file to cache storage"))
//		storageService.EXPECT().Download(gomock.Any(), "authorizedID/"+slug).Return(body, nil)
//
//		// When
//		r.ServeHTTP(w, withAuthorizedID(newDownloadRequest(slug)))
//
//		// Then
//		assert.Equal(t, "200 OK", w.Result().Status)
//	})
//
//	t.Run("When unable to save resized image file to cache storage", func(t *testing.T) {
//		// Given
//		w := httptest.NewRecorder()
//		id := primitive.NewObjectID()
//		slug := "test-" + id.Hex() + ".png"
//		var body bytes.Buffer
//		_ = png.Encode(&body, image.NewRGBA(image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 100, Y: 100}}))
//
//		file.EXPECT().FindByID(gomock.Any(), id).Return(File{Path: "authorizedID/" + slug}, nil)
//		cacheService.EXPECT().Exist("authorizedID/" + slug).Return(true)
//		cacheService.EXPECT().Retrieve("authorizedID/"+slug).Return(&body, nil)
//		cacheService.EXPECT().Exist("authorizedID/test-" + id.Hex() + "-50-0.png").Return(false)
//		resizer.EXPECT().Resize(gomock.Any(), 50, 0).Return(nil, nil)
//		cacheService.EXPECT().Store(gomock.Any(), "authorizedID/test-"+id.Hex()+"-50-0.png").Return(errors.New("test unable to save resized image file to cache storage"))
//
//		// When
//		r.ServeHTTP(w, newDownloadRequest(slug+"?width=50"))
//
//		// Then
//		assert.Equal(t, "200 OK", w.Result().Status)
//	})
//}
//
//func TestHandler_Upload(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	var (
//		resizer        = mock_image.NewMockResizer(ctrl)
//		cacheService   = mock_storage.NewMockCache(ctrl)
//		storageService = mock_storage.NewMockStorage(ctrl)
//		file           = mock_storage.NewMockFileRepository(ctrl)
//	)
//
//	r := mux.NewRouter()
//	NewHandler(cacheService, storageService, file, resizer).Register(r.PathPrefix("/v1/storage").Subrouter())
//
//	newUploadRequest := func(fileName string, body io.Reader) *http.Request {
//		buf := &bytes.Buffer{}
//		wtr := multipart.NewWriter(buf)
//		defer wtr.Close()
//
//		w, _ := wtr.CreateFormFile("file", fileName)
//		_, _ = io.Copy(w, body)
//
//		r := httptest.NewRequest(http.MethodPost, "/v1/storage/upload", buf)
//		r.Header.Set("Content-Type", wtr.FormDataContentType())
//
//		return r
//	}
//
//	withAuthorizedID := func(r *http.Request) *http.Request {
//		return r.WithContext(context.WithValue(context.Background(), auth.UserProperty, &jwt.Token{
//			Claims: jwt.MapClaims{
//				"sub": "authorizedID",
//			},
//		}))
//	}
//
//	t.Run("With successful uploading file", func(t *testing.T) {
//		// Given
//		w := httptest.NewRecorder()
//		body := bytes.NewBufferString("test")
//		fileName := "test.txt"
//
//		storageService.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, _ io.Reader, path string) error {
//			id, err := Slug(path).GetID()
//
//			assert.Nil(t, err)
//			assert.False(t, id.(primitive.ObjectID).IsZero())
//
//			file.EXPECT().Create(gomock.Any(), File{
//				ID:       id.(primitive.ObjectID),
//				Path:     path,
//				FileName: "test.txt",
//				Slug:     filepath.Base(path),
//			}).Return(File{}, nil)
//
//			return nil
//		})
//
//		// When
//		r.ServeHTTP(w, withAuthorizedID(newUploadRequest(fileName, body)))
//
//		// Then
//		assert.Equal(t, "200 OK", w.Result().Status)
//	})
//
//	t.Run("When unable to retrieve authorized ID from", func(t *testing.T) {
//		// Given
//		w := httptest.NewRecorder()
//		body := bytes.NewBufferString("test")
//		fileName := "test.txt"
//
//		// When
//		r.ServeHTTP(w, newUploadRequest(fileName, body))
//
//		// Then
//		assert.Equal(t, "401 Unauthorized", w.Result().Status)
//	})
//
//	t.Run("When unable to read form file", func(t *testing.T) {
//		// Given
//		w := httptest.NewRecorder()
//		body := bytes.NewBufferString("")
//
//		// When
//		r.ServeHTTP(w, withAuthorizedID(newUploadRequest("", body)))
//
//		// Then
//		assert.Equal(t, "500 Internal Server Error", w.Result().Status)
//	})
//
//	t.Run("When unable to upload file to the storage server", func(t *testing.T) {
//		// Given
//		w := httptest.NewRecorder()
//		body := bytes.NewBufferString("test")
//		fileName := "test.txt"
//
//		storageService.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test upload file error"))
//
//		// When
//		r.ServeHTTP(w, withAuthorizedID(newUploadRequest(fileName, body)))
//
//		// Then
//		assert.Equal(t, "500 Internal Server Error", w.Result().Status)
//	})
//
//	t.Run("When unable to create a new record on database", func(t *testing.T) {
//		// Given
//		w := httptest.NewRecorder()
//		body := bytes.NewBufferString("test")
//		fileName := "test.txt"
//
//		storageService.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
//		file.EXPECT().Create(gomock.Any(), gomock.Any()).Return(File{}, errors.New("test unable to create a new record on database"))
//
//		// When
//		r.ServeHTTP(w, withAuthorizedID(newUploadRequest(fileName, body)))
//
//		// Then
//		assert.Equal(t, "500 Internal Server Error", w.Result().Status)
//	})
//}
