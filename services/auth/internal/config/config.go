package config

import (
	"fmt"
	"os"

	"github.com/ilyas/flower/services/auth/internal/utils"
)

type HTTPConfig struct {
	Address string
}

type GRPCConfig struct {
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

type JWTConfig struct {
	Secret           string
	AccessTTLMinutes int
	RefreshTTLDays   int
}

type Config struct {
	GRPC      GRPCConfig
	HTTP      HTTPConfig
	Tarantool TarantoolConfig
	Redis     RedisConfig
	JWT       JWTConfig
}

func Load() (*Config, error) {
	httpAddr, err := getEnv("AUTH_HTTP_ADDRESS")
	if err != nil {
		return nil, err
	}
	grpcAddr := getEnvOrDefault("AUTH_GRPC_ADDRESS", ":9090")
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

	jwtSecret, err := getEnv("JWT_SECRET")
	if err != nil {
		return nil, err
	}
	jwtAccessTTL, err := getEnv("JWT_ACCESS_TTL_MIN")
	if err != nil {
		return nil, err
	}
	jwtAccessTTLInt, err := utils.ToInt(jwtAccessTTL, "JWT_ACCESS_TTL_MIN")
	if err != nil {
		return nil, err
	}
	jwtRefreshTTL, err := getEnv("JWT_REFRESH_TTL_DAYS")
	if err != nil {
		return nil, err
	}
	jwtRefreshTTLInt, err := utils.ToInt(jwtRefreshTTL, "JWT_REFRESH_TTL_DAYS")
	if err != nil {
		return nil, err
	}

	return &Config{
		GRPC: GRPCConfig{
			Address: grpcAddr,
		},
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
		JWT: JWTConfig{
			Secret:           jwtSecret,
			AccessTTLMinutes: jwtAccessTTLInt,
			RefreshTTLDays:   jwtRefreshTTLInt,
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

func getEnvOrDefault(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
