package infraconfig

import (
	"strings"
	"testing"
)

// validConfig 返回一个满足校验的配置，便于各用例做局部破坏。
func validConfig() *Config {
	return &Config{
		Port: 8080,
		JWT: JWTConfig{
			Secret:    "test-secret",
			AccessTTL: "15m",
		},
		Internal: InternalConfig{Token: "test-token"},
		Database: DatabaseConfig{
			Host: "localhost",
			Port: 3306,
			User: "root",
			Name: "PulseFeed",
		},
		Redis: RedisConfig{Addr: "localhost:6379"},
		RabbitMQ: RabbitMQConfig{
			URL:                  "amqp://guest:guest@localhost:5672/",
			InteractionExchange:  "pulsefeed.interaction",
			VideoExchange:        "pulsefeed.video",
		},
	}
}

func TestValidateOK(t *testing.T) {
	if err := validConfig().Validate(); err != nil {
		t.Fatalf("valid config should pass, got: %v", err)
	}
}

func TestValidateMissingFields(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*Config)
		want   string
	}{
		{"port zero", func(c *Config) { c.Port = 0 }, "port"},
		{"jwt secret empty", func(c *Config) { c.JWT.Secret = "  " }, "jwt.secret"},
		{"db host empty", func(c *Config) { c.Database.Host = "" }, "database.host"},
		{"db name empty", func(c *Config) { c.Database.Name = "" }, "database.name"},
		{"redis addr empty", func(c *Config) { c.Redis.Addr = "" }, "redis.addr"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := validConfig()
			tc.mutate(cfg)
			err := cfg.Validate()
			if err == nil {
				t.Fatalf("expected validation error for %s", tc.name)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error %q should mention %q", err.Error(), tc.want)
			}
		})
	}
}

func TestValidateRabbitMQOptional(t *testing.T) {
	cfg := validConfig()
	cfg.RabbitMQ.URL = "" // 禁用 RabbitMQ 是合法场景
	cfg.RabbitMQ.InteractionExchange = ""
	cfg.RabbitMQ.VideoExchange = ""
	if err := cfg.Validate(); err != nil {
		t.Fatalf("empty rabbitmq url should be allowed, got: %v", err)
	}
}

func TestValidateRabbitMQExchangeRequiredWhenEnabled(t *testing.T) {
	cfg := validConfig()
	cfg.RabbitMQ.InteractionExchange = ""
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error when rabbitmq enabled but exchange missing")
	}
	if !strings.Contains(err.Error(), "rabbitmq.interaction_exchange") {
		t.Fatalf("error should mention exchange, got: %v", err)
	}
}
