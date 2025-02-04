package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/mahabub618/minipack/config"
	log "github.com/sirupsen/logrus"
)

var DB *pgxpool.Pool

func InitDB(cfg *config.Config) error {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	var err error
	DB, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		return fmt.Errorf("failed to connect database: %v", err)
	}
	log.Println("Database connected successfully")
	return nil
}

func GetDB() *pgxpool.Pool {
	if DB == nil {
		log.Fatal("Database not initialized. Call InitDB first.")
	}
	return DB
}
