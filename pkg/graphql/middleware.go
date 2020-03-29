package graphql

import (
	"context"
	"errors"
	"github.com/nomkhonwaan/myblog/pkg/auth"
	"github.com/samsarahq/thunder/graphql"
	"net/http"
)

// AuthorizedUserIDProperty is a name of the property in the context where the authorized user ID stored
const AuthorizedUserIDProperty = "authID"

var (
	protectedResources = map[string]bool{
		"myPosts": true,
	}
)

// VerifyAuthorityMiddleware looks on the request header for the authorization token
func VerifyAuthorityMiddleware(input *graphql.ComputationInput, next graphql.MiddlewareNextFunc) *graphql.ComputationOutput {
	authID := auth.GetAuthorizedUserID(input.Ctx)

	for _, sel := range input.ParsedQuery.Selections {
		if yes := protectedResources[sel.Name]; yes {
			if authID == nil {
				return &graphql.ComputationOutput{
					Error: errors.New(http.StatusText(http.StatusUnauthorized)),
				}
			}
		}
	}

	input.Ctx = context.WithValue(input.Ctx, AuthorizedUserIDProperty, authID)
	return next(input)
}
