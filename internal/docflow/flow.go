package docflow

import (
	"doctravel/internal/docactivity"
	"go.temporal.io/sdk/workflow"
	"time"
	"unicode/utf8"
)

type ApproveSignal struct {
	Approved string
}

type EditSignal struct {
	Editing string
}

func DocProcessingFlow(ctx workflow.Context, action docactivity.DocAction, param string) (string, error) {
	logger := workflow.GetLogger(ctx)
	var resPtr string
	// Define the Activity Execution options
	// StartToCloseTimeout or ScheduleToCloseTimeout must be set
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)
	// Реєстрація сигналів (каналів)
	signChan := workflow.GetSignalChannel(ctx, ApproveSignalName)
	editChan := workflow.GetSignalChannel(ctx, EditingSignalName)
	finish := false

	// реєстрація запитів
	queryType := InfoQuery
	err := workflow.SetQueryHandler(ctx, queryType, func(prefix string, suffix string) (string, error) {
		return prefix + resPtr + suffix, nil
	})
	if err != nil {
		logger.Info("SetQueryHandler failed.", "Error", err)
		return "", err
	}

	for {
		selector := workflow.NewSelector(ctx)

		// Сигнал погодження
		selector.AddReceive(signChan, func(c workflow.ReceiveChannel, _ bool) {
			var signal ApproveSignal
			c.Receive(ctx, &signal)

			if utf8.RuneCountInString(signal.Approved) == 0 {
				return
			}
			if err := workflow.ExecuteActivity(ctx, action.ApproveActivity, param+signal.Approved).Get(ctx, &resPtr); err != nil {
				logger.Error(err.Error())
			}
			// Sleep 5 seconds
			workflow.Sleep(ctx, time.Second*5)
			finish = true
		})
		// Сигнал Початку редагування
		selector.AddReceive(editChan, func(c workflow.ReceiveChannel, _ bool) {
			var signal EditSignal
			editChan.Receive(ctx, &signal)
			if utf8.RuneCountInString(signal.Editing) == 0 {
				return
			}
			if err := workflow.ExecuteActivity(ctx, action.StartEditingActivity, param+signal.Editing).Get(ctx, &resPtr); err != nil {
				logger.Error(err.Error())
			}
			// Sleep 5 seconds
			workflow.Sleep(ctx, time.Second*5)
		})

		// Can also use `Receive()` instead of a selector, but we'll be making further
		// use of selectors in part 2 of this series.
		selector.Select(ctx)
		if finish {
			break
		}
	}
	return resPtr, nil
}
