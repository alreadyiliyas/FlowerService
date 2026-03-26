package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	tests := []struct {
		name      string
		env       map[string]string
		want      *Config
		wantError bool
	}{
		{
			name: "success",
			env: map[string]string{
				"AUTH_HTTP_ADDRESS":   "test_8080",
				"TARANTOOL_ADDR":      "test_tarantool:3301",
				"TARANTOOL_USER":      "test_user",
				"TARANTOOL_PASSWORD":  "test_password",
				"REDIS_ADDR":          "test_redis:6379",
				"REDIS_PASSWORD":      "test_redis_password",
				"REDIS_DB":            "1",
				"JWT_SECRET":          "secret",
				"JWT_ACCESS_TTL_MIN":  "30",
				"JWT_REFRESH_TTL_DAYS":"7",
			},
			want: &Config{
				GRPC: GRPCConfig{
					Address: ":9090",
				},
				HTTP: HTTPConfig{
					Address: "test_8080",
				},
				Tarantool: TarantoolConfig{
					Addr:     "test_tarantool:3301",
					User:     "test_user",
					Password: "test_password",
				},
				Redis: RedisConfig{
					Addr:                "test_redis:6379",
					Password:            "test_redis_password",
					DB:                  1,
					ConfirmationCodeTTL: 0,
				},
				JWT: JWTConfig{
					Secret:           "secret",
					AccessTTLMinutes: 30,
					RefreshTTLDays:   7,
				},
			},
			wantError: false,
		},
		{
			name: "empty env variable AUTH_HTTP_ADDRESS",
			env: map[string]string{
				"AUTH_HTTP_ADDRESS":  "",
				"TARANTOOL_ADDR":     "test_tarantool:3301",
				"TARANTOOL_USER":     "test_user",
				"TARANTOOL_PASSWORD": "test_password",
			},
			want:      nil,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}
			got, err := Load()
			if tt.wantError {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "environment variable")
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
