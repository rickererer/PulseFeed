package applicationinteraction

import (
	"context"
	"errors"
	"testing"

	domaininteraction "github.com/rickererer/PulseFeed/internal/domain/interaction"
)

// fakeRepo 最小实现 domaininteraction.Repository，仅 setActionAsync 路径用到 GetVideoStat。
type fakeRepo struct{}

func (f *fakeRepo) GetVideoStat(ctx context.Context, videoID int64) (*domaininteraction.VideoStat, error) {
	return &domaininteraction.VideoStat{VideoID: videoID}, nil
}
func (f *fakeRepo) GetVideoAuthorID(ctx context.Context, videoID int64) (int64, error)          { return 0, nil }
func (f *fakeRepo) GetUserProfile(ctx context.Context, userID int64) (*domaininteraction.UserProfile, error) {
	return &domaininteraction.UserProfile{}, nil
}
func (f *fakeRepo) SetAction(ctx context.Context, userID int64, videoID int64, actionType string, active bool, idempotencyKey string) (*domaininteraction.Action, int, int, error) {
	return nil, 0, 0, nil
}
func (f *fakeRepo) CreateComment(ctx context.Context, comment *domaininteraction.Comment) (*domaininteraction.Comment, int, int, error) {
	return nil, 0, 0, nil
}
func (f *fakeRepo) FindCommentByUserAndIdempotencyKey(ctx context.Context, userID int64, idempotencyKey string) (*domaininteraction.Comment, int, error) {
	return nil, 0, nil
}
func (f *fakeRepo) ListComments(ctx context.Context, videoID int64, cursor *domaininteraction.CommentCursor, limit int) ([]*domaininteraction.Comment, error) {
	return nil, nil
}
func (f *fakeRepo) DeleteComment(ctx context.Context, commentID int64, userID int64, role string) (*domaininteraction.Comment, int, int, error) {
	return nil, 0, 0, nil
}

// fakeStateStore 按脚本返回指定错误。
type fakeStateStore struct {
	err error
}

func (f *fakeStateStore) SetActionState(ctx context.Context, userID int64, videoID int64, actionType string, active bool, idempotencyKey string, initialStat *domaininteraction.VideoStat) (*ActionStateResult, error) {
	return nil, f.err
}

func TestSetActionConflictPropagated(t *testing.T) {
	s := &Service{
		repo:             &fakeRepo{},
		actionStateStore: &fakeStateStore{err: ErrActionConflict},
	}

	_, err := s.setActionAsync(context.Background(), 1, 2, "LIKE", true, "key-1")
	if !errors.Is(err, ErrActionConflict) {
		t.Fatalf("expected ErrActionConflict, got: %v", err)
	}
}

func TestSetActionOtherErrorWrapped(t *testing.T) {
	s := &Service{
		repo:             &fakeRepo{},
		actionStateStore: &fakeStateStore{err: errors.New("redis down")},
	}

	_, err := s.setActionAsync(context.Background(), 1, 2, "LIKE", true, "key-1")
	if errors.Is(err, ErrActionConflict) {
		t.Fatal("non-conflict error should not be ErrActionConflict")
	}
	if !errors.Is(err, ErrUpdateInteractionFailed) {
		t.Fatalf("expected ErrUpdateInteractionFailed wrapper, got: %v", err)
	}
}
