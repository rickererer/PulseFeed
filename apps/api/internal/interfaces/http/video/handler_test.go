package interfaceshttpvideo

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	domainvideo "github.com/rickererer/PulseFeed/internal/domain/video"
)

// 构造带 query 参数的 gin 上下文。
func newQueryContext(query map[string]string) *gin.Context {
	c := &gin.Context{}
	u, _ := url.Parse("http://localhost/videos")
	q := u.Query()
	for k, v := range query {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	c.Request = &http.Request{URL: u}
	return c
}

func TestParsePagination(t *testing.T) {
	cases := []struct {
		name       string
		query      map[string]string
		wantLimit  int
		wantOffset int
		wantErr    error
	}{
		{"defaults", map[string]string{}, defaultListLimit, 0, nil},
		{"valid limit and offset", map[string]string{"limit": "10", "offset": "20"}, 10, 20, nil},
		{"limit zero", map[string]string{"limit": "0"}, 0, 0, domainvideo.ErrInvalidLimit},
		{"limit negative", map[string]string{"limit": "-1"}, 0, 0, domainvideo.ErrInvalidLimit},
		{"limit over max", map[string]string{"limit": "101"}, 0, 0, domainvideo.ErrInvalidLimit},
		{"limit non numeric", map[string]string{"limit": "x"}, 0, 0, domainvideo.ErrInvalidLimit},
		{"offset negative", map[string]string{"offset": "-1"}, 0, 0, domainvideo.ErrInvalidOffset},
		{"offset non numeric", map[string]string{"offset": "x"}, 0, 0, domainvideo.ErrInvalidOffset},
		{"offset zero ok", map[string]string{"offset": "0"}, defaultListLimit, 0, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := newQueryContext(tc.query)
			limit, offset, err := parsePagination(c)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("parsePagination err = %v, want %v", err, tc.wantErr)
			}
			if limit != tc.wantLimit || offset != tc.wantOffset {
				t.Fatalf("parsePagination = (%d,%d), want (%d,%d)", limit, offset, tc.wantLimit, tc.wantOffset)
			}
		})
	}
}
