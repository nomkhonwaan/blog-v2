package graphql

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/nomkhonwaan/myblog/pkg/auth"
	"github.com/nomkhonwaan/myblog/pkg/blog"
	"github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/graphql/schemabuilder"
	"net/http"
)

// Server is our GraphQL server
type Server struct {
	service blog.Service

	// Return list of categories
	categories []blog.Category
}

// NewServer returns new GraphQL server
func NewServer(service blog.Service) *Server {
	return &Server{service: service}
}

// Schema builds the GraphQL schema
func (s *Server) Schema() *graphql.Schema {
	builder := schemabuilder.NewSchema()
	s.registerQuery(builder)
	s.registerMutation(builder)
	s.registerPost(builder)

	return builder.MustBuild()
}

func (s *Server) registerQuery(schema *schemabuilder.Schema) {
	obj := schema.Query()

	obj.FieldFunc("categories", s.makeFieldFuncCategories)
	obj.FieldFunc("latestPublishedPosts", s.makeFieldFuncLatestPublishedPosts)
}

func (s *Server) registerMutation(schema *schemabuilder.Schema) {
	obj := schema.Mutation()

	obj.FieldFunc("createPost", s.makeFieldFuncCreatePost)
}

func (s *Server) registerPost(schema *schemabuilder.Schema) {
	obj := schema.Object("Post", blog.Post{})

	obj.FieldFunc("categories", (blog.Post{}).Categories(s.service.Category()))
	obj.FieldFunc("tags", (blog.Post{}).Tags(s.service.Tag()))
}

func (s *Server) makeFieldFuncCategories(ctx context.Context) ([]blog.Category, error) {
	return s.service.Category().FindAll(ctx)
}

func (s *Server) makeFieldFuncLatestPublishedPosts(ctx context.Context, args struct{ Offset, Limit int64 }) ([]blog.Post, error) {
	return s.service.Post().FindAll(ctx,
		blog.NewPostQueryBuilder().
			WithStatus(blog.Published).
			WithOffset(args.Offset).
			WithLimit(args.Limit).
			Build(),
	)
}

func (s *Server) makeFieldFuncCreatePost(ctx context.Context) (blog.Post, error) {
	if ctx.Value(auth.UserProperty) == nil {
		return blog.Post{}, errors.New(http.StatusText(http.StatusUnauthorized))
	}

	authorID := ctx.Value(auth.UserProperty).(*jwt.Token).Claims.(jwt.MapClaims)["sub"]
	return s.service.Post().Create(ctx, authorID.(string))
}
