package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	Name         string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  int
}

var AppConfig *Config

func LoadConfig() {
	AppConfig = &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:         getEnv("DB_HOST", "localhost"),
			Port:         getEnv("DB_PORT", "3306"),
			User:         getEnv("DB_USER", "root"),
			Password:     getEnv("DB_PASSWORD", ""),
			Name:         getEnv("DB_NAME", "dday"),
			MaxOpenConns: getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: getEnvInt("DB_MAX_IDLE_CONNS", 25),
			MaxLifetime:  getEnvInt("DB_CONN_MAX_LIFETIME", 300),
		},
	}

	log.Printf("Config loaded - Port: %s, DB: %s@%s:%s/%s",
		AppConfig.Server.Port,
		AppConfig.Database.User,
		AppConfig.Database.Host,
		AppConfig.Database.Port,
		AppConfig.Database.Name)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}