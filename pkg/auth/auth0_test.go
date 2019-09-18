package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewJWTMiddleware(t *testing.T) {
	// Given
	audience := "https://api.nomkhonwaan.com"
	issuer := "https://nomkhonwaan.auth0.com"
	jwksURI := "https://nomkhonwaan.auth0.com/.well-known/jwks.json"

	// When
	jwtMiddleware := NewJWTMiddleware(audience, issuer, jwksURI)

	// Then
	assert.Equal(t, UserProperty, jwtMiddleware.Options.UserProperty)
	assert.Equal(t, true, jwtMiddleware.Options.CredentialsOptional)
	assert.Equal(t, jwt.SigningMethodRS256, jwtMiddleware.Options.SigningMethod)
}
