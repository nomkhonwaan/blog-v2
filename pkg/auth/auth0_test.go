package auth_test

import (
	"bytes"
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	. "github.com/nomkhonwaan/myblog/pkg/auth"
	mock_http "github.com/nomkhonwaan/myblog/pkg/auth/mock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestNewJWTMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var (
		audience = "https://api.nomkhonwaan.com"
		issuer   = "https://nomkhonwaan.auth0.com"
		jwksURI  = "https://nomkhonwaan.auth0.com/.well-known/jwks.json"

		transport = mock_http.NewMockRoundTripper(ctrl)
	)

	t.Run("With valid token", func(t *testing.T) {
		// Given
		transport.EXPECT().RoundTrip(gomock.Any()).DoAndReturn(func(r *http.Request) (*http.Response, error) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, jwksURI, r.URL.String())

			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString(`
				{
				  "keys": [
				    {
				      "kid": "1",
				      "x5c": [ "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0a0doNkOfm5wC56q+4ojS4KvoEgQ9OqDaMqpZqiWXdqtA5SHH0/cigu9pPaYgAfFS47NNxM5MO2ZD44pqn14dsPqTXqHDfZUgiuVaGGoWCp6ddaqc8flUYJ3ArLjUPwSGCEWV9aFiUB5NF9fd54w3tgLnB+Vu12503w+zBHtnbB1wvce6laNwSvBGyQuA5Vp9ncmbTuOvd4CLECiKmbGWlsWvKE3kJgGtT83vcn2wyzDkdJonT7NK0L5bykNXMuC9WQYYZS9V+ufZam7up2FXzhBlOclSfWjQn5OqiyBSzE4Eeaa0NlsasDT0d7YY/RNmoW84RXCcnBa73qP67+/TQIDAQAB" ]
				    }
				  ]
				}
			`)),
			}, nil
		})

		mw := NewJWTMiddleware(audience, issuer, jwksURI, transport)
		validationKeyGetter := mw.Options.ValidationKeyGetter
		token := &jwt.Token{
			Claims: jwt.MapClaims{
				"aud": audience,
			},
			Header: map[string]interface{}{
				"kid": "1",
			},
		}

		// When
		result, err := validationKeyGetter(token)

		// Then
		assert.Nil(t, err)
		assert.IsType(t, &rsa.PublicKey{}, result)
	})

	t.Run("With invalid audience", func(t *testing.T) {
		// Given
		mw := NewJWTMiddleware(audience, issuer, jwksURI, transport)
		validationKeyGetter := mw.Options.ValidationKeyGetter
		token := &jwt.Token{
			Claims: jwt.MapClaims{
				"aud": "https://www.nomkhonwaan.com",
			},
			Header: map[string]interface{}{
				"kid": "1",
			},
		}

		// When
		result, err := validationKeyGetter(token)

		// Then
		assert.Nil(t, result)
		assert.EqualError(t, err, "invalid audience")
	})

	t.Run("With invalid issuer", func(t *testing.T) {
		// Given
		mw := NewJWTMiddleware(audience, issuer, jwksURI, transport)
		validationKeyGetter := mw.Options.ValidationKeyGetter
		token := &jwt.Token{
			Claims: jwt.MapClaims{
				"aud": "https://api.nomkhonwaan.com",
				"iss": "https://nomkhonwaan.okta.com",
			},
			Header: map[string]interface{}{
				"kid": "1",
			},
		}

		// When
		result, err := validationKeyGetter(token)

		// Then
		assert.Nil(t, result)
		assert.EqualError(t, err, "invalid issuer")
	})

	t.Run("When unable to connect to the JWKS URI", func(t *testing.T) {
		// Given
		transport.EXPECT().RoundTrip(gomock.Any()).Return(nil, errors.New("test unable to connect to the JWKS URI"))

		mw := NewJWTMiddleware(audience, issuer, jwksURI, transport)
		validationKeyGetter := mw.Options.ValidationKeyGetter
		token := &jwt.Token{
			Claims: jwt.MapClaims{
				"aud": audience,
			},
			Header: map[string]interface{}{
				"kid": "1",
			},
		}

		// When
		_, err := validationKeyGetter(token)

		// Then
		assert.EqualError(t, err, "test unable to connect to the JWKS URI")
	})

	t.Run("When unable to decode the JWKS body", func(t *testing.T) {
		// Given
		transport.EXPECT().RoundTrip(gomock.Any()).DoAndReturn(func(r *http.Request) (*http.Response, error) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, jwksURI, r.URL.String())

			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString(`invalid-JWKS-body`)),
			}, nil
		})

		mw := NewJWTMiddleware(audience, issuer, jwksURI, transport)
		validationKeyGetter := mw.Options.ValidationKeyGetter
		token := &jwt.Token{
			Claims: jwt.MapClaims{
				"aud": audience,
			},
			Header: map[string]interface{}{
				"kid": "1",
			},
		}

		// When
		result, err := validationKeyGetter(token)

		// Then
		assert.Nil(t, result)
		assert.IsType(t, &json.SyntaxError{}, err)
	})

	t.Run("When unable to find appropriate key", func(t *testing.T) {
		// Given
		transport.EXPECT().RoundTrip(gomock.Any()).DoAndReturn(func(r *http.Request) (*http.Response, error) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, jwksURI, r.URL.String())

			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString(`
				{
				  "keys": [
				    {
				      "kid": "1",
				      "x5c": [ "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0a0doNkOfm5wC56q+4ojS4KvoEgQ9OqDaMqpZqiWXdqtA5SHH0/cigu9pPaYgAfFS47NNxM5MO2ZD44pqn14dsPqTXqHDfZUgiuVaGGoWCp6ddaqc8flUYJ3ArLjUPwSGCEWV9aFiUB5NF9fd54w3tgLnB+Vu12503w+zBHtnbB1wvce6laNwSvBGyQuA5Vp9ncmbTuOvd4CLECiKmbGWlsWvKE3kJgGtT83vcn2wyzDkdJonT7NK0L5bykNXMuC9WQYYZS9V+ufZam7up2FXzhBlOclSfWjQn5OqiyBSzE4Eeaa0NlsasDT0d7YY/RNmoW84RXCcnBa73qP67+/TQIDAQAB" ]
				    }
				  ]
				}
			`)),
			}, nil
		})

		mw := NewJWTMiddleware(audience, issuer, jwksURI, transport)
		validationKeyGetter := mw.Options.ValidationKeyGetter
		token := &jwt.Token{
			Claims: jwt.MapClaims{
				"aud": audience,
			},
			Header: map[string]interface{}{
				"kid": "2",
			},
		}

		// When
		result, err := validationKeyGetter(token)

		// Then
		assert.Nil(t, result)
		assert.EqualError(t, err, "unable to find appropriate key")
	})
}

func TestGetAuthorizedUserID(t *testing.T) {
	t.Run("When able to retrieve authorized ID from the context", func(t *testing.T) {
		// Given
		token := &jwt.Token{
			Claims: jwt.MapClaims{
				"sub": "test-authorized-id",
			},
		}
		ctx := context.WithValue(context.Background(), UserProperty, token)

		// When
		authorizedID := GetAuthorizedUserID(ctx)

		// Then
		assert.Equal(t, "test-authorized-id", authorizedID.(string))
	})

	t.Run("When unable to retrieve authorized ID from the context", func(t *testing.T) {
		// Given
		ctx := context.Background()

		// When
		authorizedID := GetAuthorizedUserID(ctx)

		// Then
		assert.Nil(t, authorizedID)
	})
}
