package handler

/*
All the endpoints will be implemented here.
*/

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"time"
	"web-analyzer/internal/model"
	"web-analyzer/internal/service"
	"web-analyzer/pkg/utils"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/static/login.html")
}

func AnalyzePageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/static/analyze.html")
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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Username == "" || req.Password == "" {
		http.Error(w, "Invalid request: both username and password are required", http.StatusBadRequest)
		return
	}

	loginService := service.GetLoginService()
	response := model.LoginResponse{}

	if loginService.ValidateCredentials(req.Username, req.Password) {
		http.SetCookie(w, &http.Cookie{
			Name:     "web_analyzer_session_token",
			Value:    uuid.New().String(),
			Expires:  time.Now().Add(1 * time.Hour),
			HttpOnly: true,
		})
		response.Status = "success"
		w.WriteHeader(http.StatusOK)
	} else {
		response.Status = "unauthorized"
		w.WriteHeader(http.StatusUnauthorized)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response. "+err.Error(), http.StatusInternalServerError)
	}
}
