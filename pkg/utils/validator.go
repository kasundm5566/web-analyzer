package utils

import (
	"regexp"
)

func IsValidURL(urlStr string) bool {
	regex := `^https?:\/\/[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(:[0-9]+)?(\/.*)?$` // https://regex101.com/r/zxsntB
	re := regexp.MustCompile(regex)
	return re.MatchString(urlStr)
}
