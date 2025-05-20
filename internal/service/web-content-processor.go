package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
	"web-analyzer/internal/model"
	"web-analyzer/pkg/logger"
)

func AnalyzeWebPage(urlStr string) (*model.UrlResponse, error) {
	log := logger.Log

	log.Infof("Analyzing the url: %s", urlStr)

	// Fetch HTML content as string along with the DOCTYPE string. ChromeDP only gives the HTML content.
	htmlStr, err := FetchContentAsString(urlStr)
	if err != nil {
		log.Errorf("Analyzing failed for the url: %s, error: %s", urlStr, err)
		return nil, err
	}
	response, err := AnalyzeContent(htmlStr, urlStr)
	if err != nil {
		log.Errorf("Analyzing failed for the url: %s, error: %s", urlStr, err)
		return nil, err
	}
	log.Infof("Analyzing completed successfully for the url: %s", urlStr)
	return response, nil
}

func AnalyzeContent(htmlStr string, urlStr string) (*model.UrlResponse, error) {
	log := logger.Log

	// Load and parse
	document, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		return nil, err
	}

	result := &model.UrlResponse{
		InaccessibleLinks: make([]string, 0),
	}

	// Page title
	log.Info("Extracting the title.")
	result.PageTitle = FindPageTitle(*document)

	// Headings count
	log.Info("Extracting the headings count.")
	result.NumberOfHeadings = FindHeadingsCount(*document)

	// Login form detection
	log.Info("Extracting the login form.")
	result.ContainsLoginForm = ContainsLoginForm(*document)

	// Analyze links
	log.Info("Analyzing links.")
	result, err = AnalyzeLinks(urlStr, *document, result)

	// Detect HTML version (from raw HTML string)
	log.Info("Finding the HTML version.")
	result.HtmlVersion = DetectHtmlVersion(htmlStr)

	return result, nil
}

func FetchContentAsString(urlStr string) (string, error) {
	log := logger.Log

	urlHash := md5.Sum([]byte(urlStr))
	fileName := hex.EncodeToString(urlHash[:])

	if _, err := os.Stat(fileName); err == nil {
		log.Infof("Found cache file: %s", fileName)
		content, err := os.ReadFile(fileName)
		if err != nil {
			return "", err
		}
		return string(content), nil
	}

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	var renderedHtml string
	var docType string
	err := chromedp.Run(ctx,
		chromedp.Navigate(urlStr),
		chromedp.WaitReady("body"),
		chromedp.Evaluate(`document.doctype ? '<!DOCTYPE ' + document.doctype.name + '>' : ''`, &docType),
		chromedp.OuterHTML("html", &renderedHtml),
	)
	if err != nil {
		return "", err
	}

	fullHtml := docType + "\n" + renderedHtml

	err = os.WriteFile(fileName, []byte(fullHtml), 0644)
	if err != nil {
		return "", err
	}
	log.Infof("Created cache file: %s", fileName)

	return fullHtml, nil
}

func FindPageTitle(document goquery.Document) string {
	return document.Find("title").Text()
}

func FindHeadingsCount(document goquery.Document) int {
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

	var wg sync.WaitGroup
	var mu sync.Mutex

	client := http.Client{Timeout: 3 * time.Second}

	semaphore := make(chan struct{}, 10)

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

		wg.Add(1)
		go func(resolved *url.URL) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			req, _ := http.NewRequest("HEAD", resolved.String(), nil)
			req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; WebAnalyzer/1.0)")
			logger.Log.Infof("Accessing URL: %s", resolved)
			resp, err := client.Do(req)
			if err != nil || resp.StatusCode >= 400 {
				logger.Log.Warnf("Failed accessing URL: %s", resolved)
				mu.Lock()
				result.NumberOfInaccessibleLinks++
				result.InaccessibleLinks = append(result.InaccessibleLinks, resolved.String())
				mu.Unlock()
			}
			logger.Log.Infof("Accessed URL successfully: %s", resolved)
		}(resolved)
	})
	wg.Wait()
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

func DetectHtmlVersion(docType string) string {
	htmlLower := strings.ToLower(docType)
	switch {
	case strings.Contains(htmlLower, "<!doctype html>"):
		return "HTML5"
	case strings.Contains(htmlLower, "html 4.01"):
		return "HTML 4.01"
	case strings.Contains(htmlLower, "xhtml 1.0"):
		return "XHTML 1.0"
	default:
		return "Unknown"
	}
}

func ContainsLoginForm(document goquery.Document) bool {
	hasLoginForm := false
	document.Find("form").Each(func(i int, s *goquery.Selection) {
		if s.Find("input[type='password']").Length() > 0 {
			hasLoginForm = true
		}
	})
	return hasLoginForm
}
