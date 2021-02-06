package app

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-chi/chi"
	"github.com/nomkhonwaan/myblog/internal/blob"
	"github.com/nomkhonwaan/myblog/pkg/auth"
	"github.com/nomkhonwaan/myblog/pkg/blog"
	"github.com/nomkhonwaan/myblog/pkg/data"
	"github.com/nomkhonwaan/myblog/pkg/facebook"
	"github.com/nomkhonwaan/myblog/pkg/github"
	"github.com/nomkhonwaan/myblog/pkg/graphql"
	"github.com/nomkhonwaan/myblog/pkg/image"
	"github.com/nomkhonwaan/myblog/pkg/opengraph"
	"github.com/nomkhonwaan/myblog/pkg/server"
	"github.com/nomkhonwaan/myblog/pkg/sitemap"
	"github.com/nomkhonwaan/myblog/pkg/storage"
	"github.com/nomkhonwaan/myblog/pkg/web"
	"github.com/samsarahq/thunder/graphql/introspection"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gocloud.dev/blob/s3blob"
	_ "gocloud.dev/blob/s3blob"
)

var (
	// Cmd is a root command of "serve" for serving HTTP server with options.
	Cmd = &cobra.Command{
		Use:     "serve",
		Short:   "Listen and serve HTTP server insecurely",
		PreRunE: preRunE,
		RunE:    runE,
	}
)

func init() {
	workingDirectory, _ := os.Getwd()

	Cmd.Flags().Bool("allow-cors", false, "")
	Cmd.Flags().String("listen-address", "0.0.0.0:8080", "")
	Cmd.Flags().String("base-url", "https://www.nomkhonwaan.com", "")
	Cmd.Flags().String("cache-file-path", path.Join(workingDirectory, ".cache"), "")
	Cmd.Flags().String("static-file-path", path.Join(workingDirectory, "dist", "web"), "")
	Cmd.Flags().String("mongodb-uri", "mongodb://localhost/nomkhonwaan_com", "")
	Cmd.Flags().String("db-name", "nomkhonwaan_com", "")
	Cmd.Flags().String("storage-driver", "s3", "")
	Cmd.Flags().String("amazon-s3-region", "ap-southeast-1", "")
	Cmd.Flags().String("amazon-s3-access-key", "", "")
	Cmd.Flags().String("amazon-s3-secret-key", "", "")
	Cmd.Flags().String("amazon-s3-bucket-name", "", "")
	Cmd.Flags().String("auth0-audience", "https://www.nomkhonwaan.com", "")
	Cmd.Flags().String("auth0-issuer", "https://nomkhonwaan.auth0.com/", "")
	Cmd.Flags().String("auth0-jwks-uri", "https://nomkhonwaan.auth0.com/.well-known/jwks.json", "")
	Cmd.Flags().String("facebook-app-access-token", "", "")

	_ = viper.BindPFlag("allow-cors", Cmd.Flags().Lookup("allow-cors"))
	_ = viper.BindPFlag("listen-address", Cmd.Flags().Lookup("listen-address"))
	_ = viper.BindPFlag("base-url", Cmd.Flags().Lookup("base-url"))
	_ = viper.BindPFlag("cache-file-path", Cmd.Flags().Lookup("cache-file-path"))
	_ = viper.BindPFlag("static-file-path", Cmd.Flags().Lookup("static-file-path"))
	_ = viper.BindPFlag("mongodb-uri", Cmd.Flags().Lookup("mongodb-uri"))
	_ = viper.BindPFlag("db-name", Cmd.Flags().Lookup("db-name"))
	_ = viper.BindPFlag("storage-driver", Cmd.Flags().Lookup("storage-driver"))
	_ = viper.BindPFlag("amazon-s3-region", Cmd.Flags().Lookup("amazon-s3-region"))
	_ = viper.BindPFlag("amazon-s3-access-key", Cmd.Flags().Lookup("amazon-s3-access-key"))
	_ = viper.BindPFlag("amazon-s3-secret-key", Cmd.Flags().Lookup("amazon-s3-secret-key"))
	_ = viper.BindPFlag("amazon-s3-bucket-name", Cmd.Flags().Lookup("amazon-s3-bucket-name"))
	_ = viper.BindPFlag("auth0-audience", Cmd.Flags().Lookup("auth0-audience"))
	_ = viper.BindPFlag("auth0-issuer", Cmd.Flags().Lookup("auth0-issuer"))
	_ = viper.BindPFlag("auth0-jwks-uri", Cmd.Flags().Lookup("auth0-jwks-uri"))
	_ = viper.BindPFlag("facebook-app-access-token", Cmd.Flags().Lookup("facebook-app-access-token"))
}

func preRunE(cmd *cobra.Command, _ []string) error {
	switch viper.GetString("storage-driver") {
	case "s3":
		_ = cmd.MarkFlagRequired("amazon-s3-region")
		_ = cmd.MarkFlagRequired("amazon-s3-access-key")
		_ = cmd.MarkFlagRequired("amazon-s3-secret-key")
		_ = cmd.MarkFlagRequired("amazon-s3-bucket-name")
	}

	return nil
}

