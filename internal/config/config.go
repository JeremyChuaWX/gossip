package config

import (
	"fmt"
	"os"
	"reflect"
)

type config struct {
	PostgresURL   string `env:"POSTGRES_URL"`
	ServerAddress string `env:"SERVER_ADDRESS"`
}

func Init() (config, error) {
	var env config
	queryStructType := reflect.TypeOf(env)
	queryStructValue := reflect.ValueOf(&env).Elem()
	for _, field := range reflect.VisibleFields(queryStructType) {
		key := field.Tag.Get("env")
		value := os.Getenv(key)
		if value == "" {
			return env, missingConfigError(key)
		}
		queryStructValue.FieldByIndex(field.Index).SetString(value)
	}
	return env, nil
}

func missingConfigError(key string) error {
	return fmt.Errorf("missing environment variable: %s", key)
}
