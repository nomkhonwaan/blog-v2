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

// UserProperty is a name of the property in the request where the user information stored
const UserProperty = "user"

// GetAuthorizedUserID returns a user ID which generated from the auth server
func GetAuthorizedUserID(ctx context.Context) interface{} {
	if claims := ctx.Value(UserProperty); claims != nil {
		return claims.(*jwt.Token).Claims.(jwt.MapClaims)["sub"]
	}
	return nil
}

// NewJWTMiddleware returns a new jwtmiddleware.JWTMiddleware instance.
// This middleware uses to looking for the "access_token" in the request header and call to the auth server for validating it.
func NewJWTMiddleware(audience, issuer, jwksURI string, transport http.RoundTripper) *jwtmiddleware.JWTMiddleware {
	return jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			if !token.Claims.(jwt.MapClaims).VerifyAudience(audience, false) {
				return nil, errors.New("invalid audience")
			}

			if !token.Claims.(jwt.MapClaims).VerifyIssuer(issuer, false) {
				return nil, errors.New("invalid issuer")
			}

			cert, err := getPEMCertificate(token, jwksURI, &http.Client{Transport: transport})
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

func getPEMCertificate(token *jwt.Token, jwksURI string, c *http.Client) (string, error) {
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
