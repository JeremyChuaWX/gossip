package environment

import (
	"fmt"
	"os"
	"reflect"
)

type environment struct {
	PostgresURL   string `env:"SERVER_POSTGRES_URL"`
	RedisURL      string `env:"SERVER_REDIS_URL"`
	ServerAddress string `env:"SERVER_ADDRESS"`
}

func Init() (environment, error) {
	var env environment
	queryStructType := reflect.TypeOf(env)
	queryStructValue := reflect.ValueOf(&env).Elem()
	for _, field := range reflect.VisibleFields(queryStructType) {
		key := field.Tag.Get("env")
		value := os.Getenv(key)
		if value == "" {
			return env, missingEnvError(key)
		}
		queryStructValue.FieldByIndex(field.Index).SetString(value)
	}
	return env, nil
}

func missingEnvError(key string) error {
	return fmt.Errorf("missing environment variable: %s", key)
}
