package config

import (
	"fmt"
	"os"
	"strconv"
)

func getEnv(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return valueInt
}

func LoadConfig() (*AppConfig, error) {
	cfg := &AppConfig{
		MySQL: MySQLConfig{
			Username: getEnv("MYSQL_USERNAME", ""),
			Password: getEnv("MYSQL_PASSWORD", ""),
			Host:     getEnv("MYSQL_HOST", "127.0.0.1"),
			Port:     getEnv("MYSQL_PORT", "3306"),
			Database: getEnv("MYSQL_DATABASE", ""),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "127.0.0.1"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		Jwt: JwtConfig{
			Secret: getEnv("JWT_SECRET", "mysecret"),
		},
		RunPort: getEnv("RUN_PORT", "8080"),
	}

	if cfg.MySQL.Username == "" {
		return nil, fmt.Errorf("MYSQL_USERNAME is required")
	}
	if cfg.MySQL.Password == "" {
		return nil, fmt.Errorf("MYSQL_PASSWORD is required")
	}
	if cfg.MySQL.Database == "" {
		return nil, fmt.Errorf("MYSQL_DATABASE is required")
	}

	return cfg, nil
}
