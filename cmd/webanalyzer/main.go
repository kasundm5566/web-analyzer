package main

import (
	"net/http"
	"strconv"
	"web-analyzer/internal/handler"
	"web-analyzer/pkg/config"
	"web-analyzer/pkg/logger"
	"web-analyzer/pkg/utils"
)

type Server struct {
	httpServer *http.Server
}

// NewServer Creates a new server instance with the specified port and logger.
func NewServer(port string) *Server {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(handler.RootHandler))
	mux.Handle("/analyze.html", utils.AuthMiddleware(http.HandlerFunc(handler.AnalyzePageHandler)))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	mux.Handle("/analyze-url", utils.LoggingMiddleware(http.HandlerFunc(handler.WebPageAnalyzingHandler)))
	mux.Handle("/login", utils.LoggingMiddleware(http.HandlerFunc(handler.LoginHandler)))

	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + port,
			Handler: mux,
		},
	}
}

func main() {
	logger.ConfigureLogger()
	configurations, err := config.LoadConfig()
	if err != nil {
		logger.Log.Errorf("Error loading config: %v", err)
	}

	port := configurations.ServerPort
	if port == 0 {
		port = 8080
	}

	server := NewServer(strconv.Itoa(port))

	logger.Log.Infof("Server starting on port %s...", strconv.Itoa(port))
	if err := server.httpServer.ListenAndServe(); err != nil {
		logger.Log.WithError(err).Fatal("Failed to start the server")
	}
}
