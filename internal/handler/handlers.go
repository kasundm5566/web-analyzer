package handler

/*
All the endpoints will be implemented here.
*/

import (
	"encoding/json"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/static/index.html")
}

type URLRequest struct {
	URL string `json:"url"`
}

type URLResponse struct {
	HtmlVersion               int      `json:"htmlVersion"`
	PageTitle                 string   `json:"pageTitle"`
	NumberOfHeadings          int      `json:"numberOfHeadings"`
	NumberOfInternalLinks     int      `json:"numberOfInternalLinks"`
	NumberOfExternalLinks     int      `json:"numberOfExternalLinks"`
	NumberOfInaccessibleLinks int      `json:"numberOfInaccessibleLinks"`
	InaccessibleLinks         []string `json:"inaccessibleLinks"`
	ContainsLoginForm         bool     `json:"containsLoginForm"`
}

func PostURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req URLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
		http.Error(w, "Invalid request: missing or empty 'url' field", http.StatusBadRequest)
		return
	}

	response := URLResponse{
		HtmlVersion:               5,
		PageTitle:                 "Sample Title",
		NumberOfHeadings:          5,
		NumberOfInternalLinks:     10,
		NumberOfExternalLinks:     5,
		NumberOfInaccessibleLinks: 2,
		InaccessibleLinks:         []string{"http://example.com/broken-link1", "http://example.com/broken-link2"},
		ContainsLoginForm:         true,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
