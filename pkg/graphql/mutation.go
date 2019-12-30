package graphql

import (
	"context"
	"fmt"
	"github.com/nomkhonwaan/myblog/pkg/blog"
	slugify "github.com/nomkhonwaan/myblog/pkg/slug"
	"github.com/nomkhonwaan/myblog/pkg/storage"
	"github.com/russross/blackfriday/v2"
	"github.com/samsarahq/thunder/graphql/schemabuilder"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// RegisterMutation registers pre-defined mutation fields to the provided schema
func (s *Server) RegisterMutation(schema *schemabuilder.Schema) {
	obj := schema.Mutation()

	obj.FieldFunc("createPost", s.createPostMutation)
	obj.FieldFunc("updatePostTitle", s.updatePostTitleMutation)
	obj.FieldFunc("updatePostStatus", s.updatePostStatus)
	obj.FieldFunc("updatePostContent", s.updatePostContentMutation)
	obj.FieldFunc("updatePostCategories", s.updatePostCategoriesMutation)
	obj.FieldFunc("updatePostTags", s.updatePostTagsMutation)
	obj.FieldFunc("updatePostFeaturedImage", s.updatePostFeaturedImageMutation)
	obj.FieldFunc("updatePostAttachments", s.updatePostAttachmentsMutation)
}

// mutation {
//	createPost {
//		...
//	}
// }
func (s *Server) createPostMutation(ctx context.Context) (blog.Post, error) {
	authorizedID, err := s.getAuthorizedID(ctx)
	if err != nil {
		return blog.Post{}, err
	}

	return s.service.Post().Create(ctx, authorizedID.(string))
}

// mutation {
//	updatePostTitle(slug: string!, title: string!) {
//		...
//	}
// }
func (s *Server) updatePostTitleMutation(ctx context.Context, args struct {
	Slug  Slug
	Title string
}) (blog.Post, error) {
	id := args.Slug.MustGetID()

	err := s.validateAuthority(ctx, id)
	if err != nil {
		return blog.Post{}, err
	}

	slug := fmt.Sprintf("%s-%s", slugify.Make(args.Title), id.(primitive.ObjectID).Hex())
	return s.service.Post().Save(ctx, id, blog.NewPostQueryBuilder().WithTitle(args.Title).WithSlug(slug).Build())
}

// mutation {
//	updatePostStatus(slug: string!, status: string!) {
//		...
//	}
// }
func (s *Server) updatePostStatus(ctx context.Context, args struct {
	Slug   Slug
	Status blog.Status
}) (blog.Post, error) {
	id := args.Slug.MustGetID()

	err := s.validateAuthority(ctx, id)
	if err != nil {
		return blog.Post{}, err
	}

	qb := blog.NewPostQueryBuilder().WithStatus(args.Status)

	if args.Status == blog.Published {
		post, _ := s.service.Post().FindByID(ctx, id)
		if post.PublishedAt.IsZero() {
			qb = qb.WithPublishedAt(time.Now())
		}
	}

	return s.service.Post().Save(ctx, id, qb.Build())
}

// mutation {
//	updatePostContent(slug: string!, markdown: string!) {
//		...
//	}
// }
func (s *Server) updatePostContentMutation(ctx context.Context, args struct {
	Slug     Slug
	Markdown string
}) (blog.Post, error) {
	id := args.Slug.MustGetID()

	err := s.validateAuthority(ctx, id)
	if err != nil {
		return blog.Post{}, err
	}
	html := blackfriday.Run([]byte(args.Markdown))

	return s.service.Post().Save(ctx, id, blog.NewPostQueryBuilder().WithMarkdown(args.Markdown).WithHTML(string(html)).Build())
}

// mutation {
//	updatePostCategories(slug: string!, categorySlugs: [string!]!) {
//		...
//	}
// }
func (s *Server) updatePostCategoriesMutation(ctx context.Context, args struct {
	Slug          Slug
	CategorySlugs []Slug
}) (blog.Post, error) {
	id := args.Slug.MustGetID()

	err := s.validateAuthority(ctx, id)
	if err != nil {
		return blog.Post{}, err
	}

	var ids []primitive.ObjectID
	for _, slug := range args.CategorySlugs {
		ids = append(ids, slug.MustGetID().(primitive.ObjectID))
	}

	categories, err := s.service.Category().FindAllByIDs(ctx, ids)
	if err != nil {
		return blog.Post{}, err
	}

	return s.service.Post().Save(ctx, id, blog.NewPostQueryBuilder().WithCategories(categories).Build())
}

// mutation {
//	updatePostTags(slug: string!, tags: [string!]!) {
//		...
//	}
// }
func (s *Server) updatePostTagsMutation(ctx context.Context, args struct {
	Slug     Slug
	TagSlugs []Slug
}) (blog.Post, error) {
	id := args.Slug.MustGetID()

	err := s.validateAuthority(ctx, id)
	if err != nil {
		return blog.Post{}, err
	}

	var ids []primitive.ObjectID
	for _, slug := range args.TagSlugs {
		ids = append(ids, slug.MustGetID().(primitive.ObjectID))
	}

	tags, err := s.service.Tag().FindAllByIDs(ctx, ids)
	if err != nil {
		return blog.Post{}, err
	}

	return s.service.Post().Save(ctx, id, blog.NewPostQueryBuilder().WithTags(tags).Build())
}

// mutation {
//	updatePostFeaturedImage(slug: string!, featuredImageSlug: string!) {
//		...
//	}
// }
func (s *Server) updatePostFeaturedImageMutation(ctx context.Context, args struct {
	Slug              Slug
	FeaturedImageSlug storage.Slug `graphql:",optional"`
}) (blog.Post, error) {
	id := args.Slug.MustGetID()

	err := s.validateAuthority(ctx, id)
	if err != nil {
		return blog.Post{}, err
	}

	if args.FeaturedImageSlug == "" {
		return s.service.Post().Save(ctx, id, blog.NewPostQueryBuilder().WithFeaturedImage(storage.File{}).Build())
	}

	file, err := s.service.File().FindByID(ctx, args.FeaturedImageSlug.MustGetID())
	if err != nil {
		return blog.Post{}, err
	}

	return s.service.Post().Save(ctx, id, blog.NewPostQueryBuilder().WithFeaturedImage(file).Build())
}

// mutation {
//	updatePostAttachments(slug: string!, attachmentSlugs: [string!]!) {
//		...
//	}
// }
func (s *Server) updatePostAttachmentsMutation(ctx context.Context, args struct {
	Slug            Slug
	AttachmentSlugs []storage.Slug
}) (blog.Post, error) {
	id := args.Slug.MustGetID()

	err := s.validateAuthority(ctx, id)
	if err != nil {
		return blog.Post{}, err
	}

	var ids []primitive.ObjectID
	for _, slug := range args.AttachmentSlugs {
		ids = append(ids, slug.MustGetID().(primitive.ObjectID))
	}

	files, err := s.service.File().FindAllByIDs(ctx, ids)
	if err != nil {
		return blog.Post{}, err
	}

	return s.service.Post().Save(ctx, id, blog.NewPostQueryBuilder().WithAttachments(files).Build())
}
