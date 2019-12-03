package myblog

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nomkhonwaan/myblog/pkg/auth"
	"github.com/nomkhonwaan/myblog/pkg/blog"
	"github.com/nomkhonwaan/myblog/pkg/data"
	"github.com/nomkhonwaan/myblog/pkg/facebook"
	"github.com/nomkhonwaan/myblog/pkg/github"
	"github.com/nomkhonwaan/myblog/pkg/graphql"
	"github.com/nomkhonwaan/myblog/pkg/graphql/playground"
	"github.com/nomkhonwaan/myblog/pkg/image"
	"github.com/nomkhonwaan/myblog/pkg/log"
	"github.com/nomkhonwaan/myblog/pkg/mongo"
	"github.com/nomkhonwaan/myblog/pkg/server"
	"github.com/nomkhonwaan/myblog/pkg/sitemap"
	"github.com/nomkhonwaan/myblog/pkg/storage"
	"github.com/nomkhonwaan/myblog/pkg/web"
	"github.com/samsarahq/thunder/graphql/introspection"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	// Version refers to the latest Git tag
	Version string

	// Revision refers to the latest Git commit hash
	Revision string

	cmd = &cobra.Command{
		Use:     "myblog",
		Short:   "Personal blog website written in Go with Angular 2+",
		Version: fmt.Sprintf("%s %s\n", Version, Revision),
		RunE:    action,
	}
)

func init() {
	workingDirectory, _ := os.Getwd()

	flags := cmd.Flags()
	flags.Bool("allow-cors", false, "")
	flags.String("listen-address", "0.0.0.0:8080", "")
	flags.String("cache-files-path", path.Join(workingDirectory, ".cache"), "")
	flags.String("static-files-path", path.Join(workingDirectory, "dist", "web"), "")
	flags.String("mongodb-uri", "", "")
	flags.String("amazon-s3-access-key", "", "")
	flags.String("amazon-s3-secret-key", "", "")
	flags.String("auth0-audience", baseURL, "")
	flags.String("auth0-issuer", "https://nomkhonwaan.auth0.com/", "")
	flags.String("auth0-jwks-uri", "https://nomkhonwaan.auth0.com/.well-known/jwks.json", "")
	flags.String("facebook-app-access-token", "", "")
}

// Execute proxies to the Cobra command execution function
func Execute() error {
	return cmd.Execute()
}

func action(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()

	// Create a connection to MongoDB server without time limitation
	uri, _ := flags.GetString("mongodb-uri")
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	db := client.Database("nomkhonwaan_com")

	// Create all MongoDB repositories
	file := storage.NewFileRepository(mongo.NewCustomCollection(db.Collection("files")))
	category := blog.NewCategoryRepository(mongo.NewCustomCollection(db.Collection("categories")))
	post := blog.NewPostRepository(mongo.NewCustomCollection(db.Collection("posts")), log.NewDefaultTimer())
	tag := blog.NewTagRepository(mongo.NewCustomCollection(db.Collection("tags")))

	// Create all services; as well as repositories
	blogService := blog.Service{CategoryRepository: category, PostRepository: post, TagRepository: tag}
	cacheFilesPath, _ := flags.GetString("cache-files-path")
	cacheService, err := storage.NewDiskCache(cacheFilesPath)
	if err != nil {
		return err
	}

	// Create new Amazon S3 client which provides uploader and downloader functions
	accessKey, _ := flags.GetString("amazon-s3-access-key")
	secretKey, _ := flags.GetString("amazon-s3-secret-key")
	s3, err := storage.NewCustomizedAmazonS3Client(accessKey, secretKey)
	if err != nil {
		return err
	}

	// Create new Facebook client which provides crawler bot handling and Graph API client
	appAccessToken, _ := flags.GetString("facebook-app-access-token")
	ogTemplate, _ := unzip(data.MustGzipAsset("data/facebook-opengraph-template.html"))
	fbClient, err := facebook.NewClient(baseURL, appAccessToken, string(ogTemplate), blogService, file, http.DefaultTransport)
	if err != nil {
		return err
	}

	// Create Auth0 JWT middleware for checking an authorization header
	audience, _ := flags.GetString("auth0-audience")
	issuer, _ := flags.GetString("auth0-issuer")
	jwksURI, _ := flags.GetString("auth0-jwks-uri")
	authMiddleware := auth.NewJWTMiddleware(audience, issuer, jwksURI, http.DefaultTransport)

	// Create all HTTP handlers
	ghHandler := github.NewHandler(cacheService, http.DefaultTransport)
	storageHandler := storage.NewHandler(cacheService, file, s3, s3, image.NewLanczosResizer())
	sitemapHandler := sitemap.NewHandler(baseURL, cacheService, blogService)
	schema := graphql.NewServer(blogService, fbClient, file).Schema()
	introspection.AddIntrospectionToSchema(schema)

	// Register all routes with Gorilla
	r := mux.NewRouter()

	if yes, _ := flags.GetBool("allow-cors"); yes {
		r.Use(allowCORS)
	}
	r.Use(logRequest)
	r.Use(authMiddleware.Handler)

	r.Handle("/graphiql", playground.Handler(data.MustGzipAsset("data/graphql-playground.html")))
	r.Handle("/graphql", graphql.Handler(schema))

	ghHandler.Register(r.PathPrefix("/api/v2.1/github").Subrouter())
	storageHandler.Register(r.PathPrefix("/api/v2.1/storage").Subrouter())
	sitemapHandler.Register(r.PathPrefix("/sitemap.xml").Subrouter())

	staticFilesPath, _ := flags.GetString("static-files-path")
	r.PathPrefix("/").Handler(fbClient.CrawlerHandler(web.NewSPAHandler(staticFilesPath)))

	s := server.InsecureServer{Handler: r, ShutdownTimeout: time.Minute * 5}
	stopCh := handleSignals()

	listenAddress, _ := flags.GetString("listen-address")
	err = s.ListenAndServe(listenAddress, stopCh)
	if err != nil {
		return err
	}

	<-stopCh

	return nil
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

func logRequest(h http.Handler) http.Handler {
	return log.NewLoggingInterceptor(log.NewDefaultTimer(), logrus.New()).Handler(h)
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
