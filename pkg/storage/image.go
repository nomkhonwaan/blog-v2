package storage

import "io"

// Resizer uses to resize an image file to the given width and height.
// The resize function should treat aspect ratio if one of them equal to zero,
// also reject when none of them greater than zero.
type Resizer interface {
	Resize(img io.Reader, width, height int) (io.Reader, error)
}
