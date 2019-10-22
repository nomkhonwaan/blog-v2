package graphql

import (
	"context"
	"errors"
	"github.com/nomkhonwaan/myblog/pkg/auth"
	"github.com/nomkhonwaan/myblog/pkg/blog"
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

// Server is our GraphQL server
type Server struct {
	service blog.Service

	// GraphQL schema schema object
	schema *schemabuilder.Schema
}

// NewServer returns new GraphQL server
func NewServer(service blog.Service) *Server {
	return &Server{
		service: service,
		schema:  schemabuilder.NewSchema(),
	}
}

// Schema builds the GraphQL schema
func (s *Server) Schema() *graphql.Schema {
	s.registerQuery(s.schema)
	s.registerMutation(s.schema)
	s.registerPost(s.schema)

	return s.schema.MustBuild()
}

func (s *Server) registerQuery(schema *schemabuilder.Schema) {
	obj := schema.Query()

	obj.FieldFunc("categories", s.makeFieldFuncCategories)
	obj.FieldFunc("latestPublishedPosts", s.makeFieldFuncLatestPublishedPosts)
	obj.FieldFunc("post", s.makeFieldFuncPost)
}

func (s *Server) registerMutation(schema *schemabuilder.Schema) {
	obj := schema.Mutation()

	obj.FieldFunc("createPost", s.makeFieldFuncCreatePost)
	obj.FieldFunc("updatePostTitle", s.makeFieldFuncUpdatePostTitle)
	obj.FieldFunc("updatePostContent", s.makeFieldFuncUpdatePostContent)
}

func (s *Server) registerPost(schema *schemabuilder.Schema) {
	obj := schema.Object("Post", blog.Post{})

	obj.FieldFunc("categories", (blog.Post{}).BelongToCategories(s.service.Category()))
	obj.FieldFunc("tags", (blog.Post{}).BelongToTags(s.service.Tag()))
}

func (s *Server) makeFieldFuncCategories(ctx context.Context) ([]blog.Category, error) {
	return s.service.Category().FindAll(ctx)
}

func (s *Server) makeFieldFuncLatestPublishedPosts(ctx context.Context, args struct{ Offset, Limit int64 }) ([]blog.Post, error) {
	return s.service.Post().FindAll(ctx, blog.NewPostQueryBuilder().WithStatus(blog.Published).WithOffset(args.Offset).WithLimit(args.Limit).Build())
}

func (s *Server) makeFieldFuncPost(ctx context.Context, args struct {
	Slug Slug `graphql:"slug"`
}) (blog.Post, error) {
	id := args.Slug.MustGetID()

	p, err := s.service.Post().FindByID(ctx, id)
	if err != nil {
		return blog.Post{}, err
	}

	// do not check authority if the post had published
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

func (s *Server) makeFieldFuncCreatePost(ctx context.Context) (blog.Post, error) {
	authorizedID, err := s.getAuthorizedIDOrFailed(ctx)
	if err != nil {
		return blog.Post{}, err
	}

	return s.service.Post().Create(ctx, authorizedID.(string))
}

func (s *Server) makeFieldFuncUpdatePostTitle(ctx context.Context, args struct {
	Slug  Slug   `graphql:"slug"`
	Title string `graphql:"title"`
}) (blog.Post, error) {
	id := args.Slug.MustGetID()

	err := s.validateAuthority(ctx, id)
	if err != nil {
		return blog.Post{}, err
	}

	return s.service.Post().Save(ctx, id, blog.NewPostQueryBuilder().WithTitle(args.Title).Build())
}

func (s *Server) makeFieldFuncUpdatePostContent(ctx context.Context, args struct {
	Slug     Slug   `graphql:"slug"`
	Markdown string `graphql:"markdown"`
}) (blog.Post, error) {
	id := args.Slug.MustGetID()

	err := s.validateAuthority(ctx, id)
	if err != nil {
		return blog.Post{}, err
	}

	return blog.Post{}, nil
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
