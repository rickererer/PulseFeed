package interfaceshttpmiddleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// initGin 确保 gin 处于测试模式。
func init() {
	gin.SetMode(gin.TestMode)
}

// TestJWTAuthErrorFormat 验证鉴权失败响应使用统一的 {"error": "..."} 结构。
func TestJWTAuthErrorFormat(t *testing.T) {
	// jwtManager 为 nil 时只测缺少 Authorization 头的分支（不触发 ParseAndValidateToken）。
	handler := NewJWTAuth(nil)

	r := gin.New()
	r.GET("/protected", handler, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil) // 无 Authorization 头
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	if _, ok := body["error"]; !ok {
		t.Fatalf("expected {\"error\":...} shape, got: %s", w.Body.String())
	}
	if _, ok := body["message"]; ok {
		t.Fatalf("legacy \"message\" field should not appear: %s", w.Body.String())
	}
}
