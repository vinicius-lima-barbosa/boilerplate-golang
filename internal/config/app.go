package config

import (
	"log"
	"os"
	"strconv"

	"github.com/lpernett/godotenv"
)

type Env struct {
	PORT                string
	DATABASE_URL        string
	Timezone            string
	JWTAccessSecret     string
	JWTRefreshSecret    string
	JWTAccessTTLMinutes int
	JWTRefreshTTLHours  int
}

func LoadEnv() Env {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	env := Env{
		PORT:                os.Getenv("PORT"),
		DATABASE_URL:        os.Getenv("DATABASE_URL"),
		Timezone:            os.Getenv("TIMEZONE"),
		JWTAccessSecret:     os.Getenv("JWT_ACCESS_SECRET"),
		JWTRefreshSecret:    os.Getenv("JWT_REFRESH_SECRET"),
		JWTAccessTTLMinutes: getEnvAsInt("JWT_ACCESS_TTL_MINUTES", 15),
		JWTRefreshTTLHours:  getEnvAsInt("JWT_REFRESH_TTL_HOURS", 168),
	}

	if env.JWTAccessSecret == "" {
		log.Fatal("JWT_ACCESS_SECRET is not set in environment variables")
	}

	if env.JWTRefreshSecret == "" {
		log.Fatal("JWT_REFRESH_SECRET is not set in environment variables")
	}

	return env
}

func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("Invalid value for %s: %v", key, err)
	}

	return parsedValue
}
