package config

import (
	"log"
	"os"
	"strconv"
	"sync"
)

var once sync.Once
var instance *Config

type Config struct {
	HOST              string
	PORT              int
	MB_HOST           string
	MB_PORT           int
	MB_USERNAME       string
	MB_PASSWORD       string
	USER_SERVICE_HOST string
	USER_SERVICE_PORT int
}

func Load() *Config {
	once.Do(func() {
		instance = &Config{
			HOST:              getEnv("HOST"),
			PORT:              getEnvAsInt("PORT"),
			MB_HOST:           getEnv("MB_HOST"),
			MB_PORT:           getEnvAsInt("MB_PORT"),
			MB_USERNAME:       getEnv("MB_USERNAME"),
			MB_PASSWORD:       getEnv("MB_PASSWORD"),
			USER_SERVICE_HOST: getEnv("USER_SERVICE_HOST"),
			USER_SERVICE_PORT: getEnvAsInt("USER_SERVICE_PORT"),
		}
	})
	return instance
}

func getEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	log.Fatalf("No value for %s", key)
	return ""
}

func getEnvAsInt(key string) int {
	valueS := getEnv(key)
	if valueI, err := strconv.Atoi(valueS); err == nil {
		return valueI
	}
	log.Fatalf("%s is not an int", key)
	return 0
}
