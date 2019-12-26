package generate_indices

import (
	"context"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/nomkhonwaan/myblog/pkg/blog"
	"github.com/nomkhonwaan/myblog/pkg/mongo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
)

var (
	Command = &cobra.Command{
		Use:   "generate-indices",
		Short: "Generate all indices which will be used on Algoria",
		RunE:  action,
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	flags := Command.Flags()

	flags.String("mongodb-uri", "mongodb://localhost/nomkhonwaan_com", "")
	flags.String("algoria-client-id", "", "")
	flags.String("algoria-client-secret", "", "")

	_ = viper.BindPFlag("mongodb-uri", flags.Lookup("mongodb-uri"))
	_ = viper.BindPFlag("algoria-client-id", flags.Lookup("algoria-client-id"))
	_ = viper.BindPFlag("algoria-client-secret", flags.Lookup("algoria-client-secret"))

}

func initConfig() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
}

func action(_ *cobra.Command, _ []string) error {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(viper.GetString("mongodb-uri")))
	if err != nil {
		return err
	}
	db := client.Database("nomkhonwaan_com")

	category := blog.NewCategoryRepository(mongo.NewCustomCollection(db.Collection("categories")))
	tag := blog.NewTagRepository(mongo.NewCustomCollection(db.Collection("tags")))

	blogService := blog.Service{CategoryRepository: category, TagRepository: tag}

	index := search.
		NewClient(viper.GetString("algoria-client-id"), viper.GetString("algoria-client-secret")).
		InitIndex("prod_www-nomkhonwaan-com")

	err = clearAllObjects(index)
	if err != nil {
		return err
	}

	err = createCategoryIndices(blogService, index)
	if err != nil {
		return err
	}

	err = createTagIndices(blogService, index)
	if err != nil {
		return err
	}

	err = createPostIndices(blogService, index)
	if err != nil {
		return err
	}

	return nil
}

func clearAllObjects(index *search.Index) error {
	res, err := index.ClearObjects()
	if err != nil {
		return err
	}

	return res.Wait()
}

func createCategoryIndices(blogService blog.Service, index *search.Index) error {
	categories, err := blogService.Category().FindAll(context.Background())
	if err != nil {
		return err
	}

	res, err := index.SaveObjects(categories, opt.AutoGenerateObjectIDIfNotExist(true))
	if err != nil {
		return err
	}

	return res.Wait()
}

func createTagIndices(blogService blog.Service, index *search.Index) error {
	tags, err := blogService.Tag().FindAll(context.Background())
	if err != nil {
		return err
	}

	res, err := index.SaveObjects(tags, opt.AutoGenerateObjectIDIfNotExist(true))
	if err != nil {
		return err
	}

	return res.Wait()
}

func createPostIndices(blogService blog.Service, index *search.Index) error {
	posts, err := blogService.Post().FindAll(context.Background(), blog.NewPostQueryBuilder().
		WithStatus(blog.Published).
		WithOffset(0).
		WithLimit(99).
		Build())
	if err != nil {
		return err
	}

	res, err := index.SaveObjects(posts, opt.AutoGenerateObjectIDIfNotExist(true))
	if err != nil {
		return err
	}

	return res.Wait()
}