func runE(_ *cobra.Command, _ []string) error {
	var (
		baseURL = viper.GetString("base-url")
	)
	dbName := viper.GetString("db-name")
	db, err := provideMongoDB(viper.GetString("mongodb-uri"), dbName)
	if err != nil {
		return err
	}

	var (
		fileRepository     = storage.NewFileRepository(db)
		categoryRepository = blog.NewCategoryRepository(db)
		postRepository     = blog.NewPostRepository(db)
		tagRepository      = blog.NewTagRepository(db)
	)

	cache, err := storage.NewDiskCache(afero.NewOsFs(), viper.GetString("cache-file-path"))
	if err != nil {
		return err
	}
	defer cache.Close()

	bucket, err := newBlobStorage()
	if err != nil {
		return err
	}
	defer bucket.Close()

	ogTmplData, _ := unzip(data.MustGzipAsset("data/opengraph-template.html"))
	ogTmpl := template.Must(template.New("data/opengraph-template.html").Parse(string(ogTmplData)))

	schema, err := graphql.BuildSchema(
		graphql.BuildCategorySchema(categoryRepository),
		graphql.BuildTagSchema(tagRepository),
		graphql.BuildPostSchema(postRepository),
		graphql.BuildFileSchema(fileRepository),
		graphql.BuildGraphAPISchema(baseURL, facebook.NewClient(
			viper.GetString("facebook-app-access-token"), http.DefaultTransport)),
	)
	if err != nil {
		return err
	}
	introspection.AddIntrospectionToSchema(schema)

	r := chi.NewRouter()

	if viper.GetBool("allow-cors") {
		r.Use(allowCORS)
	}
	r.Use(auth.NewJWTMiddleware(
		viper.GetString("auth0-audience"),
		viper.GetString("auth0-issuer"),
		viper.GetString("auth0-jwks-uri"),
		http.DefaultTransport,
	).Handler)

	r.Route("/api/v2.1", func(r chi.Router) {
		r.Route("/github", func(r chi.Router) {
			r.Get("/gist", github.GetGistHandlerFunc(cache, http.DefaultTransport))
		})
		r.Route("/storage", func(r chi.Router) {
			r.Get("/{slug}", storage.DownloadHandlerFunc(bucket, cache, image.NewLanczosResizer(), fileRepository))
			r.Delete("/{slug}/delete", storage.DeleteHandlerFunc(bucket, fileRepository))
			r.Post("/upload", storage.UploadHandlerFunc(bucket, fileRepository))
		})
	})
	r.With(opengraph.ServeStaticSinglePageMiddleware(baseURL, ogTmpl, postRepository, fileRepository)).
		Get("/*", web.ServeStaticHandlerFunc(viper.GetString("static-file-path")))
	r.Get("/graphiql", graphql.ServeGraphiqlHandlerFunc(data.MustGzipAsset("data/graphql-playground.html")))
	r.Handle("/graphql", graphql.Handler(schema, graphql.VerifyAuthorityMiddleware))
	r.Get("/sitemap.xml", sitemap.ServeSiteMapHandlerFunc(cache,
		sitemap.GenerateFixedURLs(baseURL),
		sitemap.GeneratePostURLs(baseURL, postRepository),
		sitemap.GenerateCategoryURLs(baseURL, categoryRepository),
		sitemap.GenerateTagURLs(baseURL, tagRepository),
	))

	s := server.InsecureServer{
		Handler:         r,
		ShutdownTimeout: time.Minute * 5,
	}
	stopCh := handleSignals()

	err = s.ListenAndServe(viper.GetString("listen-address"), stopCh)
	if err != nil {
		return err
	}

	<-stopCh

	return nil
}

func newBlobStorage() (*blob.Bucket, error) {
	switch viper.GetString("storage-driver") {
	case "s3":
		sess, err := session.NewSessionWithOptions(session.Options{
			Config: aws.Config{
				Credentials: credentials.NewStaticCredentials(
					viper.GetString("amazon-s3-access-key"),
					viper.GetString("amazon-s3-secret-key"), ""),
				Region: aws.String(viper.GetString("amazon-s3-region")),
			},
		})
		if err != nil {
			return nil, err
		}

		bucket, err := s3blob.OpenBucket(context.Background(), sess, viper.GetString("amazon-s3-bucket-name"), nil)
		return &blob.Bucket{Bucket: bucket}, err
	default:
		return nil, errors.New("unsupported storage")
	}
}

func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", strings.Join([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}, ","))
		w.Header().Set("Access-Control-Allow-Headers", strings.Join([]string{"Accept", "Accept-Encoding", "Accept-Language", "Authorization", "Content-Length", "Content-Type"}, ","))

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func handleSignals() <-chan struct{} {
	stopCh := make(chan struct{})
	sigsCh := make(chan os.Signal, 2)

	signal.Notify(sigsCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigsCh
		close(stopCh)

		<-sigsCh
		os.Exit(1)
	}()

	return stopCh
}

func unzip(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return ioutil.ReadAll(r)
}
