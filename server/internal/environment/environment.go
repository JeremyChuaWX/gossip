package environment

import (
	"errors"
	"os"
	"reflect"
)

var missingEnvError = errors.New("missing environment variable")

type environment struct {
	PostgresURL   string `env:"POSTGRES_URL"`
	RedisURL      string `env:"REDIS_URL"`
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
			return env, missingEnvError
		}
		queryStructValue.FieldByIndex(field.Index).SetString(value)
	}
	return env, nil
}
