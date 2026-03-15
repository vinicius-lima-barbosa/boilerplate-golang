package config

import (
	"log"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	Env "github.com/vinicius-lima-barbosa/boilerplate-golang/internal/config"

	gormLogger "gorm.io/gorm/logger"
)

type Database struct {
	DB *gorm.DB
}

var DB Database

func Connection() {
	env := Env.LoadEnv()

	databaseUrl := env.DATABASE_URL
	if databaseUrl == "" {
		log.Fatal("Database URL is not set in environment variables")
	}

	if !strings.Contains(strings.ToLower(databaseUrl), "sslmode=") {
		separator := "?"
		if strings.Contains(databaseUrl, "?") {
			separator = "&"
		}
		databaseUrl += separator + "sslmode=disable"
	}

	m, err := migrate.New(
		"file://internal/database/migrations",
		databaseUrl,
	)

	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}

	handleMigrateUp(m)

	db, err := gorm.Open(postgres.Open(databaseUrl), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db.Logger = gormLogger.Default.LogMode(gormLogger.Info)

	DB = Database{
		DB: db,
	}
}

func handleMigrateUp(m *migrate.Migrate) {
	if m == nil {
		log.Fatal("migration instance is nil")
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("no change")
			return
		}
		log.Fatalf("failed to apply migration: %v", err)
	}
}
