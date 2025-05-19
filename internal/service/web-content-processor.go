package service

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"web-analyzer/internal/model"
	"web-analyzer/pkg/logger"
)

func AnalyzeUrl(urlStr string) (*model.UrlResponse, error) {
	log := logger.Log

	log.Infof("Analyzing the url: %s", urlStr)

	// Fetch HTML content from the URL
	response, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	htmlStr := string(bodyBytes)

	// Load and parse
	document, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		return nil, err
	}

	result := &model.UrlResponse{
		InaccessibleLinks: make([]string, 0),
	}

	// Page title
	result.PageTitle = FindPageTitle(*document)

	// Headings count
	result.NumberOfHeadings = findHeadingsCount(*document)

	// Login form detection
	document.Find("form").Each(func(i int, s *goquery.Selection) {
		if s.Find("input[type='password']").Length() > 0 {
			result.ContainsLoginForm = true
		}
	})

	// Analyze links
	result, err = AnalyzeLinks(urlStr, *document, result)

	// Detect HTML version (from raw HTML string)
	result.HtmlVersion = DetectHtmlVersion(htmlStr)

	log.Infof("Analyzing complete for the url: %s", urlStr)

	return result, nil
}

func FindPageTitle(document goquery.Document) string {
	return document.Find("title").Text()
}

func findHeadingsCount(document goquery.Document) int {
	totalHeadings := 0
	for i := 1; i <= 6; i++ {
		tag := "h" + string('0'+i)
		totalHeadings += document.Find(tag).Length()
	}
	return totalHeadings
}

func AnalyzeLinks(urlStr string, document goquery.Document, result *model.UrlResponse) (*model.UrlResponse, error) {
	// Parse base Url
	base, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	// Link analysis
	linksFound := []string{} //  This is to keep track of links already found to avoid duplicate counts.
	document.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists || href == "" || ContainsString(linksFound, href) {
			return
		}
		link, err := url.Parse(href)
		if err != nil {
			return
		}
		resolved := base.ResolveReference(link)
		linksFound = append(linksFound, href)

		if resolved.Host == base.Host || strings.HasPrefix(href, "/") || strings.HasPrefix(href, "#") {
			result.NumberOfInternalLinks++
		} else {
			result.NumberOfExternalLinks++
		}

		client := http.Client{Timeout: 3 * time.Second}
		req, _ := http.NewRequest("HEAD", resolved.String(), nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; WebAnalyzer/1.0)")

		resp, err := client.Do(req)
		if err != nil || resp.StatusCode >= 400 {
			result.NumberOfInaccessibleLinks++
			result.InaccessibleLinks = append(result.InaccessibleLinks, resolved.String())
		}
	})
	return result, nil
}

func ContainsString(list []string, target string) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}

func DetectHtmlVersion(html string) int {
	htmlLower := strings.ToLower(html)
	switch {
	case strings.Contains(htmlLower, "<!doctype html>"):
		return 5
	case strings.Contains(htmlLower, "xhtml 1.0"):
		return 10
	case strings.Contains(htmlLower, "html 4.01"):
		return 4
	default:
		return 0
	}
}
