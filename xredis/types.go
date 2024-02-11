package xredis

import (
	"github.com/redis/go-redis/v9"
)

type (
	Z        = redis.Z
	ZWithKey = redis.ZWithKey
	ZStore   = redis.ZStore
	ZAddArgs = redis.ZAddArgs
	ZRangeBy = redis.ZRangeBy
)

type Client struct {
	prefix *string
	client *redis.Client
}

type Config struct {
	Id              string `yaml:"id" json:"id"`
	Address         string `yaml:"address" json:"address"`
	Password        string `yaml:"password" json:"password"`
	DB              int    `yaml:"db" json:"db"`
	DialTimeout     int    `yaml:"dialTimeout" json:"dial_timeout"`
	ReadTimeout     int    `yaml:"readTimeout" json:"read_timeout"`
	WriteTimeout    int    `yaml:"writeTimeout" json:"write_timeout"`
	PoolSize        int    `yaml:"poolSize" json:"pool_size"`
	MinIdleConns    int    `yaml:"minIdleConns" json:"min_idle_conns"`
	ConnMaxIdleTime int    `yaml:"connMaxIdleTime" json:"conn_max_idle_time"`
}
