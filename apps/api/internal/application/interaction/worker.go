package applicationinteraction

import (
	domaininteraction "github.com/rickererer/PulseFeed/internal/domain/interaction"
	inframetrics "github.com/rickererer/PulseFeed/internal/infra/metrics"
	"context"
	"log"
	"time"
)

type ActionEventConsumer interface {
	ConsumeActionChanged(ctx context.Context, handler func(context.Context, *ActionChangedEvent) error) error
}

type ActionWorker struct {
	repo     domaininteraction.Repository
	consumer ActionEventConsumer
}

func NewActionWorker(repo domaininteraction.Repository, consumer ActionEventConsumer) *ActionWorker {
	return &ActionWorker{
		repo:     repo,
		consumer: consumer,
	}
}

func (w *ActionWorker) Start(ctx context.Context) error {
	if w == nil || w.consumer == nil {
		return nil
	}
	return w.consumer.ConsumeActionChanged(ctx, w.HandleActionChanged)
}

func (w *ActionWorker) HandleActionChanged(ctx context.Context, event *ActionChangedEvent) error {
	start := time.Now()
	var err error
	defer func() {
		inframetrics.ObserveWorkerJob("interaction_action_changed", time.Since(start), err)
		if err != nil && event != nil {
			// 互动落库失败：用户已看到 Redis 快写反馈，这里必须留痕供对账，
			// 否则错误只会反映在指标上而无法定位原因。
			log.Printf("interaction action persist failed: user=%d video=%d type=%s active=%v err=%v",
				event.UserID, event.VideoID, event.ActionType, event.Active, err)
		}
	}()

	if event == nil {
		return nil
	}
	_, _, _, err = w.repo.SetAction(ctx, event.UserID, event.VideoID, event.ActionType, event.Active, event.IdempotencyKey)
	return err
}
