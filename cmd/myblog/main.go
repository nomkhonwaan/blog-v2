package main

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
	"github.com/nomkhonwaan/myblog/pkg/graphql"
	"github.com/nomkhonwaan/myblog/pkg/graphql/playground"
	"github.com/nomkhonwaan/myblog/pkg/log"
	"github.com/nomkhonwaan/myblog/pkg/mongo"
	"github.com/nomkhonwaan/myblog/pkg/server"
	"github.com/nomkhonwaan/myblog/pkg/storage"
	"github.com/nomkhonwaan/myblog/pkg/web"
	"github.com/samsarahq/thunder/graphql/introspection"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	version, revision string
)

func init() {
	cli.VersionPrinter = func(ctx *cli.Context) {
		fmt.Println(ctx.App.Name, ctx.App.Version, revision)
	}
}

func main() {
	app := cli.NewApp()
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "listen-address",
			EnvVar: "LISTEN_ADDRESS",
			Value:  "0.0.0.0:8080",
		},
		cli.StringFlag{
			Name:   "static-files-path",
			EnvVar: "STATIC_FILES_PATH",
			Value:  "./dist/web",
		},
		cli.StringFlag{
			Name:   "mongodb-uri",
			EnvVar: "MONGODB_URI",
			Value:  "mongodb://localhost/nomkhonwaan_com",
		},
		cli.StringFlag{
			Name:   "auth0-audience",
			EnvVar: "AUTH0_AUDIENCE",
			Value:  "https://api.nomkhonwaan.com",
		},
		cli.StringFlag{
			Name:   "auth0-issuer",
			EnvVar: "AUTH0_ISSUER",
			Value:  "https://nomkhonwaan.auth0.com/",
		},
		cli.StringFlag{
			Name:   "auth0-jwks-uri",
			EnvVar: "AUTH0_JWKS_URI",
			Value:  "https://nomkhonwaan.auth0.com/.well-known/jwks.json",
		},
		cli.StringFlag{
			Name:   "amazon-s3-access-key",
			EnvVar: "AMAZON_S3_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   "amazon-s3-secret-key",
			EnvVar: "AMAZON_S3_SECRET_KEY",
		},
		cli.BoolFlag{
			Name:   "allow-cors",
			EnvVar: "ALLOW_CORS",
		},
	}
	app.Action = action

	if err := app.Run(os.Args); err != nil {
		logrus.Fatalf("myblog: %v", err)
	}
}

func action(ctx *cli.Context) error {
	/* Connect to MongoDB */
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(ctx.String("mongodb-uri")))
	if err != nil {
		return err
	}
	db := client.Database("nomkhonwaan_com")

	catRepo := blog.NewCategoryRepository(mongo.NewCustomCollection(db.Collection("categories")))
	fileRepo := storage.NewFileRepository(mongo.NewCustomCollection(db.Collection("files")))
	postRepo := blog.NewPostRepository(mongo.NewCustomCollection(db.Collection("posts")))
	tagRepo := blog.NewTagRepository(mongo.NewCustomCollection(db.Collection("tags")))

	/* Connect to Amazon S3 */
	uploader, err := storage.NewAmazonS3(ctx.String("amazon-s3-access-key"), ctx.String("amazon-s3-secret-key"), fileRepo)
	if err != nil {
		return err
	}

	/* Build GraphQL schema */
	schema := graphql.NewServer(catRepo, fileRepo, postRepo, tagRepo).Schema()
	introspection.AddIntrospectionToSchema(schema)

	/* New authentication middleware */
	jwtMiddleware := auth.NewJWTMiddleware(ctx.String("auth0-audience"), ctx.String("auth0-issuer"), ctx.String("auth0-jwks-uri"), http.DefaultTransport)

	/* Facebook's crawler middleware */
	openGraphTemplate, err := uncomressGzipData(data.MustGzipAsset("data/facebook-opengraph-template.html"))
	if err != nil {
		return err
	}

	fbCrawlerMiddleware, err := facebook.NewCrawlerMiddleware(string(openGraphTemplate), postRepo)
	if err != nil {
		return err
	}

	/* Define HTTP routes */
	r := mux.NewRouter()

	r.Handle("/v1/storage/upload", jwtMiddleware.Handler(storage.Handler(uploader)))
	r.HandleFunc("/graphiql", playground.HandlerFunc(data.MustGzipAsset("data/graphql-playground.html")))
	r.Handle("/graphql", jwtMiddleware.Handler(graphql.Handler(schema)))
	r.PathPrefix("/").Handler(fbCrawlerMiddleware.Handler(web.NewSPAHandler(ctx.String("static-files-path"))))

	/* Instantiate an HTTP server */
	handler := logRequest(r)
	if ctx.Bool("allow-cors") {
		logrus.Info("the Cross-Origin Resource Sharing (CORS) is allowed for all sites (*)")
		handler = allowCORS(handler)
	}

	s := server.InsecureServer{
		Handler:         handler,
		ShutdownTimeout: time.Minute * 5,
	}

	stopCh := handleSignals()

	err = s.ListenAndServe(ctx.String("listen-address"), stopCh)
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

func uncomressGzipData(compressed []byte) ([]byte, error) {
	rdr, err := gzip.NewReader(bytes.NewBuffer(compressed))
	if err != nil {
		return nil, err
	}
	defer rdr.Close()

	uncompressed, _ := ioutil.ReadAll(rdr)
	return uncompressed, nil
}
