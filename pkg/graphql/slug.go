package graphql

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
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
