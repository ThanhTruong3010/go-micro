package utils

import (
	"fmt"
	"os"
)

var serviceConfig = map[string]struct {
	host    string
	portEnv string
	portDef string
}{
	"auth":     {"authentication-service", "AUTH_PORT", "8081"},
	"logger":   {"logger-service", "LOGGER_PORT", "8082"},
	"mailer":   {"mailer-service", "MAILER_PORT", "8083"},
	"broker":   {"broker-service", "BROKER_PORT", "8080"},
	"rabbitmq": {"rabbitmq", "RABBITMQ_PORT", "5672"},
}

func GetServiceURL(service string) string {
	cfg, ok := serviceConfig[service]
	if !ok {
		cfg = serviceConfig["broker"]
	}

	host := "localhost"
	envMode := GetEnv("MODE", "development")
	if envMode != "development" {
		host = cfg.host
	}

	return fmt.Sprintf("http://%s:%s", host, GetEnv(cfg.portEnv, cfg.portDef))
}

func GetEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func GetRabbitMQURL() string {
	host := "localhost"
	envMode := GetEnv("MODE", "development")
	if envMode != "development" {
		host = "rabbitmq"
	}
	return fmt.Sprintf("amqp://guest:guest@%s", host)
}
