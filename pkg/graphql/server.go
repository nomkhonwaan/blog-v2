package graphql

import (
	"context"
	"errors"
	"github.com/nomkhonwaan/myblog/pkg/auth"
	"github.com/nomkhonwaan/myblog/pkg/blog"
	"github.com/nomkhonwaan/myblog/pkg/facebook"
	"github.com/nomkhonwaan/myblog/pkg/storage"
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

type blogService struct{ blog.Service }

// MustGetID always return ID from the slug string
func (s Slug) MustGetID() interface{} {
	if id, err := s.GetID(); err == nil {
		return id
	}
	return primitive.NewObjectID()
}

// Service helps co-working between data-layer and control-layer
type Service interface {
	/* Facebook Client */
	FBClient() facebook.Client

	/* Storage Service */
	File() storage.FileRepository

	/* Blog Service */
	Category() blog.CategoryRepository
	Post() blog.PostRepository
	Tag() blog.TagRepository
}

type service struct {
	blogService

	fbClient facebook.Client
	fileRepo storage.FileRepository
}

func (s service) FBClient() facebook.Client {
	return s.fbClient
}

func (s service) File() storage.FileRepository {
	return s.fileRepo
}

// Server is our GraphQL server
type Server struct {
	service Service

	// GraphQL schema schema object
	schema *schemabuilder.Schema
}

// NewServer returns new GraphQL server
func NewServer(blogSvc blog.Service, fbClient facebook.Client, fileRepo storage.FileRepository) *Server {
	return &Server{
		service: service{
			blogService: blogService{blogSvc},
			fbClient:    fbClient,
			fileRepo:    fileRepo,
		},
		schema: schemabuilder.NewSchema(),
	}
}

// Schema builds the GraphQL schema
func (s *Server) Schema() *graphql.Schema {
	s.RegisterQuery(s.schema)
	s.RegisterMutation(s.schema)
	s.registerCategory(s.schema)
	s.registerTag(s.schema)
	s.registerPost(s.schema)

	return s.schema.MustBuild()
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
