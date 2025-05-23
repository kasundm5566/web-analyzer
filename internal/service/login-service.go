package service

import (
	"crypto/sha1"
	"encoding/hex"
	"sync"
	"web-analyzer/pkg/utils"
)

type LoginService struct {
	userDB map[string]string
	mu     sync.RWMutex
}

var instance *LoginService
var once sync.Once

func GetLoginService() *LoginService {
	once.Do(func() {
		hashAlgo := sha1.New()
		hashAlgo.Write([]byte("password"))
		hashedPassword := hex.EncodeToString(hashAlgo.Sum(nil))

		instance = &LoginService{
			userDB: map[string]string{
				"user": hashedPassword,
			},
		}
	})
	return instance
}

func (s *LoginService) ValidateCredentials(username, password string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if storedPassword, exists := s.userDB[username]; exists && storedPassword == utils.HashPassword(password) {
		return true
	}
	return false
}
