package graphql

import (
	"context"
	"github.com/nomkhonwaan/myblog/pkg/blog"
	"github.com/nomkhonwaan/myblog/pkg/storage"
	"github.com/samsarahq/thunder/graphql/schemabuilder"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Server) registerPost(schema *schemabuilder.Schema) {
	obj := schema.Object("Post", blog.Post{})

	obj.FieldFunc("categories", s.postCategoriesFieldFunc)
	obj.FieldFunc("tags", s.postTagsFieldFunc)
	obj.FieldFunc("featuredImage", s.postFeaturedImageFieldFunc)
	obj.FieldFunc("attachments", s.postAttachmentsFieldFunc)
	//obj.FieldFunc("engagement", s.postEngagementFieldFunc)
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

func (s *Server) postFeaturedImageFieldFunc(ctx context.Context, p blog.Post) storage.File {
	file, _ := s.service.File().FindByID(ctx, p.FeaturedImage.ID)
	return file
}

func (s *Server) postAttachmentsFieldFunc(ctx context.Context, p blog.Post) ([]storage.File, error) {
	ids := make([]primitive.ObjectID, len(p.Attachments))

	for i, atm := range p.Attachments {
		ids[i] = atm.ID
	}

	return s.service.File().FindAllByIDs(ctx, ids)
}

//func (s *Server) postEngagementFieldFunc(ctx context.Context, p blog.Post) blog.Engagement {
//	engagement := blog.Engagement{}
//
//	// Get engagement data from Facebook Graph API
//	id := "/" + p.PublishedAt.In(timeutil.TimeZoneAsiaBangkok).Format("2006/1/2") + "/" + p.Slug
//	url, err := s.service.FBClient().GetURL(id)
//	if err != nil {
//		logrus.Errorf("an error has occurred while getting URLNode from Facebook Graph API: %s", err)
//	}
//	engagement.ShareCount += url.Engagement.ShareCount
//
//	// Get engagement data from Twitter Search API
//	// TODO: the Twitter client is not implement yet
//
//	return engagement
//}
