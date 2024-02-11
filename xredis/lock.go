package xredis

import (
	"context"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"time"
)

// lock

type Lock struct {
	name  string
	mutex *redsync.Mutex
	opt   *LockOption
}

type LockOption struct {
	Expiry time.Duration

	Tries     int
	DelayTime time.Duration
}

func (t *Client) NewLock(ctx context.Context, name string, opt *LockOption) (l *Lock) {
	if opt == nil {
		opt = &LockOption{}
	}

	rs := redsync.New(goredis.NewPool(t.client))
	mutexOptList := make([]redsync.Option, 0)
	if opt.Expiry > 0 {
		mutexOptList = append(mutexOptList, redsync.WithExpiry(opt.Expiry))
	}
	if opt.Tries > 0 {
		mutexOptList = append(mutexOptList, redsync.WithTries(opt.Tries))
	}
	if opt.DelayTime > 0 {
		mutexOptList = append(mutexOptList, redsync.WithRetryDelay(opt.DelayTime))
	}
	mutex := rs.NewMutex(name, mutexOptList...)
	l = &Lock{
		name:  name,
		mutex: mutex,
		opt:   opt,
	}
	return
}

func (l *Lock) Lock(ctx context.Context) error {
	return l.mutex.LockContext(ctx)
}

func (l *Lock) UnLock(ctx context.Context) (bool, error) {
	return l.mutex.UnlockContext(ctx)
}
