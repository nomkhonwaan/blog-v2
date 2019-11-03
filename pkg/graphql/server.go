package graphql

import (
	"context"
	"errors"
	"github.com/nomkhonwaan/myblog/pkg/auth"
	"github.com/nomkhonwaan/myblog/pkg/blog"
	"github.com/nomkhonwaan/myblog/pkg/storage"
	"github.com/russross/blackfriday/v2"
	"github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/graphql/schemabuilder"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strings"
)

// Slug is a valid URL string composes with title and ID
type Slug string

// GetID returns an ID from the slug string
func (s Slug) GetID() (interface{}, error) {
	sl := strings.Split(string(s), "-")
	return primitive.ObjectIDFromHex(sl[len(sl)-1])
}

// MustGetID always return ID from the slug string
func (s Slug) MustGetID() interface{} {
	if id, err := s.GetID(); err == nil {
		return id
	}
	return primitive.NewObjectID()
}

// Service helps co-working between data-layer and control-layer
type Service interface {
	// A Category repository
	Category() blog.CategoryRepository

	// A File repository
	File() storage.FileRepository

	// A Post repository
	Post() blog.PostRepository

	// A Tag repository
	Tag() blog.TagRepository
}

type service struct {
	catRepo  blog.CategoryRepository
	fileRepo storage.FileRepository
	postRepo blog.PostRepository
	tagRepo  blog.TagRepository
}

func (s service) Category() blog.CategoryRepository {
	return s.catRepo
}

func (s service) File() storage.FileRepository {
	return s.fileRepo
}

func (s service) Post() blog.PostRepository {
	return s.postRepo
}

func (s service) Tag() blog.TagRepository {
	return s.tagRepo
}

// Server is our GraphQL server
type Server struct {
	service Service

	// GraphQL schema schema object
	schema *schemabuilder.Schema
}

// NewServer returns new GraphQL server
func NewServer(catRepo blog.CategoryRepository, fileRepo storage.FileRepository, postRepo blog.PostRepository, tagRepo blog.TagRepository) *Server {
	return &Server{
		service: service{
			catRepo:  catRepo,
			fileRepo: fileRepo,
			postRepo: postRepo,
			tagRepo:  tagRepo,
		},
		schema: schemabuilder.NewSchema(),
	}
}

// Schema builds the GraphQL schema
func (s *Server) Schema() *graphql.Schema {
	s.registerQuery(s.schema)
	s.registerMutation(s.schema)
	s.registerCategory(s.schema)
	s.registerPost(s.schema)

	return s.schema.MustBuild()
}

func (s *Server) registerQuery(schema *schemabuilder.Schema) {
	obj := schema.Query()

	obj.FieldFunc("category", s.findCategoryBySlugQuery)
	obj.FieldFunc("categories", s.findAllCategoriesQuery)
	obj.FieldFunc("tag", s.findTagBySlugQuery)
	obj.FieldFunc("tags", s.findAllTagsQuery)
	obj.FieldFunc("latestPublishedPosts", s.findLatestPublishedPostsQuery)
	obj.FieldFunc("post", s.findPostBySlugQuery)
}

func (s *Server) registerMutation(schema *schemabuilder.Schema) {
	obj := schema.Mutation()

	obj.FieldFunc("createPost", s.createPostMutation)
	obj.FieldFunc("updatePostTitle", s.updatePostTitleMutation)
	obj.FieldFunc("updatePostContent", s.updatePostContentMutation)
	obj.FieldFunc("updatePostTags", s.updatePostTagsMutation)
}

func (s *Server) registerCategory(schema *schemabuilder.Schema) {
	obj := schema.Object("Category", blog.Category{})

	obj.FieldFunc("latestPublishedPosts", s.categoryLatestPublishedPostsFieldFunc)
}

func (s *Server) registerPost(schema *schemabuilder.Schema) {
	obj := schema.Object("Post", blog.Post{})

	obj.FieldFunc("categories", s.postCategoriesFieldFunc)
	obj.FieldFunc("tags", s.postTagsFieldFunc)
}

