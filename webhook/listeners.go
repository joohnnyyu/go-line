package webhook

import "context"

// 需求/任务/缺陷类
type (
	FollowListener interface {
		OnFollow(ctx context.Context, event *FollowEvent) error
	}

	UnFollowListener interface {
		OnUnFollow(ctx context.Context, event *UnFollowEvent) error
	}
)
