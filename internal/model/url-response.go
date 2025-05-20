package model

type UrlResponse struct {
	HtmlVersion               string   `json:"htmlVersion"`
	PageTitle                 string   `json:"pageTitle"`
	NumberOfHeadings          int      `json:"numberOfHeadings"`
	NumberOfInternalLinks     int      `json:"numberOfInternalLinks"`
	NumberOfExternalLinks     int      `json:"numberOfExternalLinks"`
	NumberOfInaccessibleLinks int      `json:"numberOfInaccessibleLinks"`
	InaccessibleLinks         []string `json:"inaccessibleLinks"`
	ContainsLoginForm         bool     `json:"containsLoginForm"`
}
