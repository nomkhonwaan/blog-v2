package server

import (
	"context"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"time"
)

// InsecureServer is a wrapper of net/http.Server
// which embeds net/http.Handler for handling incoming HTTP requests.
type InsecureServer struct {
	// An HTTP handler
	Handler http.Handler

	// Timeout to be waited on shutting-down the server, default is: 5 minutes
	ShutdownTimeout time.Duration
}

// ListenAndServe will listen on a specific address
func (s *InsecureServer) ListenAndServe(addr string, stopCh <-chan struct{}) error {
	httpServer := http.Server{
		Handler: s.Handler,
	}

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	go func() {
		logrus.Infof("server is listening on address: %s", l.Addr().String())

		// Unable to caught an error here,
		// I'm using a simple HTTP server an validate TCP listener before.
		_ = httpServer.Serve(l)
	}()

	<-stopCh

	logrus.Info("server is shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), s.ShutdownTimeout)
	defer cancel()

	// I have no idea which case the shutdown function will return an error :p
	_ = httpServer.Shutdown(ctx)

	logrus.Info("server has been stopped gracefully")
	return nil
}
