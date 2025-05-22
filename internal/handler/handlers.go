package handler

/*
All the endpoints will be implemented here.
*/

import (
	"encoding/json"
	"net/http"
	"web-analyzer/internal/model"
	"web-analyzer/internal/service"
	"web-analyzer/pkg/utils"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/static/index.html")
}

func WebPageAnalyzingHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req model.AnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Url == "" || !utils.IsValidURL(req.Url) {
		http.Error(w, "Invalid request: invalid url. It should start with http:// or https://.", http.StatusBadRequest)
		return
	}

	response, err := service.AnalyzeWebPage(req.Url)
	if err != nil {
		http.Error(w, "Error processing URL. "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encodeError := json.NewEncoder(w).Encode(response)
	if encodeError != nil {
		http.Error(w, "Error processing URL. "+encodeError.Error(), http.StatusInternalServerError)
		return
	}
}
