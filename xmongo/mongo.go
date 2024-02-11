package xmongo

import (
	"context"
	"fmt"
	"sync"

	logger "github.com/blackRice-Tu/golib/xlogger/default"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	instanceMap sync.Map
)

type Client struct {
	Cli *mongo.Client
	Db  *mongo.Database
}

type Config struct {
	Id           string `yaml:"id" json:"id"`
	Database     string `yaml:"database" json:"database"`
	User         string `yaml:"user" json:"user"`
	Password     string `yaml:"password" json:"password"`
	Host         string `yaml:"host" json:"host"`
	Port         string `yaml:"port" json:"port"`
	AuthDatabase string `yaml:"authDatabase" json:"auth_database"`
}

func NewClient(conf *Config) (client *Client, e error) {
	ctx := context.TODO()
	defer func() {
		if e != nil {
			logger.Error(ctx, e.Error())
		}
	}()

	cred := options.Credential{
		AuthSource:  conf.AuthDatabase,
		Username:    conf.User,
		Password:    conf.Password,
		PasswordSet: true,
	}
	opt := &options.ClientOptions{
		Auth:  &cred,
		Hosts: []string{fmt.Sprintf("%s:%s", conf.Host, conf.Port)},
	}
	cli, err := mongo.Connect(ctx, opt)
	if err != nil {
		e = errors.WithMessagef(err, "Connect")
		return
	}
	client = &Client{
		Cli: cli,
		Db:  cli.Database(conf.Database),
	}
	return
}

func LoadOrNewClient(id string, config *Config) (*Client, error) {
	instance, ok := instanceMap.Load(id)
	if ok {
		return instance.(*Client), nil
	}
	client, err := NewClient(config)
	if err != nil {
		return nil, err
	}
	instanceMap.Store(id, client)
	return client, nil
}

func Clean() {
	ctx := context.TODO()

	instanceMap.Range(func(key, value any) bool {
		client := value.(*Client)
		err := client.Cli.Disconnect(ctx)
		if err != nil {
			logger.Errorf(ctx, err.Error())
		}
		instanceMap.Delete(key)
		return true
	})
}
