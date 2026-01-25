package config

import (
	"os"
)


type Config struct {
	API_PORT string
}

func LoadConfig() Config {
	return Config{
		API_PORT: os.Getenv("API_PORT"),
	}
}
