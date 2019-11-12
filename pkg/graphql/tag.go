package graphql

import (
	"context"
	"github.com/nomkhonwaan/myblog/pkg/blog"
	"github.com/samsarahq/thunder/graphql/schemabuilder"
)

func (s *Server) registerTag(schema *schemabuilder.Schema) {
	obj := schema.Object("Tag", blog.Tag{})

	obj.FieldFunc("latestPublishedPosts", s.tagLatestPublishedPostsFieldFunc)
}

func (s *Server) tagLatestPublishedPostsFieldFunc(ctx context.Context, tag blog.Tag, args struct{ Offset, Limit int64 }) ([]blog.Post, error) {
	return s.service.Post().FindAll(ctx, blog.NewPostQueryBuilder().WithTag(tag).WithOffset(args.Offset).WithLimit(args.Limit).Build())
}
