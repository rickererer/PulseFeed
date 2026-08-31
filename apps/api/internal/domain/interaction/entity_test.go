package domaininteraction

import (
	"errors"
	"strings"
	"testing"
	"time"
)

// TestNormalizeActionType 验证行为类型归一化与非法类型拒绝。
func TestNormalizeActionType(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    string
		wantErr error
	}{
		{"like uppercase", "LIKE", ActionTypeLike, nil},
		{"like lowercase", "like", ActionTypeLike, nil},
		{"like mixed case", "LiKe", ActionTypeLike, nil},
		{"favorite with spaces", "  favorite  ", ActionTypeFavorite, nil},
		{"unknown type", "SHARE", "", ErrInvalidActionType},
		{"empty", "", "", ErrInvalidActionType},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NormalizeActionType(tc.input)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("NormalizeActionType(%q) err = %v, want %v", tc.input, err, tc.wantErr)
			}
			if got != tc.want {
				t.Fatalf("NormalizeActionType(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

// TestNewComment 验证评论创建校验：视频/用户/内容/幂等键。
func TestNewComment(t *testing.T) {
	longContent := strings.Repeat("a", MaxCommentContentLength+1)
	tooLongIDKey := strings.Repeat("k", MaxIdempotencyKeyLength+1)

	cases := []struct {
		name          string
		videoID       int64
		userID        int64
		content       string
		idempotency   string
		wantErr       error
		wantContent   string
		wantStatus    int
	}{
		{"valid comment", 1, 2, "  hello  ", "key-1", nil, "hello", CommentStatusNormal},
		{"invalid video id", 0, 2, "hello", "key-1", ErrInvalidVideoID, "", 0},
		{"invalid user id", 1, 0, "hello", "key-1", ErrInvalidUserID, "", 0},
		{"empty content", 1, 2, "   ", "key-1", ErrEmptyCommentContent, "", 0},
		{"content too long", 1, 2, longContent, "key-1", ErrCommentContentTooLong, "", 0},
		{"idempotency too long", 1, 2, "hello", tooLongIDKey, ErrIdempotencyKeyTooLong, "", 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			comment, err := NewComment(tc.videoID, tc.userID, tc.content, tc.idempotency)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("NewComment err = %v, want %v", err, tc.wantErr)
			}
			if err != nil {
				return
			}
			if comment.Content != tc.wantContent {
				t.Fatalf("content = %q, want %q (should be trimmed)", comment.Content, tc.wantContent)
			}
			if comment.Status != tc.wantStatus {
				t.Fatalf("status = %d, want %d", comment.Status, tc.wantStatus)
			}
			if comment.VideoID != tc.videoID || comment.UserID != tc.userID {
				t.Fatalf("ids not preserved: %+v", comment)
			}
		})
	}
}

// TestRestoreAction 验证行为恢复：类型归一化与状态默认。
func TestRestoreAction(t *testing.T) {
	now := time.Now()
	action := RestoreAction(1, 2, 3, "  like  ", 0, "  key-1  ", now, now)

	if action.ActionType != ActionTypeLike {
		t.Fatalf("ActionType = %q, want LIKE (normalized)", action.ActionType)
	}
	if action.Status != ActionStatusActive {
		t.Fatalf("Status = %d, want ActionStatusActive (default)", action.Status)
	}
	if action.IdempotencyKey != "key-1" {
		t.Fatalf("IdempotencyKey = %q, want trimmed", action.IdempotencyKey)
	}
	if !action.Active() {
		t.Fatal("Active() should be true for restored active action")
	}
}

// TestRestoreComment 验证评论恢复：字段清洗与状态默认。
func TestRestoreComment(t *testing.T) {
	now := time.Now()
	comment := RestoreComment(1, 2, 3, "  nick  ", "  avatar  ", "  content  ", 0, "  key-1  ", now, now)

	if comment.Content != "content" || comment.UserNickname != "nick" || comment.UserAvatarURL != "avatar" {
		t.Fatalf("fields not trimmed: %+v", comment)
	}
	if comment.Status != CommentStatusNormal {
		t.Fatalf("Status = %d, want CommentStatusNormal (default)", comment.Status)
	}
	if comment.Deleted() {
		t.Fatal("Deleted() should be false for normal comment")
	}
}

// TestStatusPredicates 验证状态判断方法。
func TestStatusPredicates(t *testing.T) {
	canceled := &Action{Status: ActionStatusCanceled}
	if canceled.Active() {
		t.Fatal("canceled action should not be Active")
	}
	deleted := &Comment{Status: CommentStatusDeleted}
	if !deleted.Deleted() {
		t.Fatal("deleted comment should report Deleted")
	}
}
