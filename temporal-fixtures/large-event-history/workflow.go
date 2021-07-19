package largeeventhistory

import (
	"errors"
	"time"

	"context"

	"go.temporal.io/sdk/workflow"
)

// LargeEventHistoryWorkflow workflow definition
func LargeEventHistoryWorkflow(ctx workflow.Context, LengthOfHistory int, WillFailOrNot bool) (err error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var data []byte
	i := 1

	countActivities := LengthOfHistory / 6
	for i <= countActivities {
		err = workflow.ExecuteActivity(ctx, Activity, data).Get(ctx, nil)
		i++
	}
	if err != nil {
		return errors.New("unexpected Activity failure")
	}

	if WillFailOrNot {
		return errors.New("intentional workflow failure due to WillFailOrNot parameter")
	}
	return nil
}

func Activity(ctx context.Context) error {
	return nil
}
