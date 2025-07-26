package internal

import "os"

func GetEnvElse(key, fallback string) string {
	value, defined := os.LookupEnv(key)
	if !defined {
		return fallback
	}
	return value
}

func GetRabbitMQURL() string {
	return GetEnvElse("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
}
