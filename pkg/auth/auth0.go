package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

// The name of the property in the request where the user information from the JWT will be stored.
// Default value: "user"
const UserProperty = "user"

// GetAuthorizedUserID returns an authorized user ID (which is generated from the authentication server)
func GetAuthorizedUserID(ctx context.Context) interface{} {
	if ctx.Value(UserProperty) == nil {
		return nil
	}

	return ctx.Value(UserProperty).(*jwt.Token).Claims.(jwt.MapClaims)["sub"]
}

// JWTMiddleware is a wrapped struct to the `jwtmiddleware.JWTMiddleware` for testing purpose
type JWTMiddleware struct {
	// The original JWTMiddleware
	*jwtmiddleware.JWTMiddleware

	// An HTTP client for retrieving JSON Web Key Set (JWKS)
	transport http.RoundTripper
}

// NewJWTMiddleware creates new JWT middleware function from the given configuration.
// The JWT middleware will looking for an access_token in the incoming request header (Authorization: Bearer),
// then call to Auth0 service for checking the access_token.
func NewJWTMiddleware(audience, issuer, jwksURI string, transport http.RoundTripper) *jwtmiddleware.JWTMiddleware {
	return jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			if !token.Claims.(jwt.MapClaims).VerifyAudience(audience, false) {
				return nil, errors.New("invalid audience")
			}

			if !token.Claims.(jwt.MapClaims).VerifyIssuer(issuer, false) {
				return nil, errors.New("invalid issuer")
			}

			cert, err := getPEMCertificate(token, jwksURI, transport)
			if err != nil {
				return nil, err
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result, nil
		},
		UserProperty:        UserProperty,
		CredentialsOptional: true,
		SigningMethod:       jwt.SigningMethodRS256,
	})
}

func getPEMCertificate(token *jwt.Token, jwksURI string, transport http.RoundTripper) (string, error) {
	c := &http.Client{Transport: transport}
	req, _ := http.NewRequest(http.MethodGet, jwksURI, nil)
	res, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var jwks struct {
		Keys []struct {
			Kty string   `json:"kty"`
			Kid string   `json:"kid"`
			Use string   `json:"use"`
			N   string   `json:"n"`
			E   string   `json:"e"`
			X5c []string `json:"x5c"`
		} `json:"keys"`
	}
	if err = json.NewDecoder(res.Body).Decode(&jwks); err != nil {
		return "", err
	}

	var cert string
	for k, _ := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = fmt.Sprintf("-----BEGIN CERTIFICATE-----\n%s\n-----END CERTIFICATE-----", jwks.Keys[k].X5c[0])
		}
	}

	if cert == "" {
		return "", errors.New("unable to find appropriate key")
	}

	return cert, nil
}
