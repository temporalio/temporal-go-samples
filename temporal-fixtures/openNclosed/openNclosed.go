package openNclosed

import (
	"context"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"
)

// OpenClosedFixtureWorkflow is a basic Hello World workflow definition.
func OpenClosedFixtureWorkflow(ctx workflow.Context, name string, keep bool) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("OpenClosedFixtureWorkflow started", "name", name)

	var result string
	err := workflow.ExecuteActivity(ctx, OpenClosedFixtureActivity, name, keep).Get(ctx, &result)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return "", err
	}

	logger.Info("OpenClosedFixtureWorkflow completed.", "result", result)

	return result, nil
}

func OpenClosedFixtureActivity(ctx context.Context, name string, keep bool) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Activity", "name", name)
	if keep {
		time.Sleep(10 * time.Minute)
	}
	return "Hello " + name + "!", nil
}
