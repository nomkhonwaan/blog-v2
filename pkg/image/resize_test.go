package image_test

import (
	"bytes"
	. "github.com/nomkhonwaan/myblog/pkg/image"
	"github.com/stretchr/testify/assert"
	"image"
	"image/png"
	"testing"
)

func TestLanczosResizer_Resize(t *testing.T) {
	t.Run("With successful resizing image", func(t *testing.T) {
		// Given
		var buf bytes.Buffer
		_ = png.Encode(&buf, image.NewRGBA(image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 100, Y: 100}}))

		// When
		result, err := NewLanczosResizer().Resize(&buf, 10, 10)
		img, _, _ := image.Decode(result)

		// Then
		assert.Nil(t, err)
		assert.Equal(t, image.Point{X: 0, Y: 0}, img.Bounds().Min)
		assert.Equal(t, image.Point{X: 10, Y: 10}, img.Bounds().Max)
	})

	t.Run("When an error has occurred while decoding image", func(t *testing.T) {
		// Given
		img := bytes.NewBufferString("invalid image content")

		// When
		_, err := NewLanczosResizer().Resize(img, 10, 10)

		// Then
		assert.EqualError(t, err, "image: unknown format")
	})
}
