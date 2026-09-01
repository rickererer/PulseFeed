package infracache

import (
	applicationfeed "github.com/rickererer/PulseFeed/internal/application/feed"
	domainfeed "github.com/rickererer/PulseFeed/internal/domain/feed"
	domaininteraction "github.com/rickererer/PulseFeed/internal/domain/interaction"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

type actionStatFakeRedis struct {
	hashes map[string]map[string]string
	values map[string]string
}

func newActionStatFakeRedis() *actionStatFakeRedis {
	return &actionStatFakeRedis{
		hashes: map[string]map[string]string{},
		values: map[string]string{},
	}
}

func (r *actionStatFakeRedis) HGetAll(ctx context.Context, key string) *redis.MapStringStringCmd {
	values := r.hashes[key]
	if values == nil {
		values = map[string]string{}
	}
	return redis.NewMapStringStringResult(values, nil)
}

func (r *actionStatFakeRedis) Get(ctx context.Context, key string) *redis.StringCmd {
	value, ok := r.values[key]
	if !ok {
		return redis.NewStringResult("", redis.Nil)
	}
	return redis.NewStringResult(value, nil)
}

func (r *actionStatFakeRedis) Set(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd {
	switch typed := value.(type) {
	case string:
		r.values[key] = typed
	case []byte:
		r.values[key] = string(typed)
	default:
		content, _ := json.Marshal(typed)
		r.values[key] = string(content)
	}
	return redis.NewStatusResult("OK", nil)
}

func (r *actionStatFakeRedis) MGet(ctx context.Context, keys ...string) *redis.SliceCmd {
	values := make([]any, 0, len(keys))
	for _, key := range keys {
		if value, ok := r.values[key]; ok {
			values = append(values, value)
			continue
		}
		values = append(values, nil)
	}
	return redis.NewSliceResult(values, nil)
}

func TestActionStatAggregatesCounterShards(t *testing.T) {
	ctx := context.Background()
	videoID := int64(1001)
	redisClient := newActionStatFakeRedis()
	redisClient.hashes[interactionStatCounterBaseKey(videoID)] = map[string]string{
		"like_count":     "10",
		"comment_count":  "3",
		"favorite_count": "4",
	}
	redisClient.hashes[interactionStatCounterShardKey(videoID, interactionStatCounterShardIndex(42))] = map[string]string{
		"like_count":     "1",
		"favorite_count": "1",
	}
	redisClient.hashes[interactionStatCounterShardKey(videoID, interactionStatCounterShardIndex(43))] = map[string]string{
		"like_count": "-1",
	}
	redisClient.hashes[interactionStatCounterShardKey(videoID, interactionStatCounterShardIndex(44))] = map[string]string{
		"like_count": "1",
	}

	stat, err := actionStat(ctx, redisClient, interactionStatCounterBaseKey(videoID), interactionStatCounterShardKeys(videoID), feedStatKey(videoID), videoID, nil)
	if err != nil {
		t.Fatalf("actionStat: %v", err)
	}
	if stat.LikeCount != 11 || stat.FavoriteCount != 5 || stat.CommentCount != 3 {
		t.Fatalf("unexpected stat: %+v", stat)
	}
}

func TestActionStatFallsBackToInitialStat(t *testing.T) {
	ctx := context.Background()
	videoID := int64(1002)
	redisClient := newActionStatFakeRedis()
	initial := &domaininteraction.VideoStat{
		VideoID:       videoID,
		LikeCount:     7,
		CommentCount:  2,
		FavoriteCount: 1,
	}
	redisClient.hashes[interactionStatCounterShardKey(videoID, interactionStatCounterShardIndex(42))] = map[string]string{
		"like_count":     "1",
		"favorite_count": "-1",
	}

	stat, err := actionStat(ctx, redisClient, interactionStatCounterBaseKey(videoID), interactionStatCounterShardKeys(videoID), feedStatKey(videoID), videoID, initial)
	if err != nil {
		t.Fatalf("actionStat: %v", err)
	}
	if stat.LikeCount != 8 || stat.FavoriteCount != 0 || stat.CommentCount != 2 {
		t.Fatalf("unexpected stat: %+v", stat)
	}
}

func TestGetStatsReadsShardedCountersOnJSONMiss(t *testing.T) {
	ctx := context.Background()
	videoID := int64(1003)
	redisClient := newActionStatFakeRedis()
	redisClient.hashes[interactionStatCounterBaseKey(videoID)] = map[string]string{
		"like_count":     "2",
		"comment_count":  "1",
		"favorite_count": "0",
	}
	redisClient.hashes[interactionStatCounterShardKey(videoID, interactionStatCounterShardIndex(42))] = map[string]string{
		"like_count":     "1",
		"favorite_count": "1",
	}
	stats, err := getStats(ctx, redisClient, []int64{videoID})
	if err != nil {
		t.Fatalf("GetStats: %v", err)
	}
	stat := stats[videoID]
	if stat == nil || stat.LikeCount != 3 || stat.FavoriteCount != 1 || stat.CommentCount != 1 {
		t.Fatalf("unexpected stats: %+v", stats)
	}
	if _, ok := redisClient.values[feedStatKey(videoID)]; !ok {
		t.Fatalf("expected sharded stat to be written back to JSON cache")
	}
}

func TestSetVideoStatWritesJSONCache(t *testing.T) {
	ctx := context.Background()
	videoID := int64(1005)
	redisClient := newActionStatFakeRedis()

	err := setActionStatJSON(ctx, redisClient, feedStatKey(videoID), videoStatToFeedStat(&domaininteraction.VideoStat{
		VideoID:       videoID,
		LikeCount:     2,
		CommentCount:  3,
		FavoriteCount: 1,
	}))
	if err != nil {
		t.Fatalf("SetVideoStat: %v", err)
	}

	stats, err := getStats(ctx, redisClient, []int64{videoID})
	if err != nil {
		t.Fatalf("getStats: %v", err)
	}
	stat := stats[videoID]
	if stat == nil || stat.LikeCount != 2 || stat.CommentCount != 3 || stat.FavoriteCount != 1 {
		t.Fatalf("unexpected stat: %+v", stat)
	}
}

func TestActionStatBaseInitUsesInitialStat(t *testing.T) {
	videoID := int64(1004)
	initial := &domaininteraction.VideoStat{
		VideoID:       videoID,
		LikeCount:     1,
		CommentCount:  1,
		FavoriteCount: 1,
	}

	stat := actionStatBaseInit(videoID, initial)
	if stat != initial {
		t.Fatalf("unexpected stat: %+v", stat)
	}
}

// TestPageCacheRoundTrip 验证页缓存写入后命中返回相同内容，未命中返回 miss 而非错误。
func TestPageCacheRoundTrip(t *testing.T) {
	client := newActionStatFakeRedis()
	ctx := context.Background()
	key := "feed:page:test:1"
	want := &applicationfeed.FeedPage{
		Scene:      domainfeed.SceneTimeline,
		Items:      []*domainfeed.FeedPageItem{{VideoID: 1, AuthorID: 2}},
		NextCursor: "cursor-1",
	}

	if err := setPage(ctx, client, key, want, time.Minute); err != nil {
		t.Fatalf("setPage failed: %v", err)
	}

	got, hit, err := getPage(ctx, client, key)
	if err != nil || !hit {
		t.Fatalf("getPage hit=%v err=%v, want hit=true err=nil", hit, err)
	}
	if got.Scene != want.Scene || got.NextCursor != want.NextCursor {
		t.Fatalf("page mismatch: %+v", got)
	}
	if len(got.Items) != 1 || got.Items[0].VideoID != 1 {
		t.Fatalf("items mismatch: %+v", got.Items)
	}

	if _, hit, err := getPage(ctx, client, "feed:page:missing"); err != nil || hit {
		t.Fatalf("missing key: hit=%v err=%v, want hit=false err=nil", hit, err)
	}
}

// TestGetCardsPartialHit 验证卡片缓存部分命中：命中的返回，未命中的跳过。
func TestGetCardsPartialHit(t *testing.T) {
	client := newActionStatFakeRedis()
	ctx := context.Background()
	publishedAt := time.Date(2026, 9, 1, 10, 0, 0, 0, time.UTC)

	// 预填两个卡片 JSON，模拟 SetCards 已写入后的缓存状态。
	for _, card := range []*domainfeed.FeedCard{
		{VideoID: 1, AuthorID: 10, Title: "one", PublishedAt: publishedAt},
		{VideoID: 2, AuthorID: 10, Title: "two", PublishedAt: publishedAt},
	} {
		content, err := json.Marshal(card)
		if err != nil {
			t.Fatalf("marshal card: %v", err)
		}
		client.values[feedCardKey(card.VideoID)] = string(content)
	}

	got, err := getCards(ctx, client, []int64{1, 2, 3})
	if err != nil {
		t.Fatalf("getCards failed: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("cards = %d entries, want 2 (partial hit)", len(got))
	}
	if got[1].Title != "one" || got[2].Title != "two" {
		t.Fatalf("card content mismatch: %+v %+v", got[1], got[2])
	}
	if _, ok := got[3]; ok {
		t.Fatal("missing video should not appear in result")
	}
}

// TestGetCardsEmptyInput 验证空入参直接返回空结果。
func TestGetCardsEmptyInput(t *testing.T) {
	client := newActionStatFakeRedis()
	got, err := getCards(context.Background(), client, nil)
	if err != nil {
		t.Fatalf("getCards failed: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("cards = %d entries, want 0", len(got))
	}
}
