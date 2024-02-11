package xtask

import (
	"context"
	"fmt"
	"github.com/blackRice-Tu/golib/utils/xcontext"
	"time"

	logger "github.com/blackRice-Tu/golib/xlogger/default"

	"github.com/hibiken/asynq"
)

func MaxRetry(n int) Option {
	return asynq.MaxRetry(n)
}

func Queue(name string) Option {
	return asynq.Queue(name)
}

func TaskId(id string) Option {
	return asynq.TaskID(id)
}

func Timeout(d time.Duration) Option {
	return asynq.Timeout(d)
}

func Deadline(t time.Time) Option {
	return asynq.Deadline(t)
}

func Unique(ttl time.Duration) Option {
	return asynq.Unique(ttl)
}

func ProcessAt(t time.Time) Option {
	return asynq.ProcessAt(t)
}

func ProcessIn(d time.Duration) Option {
	return asynq.ProcessIn(d)
}

//func Retention(d time.Duration) Option {
//	return asynq.Retention(d)
//}

func Group(name string) Option {
	return asynq.Group(name)
}

func defaultPreEnqueueFunc(task *Task, opts []Option) {
	if task == nil {
		return
	}
}

func defaultPostEnqueueFunc(info *TaskInfo, err error) {
	if info == nil {
		return
	}
	ctx := context.TODO()
	msg := fmt.Sprintf("ScheduleTaskDone: %s", info.Type)
	if err == nil {
		logger.Infof(ctx, msg)
	} else {
		msg = fmt.Sprintf("%s, err: %s", msg, err.Error())
		logger.Errorf(ctx, msg)
	}
}

func defaultErrorHandler(ctx context.Context, task *Task, err error) {
	if err != nil {
		logger.Errorf(ctx, "TaskError: %s, %s", task.Type(), err.Error())
	}
}

func defaultBaseContext() context.Context {
	ctx, _ := xcontext.NewTraceContext()
	return ctx
}
