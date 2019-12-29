//go:generate mockgen -destination=./mock/server_mock.go github.com/nomkhonwaan/myblog/pkg/graphql Service

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

// MustGetID always return ID from the slug string
func (s Slug) MustGetID() interface{} {
	if id, err := s.GetID(); err == nil {
		return id
	}
	return primitive.NewObjectID()
}

// Service helps co-working between data-layer and control-layer
type Service interface {
	// The Facebook client
	FBClient() facebook.Client
	// A file repository
	File() storage.FileRepository
	// A category repository
	Category() blog.CategoryRepository
	// A post repository
	Post() blog.PostRepository
	// A tag repository
	Tag() blog.TagRepository
}

type service struct {
	blog.Service

	fbClient facebook.Client
	file     storage.FileRepository
}

func (s service) FBClient() facebook.Client {
	return s.fbClient
}

func (s service) File() storage.FileRepository {
	return s.file
}

// Server is our GraphQL server
type Server struct {
	service Service

	// GraphQL schema schema object
	schema *schemabuilder.Schema
}

// NewServer returns new GraphQL server
func NewServer(blogService blog.Service, fbClient facebook.Client, file storage.FileRepository) *Server {
	return &Server{
		service: service{
			Service:  blogService,
			fbClient: fbClient,
			file:     file,
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
func (s *Server) getAuthorizedID(ctx context.Context) (interface{}, error) {
	authorizedID := auth.GetAuthorizedUserID(ctx)
	if authorizedID == nil {
		return nil, errors.New(http.StatusText(http.StatusUnauthorized))
	}

	return authorizedID, nil
}

// validateAuthority performs validation against the authorized ID and post's author ID
func (s *Server) validateAuthority(ctx context.Context, id interface{}) error {
	authorizedID, err := s.getAuthorizedID(ctx)
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
