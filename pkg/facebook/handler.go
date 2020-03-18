package facebook

import (
	"github.com/nomkhonwaan/myblog/pkg/blog"
	"github.com/nomkhonwaan/myblog/pkg/storage"
	"text/template"
)

// Handler provides Facebook HTTP handler functions
type Handler struct {
	baseURL           string
	openGraphTemplate *template.Template
	postRepository    blog.PostRepository
	fileRepository    storage.FileRepository
}

//// NewHandler returns a new Handler instance
//func NewHandler(baseURL string, postRepository blog.PostRepository, fileRepository storage.FileRepository, openGraphTemplate *template.Template) *Handler {
//	//return &Handler{
//	//	baseURL:           baseURL,
//	//	openGraphTemplate: openGraphTemplate,
//	//	postRepository:    postRepository,
//	//	fileRepository:    fileRepository,
//	//}
//}

// HandleCrawler uses to handling Facebook sharing bot crawler
// which returns static HTML rather than single page application
// for displaying on Facebook feed.
//func (h Handler) HandleCrawler(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		if IsFacebookCrawlerRequest(r.UserAgent()) {
//			if id, yes := IsSinglePage(r.URL.Path); yes {
//				postID, err := primitive.ObjectIDFromHex(id)
//				if err == nil {
//					//buf, statusCode := h.serveStaticSinglePage(r.Context(), postID)
//					//if err != nil {
//					//	return
//					//}
//					//return
//				}
//
//			}
//		}
//
//		next.ServeHTTP(w, r)
//	})
//}

//func (h Handler) serveStaticSinglePage(ctx context.Context, id interface{}) ([]byte, int) {
//	p, err := h.postRepository.FindByID(ctx, id)
//	if err != nil {
//		return nil, http.StatusNotFound
//	}
//	if p.Status != blog.Published {
//		return nil, http.StatusForbidden
//	}
//
//	defaultFeaturedImage :=
//}
