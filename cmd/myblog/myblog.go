package myblog

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nomkhonwaan/myblog/pkg/auth"
	"github.com/nomkhonwaan/myblog/pkg/aws"
	"github.com/nomkhonwaan/myblog/pkg/blog"
	"github.com/nomkhonwaan/myblog/pkg/data"
	"github.com/nomkhonwaan/myblog/pkg/facebook"
	"github.com/nomkhonwaan/myblog/pkg/gcloud"
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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	baseURL       = "https://www.nomkhonwaan.com"
	storageBucket = "www-nomkhonwaan-com"
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

	cobra.OnInitialize(initConfig)

	flags := cmd.Flags()

	flags.Bool("allow-cors", false, "")
	flags.String("listen-address", "0.0.0.0:8080", "")
	flags.String("cache-file-path", path.Join(workingDirectory, ".cache"), "")
	flags.String("web-file-path", path.Join(workingDirectory, "dist", "web"), "")
	flags.String("mongodb-uri", "mongodb://localhost/nomkhonwaan_com", "")
	flags.String("storage", "local-disk", "")
	flags.String("amazon-s3-access-key", "", "")
	flags.String("amazon-s3-secret-key", "", "")
	flags.String("gcloud-credentials-file-path", "", "")
	flags.String("auth0-audience", baseURL, "")
	flags.String("auth0-issuer", "https://nomkhonwaan.auth0.com/", "")
	flags.String("auth0-jwks-uri", "https://nomkhonwaan.auth0.com/.well-known/jwks.json", "")
	flags.String("facebook-app-access-token", "", "")

	_ = viper.BindPFlag("allow-cors", flags.Lookup("allow-cors"))
	_ = viper.BindPFlag("listen-address", flags.Lookup("listen-address"))
	_ = viper.BindPFlag("cache-file-path", flags.Lookup("cache-file-path"))
	_ = viper.BindPFlag("web-file-path", flags.Lookup("web-file-path"))
	_ = viper.BindPFlag("mongodb-uri", flags.Lookup("mongodb-uri"))
	_ = viper.BindPFlag("storage", flags.Lookup("storage"))
	_ = viper.BindPFlag("amazon-s3-access-key", flags.Lookup("amazon-s3-access-key"))
	_ = viper.BindPFlag("amazon-s3-secret-key", flags.Lookup("amazon-s3-secret-key"))
	_ = viper.BindPFlag("gcloud-credentials-file-path", flags.Lookup("gcloud-credentials-file-path"))
	_ = viper.BindPFlag("auth0-audience", flags.Lookup("auth0-audience"))
	_ = viper.BindPFlag("auth0-issuer", flags.Lookup("auth0-issuer"))
	_ = viper.BindPFlag("auth0-jwks-uri", flags.Lookup("auth0-jwks-uri"))
	_ = viper.BindPFlag("facebook-app-access-token", flags.Lookup("facebook-app-access-token"))

}

// Execute proxies to the Cobra command execution function
func Execute() error {
	return cmd.Execute()
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

	file := storage.NewFileRepository(mongo.NewCustomCollection(db.Collection("files")))
	category := blog.NewCategoryRepository(mongo.NewCustomCollection(db.Collection("categories")))
	post := blog.NewPostRepository(mongo.NewCustomCollection(db.Collection("posts")), log.NewDefaultTimer())
	tag := blog.NewTagRepository(mongo.NewCustomCollection(db.Collection("tags")))

	blogService := blog.Service{CategoryRepository: category, PostRepository: post, TagRepository: tag}

	var (
		cacheService storage.Cache
	)
	cacheService, err = storage.NewLocalDiskCache(viper.GetString("cache-file-path"))
	if err != nil {
		return err
	}

	var (
		uploader   storage.Uploader
		downloader storage.Downloader
	)
	switch viper.GetString("storage") {
	case "gcloud":
		cloudStorage := gcloud.NewCloudStorage(viper.GetString("gcloud-credentials-file-path"), storageBucket)
		uploader, downloader = cloudStorage, cloudStorage
	case "s3":
		s3, err := aws.NewS3(
			viper.GetString("amazon-s3-access-key"),
			viper.GetString("amazon-s3-secret-key"),
			storageBucket,
		)
		if err != nil {
			return err
		}
		uploader, downloader = s3, s3
	case "local-disk":
	default:
		uploader, downloader = storage.LocalDiskStorage(cacheService.(storage.LocalDiskCache)), storage.LocalDiskStorage(cacheService.(storage.LocalDiskCache))
	}

	ogTemplate, _ := unzip(data.MustGzipAsset("data/facebook-opengraph-template.html"))
	fbClient, err := facebook.NewClient(baseURL, viper.GetString("facebook-app-access-token"), string(ogTemplate), blogService, file, http.DefaultTransport)
	if err != nil {
		return err
	}

	authMiddleware := auth.NewJWTMiddleware(
		viper.GetString("auth0-audience"),
		viper.GetString("auth0-issuer"),
		viper.GetString("auth0-jwks-uri"),
		http.DefaultTransport,
	)

	gitHubHandler := github.NewHandler(cacheService, http.DefaultTransport)
	storageHandler := storage.NewHandler(cacheService, file, downloader, uploader, image.NewLanczosResizer())
	sitemapHandler := sitemap.NewHandler(baseURL, cacheService, blogService)
	schema := graphql.NewServer(blogService, fbClient, file).Schema()
	introspection.AddIntrospectionToSchema(schema)

	r := mux.NewRouter()

	if viper.GetBool("allow-cors") {
		r.Use(allowCORS)
	}
	r.Use(authMiddleware.Handler)

	r.Handle("/graphiql", playground.Handler(data.MustGzipAsset("data/graphql-playground.html")))
	r.Handle("/graphql", graphql.Handler(schema))
	gitHubHandler.Register(r.PathPrefix("/api/v2.1/github").Subrouter())
	storageHandler.Register(r.PathPrefix("/api/v2.1/storage").Subrouter())
	sitemapHandler.Register(r.PathPrefix("/sitemap.xml").Subrouter())
	r.PathPrefix("/").Handler(fbClient.CrawlerHandler(web.NewSPAHandler(viper.GetString("web-file-path"))))

	s := server.InsecureServer{Handler: r, ShutdownTimeout: time.Minute * 5}
	stopCh := handleSignals()

	err = s.ListenAndServe(viper.GetString("listen-address"), stopCh)
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
