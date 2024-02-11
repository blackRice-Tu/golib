package xtask

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	logger "github.com/blackRice-Tu/golib/xlogger/default"
	"github.com/blackRice-Tu/golib/xredis"

	"github.com/golang-module/carbon"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

const (
	defaultQueueName = "default"

	redisQueueKey     = "asynq:queues"
	redisWorkerKey    = "asynq:workers"
	redisServerKey    = "asynq:servers"
	redisSchedulerKey = "asynq:schedulers"

	redisProcessedKeyFmt = "asynq:{%s}:processed:*"
	processedBackupDays  = 2
)

type (
	RedisOpt     = xredis.Config
	ServerOpt    = asynq.Config
	SchedulerOpt = asynq.SchedulerOpts
	Task         = asynq.Task
	TaskInfo     = asynq.TaskInfo
	Option       = asynq.Option
	Handler      = func(context.Context, *Task) error
)

func getRedisOpt(redisOpt *RedisOpt) asynq.RedisClientOpt {
	if redisOpt == nil {
		return asynq.RedisClientOpt{}
	}
	return asynq.RedisClientOpt{
		Addr:         redisOpt.Address,
		Password:     redisOpt.Password,
		DB:           redisOpt.DB,
		DialTimeout:  time.Duration(redisOpt.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(redisOpt.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(redisOpt.WriteTimeout) * time.Second,
		PoolSize:     redisOpt.PoolSize,
	}
}

func getServerOpt(name string, serverOpt *ServerOpt) *ServerOpt {
	if serverOpt == nil {
		serverOpt = &ServerOpt{}
	}
	// set queue name
	name = strings.TrimSpace(name)
	if name == "" {
		name = defaultQueueName
	}
	serverOpt.Queues = map[string]int{
		name: 1,
	}
	// add default BaseContext
	serverOpt.BaseContext = defaultBaseContext
	// add default ErrorHandler
	if serverOpt.ErrorHandler == nil {
		serverOpt.ErrorHandler = asynq.ErrorHandlerFunc(defaultErrorHandler)
	}

	return serverOpt
}

func getSchedulerOpt(schedulerOpt *SchedulerOpt) *SchedulerOpt {
	if schedulerOpt == nil {
		schedulerOpt = &SchedulerOpt{}
	}
	// add default PreEnqueueFunc
	if schedulerOpt.PreEnqueueFunc == nil {
		schedulerOpt.PreEnqueueFunc = defaultPreEnqueueFunc
	} else {
		f := schedulerOpt.PreEnqueueFunc
		schedulerOpt.PreEnqueueFunc = func(task *Task, opts []Option) {
			defaultPreEnqueueFunc(task, opts)
			f(task, opts)
		}
	}
	// add default PostEnqueueFunc
	if schedulerOpt.PostEnqueueFunc == nil {
		schedulerOpt.PostEnqueueFunc = defaultPostEnqueueFunc
	} else {
		f := schedulerOpt.PostEnqueueFunc
		schedulerOpt.PostEnqueueFunc = func(info *TaskInfo, err error) {
			defaultPostEnqueueFunc(info, err)
			f(info, err)
		}
	}
	return schedulerOpt
}

type Manager struct {
	queueName string

	redisOpt *RedisOpt

	server    *asynq.Server
	scheduler *asynq.Scheduler
	client    *asynq.Client
	mux       *asynq.ServeMux
}

var (
	manager  *Manager
	initOnce sync.Once
)

func Init(name string, redisOpt *RedisOpt, serverOpt *ServerOpt, schedulerOpt *SchedulerOpt) {
	if manager != nil {
		return
	}

	initOnce.Do(func() {
		opt := getRedisOpt(redisOpt)
		serverOpt = getServerOpt(name, serverOpt)
		schedulerOpt = getSchedulerOpt(schedulerOpt)

		server := asynq.NewServer(opt, *serverOpt)
		scheduler := asynq.NewScheduler(opt, schedulerOpt)
		client := asynq.NewClient(opt)
		mux := asynq.NewServeMux()

		m := &Manager{
			redisOpt:  redisOpt,
			queueName: name,
			server:    server,
			scheduler: scheduler,
			client:    client,
			mux:       mux,
		}
		manager = m
	})
}

func notInitError() error {
	return errors.Errorf("manager is nil, please Init it first")
}

func NewTask(name string, payload []byte, opts ...Option) *Task {
	return asynq.NewTask(name, payload, opts...)
}

func AddTask(name string, handler Handler) error {
	if manager == nil {
		return notInitError()
	}
	if handler == nil {
		return errors.Errorf("handler can not be nil")
	}
	manager.mux.HandleFunc(name, handler)
	return nil
}

func AddCronTask(ctx context.Context, handler Handler, cronSpec string, task *Task, opts ...Option) (entryId string, e error) {
	// add to server
	err := AddTask(task.Type(), handler)
	if err != nil {
		e = err
		return
	}
	defaultOptions := []Option{
		asynq.Queue(manager.queueName),
		asynq.Unique(time.Minute),
		asynq.MaxRetry(0),
	}
	// add to scheduler
	opts = append(defaultOptions, opts...)
	entryId, e = manager.scheduler.Register(cronSpec, task, opts...)

	msg := fmt.Sprintf("AddCronTask: %s, cron: %s, entryId: %s", task.Type(), cronSpec, entryId)
	if e == nil {
		logger.Infof(ctx, msg)
	} else {
		msg = fmt.Sprintf("%s, err: %s", msg, e.Error())
		logger.Errorf(ctx, msg)
	}
	return
}

func CallTask(ctx context.Context, task *Task, opts ...Option) (taskInfo *TaskInfo, e error) {
	if manager == nil {
		e = notInitError()
		return
	}
	opts = append(opts, asynq.Queue(manager.queueName))
	return manager.client.EnqueueContext(ctx, task, opts...)
}

func Run(ctx context.Context) (e error) {
	defer func() {
		if e != nil {
			logger.Errorf(ctx, e.Error())
		}
	}()

	if manager == nil {
		e = notInitError()
		return
	}

	go func() {
		err := manager.server.Run(manager.mux)
		if err != nil {
			e = err
		}
	}()

	go func() {
		err := manager.scheduler.Run()
		if err != nil {
			e = err
		}
	}()

	return
}

func Clean(ctx context.Context) {
	if manager == nil {
		return
	}
	if manager.queueName == defaultQueueName {
		return
	}
	hostname, err := os.Hostname()
	if err != nil {
		return
	}
	clean(ctx, hostname)
	return
}

func clean(ctx context.Context, hostname string) {
	// connect to redis
	rds, err := xredis.NewClient(manager.redisOpt, nil)
	if err != nil {
		logger.Errorf(ctx, errors.WithMessagef(err, "xredis.NewClient").Error())
		return
	}
	defer func() {
		_ = rds.Close()
	}()

	// clean zset redis keys
	var cleanRedisZsetKey = func(key string) {
		var match = fmt.Sprintf("%s:{%s:*}", key, hostname)
		keys, _ := rds.NativeZScanAll(ctx, key, match, 1)
		if len(keys) != 0 {
			_, _ = rds.NativeZRem(ctx, key, keys)
		}
	}
	// clean workers
	cleanRedisZsetKey(redisWorkerKey)
	// clean servers
	cleanRedisZsetKey(redisServerKey)
	// clean schedulers
	cleanRedisZsetKey(redisSchedulerKey)

	// clean processed key
	current := carbon.Now()
	dateStringList := make([]string, processedBackupDays)
	for i := 0; i < processedBackupDays; i += 1 {
		dateStringList[i] = current.SubDays(i).ToDateString()
	}
	key := fmt.Sprintf(redisProcessedKeyFmt, manager.queueName)
	keys, _ := rds.NativeScanAll(ctx, key, 1)
	for _, _key := range keys {
		items := strings.Split(_key, ":")
		dateString := items[len(items)-1]
		if !lo.Contains(dateStringList, dateString) {
			_, _ = rds.Del(ctx, _key)
		}
	}
}
