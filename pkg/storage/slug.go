package storage

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"path/filepath"
	"strings"
)

// Slug is a valid URL string composes with file name and ID
type Slug string

// GetID returns an ID from the slug string
func (s Slug) GetID() (interface{}, error) {
	sl := strings.Split(string(s), "-")
	fileName := sl[len(sl)-1]
	return primitive.ObjectIDFromHex(fileName[0 : len(fileName)-len(filepath.Ext(fileName))])
}

// MustGetID always return ID from the slug string
func (s Slug) MustGetID() interface{} {
	if id, err := s.GetID(); err == nil {
		return id
	}
	return primitive.NewObjectID()
}
