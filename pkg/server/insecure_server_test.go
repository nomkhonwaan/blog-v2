package server_test

import (
	. "github.com/nomkhonwaan/myblog/pkg/server"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestInsecureServer_ListenAndServe(t *testing.T) {
	t.Run("With successful listening and serving the server", func(t *testing.T) {
		// Given
		stopCh := make(chan struct{})
		defer func() { close(stopCh) }()

		s := InsecureServer{}
		addr := ":0"

		// When
		err := s.ListenAndServe(addr, stopCh)

		// Then
		assert.Nil(t, err)
	})

	t.Run("When unable to binding TCP port", func(t *testing.T) {
		// Given
		stopCh := make(chan struct{})
		defer func() { close(stopCh) }()

		l, _ := net.Listen("tcp", ":0")
		defer l.Close()

		s := InsecureServer{}

		// When
		err := s.ListenAndServe(l.Addr().String(), stopCh)

		// Then
		assert.IsType(t, &net.OpError{}, err)
	})
}
