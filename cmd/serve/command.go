package serve

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-chi/chi"
	"github.com/nomkhonwaan/myblog/internal/blob"
	"github.com/nomkhonwaan/myblog/pkg/auth"
	"github.com/nomkhonwaan/myblog/pkg/image"
	"github.com/nomkhonwaan/myblog/pkg/mongo"
	"github.com/nomkhonwaan/myblog/pkg/server"
	"github.com/nomkhonwaan/myblog/pkg/storage"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gocloud.dev/blob/s3blob"
	_ "gocloud.dev/blob/s3blob"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"
)

const (
	baseURL = "https://www.nomkhonwaan.com"
)

var (
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
	Cmd.Flags().String("cache-file-path", path.Join(workingDirectory, ".cache"), "")
	Cmd.Flags().String("web-file-path", path.Join(workingDirectory, "dist", "web"), "")
	Cmd.Flags().String("mongodb-uri", "mongodb://localhost/nomkhonwaan_com", "")
	Cmd.Flags().String("storage-driver", "s3", "")
	Cmd.Flags().String("amazon-s3-region", "ap-southeast-1", "")
	Cmd.Flags().String("amazon-s3-access-key", "", "")
	Cmd.Flags().String("amazon-s3-secret-key", "", "")
	Cmd.Flags().String("amazon-s3-bucket-name", "", "")
	Cmd.Flags().String("auth0-audience", baseURL, "")
	Cmd.Flags().String("auth0-issuer", "https://nomkhonwaan.auth0.com/", "")
	Cmd.Flags().String("auth0-jwks-uri", "https://nomkhonwaan.auth0.com/.well-known/jwks.json", "")
	Cmd.Flags().String("facebook-app-access-token", "", "")

	_ = viper.BindPFlag("allow-cors", Cmd.Flags().Lookup("allow-cors"))
	_ = viper.BindPFlag("listen-address", Cmd.Flags().Lookup("listen-address"))
	_ = viper.BindPFlag("cache-file-path", Cmd.Flags().Lookup("cache-file-path"))
	_ = viper.BindPFlag("web-file-path", Cmd.Flags().Lookup("web-file-path"))
	_ = viper.BindPFlag("mongodb-uri", Cmd.Flags().Lookup("mongodb-uri"))
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
	db, err := newMongoDB(viper.GetString("mongodb-uri"), "nomkhonwaan_com")
	if err != nil {
		return err
	}

	var (
		fileRepository = storage.NewFileRepository(db)
		//categoryRepository = blog.NewCategoryRepository(mongo.NewCollection(db.Collection("categories")))
		//postRepository     = blog.NewPostRepository(db, log.NewDefaultTimer())
		//tagRepository      = blog.NewTagRepository(mongo.NewCollection(db.Collection("tags")))
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

	//ogTemplate, _ := unzip(data.MustGzipAsset("data/facebook-opengraph-template.html"))
	//fbClient, err := facebook.NewClient(baseURL, viper.GetString("facebook-app-access-token"), string(ogTemplate), blogService, fileRepository, http.DefaultTransport)
	//if err != nil {
	//	return err
	//}

	authMiddleware := auth.NewJWTMiddleware(
		viper.GetString("auth0-audience"),
		viper.GetString("auth0-issuer"),
		viper.GetString("auth0-jwks-uri"),
		http.DefaultTransport,
	)

	//schema := graphql.NewServer(blogService, fbClient, fileRepository).Schema()
	//introspection.AddIntrospectionToSchema(schema)

	r := chi.NewRouter()

	if viper.GetBool("allow-cors") {
		r.Use(allowCORS)
	}
	r.Use(authMiddleware.Handler)

	r.Route("/api/v2.1/storage", func(r chi.Router) {
		r.Get("/{slug}", storage.DownloadHandlerFunc(bucket, cache, image.NewLanczosResizer(), fileRepository))
		r.Delete("/delete/{slug}", storage.DeleteHandlerFunc(bucket, fileRepository))
		r.Post("/upload", storage.UploadHandlerFunc(bucket, fileRepository))
	})

	//r.Handle("/graphiql", playground.Handler(data.MustGzipAsset("data/graphql-playground.html")))
	//r.Handle("/graphql", graphql.Handler(schema))
	//github.NewHandler(cache, http.DefaultTransport).
	//	Register(r.PathPrefix("/api/v2.1/github").Subrouter())
	//r.PathPrefix("/api/v2.1/storage").Subrouter()

	//sitemap.NewHandler(baseURL, cache, blogService).
	//	Register(r.PathPrefix("/sitemap.xml").Subrouter())
	//r.PathPrefix("/").Handler(fbClient.CrawlerHandler(web.NewSPAHandler(viper.GetString("web-`file-path"))))

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

func newMongoDB(uri, dbName string) (mongo.Database, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return client.Database(dbName), nil
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

func unzip(compressed []byte) ([]byte, error) {
	rdr, err := gzip.NewReader(bytes.NewBuffer(compressed))
	if err != nil {
		return nil, err
	}
	defer rdr.Close()

	uncompressed, _ := ioutil.ReadAll(rdr)
	return uncompressed, nil
}
