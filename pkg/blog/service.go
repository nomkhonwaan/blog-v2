package blog

// Service helps co-working between data-layer and control-layer
type Service struct {
	CategoryRepository
	PostRepository
	TagRepository
}

// Category returns a category repository
func (s Service) Category() CategoryRepository {
	return s.CategoryRepository
}

// Post returns a post repository
func (s Service) Post() PostRepository {
	return s.PostRepository
}

// Tag returns a tag repository
func (s Service) Tag() TagRepository {
	return s.TagRepository
}
