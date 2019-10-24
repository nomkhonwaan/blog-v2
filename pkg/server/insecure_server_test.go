package server_test

import (
	. "github.com/nomkhonwaan/myblog/pkg/server"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

func TestInsecureServer_ListenAndServe(t *testing.T) {
	t.Run("With successful listening and serving the server", func(t *testing.T) {
		// Given
		stopCh := make(chan struct{})
		go func() {
			// Wait for before closing the server
			time.Sleep(time.Millisecond * 100)

			close(stopCh)
		}()

		s := InsecureServer{}
		addr := ":0"

		// When
		err := s.ListenAndServe(addr, stopCh)

		// Then
		assert.Nil(t, err)
	})

	t.Run("When unable to binding TCP port", func(t *testing.T) {
		// Given
		l, _ := net.Listen("tcp", ":0")
		defer l.Close()

		s := InsecureServer{}

		// When
		err := s.ListenAndServe(l.Addr().String(), nil)

		// Then
		assert.IsType(t, &net.OpError{}, err)
	})

	t.Run("When it takes too long to shutdown", func(t *testing.T) {
		// Given
		stopCh := make(chan struct{})
		go func() {
			// Wait for before closing the server
			time.Sleep(time.Millisecond * 100)

			close(stopCh)
		}()

		s := InsecureServer{ShutdownTimeout: 0}

		// When
		err := s.ListenAndServe(":0", stopCh)

		// Then
		logrus.Error(err)
	})
}
