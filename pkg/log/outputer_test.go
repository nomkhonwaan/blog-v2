package log

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestDefaultOutputer_Printf(t *testing.T) {
	// Given
	var buf bytes.Buffer
	l := log.New(&buf, "", 0)
	outputer := DefaultOutputer{Logger: l}

	// When
	outputer.Printf("The quick brown %s jumps over the lazy %s", "fox", "dog")

	// Then
	assert.Equal(t, "The quick brown fox jumps over the lazy dog\n", buf.String())
}
