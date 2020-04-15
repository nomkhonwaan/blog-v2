package graphql

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/nomkhonwaan/myblog/pkg/auth"
	"github.com/samsarahq/thunder/graphql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVerifyAuthorityMiddleware(t *testing.T) {
	t.Run("With authorized request", func(t *testing.T) {
		// Given
		input := &graphql.ComputationInput{
			Ctx: context.WithValue(context.Background(), auth.UserProperty, &jwt.Token{Claims: jwt.MapClaims{"sub": "authorizedID"}}),
			ParsedQuery: &graphql.Query{
				SelectionSet: &graphql.SelectionSet{
					Selections: []*graphql.Selection{{Name: "createPost"}},
				},
			},
		}
		next := func(input *graphql.ComputationInput) *graphql.ComputationOutput {
			assert.Equal(t, "authorizedID", input.Ctx.Value(AuthorizedID).(string))

			return &graphql.ComputationOutput{}
		}

		// When
		VerifyAuthorityMiddleware(input, next)

		// Then
	})

	t.Run("With unauthorized request", func(t *testing.T) {
		// Given
		input := &graphql.ComputationInput{
			Ctx: context.Background(),
			ParsedQuery: &graphql.Query{
				SelectionSet: &graphql.SelectionSet{
					Selections: []*graphql.Selection{{Name: "createPost"}},
				},
			},
		}
		next := func(input *graphql.ComputationInput) *graphql.ComputationOutput {
			return &graphql.ComputationOutput{}
		}

		// When
		output := VerifyAuthorityMiddleware(input, next)

		// Then
		assert.EqualError(t, output.Error, "Unauthorized")
	})
}
