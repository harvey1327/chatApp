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
	MB_HOST     string
	MB_PORT     int
	MB_USERNAME string
	MB_PASSWORD string
	DB_HOST     string
	DB_PORT     int
	DB_USERNAME string
	DB_PASSWORD string
}

func Load() *Config {
	once.Do(func() {
		instance = &Config{
			MB_HOST:     getEnv("MB_HOST"),
			MB_PORT:     getEnvAsInt("MB_PORT"),
			MB_USERNAME: getEnv("MB_USERNAME"),
			MB_PASSWORD: getEnv("MB_PASSWORD"),
			DB_HOST:     getEnv("DB_HOST"),
			DB_PORT:     getEnvAsInt("DB_PORT"),
			DB_USERNAME: getEnv("DB_USERNAME"),
			DB_PASSWORD: getEnv("DB_PASSWORD"),
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
