package infrajwt

import (
	"errors"
	"testing"
	"time"
)

// TestNewManagerAccessTTLFromConfig 验证 access_ttl 从配置解析生效（而非固定默认值）。
func TestNewManagerAccessTTLFromConfig(t *testing.T) {
	m, err := NewManager("test-secret", "30m")
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}
	if got := m.AccessTTL(); got != 30*time.Minute {
		t.Fatalf("AccessTTL = %v, want 30m", got)
	}
}

// TestNewManagerDefaultTTL 验证空配置回退到默认 15 分钟。
func TestNewManagerDefaultTTL(t *testing.T) {
	m, err := NewManager("test-secret", "")
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}
	if got := m.AccessTTL(); got != defaultAccessTTL {
		t.Fatalf("AccessTTL = %v, want default %v", got, defaultAccessTTL)
	}
}

// TestNewManagerInvalidTTL 验证非法/非正 TTL 配置直接报错（fail-fast，不静默回退）。
func TestNewManagerInvalidTTL(t *testing.T) {
	cases := []string{"abc", "0s", "-5m", "1x"}
	for _, raw := range cases {
		_, err := NewManager("test-secret", raw)
		if !errors.Is(err, ErrParseAccessTTL) {
			t.Fatalf("NewManager(%q) err = %v, want ErrParseAccessTTL", raw, err)
		}
	}
}

// TestNewManagerEmptySecret 验证空 secret 被拒绝。
func TestNewManagerEmptySecret(t *testing.T) {
	if _, err := NewManager("", "15m"); !errors.Is(err, ErrEmptyJWTSecret) {
		t.Fatalf("expected ErrEmptyJWTSecret, got %v", err)
	}
}

// TestSignAccessTokenExpiryMatchesTTL 验证签发 token 的 exp-iat 与配置 TTL 一致。
// 这是"改配置后登录过期时间变化"的直接证据：TTL 从 30m 改为 1h，token 有效期随之变化。
func TestSignAccessTokenExpiryMatchesTTL(t *testing.T) {
	cases := []struct {
		ttlRaw string
		want   time.Duration
	}{
		{"30m", 30 * time.Minute},
		{"1h", time.Hour},
		{"45m", 45 * time.Minute},
	}
	for _, tc := range cases {
		m, err := NewManager("test-secret", tc.ttlRaw)
		if err != nil {
			t.Fatalf("NewManager(%q) failed: %v", tc.ttlRaw, err)
		}
		token, err := m.SignAccessToken(100, "user")
		if err != nil {
			t.Fatalf("SignAccessToken failed: %v", err)
		}
		claims, err := m.ParseAndValidateToken(token, TokenTypeAccess)
		if err != nil {
			t.Fatalf("ParseAndValidateToken failed: %v", err)
		}
		got := time.Duration(claims.ExpiresAt-claims.IssuedAt) * time.Second
		if got != tc.want {
			t.Fatalf("ttl=%s token lifetime = %v, want %v", tc.ttlRaw, got, tc.want)
		}
	}
}

// TestParseAndValidateTokenTypeMismatch 验证用错误 token 类型校验会被拒绝。
func TestParseAndValidateTokenTypeMismatch(t *testing.T) {
	m, _ := NewManager("test-secret", "15m")
	token, err := m.SignAccessToken(100, "user")
	if err != nil {
		t.Fatalf("SignAccessToken failed: %v", err)
	}
	if _, err := m.ParseAndValidateToken(token, "refresh"); !errors.Is(err, ErrInvalidTokenType) {
		t.Fatalf("expected ErrInvalidTokenType, got %v", err)
	}
}

// TestParseAndValidateEmptyToken 验证空 token 被拒绝。
func TestParseAndValidateEmptyToken(t *testing.T) {
	m, _ := NewManager("test-secret", "15m")
	if _, err := m.ParseAndValidateToken("", TokenTypeAccess); !errors.Is(err, ErrEmptyToken) {
		t.Fatalf("expected ErrEmptyToken, got %v", err)
	}
}
