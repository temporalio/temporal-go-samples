package openNclosed

import (
	"context"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// OpenNClosedWorkflow workflow definition.
func OpenNClosedWorkflow(ctx workflow.Context, keepOpen bool) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 1 * time.Hour,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 1,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("OpenNClosedWorkflow workflow started", "keepOpen", keepOpen)

	var result string
	err := workflow.ExecuteActivity(ctx, Activity, keepOpen).Get(ctx, &result)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return "", err
	}

	logger.Info("OpenNClosedWorkflow workflow completed.", "result", result)

	return result, nil
}

func Activity(ctx context.Context, keepOpen bool) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("OpenNClosedActivity", "keepOpen", keepOpen)
	if keepOpen {
		time.Sleep(10 * time.Minute)
	}
	return "Hello!", nil
}
