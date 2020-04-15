package graphql

import (
	"context"
	"errors"
	"github.com/nomkhonwaan/myblog/pkg/auth"
	"github.com/samsarahq/thunder/graphql"
	"net/http"
)

// AuthorizedID is a context.Context key where an authorized ID value stored
const AuthorizedID = "authID"

var (
	protectedResources = map[string]bool{
		"myPosts":                 true,
		"createPost":              true,
		"updatePostTitle":         true,
		"updatePostStatus":        true,
		"updatePostContent":       true,
		"updatePostCategories":    true,
		"updatePostTags":          true,
		"updatePostFeaturedImage": true,
		"updatePostAttachments":   true,
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

	input.Ctx = context.WithValue(input.Ctx, AuthorizedID, authID)
	return next(input)
}
