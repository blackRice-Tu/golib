package xlogger

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/blackRice-Tu/golib"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	Id            string `yaml:"id" json:"id"`
	Level         string `yaml:"level" json:"level"`
	FlushInterval int    `yaml:"flushInterval" json:"flush_interval"`
	Path          string `yaml:"path" json:"path"`
	Name          string `yaml:"name" json:"name"`
	Size          int    `yaml:"size" json:"size"`
	Age           int    `yaml:"age" json:"age"`
	Backups       int    `yaml:"backups" json:"backups"`
	AddCaller     bool   `yaml:"addCaller" json:"add_caller"`
	AddStacktrace bool   `yaml:"addStacktrace" json:"add_stacktrace"`

	KafkaHook *KafkaHookConfig `yaml:"kafkaHook" json:"kafka_hook"`
}

const (
	maxFileSize = 512
	maxAge      = 7
	maxBackups  = 10
)

var (
	loggerInstanceMap sync.Map
)

var levelMap = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
	"panic": zapcore.PanicLevel,
	"fatal": zapcore.FatalLevel,
}

func getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:    "msg",
		LevelKey:      "level",
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		TimeKey:       "time",
		StacktraceKey: "stacktrace",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		CallerKey:    "caller",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})
}

func NewLogger(conf *Config) *zap.Logger {
	if conf == nil {
		return nil
	}
	writer, err := getWriter(conf)
	if err != nil {
		err = errors.Wrapf(err, "getWriter")
		logger := golib.GetStdLogger()
		logger.Println(err.Error())
		return nil
	}
	w := &zapcore.BufferedWriteSyncer{
		WS:            zapcore.AddSync(writer),
		FlushInterval: time.Duration(conf.FlushInterval) * time.Second,
	}

	l, ok := levelMap[conf.Level]
	if !ok {
		l = zap.InfoLevel
	}

	coreList := make([]zapcore.Core, 0)
	// add file logger
	fileCore := zapcore.NewCore(getEncoder(), w, l)
	coreList = append(coreList, fileCore)
	// add console logger
	if golib.GetMode() == golib.DebugMode {
		consoleCore := zapcore.NewCore(zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()), zapcore.Lock(os.Stdout), l)
		coreList = append(coreList, consoleCore)
	}
	core := zapcore.NewTee(coreList...)
	if conf.KafkaHook != nil { // add kafka logger
		hook := NewKafkaHook(conf.KafkaHook)
		core = zapcore.RegisterHooks(core, hook)
	}
	options := make([]zap.Option, 0)
	if conf.AddCaller {
		options = append(options, zap.AddCaller(), zap.AddCallerSkip(1))
	}
	if conf.AddStacktrace {
		options = append(options, zap.AddStacktrace(zap.ErrorLevel))
	}

	return zap.New(core, options...)
}

// LoadOrNewLogger ...
func LoadOrNewLogger(id string, conf *Config) (*zap.Logger, error) {
	instance, ok := loggerInstanceMap.Load(id)
	if ok {
		return instance.(*zap.Logger), nil
	}
	logger := NewLogger(conf)
	loggerInstanceMap.Store(id, logger)
	return logger, nil
}

func getWriter(conf *Config) (io.Writer, error) {
	if _, err := os.Stat(conf.Path); os.IsNotExist(err) {
		_ = os.MkdirAll(conf.Path, os.ModePerm)
	}

	name := conf.Name
	rawPath := strings.TrimSuffix(conf.Path, "/")
	fileName := filepath.Join(rawPath, name)

	size := conf.Size
	if size == 0 {
		size = maxFileSize
	}
	age := conf.Age
	if age == 0 {
		age = maxAge
	}
	backups := conf.Backups
	if backups == 0 {
		backups = maxBackups
	}

	// log rotate
	writer := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    size,
		MaxBackups: backups,
		MaxAge:     age,
		LocalTime:  true,
	}

	return writer, nil
}

type CommonField struct {
	ServerIp   string `json:"server_ip"`
	Env        string `json:"env"`
	ProgramKey string `json:"program_key"`
	TraceId    string `json:"trace_id"`
}

func GetCommonField(traceId string) CommonField {
	field := CommonField{
		ServerIp:   golib.GetServerIp(),
		Env:        golib.GetEnv(),
		ProgramKey: golib.GetProgramKey(),
		TraceId:    traceId,
	}
	return field
}
