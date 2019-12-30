//go:generate mockgen -destination=./mock/resize_mock.go github.com/nomkhonwaan/myblog/pkg/image Resizer

package image

import (
	"bytes"
	"github.com/disintegration/imaging"
	"image/png"
	"io"
)

// Resizer uses to resize an image file to the given width and height.
// The resize function should treat aspect ratio if one of them equal to zero.
type Resizer interface {
	Resize(img io.Reader, width, height int) (io.Reader, error)
}

// LanczosResizer does resizing the image with Lanczos algorithm,
// https://en.wikipedia.org/wiki/Lanczos_algorithm
type LanczosResizer struct{}

// NewLanczosResizer returns a new Lanczos instance
func NewLanczosResizer() LanczosResizer {
	return LanczosResizer{}
}

func (r LanczosResizer) Resize(img io.Reader, width, height int) (io.Reader, error) {
	i, err := imaging.Decode(img, imaging.AutoOrientation(true))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = imaging.Encode(
		&buf,
		imaging.Resize(i, width, height, imaging.Lanczos),
		imaging.PNG,
		imaging.PNGCompressionLevel(png.NoCompression),
	)

	return &buf, err
}
