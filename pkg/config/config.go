package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	ServerPort int `json:"server_port"`
}

func LoadConfig() (*Config, error) {
	file, err := os.Open("pkg/config/config.json") // relative path from root
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	err = json.NewDecoder(file).Decode(&cfg)
	return &cfg, err
}
