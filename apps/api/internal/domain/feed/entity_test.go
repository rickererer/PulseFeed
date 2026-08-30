package domainfeed

import (
	"testing"
	"time"
)

// TestNormalizeScene 验证场景归一化：空值用默认、大小写与空白统一。
func TestNormalizeScene(t *testing.T) {
	cases := []struct {
		name  string
		input Scene
		want  Scene
	}{
		{"empty defaults to timeline", "", SceneTimeline},
		{"whitespace defaults", "  ", SceneTimeline},
		{"uppercase normalized", "TIMELINE", SceneTimeline},
		{"mixed case normalized", "Following", SceneFollowing},
		{"surrounding whitespace trimmed", "  hot  ", SceneHot},
		{"unknown scene preserved", "custom", Scene("custom")},
		{"recommend preserved", SceneRecommend, SceneRecommend},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := NormalizeScene(tc.input); got != tc.want {
				t.Fatalf("NormalizeScene(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

// TestScoreHotFeedItem 验证热榜权重：评论(5) > 收藏(4) > 点赞(3)。
func TestScoreHotFeedItem(t *testing.T) {
	cases := []struct {
		name          string
		like, comment int
		favorite      int
		want          int
	}{
		{"zero interaction", 0, 0, 0, 0},
		{"like weight 3", 1, 0, 0, 3},
		{"comment weight 5", 0, 1, 0, 5},
		{"favorite weight 4", 0, 0, 1, 4},
		{"mixed weights", 2, 1, 3, 2*3 + 1*5 + 3*4},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := ScoreHotFeedItem(tc.like, tc.comment, tc.favorite); got != tc.want {
				t.Fatalf("ScoreHotFeedItem(%d,%d,%d) = %d, want %d",
					tc.like, tc.comment, tc.favorite, got, tc.want)
			}
		})
	}
}

// TestRestoreFeedItem 验证展示字段清洗与热度联动。
func TestRestoreFeedItem(t *testing.T) {
	publishedAt := time.Date(2026, 8, 30, 10, 0, 0, 0, time.UTC)
	item := RestoreFeedItem(
		101, 202, "  nickname  ", " http://avatar ", "  title  ",
		"  description  ", " http://media ", " http://cover ",
		2, 1, 3, publishedAt,
	)

	if item.VideoID != 101 || item.AuthorID != 202 {
		t.Fatalf("ids not preserved: %+v", item)
	}
	for field, got := range map[string]string{
		"authorNickname": item.AuthorNickname,
		"authorAvatar":   item.AuthorAvatarURL,
		"title":          item.Title,
		"description":    item.Description,
		"mediaURL":       item.MediaURL,
		"coverURL":       item.CoverURL,
	} {
		if len(got) > 0 && (got[0] == ' ' || got[len(got)-1] == ' ') {
			t.Fatalf("%s not trimmed: %q", field, got)
		}
	}
	if item.HotScore != ScoreHotFeedItem(2, 1, 3) {
		t.Fatalf("HotScore = %d, want %d", item.HotScore, ScoreHotFeedItem(2, 1, 3))
	}
	if !item.PublishedAt.Equal(publishedAt) {
		t.Fatalf("PublishedAt not preserved: %v", item.PublishedAt)
	}
}

// TestFeedDomainConstants 锁定领域常量语义，防止未来误改破坏业务约束。
func TestFeedDomainConstants(t *testing.T) {
	if MaxLimit != 100 {
		t.Fatalf("MaxLimit = %d, want 100", MaxLimit)
	}
	if BigCreatorFollowerThreshold != 10000 {
		t.Fatalf("BigCreatorFollowerThreshold = %d, want 10000", BigCreatorFollowerThreshold)
	}
}
