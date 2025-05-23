package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"regexp"
)

func IsValidURL(urlStr string) bool {
	regex := `^https?:\/\/[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(:[0-9]+)?(\/.*)?$` // https://regex101.com/r/zxsntB
	re := regexp.MustCompile(regex)
	return re.MatchString(urlStr)
}

func HashPassword(password string) string {
	hashAlgo := sha1.New()
	hashAlgo.Write([]byte(password))
	return hex.EncodeToString(hashAlgo.Sum(nil))
}
