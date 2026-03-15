package config

import (
	"log"
	"os"

	"github.com/lpernett/godotenv"
)

type Env struct {
	PORT         string
	DATABASE_URL string
	Timezone     string
}

func LoadEnv() Env {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	env := Env{
		PORT:         os.Getenv("PORT"),
		DATABASE_URL: os.Getenv("DATABASE_URL"),
		Timezone:     os.Getenv("TIMEZONE"),
	}

	return env
}
