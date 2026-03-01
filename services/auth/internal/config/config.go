package config

import (
	"fmt"
	"os"

	"github.com/ilyas/flower/services/auth/internal/utils"
)

type HTTPConfig struct {
	Address string
}

type TarantoolConfig struct {
	Addr     string
	User     string
	Password string
}

type RedisConfig struct {
	Addr                string
	Password            string
	DB                  int
	ConfirmationCodeTTL int
}

// общий конфиг сервиса.
type Config struct {
	HTTP      HTTPConfig
	Tarantool TarantoolConfig
	Redis     RedisConfig
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

	redisAddr, err := getEnv("REDIS_ADDR")
	if err != nil {
		return nil, err
	}
	redisPassword, err := getEnv("REDIS_PASSWORD")
	if err != nil {
		return nil, err
	}
	redisDB, err := getEnv("REDIS_DB")
	if err != nil {
		return nil, err
	}
	redisDBInt, err := utils.ToInt(redisDB, "REDIS_DB")
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
		Redis: RedisConfig{
			Addr:                redisAddr,
			Password:            redisPassword,
			DB:                  redisDBInt,
			ConfirmationCodeTTL: 0,
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
