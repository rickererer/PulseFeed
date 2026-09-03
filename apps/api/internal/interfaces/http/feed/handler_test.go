package interfaceshttpfeed

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	applicationfeed "github.com/rickererer/PulseFeed/internal/application/feed"
	domainfeed "github.com/rickererer/PulseFeed/internal/domain/feed"

	"github.com/gin-gonic/gin"
)

func TestParseLimit(t *testing.T) {
	cases := []struct {
		name    string
		raw     string
		want    int
		wantErr error
	}{
		{"empty default", "", 0, nil},
		{"valid lower bound", "1", 1, nil},
		{"valid upper bound", "100", 100, nil},
		{"valid middle", "20", 20, nil},
		{"zero", "0", 0, domainfeed.ErrInvalidLimit},
		{"negative", "-5", 0, domainfeed.ErrInvalidLimit},
		{"over max", "101", 0, domainfeed.ErrInvalidLimit},
		{"way over max", "10000", 0, domainfeed.ErrInvalidLimit},
		{"non numeric", "abc", 0, domainfeed.ErrInvalidLimit},
		{"whitespace", "  10  ", 10, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseLimit(tc.raw)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("parseLimit(%q) err = %v, want %v", tc.raw, err, tc.wantErr)
			}
			if got != tc.want {
				t.Fatalf("parseLimit(%q) = %d, want %d", tc.raw, got, tc.want)
			}
		})
	}
}

// feedRepoStub 按 scene 返回预置数据或错误，用于 handler 请求-响应断言。
type feedRepoStub struct {
	items      []*domainfeed.FeedPageItem
	listErr    error
	lastLimit  int
	lastCursor *domainfeed.TimelineCursor
}

func (r *feedRepoStub) ListTimelinePage(ctx context.Context, cursor *domainfeed.TimelineCursor, limit int) ([]*domainfeed.FeedPageItem, error) {
	r.lastLimit = limit
	r.lastCursor = cursor
	if r.listErr != nil {
		return nil, r.listErr
	}
	return r.items, nil
}

func (r *feedRepoStub) ListHotPage(ctx context.Context, cursor *domainfeed.HotCursor, limit int) ([]*domainfeed.FeedPageItem, error) {
	return nil, nil
}

func (r *feedRepoStub) ListFollowingPage(ctx context.Context, viewerID int64, cursor *domainfeed.TimelineCursor, limit int) ([]*domainfeed.FeedPageItem, error) {
	return nil, nil
}

func (r *feedRepoStub) ListFollowingPullAuthorIDs(ctx context.Context, viewerID int64) ([]int64, error) {
	return nil, nil
}

func (r *feedRepoStub) BatchGetFeedCards(ctx context.Context, videoIDs []int64) (map[int64]*domainfeed.FeedCard, error) {
	cards := map[int64]*domainfeed.FeedCard{}
	for _, id := range videoIDs {
		cards[id] = &domainfeed.FeedCard{
			VideoID: id, AuthorID: 10, AuthorNickname: " author ",
			Title: " title ", PublishedAt: time.Date(2026, 9, 1, 10, 0, 0, 0, time.UTC),
		}
	}
	return cards, nil
}

func (r *feedRepoStub) BatchGetFeedStats(ctx context.Context, videoIDs []int64) (map[int64]*domainfeed.FeedStat, error) {
	stats := map[int64]*domainfeed.FeedStat{}
	for _, id := range videoIDs {
		stats[id] = &domainfeed.FeedStat{VideoID: id, LikeCount: 3, CommentCount: 2, FavoriteCount: 1}
	}
	return stats, nil
}

func (r *feedRepoStub) BatchGetViewerActionStates(ctx context.Context, viewerID int64, videoIDs []int64) (map[int64]*domainfeed.ViewerActionState, error) {
	return map[int64]*domainfeed.ViewerActionState{}, nil
}

// newFeedTestRouter 构造 handler + 真实 Service + stub 仓储的测试路由。
func newFeedTestRouter(repo *feedRepoStub) *gin.Engine {
	gin.SetMode(gin.TestMode)
	service := applicationfeed.New(repo)
	handler := New(service)
	r := gin.New()
	r.GET("/api/feed-items", handler.ListFeedItems)
	return r
}

// TestListFeedItemsOK 验证正常路径：200 + JSON 字段断言。
func TestListFeedItemsOK(t *testing.T) {
	publishedAt := time.Date(2026, 9, 1, 10, 0, 0, 0, time.UTC)
	repo := &feedRepoStub{items: []*domainfeed.FeedPageItem{
		{VideoID: 101, AuthorID: 10, PublishedAt: publishedAt},
	}}
	r := newFeedTestRouter(repo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/feed-items?scene=timeline&limit=5", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", w.Code, w.Body.String())
	}

	var body struct {
		Scene      string `json:"scene"`
		NextCursor string `json:"next_cursor"`
		HasMore    bool   `json:"has_more"`
		Items      []struct {
			VideoID        int64  `json:"video_id"`
			AuthorID       int64  `json:"author_id"`
			AuthorNickname string `json:"author_nickname"`
			Title          string `json:"title"`
			LikeCount      int    `json:"like_count"`
		} `json:"items"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON: %v; body: %s", err, w.Body.String())
	}
	if body.Scene != "timeline" {
		t.Fatalf("scene = %q, want timeline", body.Scene)
	}
	if len(body.Items) != 1 {
		t.Fatalf("items = %d, want 1", len(body.Items))
	}
	item := body.Items[0]
	if item.VideoID != 101 || item.AuthorID != 10 {
		t.Fatalf("item ids mismatch: %+v", item)
	}
	if item.AuthorNickname != "author" || item.Title != "title" {
		t.Fatalf("display fields not trimmed: %+v", item)
	}
	if item.LikeCount != 3 {
		t.Fatalf("like_count = %d, want 3", item.LikeCount)
	}
}

// TestListFeedItemsInvalidLimit 验证参数校验：非法 limit 返回 400 与统一错误结构。
func TestListFeedItemsInvalidLimit(t *testing.T) {
	r := newFeedTestRouter(&feedRepoStub{})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/feed-items?scene=timeline&limit=101", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := body["error"]; !ok {
		t.Fatalf("expected error field, got: %s", w.Body.String())
	}
}

// TestListFeedItemsRepoError 验证服务端错误：500 + 固定文案（不泄漏内部细节）。
func TestListFeedItemsRepoError(t *testing.T) {
	repo := &feedRepoStub{listErr: errors.New("db connection lost")}
	r := newFeedTestRouter(repo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/feed-items?scene=timeline&limit=10", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if body["error"] != "internal server error" {
		t.Fatalf("error = %q, want fixed internal server error text", body["error"])
	}
}
