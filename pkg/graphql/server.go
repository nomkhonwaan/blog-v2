package graphql

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/nomkhonwaan/myblog/pkg/auth"
	"github.com/nomkhonwaan/myblog/pkg/blog"
	"github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/graphql/schemabuilder"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strings"
)

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
	return s.service.Post().FindAll(ctx,
		blog.NewPostQueryBuilder().
			WithStatus(blog.Published).
			WithOffset(args.Offset).
			WithLimit(args.Limit).
			Build(),
	)
}

func (s *Server) makeFieldFuncPost(ctx context.Context, args struct {
	IDOrSlug string `graphql:"idOrSlug"`
}) (blog.Post, error) {
	sl := strings.Split(args.IDOrSlug, "-")

	id, err := primitive.ObjectIDFromHex(sl[len(sl)-1])
	if err != nil {
		return blog.Post{}, err
	}

	p, err := s.service.Post().FindByID(ctx, id)
	if err != nil {
		return blog.Post{}, err
	}

	// do not check authority if the post had published
	if p.Status == blog.Published {
		return p, nil
	}

	authorID, err := s.getAuthorizedUserID(ctx)
	if err != nil {
		return blog.Post{}, err
	}

	if p.AuthorID == authorID.(string) {
		return p, nil
	}

	return blog.Post{}, errors.New(http.StatusText(http.StatusForbidden))
}

func (s *Server) makeFieldFuncCreatePost(ctx context.Context) (blog.Post, error) {
	authorID, err := s.getAuthorizedUserID(ctx)
	if err != nil {
		return blog.Post{}, err
	}

	return s.service.Post().Create(ctx, authorID.(string))
}

func (s *Server) makeFieldFuncUpdatePostTitle(ctx context.Context, args struct {
	IDOrSlug string `graphql:"idOrSlug"`
	Title    string `graphql:"title"`
}) (blog.Post, error) {
	authorID, err := s.getAuthorizedUserID(ctx)
	if err != nil {
		return blog.Post{}, err
	}

	sl := strings.Split(args.IDOrSlug, "-")

	id, err := primitive.ObjectIDFromHex(sl[len(sl)-1])
	if err != nil {
		return blog.Post{}, err
	}

	p, err := s.service.Post().FindByID(ctx, id)
	if err != nil {
		return blog.Post{}, err
	}

	if p.AuthorID != authorID.(string) {
		return blog.Post{}, errors.New(http.StatusText(http.StatusForbidden))
	}

	return s.service.Post().Save(ctx, id,
		blog.NewPostQueryBuilder().
			WithTitle(args.Title).
			Build(),
	)
}

// getAuthorizedUserID returns an authorized user ID (which generated from the authentication server),
// an error unauthorized will be returned if the context is nil
func (s *Server) getAuthorizedUserID(ctx context.Context) (interface{}, error) {
	if ctx.Value(auth.UserProperty) == nil {
		return nil, errors.New(http.StatusText(http.StatusUnauthorized))
	}

	return ctx.Value(auth.UserProperty).(*jwt.Token).Claims.(jwt.MapClaims)["sub"], nil
}
