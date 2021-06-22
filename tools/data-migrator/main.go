package main

import (
	"context"
	"flag"
	"github.com/nomkhonwaan/myblog/pkg/blog"
	"github.com/nomkhonwaan/myblog/pkg/mongo"
	"github.com/nomkhonwaan/myblog/pkg/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strings"
	"time"
)

var (
	uri string
)

func init() {
	flag.StringVar(&uri, "mongodb-uri", "", "")
	flag.Parse()
}

func main() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	fatal(err)

	source := client.Database("nomkhonwaan_com")
	destination := client.Database("beta_nomkhonwaan_com")

	fatal(migrateCategoriesCollection(source, destination))
	fatal(migrateTagsCollection(source, destination))
	fatal(migratePostsCollection(source, destination))
}

func migrateCategoriesCollection(source, destination mongo.Database) error {
	repo := blog.NewCategoryRepository(source)
	col := destination.Collection("taxonomies")

	if err := col.Drop(context.Background()); err != nil {
		return err
	}

	cats, err := repo.FindAll(context.Background())

	for _, c := range cats {
		doc := bson.M{
			"_id":  c.ID,
			"name": c.Name,
			"slug": c.Slug,
			"type": 0,
		}

		_, err = col.InsertOne(context.Background(), doc)
		if err != nil {
			return err
		}
	}

	return nil
}

func migrateTagsCollection(source, destination mongo.Database) error {
	repo := blog.NewTagRepository(source)
	col := destination.Collection("taxonomies")

	cats, err := repo.FindAll(context.Background())

	for _, c := range cats {
		doc := bson.M{
			"_id":  c.ID,
			"name": c.Name,
			"slug": c.Slug,
			"type": 1,
		}

		_, err = col.InsertOne(context.Background(), doc)
		if err != nil {
			return err
		}
	}

	return nil
}

func migrateFilesCollection(source, destination mongo.Database, ids interface{}) error {
	repo := storage.NewFileRepository(source)
	col := destination.Collection("files")

	files, err := repo.FindAllByIDs(context.Background(), ids)
	if err != nil {
		return err
	}

	for _, f := range files {
		r := col.FindOne(context.Background(), bson.M{"_id": f.ID})

		var h storage.File
		err := r.Decode(&h)
		if err == nil {
			continue
		}

		doc := bson.M{
			"_id":              f.ID,
			"fileName":         f.FileName,
			"slug":             f.Slug,
			"uploadedFilePath": f.Path,
			"mimeType":         getMimeTypeFromFileName(f.FileName),
			"provider":         "s3",
			"region":           "ap-southeast-1",
			"bucket":           "nomkhonwaan-com",
			"uploadedAt":       f.CreatedAt,
		}

		_, err = col.InsertOne(context.Background(), doc)
		if err != nil {
			return err
		}
	}

	return nil
}

func migratePostsCollection(source, destination mongo.Database) error {
	repo := blog.NewPostRepository(source)
	col := destination.Collection("posts")

	if err := col.Drop(context.Background()); err != nil {
		return err
	}
	if err := destination.Collection("users").Drop(context.Background()); err != nil {
		return err
	}

	if err := destination.Collection("files").Drop(context.Background()); err != nil {
		return err
	}

	var (
		posts []blog.Post
		err   error
	)
	for i := 0; i == 0 || len(posts) > 0; i++ {
		for _, p := range posts {
			var status int
			if p.Status == blog.StatusPublished {
				status = 1
			}

			authorID, err := migrateUsersCollection(p.AuthorID, destination)
			if err != nil {
				return err
			}

			doc := bson.M{
				"_id":         p.ID,
				"title":       p.Title,
				"slug":        p.Slug,
				"status":      status,
				"markdown":    p.Markdown,
				"html":        p.HTML,
				"publishedAt": p.PublishedAt,
				"author":      authorID,
			}

			cats := make([]primitive.ObjectID, 0)
			for _, c := range p.Categories {
				cats = append(cats, c.ID)
			}
			if len(cats) > 0 {
				doc["categories"] = cats
			}

			tags := make([]primitive.ObjectID, 0)
			for _, t := range p.Tags {
				tags = append(tags, t.ID)
			}
			if len(tags) > 0 {
				doc["tags"] = tags
			}

			err = migrateFilesCollection(source, destination, []primitive.ObjectID{p.FeaturedImage.ID})
			if err != nil {
				return err
			}
			doc["featuredImage"] = p.FeaturedImage.ID

			attachments := make([]primitive.ObjectID, 0)
			for _, a := range p.Attachments {
				attachments = append(attachments, a.ID)
			}
			err = migrateFilesCollection(source, destination, attachments)
			if err != nil {
				return err
			}
			doc["attachments"] = attachments
			doc["createdAt"] = p.CreatedAt
			doc["updatedAt"] = p.UpdatedAt

			_, err = col.InsertOne(context.Background(), doc)
			if err != nil {
				return err
			}
		}

		posts, err = repo.FindAll(context.Background(), blog.NewPostQueryBuilder().
			WithOffset(int64(i*5)).
			WithLimit(5).
			Build())
		if err != nil {
			return err
		}
	}

	return nil
}

func migrateUsersCollection(user string, destination mongo.Database) (primitive.ObjectID, error) {
	col := destination.Collection("users")

	r := col.FindOne(context.Background(), bson.M{"user": user})

	var u User
	err := r.Decode(&u)
	if err != nil {
		if err.Error() != "mongo: no documents in result" {
			return primitive.ObjectID{}, err
		}

		u.ID = primitive.NewObjectID()
		u.User = user
		u.DisplayName = "Natcha Luangaroonchai"
		u.ProfilePicture = "https://avatars.githubusercontent.com/u/303589?v=4"
		u.CreatedAt = time.Date(2017, 11, 7, 2, 30, 31, 864000, time.UTC)
		u.UpdatedAt = time.Date(2021, 3, 13, 2, 11, 35, 864000, time.UTC)

		doc, _ := bson.Marshal(u)
		_, err = col.InsertOne(context.Background(), doc)
		if err != nil {
			return primitive.ObjectID{}, err
		}

		return u.ID, nil
	}

	return u.ID, nil
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type User struct {
	ID             primitive.ObjectID `bson:"_id" json:"id"`
	User           string             `bson:"user" json:"user"`
	DisplayName    string             `bson:"displayName" json:"displayName"`
	ProfilePicture string             `bson:"profilePicture" json:"profilePicture`
	CreatedAt      time.Time          `bson:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt"`
}

func getMimeTypeFromFileName(fileName string) string {
	if strings.Contains(strings.ToLower(fileName), ".jpeg") {
		return "image/jpeg"
	}
	if strings.Contains(strings.ToLower(fileName), ".jpg") {
		return "image/jpg"
	}
	if strings.Contains(strings.ToLower(fileName), ".png") {
		return "image/png"
	}
	return ""
}
