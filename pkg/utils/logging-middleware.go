package utils

/*
LoggingMiddleware is a middleware that logs incoming HTTP requests and their responses.
*/

import (
	"bytes"
	"io"
	"net/http"
	"time"
	"web-analyzer/pkg/logger"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.Log

		// Log the incoming request
		bodyBytes, _ := io.ReadAll(r.Body)
		log.Infof("Request: Method=[%s], Url=[%s], Body=[%s], RemoteAddr=[%s]", r.Method, r.URL.String(), string(bodyBytes), r.RemoteAddr)
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Restore the body for further processing

		// Capture the response
		rec := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK, body: &bytes.Buffer{}}
		start := time.Now()

		next.ServeHTTP(rec, r)

		// Log the response
		duration := time.Since(start)
		log.Infof("Response: Status=[%d], Duration=[%s], Body=[%s]", rec.statusCode, duration, rec.body.String())
	})
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

// WriteHeader captures the status code and writes it to the response.
func (rec *responseRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

// Write captures the response body and writes it to the response.
func (rec *responseRecorder) Write(p []byte) (int, error) {
	rec.body.Write(p)
	return rec.ResponseWriter.Write(p)
}
