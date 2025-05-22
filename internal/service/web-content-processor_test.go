package service

import (
	"strings"
	"testing"
	"web-analyzer/internal/model"

	"github.com/PuerkitoBio/goquery"
)

func TestFindPageTitle(t *testing.T) {
	html := `<html><head><title>Test Page</title></head><body></body></html>`
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	title := FindPageTitle(*doc)
	if title != "Test Page" {
		t.Errorf("FindPageTitle() = %q; want %q", title, "Test Page")
	}
}

func TestFindHeadingsCount(t *testing.T) {
	html := `<html><body><h1>Heading 1</h1><h2>Heading 2</h2><h3>Heading 3</h3></body></html>`
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	count := FindHeadingsCount(*doc)
	if count != 3 {
		t.Errorf("FindHeadingsCount() = %d; want %d", count, 3)
	}
}

func TestAnalyzeLinks(t *testing.T) {
	html := `<html><body>
		<a href="https://example.com/test.html">External Link</a>
		<a href="/internal">Internal Link</a>
		<a href="https://abc.com">Internal Link</a>
		<a href="mailto:abc@example.com">Internal Link</a>
		<a href="tel:123456789">Internal Link</a>
		<a href="ftp://ftp.example.com/file.txt">Internal Link</a>
	</body></html>`
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	result := &model.AnalyzeResponse{
		InaccessibleLinks: make([]string, 0),
	}
	urlStr := "https://example.com"
	result, err := AnalyzeLinks(urlStr, *doc, result)
	if err != nil {
		t.Errorf("AnalyzeLinks() failed: %v", err)
	}

	if result.InternalLinksCount != 2 {
		t.Errorf("AnalyzeLinks() internal links = %d; want %d", result.InternalLinksCount, 2)
	}
	if result.ExternalLinksCount != 1 {
		t.Errorf("AnalyzeLinks() external links = %d; want %d", result.ExternalLinksCount, 1)
	}
}

func TestDetectHtmlVersion(t *testing.T) {
	tests := []struct {
		docType  string
		expected string
	}{
		{"<!DOCTYPE html>", "HTML5"},
		{"<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 4.01 Transitional//EN\">", "HTML 4.01"},
		{"<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Strict//EN\">", "XHTML 1.0"},
		{"", "Unknown"},
	}

	for _, test := range tests {
		result := DetectHtmlVersion(test.docType)
		if result != test.expected {
			t.Errorf("DetectHtmlVersion(%q) = %q; want %q", test.docType, result, test.expected)
		}
	}
}

func TestContainsLoginForm(t *testing.T) {
	html := `<html><body><form><input type="password"></form></body></html>`
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	hasLoginForm := ContainsLoginForm(*doc)
	if !hasLoginForm {
		t.Errorf("ContainsLoginForm() = %v; want %v", hasLoginForm, true)
	}
}

func TestContainsString(t *testing.T) {
	list := []string{"car", "van", "bus"}
	if !ContainsString(list, "bus") {
		t.Errorf("ContainsString() = false; want true")
	}
	if ContainsString(list, "lorry") {
		t.Errorf("ContainsString() = true; want false")
	}
}
