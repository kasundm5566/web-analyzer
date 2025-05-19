package main

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"web-analyzer/internal/handler"
	"web-analyzer/pkg/logger"
	"web-analyzer/pkg/utils"
)

type Server struct {
	httpServer *http.Server
}

// NewServer Creates a new server instance with the specified port and logger.
func NewServer(port string, log *logrus.Logger) *Server {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(handler.RootHandler))
	mux.Handle("/analyze-url", utils.LoggingMiddleware(http.HandlerFunc(handler.UrlAnalyzingHandler)))

	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + port,
			Handler: mux,
		},
	}
}

func main() {
	logger.ConfigureLogger()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := NewServer(port, logger.Log)

	logger.Log.Infof("Server starting on port %s...", port)
	if err := server.httpServer.ListenAndServe(); err != nil {
		logger.Log.WithError(err).Fatal("Failed to start the server")
	}
}
