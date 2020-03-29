package graphql

import (
	"context"
	"errors"
	"fmt"
	"github.com/nomkhonwaan/myblog/pkg/blog"
	"github.com/nomkhonwaan/myblog/pkg/facebook"
	"github.com/nomkhonwaan/myblog/pkg/log"
	slugify "github.com/nomkhonwaan/myblog/pkg/slug"
	"github.com/nomkhonwaan/myblog/pkg/storage"
	"github.com/nomkhonwaan/myblog/pkg/timeutil"
	"github.com/russross/blackfriday/v2"
	"github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/graphql/schemabuilder"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

// BuildSchema accepts build schema function(s) for applying to the schemabuilding.Schema object
func BuildSchema(buildSchemaFunc ...func(object *schemabuilder.Schema)) (*graphql.Schema, error) {
	s := schemabuilder.NewSchema()
	for _, f := range buildSchemaFunc {
		f(s)
	}
	return s.Build()
}

// BuildCategorySchema builds all category related schemas
func BuildCategorySchema(repository blog.CategoryRepository) func(*schemabuilder.Schema) {
	return func(s *schemabuilder.Schema) {
		q := s.Query()
		q.FieldFunc("category", FindCategoryBySlugFieldFunc(repository))
		q.FieldFunc("categories", FindAllCategoriesFieldFunc(repository))
		p := s.Object("Post", blog.Post{})
		p.FieldFunc("categories", FindAllCategoriesBelongedToPostFieldFunc(repository))
	}
}

// BuildTagSchema builds all tag related schemas
func BuildTagSchema(repository blog.TagRepository) func(*schemabuilder.Schema) {
	return func(s *schemabuilder.Schema) {
		q := s.Query()
		q.FieldFunc("tag", FindTagBySlugFieldFunc(repository))
		q.FieldFunc("tags", FindAllTagsFieldFunc(repository))
		p := s.Object("Post", blog.Post{})
		p.FieldFunc("tags", FindAllTagsBelongedToPostFieldFunc(repository))
	}
}

// BuildPostSchema builds all post related schemas
func BuildPostSchema(repository blog.PostRepository, timer log.Timer) func(*schemabuilder.Schema) {
	return func(s *schemabuilder.Schema) {
		q := s.Query()
		q.FieldFunc("latestPublishedPosts", FindAllLatestPublishedPostsFieldFunc(repository))
		q.FieldFunc("myPosts", FindAllMyPostsFieldFunc(repository))
		q.FieldFunc("post", FindPostBySlugFieldFunc(repository))
		m := s.Mutation()
		m.FieldFunc("createPost", CreatePostFieldFunc(repository))
		m.FieldFunc("updatePostTitle", UpdatePostTitleFieldFunc(repository))
		m.FieldFunc("updatePostStatus", UpdatePostStatusFieldFunc(repository, timer))
		m.FieldFunc("updatePostContent", UpdatePostContentFieldFunc(repository))
		m.FieldFunc("updatePostCategories", UpdatePostCategoriesFieldFunc(repository))
		m.FieldFunc("updatePostTags", UpdatePostTagsFieldFunc(repository))
		m.FieldFunc("updatePostFeaturedImage", UpdatePostFeaturedImageFieldFunc(repository))
		m.FieldFunc("updatePostAttachments", UpdatePostAttachmentsFieldFunc(repository))
		c := s.Object("Category", blog.Category{})
		c.FieldFunc("latestPublishedPosts", FindAllLPPBelongedToCategoryFieldFunc(repository))
		t := s.Object("Tag", blog.Tag{})
		t.FieldFunc("latestPublishedPosts", FindAllLPPBelongedToTagFieldFunc(repository))
	}
}

// BuildFileSchema builds all file related schemas
func BuildFileSchema(repository storage.FileRepository) func(*schemabuilder.Schema) {
	return func(s *schemabuilder.Schema) {
		p := s.Object("Post", blog.Post{})
		p.FieldFunc("featuredImage", FindFeaturedImageBelongedToPostFieldFunc(repository))
		p.FieldFunc("attachments", FindAllAttachmentsBelongedToPostFieldFunc(repository))
	}
}

// BuildGraphAPISchema builds all Facebook Graph API related schemas
func BuildGraphAPISchema(baseURL string, c facebook.Client) func(*schemabuilder.Schema) {
	return func(s *schemabuilder.Schema) {
		p := s.Object("Post", blog.Post{})
		p.FieldFunc("engagement", GetURLNodeShareCountFieldFunc(baseURL, c))
	}
}

// FindCategoryBySlugFieldFunc handles the following query
// ```graphql
// 	{
//		category(slug: string!) { ... }
// 	}
// ```
func FindCategoryBySlugFieldFunc(repository blog.CategoryRepository) interface{} {
	return func(ctx context.Context, args struct{ Slug Slug }) (blog.Category, error) {
		return repository.FindByID(ctx, args.Slug.MustGetID())
	}
}

// FindAllCategoriesFieldFunc handles the following query
// ```graphql
//	{
//		categories { ... }
//	}
// ```
func FindAllCategoriesFieldFunc(repository blog.CategoryRepository) interface{} {
	return func(ctx context.Context) ([]blog.Category, error) {
		return repository.FindAll(ctx)
	}
}

// FindAllCategoriesBelongedToPostFieldFunc handles the following query in the Post type
// ```graphql
//	{
//		Post {
//			...
//			categories { ... }
//		}
//	}
// ```
func FindAllCategoriesBelongedToPostFieldFunc(repository blog.CategoryRepository) interface{} {
	return func(ctx context.Context, p blog.Post) ([]blog.Category, error) {
		ids := make([]primitive.ObjectID, len(p.Categories))
		for i, c := range p.Categories {
			ids[i] = c.ID
		}
		return repository.FindAllByIDs(ctx, ids)
	}
}

// FindTagBySlugFieldFunc handles the following query
// ```graphql
//	{
//		tag(slug: string!) { ... }
//	}
// ```
func FindTagBySlugFieldFunc(repository blog.TagRepository) interface{} {
	return func(ctx context.Context, args struct{ Slug Slug }) (blog.Tag, error) {
		return repository.FindByID(ctx, args.Slug.MustGetID())
	}
}

// FindAllTagsFieldFunc handles the following query
// ```graphql
//	{
//		tags { ... }
//	}
// ```
func FindAllTagsFieldFunc(repository blog.TagRepository) interface{} {
	return func(ctx context.Context) ([]blog.Tag, error) {
		return repository.FindAll(ctx)
	}
}

// FindAllTagsBelongedToPostFieldFunc handles the following query in the Post type
// ```graphql
//	{
//		Post {
//			...
//			tags { ... }
//		}
//	}
// ```
func FindAllTagsBelongedToPostFieldFunc(repository blog.TagRepository) interface{} {
	return func(ctx context.Context, p blog.Post) ([]blog.Tag, error) {
		ids := make([]primitive.ObjectID, len(p.Tags))
		for i, t := range p.Tags {
			ids[i] = t.ID
		}
		return repository.FindAllByIDs(ctx, ids)
	}
}

// FindAllLatestPublishedPostsFieldFunc handles the following query
// ```graphql
//	{
//		latestPublishedPosts(offset: int!, limit: int!) { ... }
//	}
// ```
func FindAllLatestPublishedPostsFieldFunc(repository blog.PostRepository) interface{} {
	return func(ctx context.Context, args struct{ Offset, Limit int64 }) ([]blog.Post, error) {
		return repository.FindAll(ctx, blog.NewPostQueryBuilder().WithStatus(blog.StatusPublished).
			WithOffset(args.Offset).WithLimit(args.Limit).Build())
	}
}

// FindAllLPPBelongedToCategoryFieldFunc handles the following query in the Category type
// ```graphql
//	{
//		Category {
//			...
//			latestPublishedPosts(offset: int!, limit: int!) { ... }
//		}
//	}
// ```
func FindAllLPPBelongedToCategoryFieldFunc(repository blog.PostRepository) interface{} {
	return func(ctx context.Context, c blog.Category, args struct{ Offset, Limit int64 }) ([]blog.Post, error) {
		return repository.FindAll(ctx, blog.NewPostQueryBuilder().WithCategory(c).WithStatus(blog.StatusPublished).
			WithOffset(args.Offset).WithLimit(args.Limit).Build())
	}
}

// FindAllLPPBelongedToTagFieldFunc handles the following query in the Tag type
// ```graphql
//	{
//		Tag {
//			...
//			latestPublishedPosts(offset: int!, limit: int!) { ... }
//		}
//	}
// ```
func FindAllLPPBelongedToTagFieldFunc(repository blog.PostRepository) interface{} {
	return func(ctx context.Context, t blog.Tag, args struct{ Offset, Limit int64 }) ([]blog.Post, error) {
		return repository.FindAll(ctx, blog.NewPostQueryBuilder().WithTag(t).WithStatus(blog.StatusPublished).
			WithOffset(args.Offset).WithLimit(args.Limit).Build())
	}
}

// FindAllMyPostsFieldFunc handles the following query
// ```graphql
//	{
//		myPosts(offset: int!, limit: int!) { ... }
//	}
// ```
func FindAllMyPostsFieldFunc(repository blog.PostRepository) interface{} {
	return func(ctx context.Context, args struct{ Offset, Limit int64 }) ([]blog.Post, error) {
		return repository.FindAll(ctx, blog.NewPostQueryBuilder().WithAuthorID(ctx.Value(AuthorizedID).(string)).
			WithOffset(args.Offset).WithLimit(args.Limit).Build())
	}
}

// FindPostBySlugFieldFunc handles the following query
// ```graphql
//	{
//		post(slug: string!) { ... }
//	}
// ```
func FindPostBySlugFieldFunc(repository blog.PostRepository) interface{} {
	return func(ctx context.Context, args struct{ Slug Slug }) (blog.Post, error) {
		id := args.Slug.MustGetID()

		p, err := repository.FindByID(ctx, id)
		if err != nil {
			return blog.Post{}, errors.New(http.StatusText(http.StatusNotFound))
		}
		if p.Status == blog.StatusPublished {
			return p, nil
		}

		if authID := ctx.Value(AuthorizedID); authID != nil {
			if p.AuthorID == authID.(string) {
				return p, nil
			}
		}

		return blog.Post{}, errors.New(http.StatusText(http.StatusForbidden))
	}
}

// CreatePostFieldFunc handles the following mutation
// ```graphql
//	mutation {
//		createPost { ... }
//	}
// ```
func CreatePostFieldFunc(repository blog.PostRepository) interface{} {
	return func(ctx context.Context) (blog.Post, error) {
		return repository.Create(ctx, ctx.Value(AuthorizedID).(string))
	}
}

// UpdatePostTitleFieldFunc handles the following mutation
// ```graphql
//	mutation {
//		updatePostTitle(slug: string!, title: string!) { ... }
//	}
// ```
func UpdatePostTitleFieldFunc(repository blog.PostRepository) interface{} {
	return func(ctx context.Context, args struct {
		Slug  Slug
		Title string
	}) (blog.Post, error) {
		id := args.Slug.MustGetID()

		p, err := repository.FindByID(ctx, id)
		if err != nil {
			return blog.Post{}, errors.New(http.StatusText(http.StatusNotFound))
		}

		if authID := ctx.Value(AuthorizedID); authID != nil {
			if p.AuthorID == authID.(string) {
				slug := fmt.Sprintf("%s-%s", slugify.Make(args.Title), id.(primitive.ObjectID).Hex())
				return repository.Save(ctx, id, blog.NewPostQueryBuilder().WithTitle(args.Title).WithSlug(slug).
					Build())
			}
		}

		return blog.Post{}, errors.New(http.StatusText(http.StatusForbidden))
	}
}

// UpdatePostStatusFieldFunc handles the following mutation
// ```graphql
//	mutation {
//		updatePostStatus(slug: string!, status: Status!) { ... }
//	}
// ```
func UpdatePostStatusFieldFunc(repository blog.PostRepository, timer log.Timer) interface{} {
	return func(ctx context.Context, args struct {
		Slug   Slug
		Status blog.Status
	}) (blog.Post, error) {
		id := args.Slug.MustGetID()

		p, err := repository.FindByID(ctx, id)
		if err != nil {
			return blog.Post{}, errors.New(http.StatusText(http.StatusNotFound))
		}

		if authID := ctx.Value(AuthorizedID); authID != nil {
			if p.AuthorID == authID.(string) {
				qb := blog.NewPostQueryBuilder().WithStatus(args.Status)
				if args.Status.IsPublished() && p.PublishedAt.IsZero() {
					qb.WithPublishedAt(timer.Now())
				}
				return repository.Save(ctx, id, qb.Build())
			}
		}

		return blog.Post{}, errors.New(http.StatusText(http.StatusForbidden))
	}
}

// UpdatePostContentFieldFunc handles the following mutation
// ```graphql
//	mutation {
//		updatePostContent(slug: string!, markdown: string!) { ... }
//	}
// ```
func UpdatePostContentFieldFunc(repository blog.PostRepository) interface{} {
	return func(ctx context.Context, args struct {
		Slug     Slug
		Markdown string
	}) (blog.Post, error) {
		id := args.Slug.MustGetID()

		p, err := repository.FindByID(ctx, id)
		if err != nil {
			return blog.Post{}, errors.New(http.StatusText(http.StatusNotFound))
		}

		if authID := ctx.Value(AuthorizedID); authID != nil {
			if p.AuthorID == authID.(string) {
				html := blackfriday.Run([]byte(args.Markdown), blackfriday.
					WithExtensions(blackfriday.CommonExtensions+blackfriday.Footnotes))
				return repository.Save(ctx, id, blog.NewPostQueryBuilder().WithMarkdown(args.Markdown).
					WithHTML(string(html)).Build())
			}
		}

		return blog.Post{}, errors.New(http.StatusText(http.StatusForbidden))
	}
}

// UpdatePostCategoriesFieldFunc handles the following mutation
// ```graphql
//	mutation {
//		updatePostCategories(slug: string!, categorySlugs: [string!]!) { ... }
//	}
// ```
func UpdatePostCategoriesFieldFunc(repository blog.PostRepository) interface{} {
	return func(ctx context.Context, args struct {
		Slug          Slug
		CategorySlugs []Slug
	}) (blog.Post, error) {
		id := args.Slug.MustGetID()

		p, err := repository.FindByID(ctx, id)
		if err != nil {
			return blog.Post{}, errors.New(http.StatusText(http.StatusNotFound))
		}

		if authID := ctx.Value(AuthorizedID); authID != nil {
			if p.AuthorID == authID.(string) {
				var cats []blog.Category
				for _, slug := range args.CategorySlugs {
					cats = append(cats, blog.Category{ID: slug.MustGetID().(primitive.ObjectID)})
				}
				return repository.Save(ctx, id, blog.NewPostQueryBuilder().WithCategories(cats).Build())
			}
		}

		return blog.Post{}, errors.New(http.StatusText(http.StatusForbidden))
	}
}

// UpdatePostTagsFieldFunc handles the following mutation
// ```graphql
//	mutation {
//		updatePostTags(slug: string!, tags: [string!]!) { ... }
//	}
// ```
func UpdatePostTagsFieldFunc(repository blog.PostRepository) interface{} {
	return func(ctx context.Context, args struct {
		Slug     Slug
		TagSlugs []Slug
	}) (blog.Post, error) {
		id := args.Slug.MustGetID()

		p, err := repository.FindByID(ctx, id)
		if err != nil {
			return blog.Post{}, errors.New(http.StatusText(http.StatusNotFound))
		}

		if authID := ctx.Value(AuthorizedID); authID != nil {
			if p.AuthorID == authID.(string) {
				var tags []blog.Tag
				for _, slug := range args.TagSlugs {
					tags = append(tags, blog.Tag{ID: slug.MustGetID().(primitive.ObjectID)})
				}
				return repository.Save(ctx, id, blog.NewPostQueryBuilder().WithTags(tags).Build())
			}
		}

		return blog.Post{}, errors.New(http.StatusText(http.StatusForbidden))
	}
}

// UpdatePostFeaturedImageFieldFunc handles the following mutation
// ```graphql
//	mutation {
//		updatePostFeaturedImage(slug: string!, featuredImageSlug: string!) { ... }
//	}
// ```
func UpdatePostFeaturedImageFieldFunc(repository blog.PostRepository) interface{} {
	return func(ctx context.Context, args struct {
		Slug              Slug
		FeaturedImageSlug storage.Slug `graphql:",optional"`
	}) (blog.Post, error) {
		id := args.Slug.MustGetID()

		p, err := repository.FindByID(ctx, id)
		if err != nil {
			return blog.Post{}, errors.New(http.StatusText(http.StatusNotFound))
		}

		if authID := ctx.Value(AuthorizedID); authID != nil {
			if p.AuthorID == authID.(string) {
				if args.FeaturedImageSlug == "" {
					return repository.Save(ctx, id, blog.NewPostQueryBuilder().
						WithFeaturedImage(storage.File{}).Build())
				}
				return repository.Save(ctx, id, blog.NewPostQueryBuilder().
					WithFeaturedImage(storage.File{
						ID: args.FeaturedImageSlug.MustGetID().(primitive.ObjectID),
					}).Build())
			}
		}

		return blog.Post{}, errors.New(http.StatusText(http.StatusForbidden))
	}
}

// UpdatePostAttachmentsFieldFunc handles the following mutation
// ```graphql
//	mutation {
//		updatePostAttachments(slug: string!, attachmentSlugs: [string!]!) { ... }
//	}
// ```
func UpdatePostAttachmentsFieldFunc(repository blog.PostRepository) interface{} {
	return func(ctx context.Context, args struct {
		Slug            Slug
		AttachmentSlugs []storage.Slug
	}) (blog.Post, error) {
		id := args.Slug.MustGetID()

		p, err := repository.FindByID(ctx, id)
		if err != nil {
			return blog.Post{}, errors.New(http.StatusText(http.StatusNotFound))
		}

		if authID := ctx.Value(AuthorizedID); authID != nil {
			if p.AuthorID == authID.(string) {
				var attachments []storage.File
				for _, slug := range args.AttachmentSlugs {
					attachments = append(attachments, storage.File{ID: slug.MustGetID().(primitive.ObjectID)})
				}
				return repository.Save(ctx, id, blog.NewPostQueryBuilder().WithAttachments(attachments).Build())
			}
		}

		return blog.Post{}, errors.New(http.StatusText(http.StatusForbidden))
	}
}

// FindFeaturedImageBelongedToPostFieldFunc handles the following query in the Post type
// ```graphql
//	{
//		Post {
//			...
//			featuredImage { .. }
//		}
//	}
// ```
func FindFeaturedImageBelongedToPostFieldFunc(repository storage.FileRepository) interface{} {
	return func(ctx context.Context, p blog.Post) storage.File {
		file, _ := repository.FindByID(ctx, p.FeaturedImage.ID)
		return file
	}
}

// FindAllAttachmentsBelongedToPostFieldFunc handles the following query in the Post type
// ```graphql
//	{
//		Post {
//			...
//			attachments { ... }
//		}
//	}
// ```
func FindAllAttachmentsBelongedToPostFieldFunc(repository storage.FileRepository) interface{} {
	return func(ctx context.Context, p blog.Post) ([]storage.File, error) {
		ids := make([]primitive.ObjectID, len(p.Attachments))
		for i, f := range p.Attachments {
			ids[i] = f.ID
		}
		return repository.FindAllByIDs(ctx, ids)
	}
}

// GetURLNodeShareCountFieldFunc handles the following query in the Post type
// ```graphql
//	{
//		Post {
//			...
//			engagement { ... }
//		}
//	}
// ```
func GetURLNodeShareCountFieldFunc(baseURL string, c facebook.Client) interface{} {
	return func(ctx context.Context, p blog.Post) (engagement blog.Engagement) {
		id := baseURL + "/" + p.PublishedAt.In(timeutil.TimeZoneAsiaBangkok).Format("2006/1/2") + "/" + p.Slug
		urlNode, err := c.GetURLNodeFields(id)
		if err != nil {
			logrus.Errorf("unable to retrieve URLNode on ID: %s", id)
			return
		}
		engagement.ShareCount = urlNode.Engagement.ShareCount
		return
	}
}
