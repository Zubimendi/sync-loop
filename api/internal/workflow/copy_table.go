package workflow

import (
	"time"
	"go.temporal.io/sdk/workflow"
	"github.com/Zubimendi/sync-loop/api/internal/activity"
)

func CopyTableWorkflow(ctx workflow.Context, table string) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	for {
		if err := workflow.ExecuteActivity(ctx, activity.CopyTableActivity, table).Get(ctx, nil); err != nil {
			return err
		}
		if err := workflow.Sleep(ctx, 1*time.Minute); err != nil {
			return err
		}
	}
}