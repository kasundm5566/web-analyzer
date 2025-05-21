# web-analyzer

This is a web app used to analyze the content of a webpage.
It reads a webpage content and extracts the following information.

1. HTML version
2. Page title
3. Number of headings
4. Number of internal and external links
5. Count of the inaccessible links
6. List of the inaccessible links
7. Checks whether the page contains a login form

The app is built using GoLang.

## Prerequisites

- Go 1.24.3 or later
- Latest version of Google Chrome

---

## External dependencies

- chromedp: For rendering and fetching webpage content. `github.com/chromedp/chromedp`
- goquery: For parsing and querying HTML documents. `github.com/PuerkitoBio/goquery`
- logger: For structured logging. `github.com/sirupsen/logrus`

---

## Setup

### Clone the repository:

`git clone https://github.com/kasundm5566/web-analyzer.git`

`cd web-analyzer`

### Install Go dependencies:

`go mod tidy`

### Run the application:

`go run main.go`

Access the web interface at http://localhost:8080.

---

## Limitations

- The application relies on Google Chrome for rendering webpages, which may increase resource usage.
- The analysis of inaccessible links is limited to a HEAD request with a short timeout, which may not work for all
  servers.
- The detection of the HTML version is based on the DOCTYPE string and may not be accurate for non-standard pages.
- The results may not be accurate for complex webpages with dynamic content or heavy JavaScript usage or iframes.

---

## Challenges faced

- Getting use to the GoLang syntax and libraries.
- Finding a better library to extract the webpage content.
- Implementation of wait groups.

---

## Possible improvements

- Implement a more robust method for checking link accessibility.
- Implement a more user-friendly interface.
- Implement a caching mechanism to store results for previously analyzed page results.
- Add support for analyzing multiple pages in parallel.
- Implement a more sophisticated method for detecting login forms.
- Add support for analyzing other types of content (e.g., images, videos, etc.).

---

## Demo
