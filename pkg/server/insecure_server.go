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
		<-stopCh

		logrus.Info("server is shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), s.ShutdownTimeout)
		defer cancel()

		if err := httpServer.Shutdown(ctx); err != nil {
			logrus.Fatalf("an error has occurred while shutting-down the server: %v", err)
		}

		// A server has been stopped gracefully
		logrus.Info("server has been stopped")
	}()

	go func() {
		_ = httpServer.Serve(l)
	}()

	logrus.Infof("server is listening on address: %s", l.Addr().String())
	return nil
}
