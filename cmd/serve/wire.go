//+build wireinject

package serve

import (
	"github.com/google/wire"
	"github.com/nomkhonwaan/myblog/pkg/facebook"
	"github.com/nomkhonwaan/myblog/pkg/mongo"
	"github.com/nomkhonwaan/myblog/pkg/storage"
	"github.com/spf13/viper"
)

//func ProvideMongoDB() (mongo.Database, error) {
//	return NewMongoDB()
//}
//
func ProvideMongoDB() (mongo.Database, error) {
	return NewMongoDB(
		viper.GetString("mongodb-uri"),
		"nomkhonwaan_com",
	)
}

func ProvideFileRepository(db mongo.Database) storage.FileRepository {
	return storage.NewFileRepository(db)
}

func InitFacebookHandler() *facebook.Handler {
	wire.Build(facebook.NewHandler, ProvideFileRepository, ProvideMongoDB)
	return &facebook.Handler{}
}

//func InitPostRepository() blog.PostRepository {
//	wire.Build(NewMongoDB, )
//	//db, err := NewMongoDB()
//	//if err != nil {
//	//	panic(err)
//	//}
//	//return blog.NewPostRepository(db, log.NewDefaultTimer())
//}

//func InitFileRepository() storage.fileRepository {
//	wire.Build(NewMongoDB)
//	return storage.MongoFileRepository{}
//}
//
//func ProvideFacebookHandler() facebook.Handler {
//	wire.Build(blog.NewPostRepository, storage.NewFileRepository)
//	return facebook.Handler{}
//}
