package infraconfig

import (
	"fmt"
	"strings"
)

// Validate 校验启动配置中的必填项，缺失返回描述性错误。
// 目标：让配置问题在启动阶段尽早暴露，而非运行期才失败。
func (c *Config) Validate() error {
	var missing []string

	if c.Port <= 0 || c.Port > 65535 {
		missing = append(missing, "port (1-65535)")
	}
	if strings.TrimSpace(c.JWT.Secret) == "" {
		missing = append(missing, "jwt.secret")
	}
	if strings.TrimSpace(c.JWT.AccessTTL) == "" {
		missing = append(missing, "jwt.access_ttl")
	}
	if strings.TrimSpace(c.Internal.Token) == "" {
		missing = append(missing, "internal.token")
	}

	db := c.Database
	if strings.TrimSpace(db.Host) == "" {
		missing = append(missing, "database.host")
	}
	if db.Port <= 0 || db.Port > 65535 {
		missing = append(missing, "database.port (1-65535)")
	}
	if strings.TrimSpace(db.User) == "" {
		missing = append(missing, "database.user")
	}
	if strings.TrimSpace(db.Name) == "" {
		missing = append(missing, "database.name")
	}

	if strings.TrimSpace(c.Redis.Addr) == "" {
		missing = append(missing, "redis.addr")
	}

	// RabbitMQ 是可选项：URL 为空时上层选择禁用异步队列。
	if strings.TrimSpace(c.RabbitMQ.URL) != "" {
		if strings.TrimSpace(c.RabbitMQ.InteractionExchange) == "" {
			missing = append(missing, "rabbitmq.interaction_exchange")
		}
		if strings.TrimSpace(c.RabbitMQ.VideoExchange) == "" {
			missing = append(missing, "rabbitmq.video_exchange")
		}
	}

	if len(missing) == 0 {
		return nil
	}
	return fmt.Errorf("missing or invalid config: %s", strings.Join(missing, ", "))
}