func (s *Server) findCategoryBySlugQuery(ctx context.Context, args struct{ Slug Slug }) (blog.Category, error) {
	id := args.Slug.MustGetID()
	return s.service.Category().FindByID(ctx, id)
}

func (s *Server) findAllCategoriesQuery(ctx context.Context) ([]blog.Category, error) {
	return s.service.Category().FindAll(ctx)
}

func (s *Server) findTagBySlugQuery(ctx context.Context, args struct{ Slug Slug }) (blog.Tag, error) {
	id := args.Slug.MustGetID()
	return s.service.Tag().FindByID(ctx, id)
}

func (s *Server) findAllTagsQuery(ctx context.Context) ([]blog.Tag, error) {
	return s.service.Tag().FindAll(ctx)
}

func (s *Server) findLatestPublishedPostsQuery(ctx context.Context, args struct{ Offset, Limit int64 }) ([]blog.Post, error) {
	return s.service.Post().FindAll(ctx, blog.NewPostQueryBuilder().WithStatus(blog.Published).WithOffset(args.Offset).WithLimit(args.Limit).Build())
}

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

func (s *Server) createPostMutation(ctx context.Context) (blog.Post, error) {
	authorizedID, err := s.getAuthorizedIDOrFailed(ctx)
	if err != nil {
		return blog.Post{}, err
	}

	return s.service.Post().Create(ctx, authorizedID.(string))
}

func (s *Server) updatePostTitleMutation(ctx context.Context, args struct {
	Slug  Slug
	Title string
}) (blog.Post, error) {
	id := args.Slug.MustGetID()

	err := s.validateAuthority(ctx, id)
	if err != nil {
		return blog.Post{}, err
	}

	return s.service.Post().Save(ctx, id, blog.NewPostQueryBuilder().WithTitle(args.Title).Build())
}

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

func (s *Server) categoryLatestPublishedPostsFieldFunc(ctx context.Context, cat blog.Category, args struct{ Offset, Limit int64 }) ([]blog.Post, error) {
	return s.service.Post().FindAll(ctx, blog.NewPostQueryBuilder().WithCategory(cat).WithOffset(args.Offset).WithLimit(args.Limit).Build())
}

func (s *Server) tagLatestPublishedPostsFieldFunc(ctx context.Context, tag blog.Tag, args struct{ Offset, Limit int64 }) ([]blog.Post, error) {
	return s.service.Post().FindAll(ctx, blog.NewPostQueryBuilder().WithTag(tag).WithOffset(args.Offset).WithLimit(args.Limit).Build())
}

func (s *Server) postCategoriesFieldFunc(ctx context.Context, p blog.Post) ([]blog.Category, error) {
	ids := make([]primitive.ObjectID, len(p.Categories))

	for i, cat := range p.Categories {
		ids[i] = cat.ID
	}

	return s.service.Category().FindAllByIDs(ctx, ids)
}

func (s *Server) postTagsFieldFunc(ctx context.Context, p blog.Post) ([]blog.Tag, error) {
	ids := make([]primitive.ObjectID, len(p.Tags))

	for i, tag := range p.Tags {
		ids[i] = tag.ID
	}

	return s.service.Tag().FindAllByIDs(ctx, ids)
}

// getAuthorizedUserID returns an authorized user ID (which generated from the authentication server),
// an error unauthorized will be returned if the context is nil
func (s *Server) getAuthorizedIDOrFailed(ctx context.Context) (interface{}, error) {
	authorizedID := auth.GetAuthorizedUserID(ctx)
	if authorizedID == nil {
		return nil, errors.New(http.StatusText(http.StatusUnauthorized))
	}

	return authorizedID, nil
}

// validateAuthority performs validation against the authorized ID and post's author ID
func (s *Server) validateAuthority(ctx context.Context, id interface{}) error {
	authorizedID, err := s.getAuthorizedIDOrFailed(ctx)
	if err != nil {
		return err
	}

	p, err := s.service.Post().FindByID(ctx, id)
	if err != nil {
		return err
	}
	if p.AuthorID != authorizedID.(string) {
		return errors.New(http.StatusText(http.StatusForbidden))
	}

	return nil
}
