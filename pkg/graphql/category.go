package graphql

import (
	"context"
	"github.com/nomkhonwaan/myblog/pkg/blog"
	"github.com/samsarahq/thunder/graphql/schemabuilder"
)

func (s *Server) registerCategory(schema *schemabuilder.Schema) {
	obj := schema.Object("Category", blog.Category{})

	obj.FieldFunc("latestPublishedPosts", s.categoryLatestPublishedPostsFieldFunc)
}

func (s *Server) categoryLatestPublishedPostsFieldFunc(ctx context.Context, cat blog.Category, args struct{ Offset, Limit int64 }) ([]blog.Post, error) {
	return s.service.Post().FindAll(ctx, blog.NewPostQueryBuilder().
		WithCategory(cat).
		WithStatus(blog.StatusPublished).
		WithOffset(args.Offset).
		WithLimit(args.Limit).
		Build(),
	)
}
