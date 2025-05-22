package model

type AnalyzeResponse struct {
	HtmlVersion            string   `json:"htmlVersion"`
	PageTitle              string   `json:"pageTitle"`
	HeadingsCount          int      `json:"headingsCount"`
	InternalLinksCount     int      `json:"internalLinksCount"`
	ExternalLinksCount     int      `json:"externalLinksCount"`
	InaccessibleLinksCount int      `json:"inaccessibleLinksCount"`
	InaccessibleLinks      []string `json:"inaccessibleLinks"`
	ContainsLoginForm      bool     `json:"containsLoginForm"`
}
