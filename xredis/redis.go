package xredis

import (
	"context"
	"runtime"
	"sync"
	"time"

	logger "github.com/blackRice-Tu/golib/xlogger/default"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

const (
	dialTimeout          = 2
	readTimeout          = 1
	writeTimeout         = 1
	pollSizeMultiple     = 20
	minIdleConnsMultiple = 10
	connMaxIdleTime      = 3
)

var (
	instanceMap sync.Map
)

// NewClient ...
func NewClient(config *Config, prefix *string) (*Client, error) {
	if config == nil {
		e := errors.Errorf("config can not be nil")
		logger.Error(context.TODO(), e.Error())
		return nil, e
	}

	opt := &redis.Options{
		Addr:     config.Address,
		Password: config.Password,
		DB:       config.DB,
	}

	opt.DialTimeout = time.Duration(dialTimeout) * time.Second
	opt.ReadTimeout = time.Duration(readTimeout) * time.Second
	opt.WriteTimeout = time.Duration(writeTimeout) * time.Second
	opt.PoolSize = pollSizeMultiple * runtime.NumCPU()
	opt.MinIdleConns = minIdleConnsMultiple * runtime.NumCPU()
	opt.ConnMaxIdleTime = time.Duration(connMaxIdleTime) * time.Minute

	if config.DialTimeout > 0 {
		opt.DialTimeout = time.Duration(config.DialTimeout) * time.Second
	}
	if config.ReadTimeout > 0 {
		opt.ReadTimeout = time.Duration(config.ReadTimeout) * time.Second
	}
	if config.WriteTimeout > 0 {
		opt.WriteTimeout = time.Duration(config.WriteTimeout) * time.Second
	}
	if config.PoolSize > 0 {
		opt.PoolSize = config.PoolSize
	}
	if config.MinIdleConns > 0 {
		opt.MinIdleConns = config.MinIdleConns
	}
	if config.MinIdleConns > 0 {
		opt.ConnMaxIdleTime = time.Duration(config.ConnMaxIdleTime) * time.Second
	}

	client := &Client{
		client: redis.NewClient(opt),
	}
	if prefix != nil {
		client.prefix = prefix
	}
	return client, nil
}

func LoadOrNewClient(id string, config *Config, prefix *string) (*Client, error) {
	instance, ok := instanceMap.Load(id)
	if ok {
		return instance.(*Client), nil
	}
	client, err := NewClient(config, prefix)
	if err != nil {
		return nil, err
	}
	instanceMap.Store(id, client)
	return client, nil
}

func Clean() {
	instanceMap.Range(func(key, value any) bool {
		client := value.(*Client)
		err := client.client.Close()
		if err != nil {
			logger.Errorf(context.TODO(), err.Error())
		}
		instanceMap.Delete(key)
		return true
	})
}
