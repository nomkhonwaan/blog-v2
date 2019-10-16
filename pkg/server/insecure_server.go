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
	Handler         http.Handler
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
			logrus.Fatalf("error: %v", err)
		}
	}()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorf("recover: %v", r)
			}
		}()

		err := httpServer.Serve(l)

		select {
		case <-stopCh:
			logrus.Info("server has been stopped")
		default:
			logrus.Errorf("error: %v", err)
		}
	}()

	logrus.Infof("server is listening on address: %s", l.Addr().String())
	return nil
}
