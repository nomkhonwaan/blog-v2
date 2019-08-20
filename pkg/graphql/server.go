package graphql

import (
	"context"

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

// Schema builts the GraphQL schema
func (s *Server) Schema() *graphql.Schema {
	builder := schemabuilder.NewSchema()
	s.registerQuery(builder)

	return builder.MustBuild()
}

func (s *Server) registerQuery(schema *schemabuilder.Schema) {
	obj := schema.Query()

	obj.FieldFunc("categories", func(ctx context.Context) ([]blog.Category, error) {
		return s.service.Category().FindAll(ctx)
	})
}
