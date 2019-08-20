package blog

// Service helps co-working between data-layer and control-layer
type Service interface {
	// Allow accessing to the category repository
	Category() CategoryRepository

	// Allow accessing to the post repository
	Post() PostRepository
}

// NewService returns blog service which embeds the following repositories
//
// CategoryRepository - For CRUD with "categories" entity or collection
// PostRepository     - For CRUD with "posts" entity or collection
func NewService(
	c CategoryRepository,
	p PostRepository,
) Service {
	return service{
		CategoryRepository: c,
		PostRepository:     p,
	}
}

type service struct {
	CategoryRepository
	PostRepository
}

func (s service) Category() CategoryRepository {
	return s.CategoryRepository
}

func (s service) Post() PostRepository {
	return s.PostRepository
}
