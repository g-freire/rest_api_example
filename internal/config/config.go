package config

import (
	_ "github.com/joho/godotenv/autoload"
	"os"
	"strings"
)

const (
	defaultPostgresURILocal = "postgres://gym:gym@localhost:5432/gym?sslmode=disable"
	defaultPort             = "5000"
)

type Config struct {
	Environment    string   `json:"environment,omitempty"`
	Port           string   `json:"port"`
	PostgresHost   string   `json:"pg_host"`
	RedisAddress   string   `json:"redis_address"`
	RedisPassword  string   `json:"redis_password"`
	EmailSender    string   `json:"email_sender,omitempty"`
	EmailReceivers []string `json:"email_receivers"`
	EmailAPIKey    string   `json:"email_api_key,omitempty"`
}

func GetConfig() *Config {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "DEV"
	}
	receiversEnv := os.Getenv("EMAIL_RECEIVERS")
	var receivers []string
	if receiversEnv == "" {
		receivers = []string{receiversEnv}
	} else {
		receivers = strings.Split(receiversEnv, ",")
	}
	port := os.Getenv("GIN_PORT")
	if port == "" {
		port = defaultPort
	}
	postgresHost := os.Getenv("POSTGRES_HOST")
	if postgresHost == "" {
		postgresHost = defaultPostgresURILocal
	}

	return &Config{
		Environment:    env,
		PostgresHost:   postgresHost,
		Port:           port,
		RedisAddress:   os.Getenv("REDIS_ADDRESS"),
		RedisPassword:  os.Getenv("REDIS_PASSWORD"),
		EmailSender:    os.Getenv("EMAIL_SENDER"),
		EmailAPIKey:    os.Getenv("EMAIL_API_KEY"),
		EmailReceivers: receivers,
	}
}
