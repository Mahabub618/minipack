package config

import (
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	DBHost             string        `envconfig: "DB_HOST"`
	DBPort             string        `envconfig: "DB_PORT"`
	DBUser             string        `envconfig: "DB_USER"`
	DBPassword         string        `envconfig: "DB_PASSWORD"`
	DBName             string        `envconfig: "DB_NAME"`
	JWTSecret          string        `envconfig: "JWT_SECRET"`
	AccessTokenExpiry  time.Duration `envconfig: "ACCESS_TOKEN_EXPIRY"`
	RefreshTokenExpiry time.Duration `envconfig: "REFRESH_TOKEN_EXPIRY"`
}

var (
	instance *Config
	once     sync.Once
)

func LoadConfig() (*Config, error) {
	var err error
	once.Do(func() {
		var config Config
		if err = envconfig.Process("", &config); err != nil {
			log.Errorf("Error parsing environment variables: %v", err)
			return
		}
		instance = &config
	})
	return instance, err
}
