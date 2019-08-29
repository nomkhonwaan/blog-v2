package graphql

import (
	"context"
	"errors"

	"github.com/nomkhonwaan/myblog/pkg/blog"
	"github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/graphql/schemabuilder"
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

	return builder.MustBuild()
}

func (s *Server) registerQuery(schema *schemabuilder.Schema) {
	obj := schema.Query()

	obj.FieldFunc("categories", s.makeFieldFuncCategories)
	obj.FieldFunc("latestPublished", s.makeFieldFuncLatestPublishedPosts)
}

func (s *Server) registerMutation(schema *schemabuilder.Schema) {
	obj := schema.Mutation()

	obj.FieldFunc("createPost", s.makeFieldFuncCreatePost)
}

func (s *Server) makeFieldFuncCategories(ctx context.Context) ([]blog.Category, error) {
	return s.service.Category().FindAll(ctx)
}

func (s *Server) makeFieldFuncLatestPublishedPosts(ctx context.Context) ([]blog.Post, error) {
	return s.service.Post().FindAll(ctx, blog.NewPostQueryBuilder().WithStatus(blog.Published).Build())
}

// TODO: implement create new post mutation which requires nothing but returns an empty post with "DRAFT" status
func (s *Server) makeFieldFuncCreatePost(ctx context.Context) (blog.Post, error) {
	return blog.Post{}, errors.New("not implement yet")
}
