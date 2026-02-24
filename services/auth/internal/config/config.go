package config

import (
	"fmt"
	"os"
)

// HTTPConfig описывает настройки HTTP-сервера.
type HTTPConfig struct {
	Address string
}

// TarantoolConfig описывает настройки подключения к Tarantool.
type TarantoolConfig struct {
	Addr     string
	User     string
	Password string
}

// Config — общий конфиг сервиса.
type Config struct {
	HTTP      HTTPConfig
	Tarantool TarantoolConfig
}

// Load загружает конфиг из переменных окружения с дефолтами.
func Load() (*Config, error) {
	httpAddr, err := getEnv("AUTH_HTTP_ADDRESS")
	if err != nil {
		return nil, err
	}
	tntAddr, err := getEnv("TARANTOOL_ADDR")
	if err != nil {
		return nil, err
	}
	tntUser, err := getEnv("TARANTOOL_USER")
	if err != nil {
		return nil, err
	}
	tntPassword, err := getEnv("TARANTOOL_PASSWORD")
	if err != nil {
		return nil, err
	}

	return &Config{
		HTTP: HTTPConfig{
			Address: httpAddr,
		},
		Tarantool: TarantoolConfig{
			Addr:     tntAddr,
			User:     tntUser,
			Password: tntPassword,
		},
	}, nil
}

func getEnv(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return "", fmt.Errorf("environment variable %s is not set", key)
	}
	return v, nil
}
