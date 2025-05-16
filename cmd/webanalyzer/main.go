package main

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web-analyzer/internal/handler"
	"web-analyzer/pkg/logger"
	"web-analyzer/pkg/utils"
)

type Server struct {
	httpServer *http.Server
	log        *logrus.Logger
}

// NewServer Creates a new server instance with the specified port and logger.
func NewServer(port string, log *logrus.Logger) *Server {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(handler.RootHandler))
	mux.Handle("/analyze-url", utils.LoggingMiddleware(http.HandlerFunc(handler.PostURLHandler)))

	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + port,
			Handler: mux,
		},
		log: log,
	}
}

/*
Start Starts the server.
*/
func (s *Server) Start() {
	go func() {
		s.log.Infof("Server started successfully on port %s.", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.WithError(err).Fatal("Failed to start the server")
		}
	}()
}

/*
Shutdown Shuts down the server gracefully.
*/
func (s *Server) Shutdown(ctx context.Context) {
	s.log.Info("Shutting down server...")
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.log.WithError(err).Error("Error during server shutdown")
	}
	s.log.Info("Server stopped gracefully.")
}

func main() {
	log := logger.ConfigureLogger()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := NewServer(port, log)

	server.Start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}
