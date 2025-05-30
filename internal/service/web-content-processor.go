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

func AnalyzeWebPage(urlStr string) (*model.AnalyzeResponse, error) {
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

func AnalyzeContent(htmlStr string, urlStr string) (*model.AnalyzeResponse, error) {
	log := logger.Log

	// Load and parse
	document, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		return nil, err
	}

	result := &model.AnalyzeResponse{
		InaccessibleLinks: make([]string, 0),
	}

	// Page title
	log.Info("Extracting the title.")
	result.PageTitle = FindPageTitle(*document)

	// Headings count
	log.Info("Extracting the headings count.")
	result.HeadingsCount = FindHeadingsCount(*document)

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

	cacheDir := "cache"
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		err = os.Mkdir(cacheDir, 0755)
		if err != nil {
			return "", err
		}
	}

	urlHash := md5.Sum([]byte(urlStr))
	fileName := cacheDir + string(os.PathSeparator) + hex.EncodeToString(urlHash[:])

	if _, err := os.Stat(fileName); err == nil {
		log.Infof("Found cache file: %s for the url: %s", fileName, urlStr)
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
	log.Infof("Created cache file: %s for the url: %s", fileName, urlStr)

	return fullHtml, nil
}

func FindPageTitle(document goquery.Document) string {
	return document.Find("title").Text()
}

func FindHeadingsCount(document goquery.Document) int {
	totalHeadings := 0
	for i := 1; i <= 6; i++ {
		tag := "h" + string(rune('0'+i))
		totalHeadings += document.Find(tag).Length()
	}
	return totalHeadings
}

func AnalyzeLinks(urlStr string, document goquery.Document, result *model.AnalyzeResponse) (*model.AnalyzeResponse, error) {

	// Parse base URL
	base, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	links, err := findAndCategorizeLinks(base, document, result)
	if err != nil {
		return nil, err
	}

	checkLinkAccessibility(links, result)

	return result, nil
}

func findAndCategorizeLinks(base *url.URL, document goquery.Document, result *model.AnalyzeResponse) ([]*url.URL, error) {
	linksFound := make(map[string]struct{}) // Use a map to avoid duplicates
	var links []*url.URL

	document.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists || href == "" {
			return
		}

		if _, found := linksFound[href]; found {
			return
		}

		link, err := url.Parse(href)
		if err != nil || (link.Scheme != "" && link.Scheme != "http" && link.Scheme != "https") {
			return
		}

		resolved := base.ResolveReference(link)
		linksFound[href] = struct{}{}
		links = append(links, resolved)

		if resolved.Host == base.Host || strings.HasPrefix(href, "/") || strings.HasPrefix(href, "#") {
			result.InternalLinksCount++
		} else {
			result.ExternalLinksCount++
		}
	})

	return links, nil
}

func checkLinkAccessibility(links []*url.URL, result *model.AnalyzeResponse) {
	log := logger.ConfigureLogger()

	var wg sync.WaitGroup
	var mu sync.Mutex
	client := http.Client{Timeout: 3 * time.Second}
	semaphore := make(chan struct{}, 10)

	for _, link := range links {
		wg.Add(1)
		go func(link *url.URL) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			req, _ := http.NewRequest("HEAD", link.String(), nil)
			req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; WebAnalyzer/1.0)")
			log.Infof("Accessing URL: %s", link)
			resp, err := client.Do(req)
			if err != nil || resp.StatusCode >= 400 {
				log.Warnf("Failed accessing URL: %s", link)
				mu.Lock()
				result.InaccessibleLinksCount++
				result.InaccessibleLinks = append(result.InaccessibleLinks, link.String())
				mu.Unlock()
			} else {
				log.Infof("Accessed URL successfully: %s", link)
			}
		}(link)
	}

	wg.Wait()
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
