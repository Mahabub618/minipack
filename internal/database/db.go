package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mahabub618/minipack/config"
)

var DB *sqlx.DB

func InitDB(cfg *config.Config) error {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBPassword, cfg.DBName)
	var err error
	DB, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect database: %v", err)
	}
	return nil
}
