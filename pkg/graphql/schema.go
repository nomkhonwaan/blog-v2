package graphql

import (
	"context"
	"github.com/nomkhonwaan/myblog/pkg/blog"
	"github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/graphql/schemabuilder"
)

func BuildSchema(categoryRepository blog.CategoryRepository) (*graphql.Schema, error) {
	s := schemabuilder.NewSchema()
	q := s.Query()
	q.FieldFunc("category", FindCategoryBySlugFieldFunc(categoryRepository))
	q.FieldFunc("categories", FindAllCategoriesFieldFunc(categoryRepository))
	return s.Build()
}

// FindCategoryBySlugFieldFunc handles the following query
// ```graphql
// 	{
//		category(slug: string!) { ... }
// 	}
// ```
func FindCategoryBySlugFieldFunc(repository blog.CategoryRepository) interface{} {
	return func(ctx context.Context, args struct{ Slug Slug }) (blog.Category, error) {
		return repository.FindByID(ctx, args.Slug.MustGetID())
	}
}

// FindAllCategoriesFieldFunc handles the following query
// ```graphql
//	{
//		categories { ... }
//	}
// ```
func FindAllCategoriesFieldFunc(repository blog.CategoryRepository) interface{} {
	return func(ctx context.Context) ([]blog.Category, error) {
		return repository.FindAll(ctx)
	}
}
