package xredis

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/ulule/limiter/v3"
	storeRedis "github.com/ulule/limiter/v3/drivers/store/redis"
)

// rate limiter

type RateLimiter struct {
	Limiter *limiter.Limiter
	opt     *RateLimiterOption
}

type RateLimiterOption struct {
	Block    bool
	Interval time.Duration
	Timeout  time.Duration
}

func (t *Client) NewRateLimiter(ctx context.Context, formatted string, opt *RateLimiterOption) (l *RateLimiter, e error) {
	rate, err := limiter.NewRateFromFormatted(formatted)
	if err != nil {
		e = errors.WithMessagef(err, "limiter.NewRateFromFormatted")
		return
	}
	store, err := storeRedis.NewStoreWithOptions(t.client, limiter.StoreOptions{
		Prefix: fmt.Sprintf("%s:limiter", *t.prefix),
	})
	if err != nil {
		e = errors.WithMessagef(err, "storeRedis.NewStoreWithOptions")
		return
	}
	// init opt
	if opt == nil {
		opt = &RateLimiterOption{}
	}
	if opt.Block && opt.Interval == 0 {
		opt.Interval = time.Second
	}
	if opt.Block && opt.Interval > 0 && opt.Timeout == 0 {
		opt.Timeout = 10 * opt.Interval
	}

	l = &RateLimiter{
		Limiter: limiter.New(store, rate),
		opt:     opt,
	}
	return
}

func (l *RateLimiter) Get(ctx context.Context, key string) (ok bool, e error) {
	result, err := l.Limiter.Get(ctx, key)
	if err != nil {
		e = errors.WithMessagef(err, "Get")
		return
	}
	ok = !result.Reached
	if !ok && l.opt.Block {
		ticker := time.NewTicker(l.opt.Interval)
		defer ticker.Stop()

		done := make(chan struct{}, 1)
		go func() {
			time.Sleep(l.opt.Timeout)
			done <- struct{}{}
		}()

		for {
			select {
			case <-ticker.C:
				result, err = l.Limiter.Get(ctx, key)
				if err != nil {
					e = errors.WithMessagef(err, "Get")
					return
				}
				ok = !result.Reached
				if ok {
					return
				}
			case <-done:
				ok = false
				return
			}
		}
	}
	return
}
