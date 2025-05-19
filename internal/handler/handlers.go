package handler

/*
All the endpoints will be implemented here.
*/

import (
	"encoding/json"
	"net/http"
	"web-analyzer/internal/model"
	"web-analyzer/internal/service"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/static/index.html")
}

func UrlAnalyzingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req model.UrlRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Url == "" {
		http.Error(w, "Invalid request: missing or empty 'url' field", http.StatusBadRequest)
		return
	}

	response, err := service.AnalyzeUrl(req.Url)
	if err != nil {
		http.Error(w, "Error processing URL", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encodeError := json.NewEncoder(w).Encode(response)
	if encodeError != nil {
		http.Error(w, "Error processing URL", http.StatusInternalServerError)
		return
	}
}
