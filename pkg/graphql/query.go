package graphql

import (
	"context"
	"errors"
	"github.com/nomkhonwaan/myblog/pkg/blog"
	"github.com/samsarahq/thunder/graphql/schemabuilder"
	"net/http"
)

// RegisterQuery registers pre-defined query fields to the provided schema
func (s *Server) RegisterQuery(schema *schemabuilder.Schema) {
	obj := schema.Query()

	obj.FieldFunc("category", s.findCategoryBySlugQuery)
	obj.FieldFunc("categories", s.findAllCategoriesQuery)
	obj.FieldFunc("tag", s.findTagBySlugQuery)
	obj.FieldFunc("tags", s.findAllTagsQuery)
	obj.FieldFunc("latestPublishedPosts", s.findLatestPublishedPostsQuery)
	obj.FieldFunc("post", s.findPostBySlugQuery)
}

// query {
//	category(slug: sting!) {
//		...
//	}
// }
func (s *Server) findCategoryBySlugQuery(ctx context.Context, args struct{ Slug Slug }) (blog.Category, error) {
	id := args.Slug.MustGetID()
	return s.service.Category().FindByID(ctx, id)
}

// query {
//	categories {
//		...
//	}
// }
func (s *Server) findAllCategoriesQuery(ctx context.Context) ([]blog.Category, error) {
	return s.service.Category().FindAll(ctx)
}

// query {
//	tag(slug: string!) {
//		...
//	}
// }
func (s *Server) findTagBySlugQuery(ctx context.Context, args struct{ Slug Slug }) (blog.Tag, error) {
	id := args.Slug.MustGetID()
	return s.service.Tag().FindByID(ctx, id)
}

// query {
//	tags {
//		...
//	}
// }
func (s *Server) findAllTagsQuery(ctx context.Context) ([]blog.Tag, error) {
	return s.service.Tag().FindAll(ctx)
}

// query {
//	latestPublishedPosts(offset: int!, limit: int!) {
//		...
//	}
// }
func (s *Server) findLatestPublishedPostsQuery(ctx context.Context, args struct{ Offset, Limit int64 }) ([]blog.Post, error) {
	return s.service.Post().FindAll(ctx, blog.NewPostQueryBuilder().WithStatus(blog.Published).WithOffset(args.Offset).WithLimit(args.Limit).Build())
}

// query {
//	post(slug: string!) {
//		...
//	}
// }
func (s *Server) findPostBySlugQuery(ctx context.Context, args struct {
	Slug Slug `graphql:"slug"`
}) (blog.Post, error) {
	id := args.Slug.MustGetID()

	p, err := s.service.Post().FindByID(ctx, id)
	if err != nil {
		return blog.Post{}, err
	}

	if p.Status == blog.Published {
		return p, nil
	}

	authorizedID, err := s.getAuthorizedIDOrFailed(ctx)
	if err != nil {
		return blog.Post{}, err
	}

	if p.AuthorID == authorizedID.(string) {
		return p, nil
	}

	return blog.Post{}, errors.New(http.StatusText(http.StatusForbidden))
}
