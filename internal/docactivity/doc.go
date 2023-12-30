package docactivity

import (
	"context"
	"go.temporal.io/sdk/activity"
	"strings"
)

type DocAction struct{}

func (DocAction) StartEditingActivity(ctx context.Context, param string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("StartEditingActivity pass")
	return strings.ToUpper(param), nil
}

func (DocAction) ApproveActivity(ctx context.Context, param string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("ApproveActivity pass")
	return strings.ToLower(param), nil
}
