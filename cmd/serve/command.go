package serve

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"github.com/go-chi/chi"
	"github.com/nomkhonwaan/myblog/pkg/auth"
	"github.com/nomkhonwaan/myblog/pkg/aws"
	"github.com/nomkhonwaan/myblog/pkg/gcloud"
	"github.com/nomkhonwaan/myblog/pkg/image"
	"github.com/nomkhonwaan/myblog/pkg/mongo"
	"github.com/nomkhonwaan/myblog/pkg/server"
	"github.com/nomkhonwaan/myblog/pkg/storage"
	"github.com/spf13/afero"
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
	baseURL = "https://www.nomkhonwaan.com"
)

var (
	Cmd = &cobra.Command{
		Use:   "serve",
		Short: "Listen and serve HTTP server insecurely",
		RunE:  runE,
	}
)

func init() {
	wd, _ := os.Getwd()

	cobra.OnInitialize(initConfig)

	flags := Cmd.Flags()

	flags.Bool("allow-cors", false, "")
	flags.String("listen-address", "0.0.0.0:8080", "")
	flags.String("cache-file-path", path.Join(wd, ".cache"), "")
	flags.String("web-file-path", path.Join(wd, "dist", "web"), "")
	flags.String("mongodb-uri", "mongodb://localhost/nomkhonwaan_com", "")
	flags.String("storage", "local-disk", "")
	flags.String("amazon-s3-access-key", "", "")
	flags.String("amazon-s3-secret-key", "", "")
	flags.String("gcloud-credentials-file-path", os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"), "")
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

func initConfig() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
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

	stg, err := initStorage()
	if err != nil {
		return err
	}

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
		r.Get("/{slug}", storage.DownloadHandlerFunc(stg, cache, image.NewLanczosResizer(), fileRepository))
		r.Delete("/delete/{slug}", storage.DeleteHandlerFunc(stg, fileRepository))
		r.Post("/upload", storage.UploadHandlerFunc(stg, fileRepository))
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

func initStorage() (storage.Storage, error) {
	storageBucket := "www-nomkhonwaan-com"

	switch viper.GetString("storage") {
	case "gcloud":
		cloudStorage, err := gcloud.NewCloudStorage(
			viper.GetString("gcloud-credentials-file-path"),
			storageBucket,
		)
		if err != nil {
			return nil, err
		}
		return cloudStorage, nil
	case "s3":
		amazonS3, err := aws.NewS3(
			viper.GetString("amazon-s3-access-key"),
			viper.GetString("amazon-s3-secret-key"),
			storageBucket,
		)
		if err != nil {
			return nil, err
		}
		return amazonS3, nil
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
