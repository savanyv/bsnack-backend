package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName string
	AppEnv string
	AppPort string

	DBHost string
	DBPort string
	DBUser string
	DBPassword string
	DBName string

	RedisAddr string
	RedisPassword string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		AppName:    loadEnv("APP_NAME"),
		AppEnv:     loadEnv("APP_ENV"),
		AppPort:    loadEnv("APP_PORT"),
		DBHost:     loadEnv("DB_HOST"),
		DBPort:     loadEnv("DB_PORT"),
		DBUser:     loadEnv("DB_USER"),
		DBPassword: loadEnv("DB_PASSWORD"),
		DBName:     loadEnv("DB_NAME"),
		RedisAddr:  loadEnv("REDIS_ADDR"),
		RedisPassword: loadEnv("REDIS_PASSWORD"),
	}
}

func loadEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic("Environment variable " + key + " not set")
	}

	return value
}
