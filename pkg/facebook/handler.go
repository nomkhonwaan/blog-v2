package facebook

import (
	"github.com/nomkhonwaan/myblog/pkg/storage"
)

// Handler provides Facebook HTTP handler functions
type Handler struct {
	//PostRepository blog.PostRepository
	fileRepository storage.FileRepository
}

func NewHandler(fileRepository storage.FileRepository) *Handler {
	return &Handler{fileRepository: fileRepository}
}

// HandlerOption is a function for applying option to the Handler
type HandlerOption func(*Handler)

//// WithPostRepository allows to setup blog.PostRepository to the Handler
//func WithPostRepository(postRepository blog.PostRepository) HandlerOption {
//	return func(h *Handler) {
//		h.PostRepository = postRepository
//	}
//}

// WithFileRepository allows to setup storage.fileRepository to the Handler
func WithFileRepository(fileRepository storage.FileRepository) HandlerOption {
	return func(h *Handler) {
		h.fileRepository = fileRepository
	}
}
