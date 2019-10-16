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

	if ctx.Value(auth.UserProperty) == nil {
		return blog.Post{}, errors.New(http.StatusText(http.StatusUnauthorized))
	}

	// only author can get their own post
	authorID := ctx.Value(auth.UserProperty).(*jwt.Token).Claims.(jwt.MapClaims)["sub"]
	if p.AuthorID == authorID {
		return p, nil
	}

	return blog.Post{}, errors.New(http.StatusText(http.StatusForbidden))
}

func (s *Server) makeFieldFuncCreatePost(ctx context.Context) (blog.Post, error) {
	if ctx.Value(auth.UserProperty) == nil {
		return blog.Post{}, errors.New(http.StatusText(http.StatusUnauthorized))
	}

	authorID := ctx.Value(auth.UserProperty).(*jwt.Token).Claims.(jwt.MapClaims)["sub"]
	return s.service.Post().Create(ctx, authorID.(string))
}
