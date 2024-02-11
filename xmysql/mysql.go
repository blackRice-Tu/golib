package xmysql

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/blackRice-Tu/golib/utils/xcommon"
	"github.com/blackRice-Tu/golib/utils/xcontext"
	"github.com/blackRice-Tu/golib/xlogger"

	logger "github.com/blackRice-Tu/golib/xlogger/default"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gLogger "gorm.io/gorm/logger"
	gUtils "gorm.io/gorm/utils"
)

type (
	Client = gorm.DB
)

const (
	defaultTimeZone = "Asia%2FShanghai"
	slowThreshold   = 200 * time.Millisecond

	OperationSelect = "SELECT"
	OperationUpdate = "UPDATE"
	OperationDelete = "DELETE"
	OperationInsert = "INSERT"
)

var (
	instanceMap sync.Map
)

type Config struct {
	Id       string `yaml:"id" json:"id"`
	Name     string `yaml:"name" json:"name"`
	User     string `yaml:"user" json:"user"`
	Password string `yaml:"password" json:"password"`
	Host     string `yaml:"host" json:"host"`
	Port     string `yaml:"port" json:"port"`
	TimeZone string `yaml:"timeZone"`
	Debug    bool   `yaml:"debug" json:"debug"`
}

// Logger ...
type Logger struct {
	Log *zap.Logger

	gLogger.Config

	Db string
}

// MsgField ...
type MsgField struct {
	xlogger.CommonField

	Db        string  `json:"db"`
	Operation string  `json:"operation"`
	Sql       string  `json:"sql"`
	Rows      int64   `json:"rows"`
	File      string  `json:"file"`
	IsSlow    int     `json:"is_slow"`
	Msg       string  `json:"msg"`
	Latency   float64 `json:"latency"`
}

// NewClient ...
func NewClient(dbConfig *Config, loggerConfig *xlogger.Config) (*Client, error) {
	if dbConfig == nil {
		e := errors.Errorf("dbConfig can not be nil")
		logger.Error(context.TODO(), e.Error())
		return nil, e
	}
	timeZone := dbConfig.TimeZone
	if timeZone == "" {
		timeZone = defaultTimeZone
	}
	path := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=%s",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name, timeZone)
	dsn := mysql.Open(path)

	cfg := gorm.Config{
		Logger: gLogger.Default,
	}
	l := NewLogger(dbConfig, loggerConfig)
	if l != nil {
		cfg.Logger = l
	}

	DB, _ := gorm.Open(dsn, &cfg)
	return DB, nil
}

// LoadOrNewClient ...
func LoadOrNewClient(id string, dbConfig *Config, loggerConfig *xlogger.Config) (*Client, error) {
	instance, ok := instanceMap.Load(id)
	if ok {
		return instance.(*Client), nil
	}
	client, err := NewClient(dbConfig, loggerConfig)
	if err != nil {
		return nil, err
	}
	instanceMap.Store(id, client)
	return client, nil
}

func NewLogger(dbConfig *Config, loggerConfig *xlogger.Config) *Logger {
	xlog := xlogger.NewLogger(loggerConfig)
	if xlog == nil {
		return nil
	}
	l := Logger{
		Log: xlog,
		Config: gLogger.Config{
			SlowThreshold:             slowThreshold,
			Colorful:                  true,
			IgnoreRecordNotFoundError: true,
			LogLevel:                  gLogger.Info,
		},
		Db: dbConfig.Id,
	}
	return &l
}

func (l *Logger) LogMode(level gLogger.LogLevel) gLogger.Interface {
	newLogger := *l
	l.LogLevel = level
	return &newLogger
}

func (l *Logger) Info(ctx context.Context, msg string, vars ...any) {
	if l.LogLevel >= gLogger.Info {
		l.Log.Info(fmt.Sprintf(msg, vars...))
	}
}

func (l *Logger) Warn(ctx context.Context, msg string, vars ...any) {
	if l.LogLevel >= gLogger.Warn {
		l.Log.Warn(fmt.Sprintf(msg, vars...))
	}
}

func (l *Logger) Error(ctx context.Context, msg string, vars ...any) {
	if l.LogLevel >= gLogger.Error {
		l.Log.Error(fmt.Sprintf(msg, vars...))
	}
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= gLogger.Silent {
		return
	}

	traceId := xcontext.GetOrNewTraceId(ctx)
	sql, rows := fc()
	elapsed := time.Since(begin)
	latency := float64(elapsed.Nanoseconds()) / 1e6

	operation := ""
	for _, op := range []string{OperationSelect, OperationDelete, OperationInsert, OperationUpdate} {
		if strings.HasPrefix(sql, op) {
			operation = op
			break
		}
	}

	msg := struct {
		xlogger.CommonField

		Db        string  `json:"db"`
		Operation string  `json:"operation"`
		Sql       string  `json:"sql"`
		Rows      int64   `json:"rows"`
		File      string  `json:"file"`
		IsSlow    int     `json:"is_slow"`
		Msg       string  `json:"msg"`
		Latency   float64 `json:"latency"`
	}{
		CommonField: xlogger.GetCommonField(traceId),
		Db:          l.Db,
		Operation:   operation,
		Sql:         sql,
		Rows:        rows,
		File:        gUtils.FileWithLineNum(),
		IsSlow:      0,
		Msg:         "",
		Latency:     latency,
	}

	switch {
	case err != nil && l.LogLevel >= gLogger.Error && (!errors.Is(err, gLogger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		msg.IsSlow = 0
		msg.Msg = err.Error()
		l.Error(ctx, xcommon.JsonMarshal(msg))
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gLogger.Warn:
		msg.IsSlow = 1
		msg.Msg = fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		l.Warn(ctx, xcommon.JsonMarshal(msg))
	case l.LogLevel == gLogger.Info:
		l.Info(ctx, xcommon.JsonMarshal(msg))
	}
}

func Clean() {
	instanceMap.Range(func(key, value any) bool {
		client := value.(*Client)
		mdb, _ := client.DB()
		mdb.Close()
		instanceMap.Delete(key)
		return true
	})
}
