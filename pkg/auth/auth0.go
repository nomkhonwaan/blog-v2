package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

// Auth0JWTMiddlewareFunc performs validation against an access_token in the incoming request
func Auth0JWTMiddlewareFunc(audience, issuer, jwksURI string) {
	jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			if ! token.Claims.(jwt.MapClaims).VerifyAudience(audience, false) {
				return nil, errors.New("invalid audience")
			}

			if ! token.Claims.(jwt.MapClaims).VerifyIssuer(issuer, false) {
				return nil, errors.New("invalid issuer")
			}

			certificate, err := getPEMCertificate(token, jwksURI)
			if err != nil {
				return nil, err
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(certificate))
			return result, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})
}

func getPEMCertificate(token *jwt.Token, jwksURI string) (string, error) {
	res, err := http.Get(jwksURI)
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

	var certificate string
	for k, _ := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			certificate = fmt.Sprintf("-----BEGIN CERTIFICATE-----\n%s\n-----END CERTIFICATE-----", jwks.Keys[k].X5c[0])
		}
	}

	if certificate == "" {
		return "", errors.New("unable to find appropriate key")
	}

	return certificate, nil
}
