package golib

import (
	"sync"

	"github.com/blackRice-Tu/golib/utils/xsys"
)

type Mode string

type Config struct {
	mode       Mode
	logPrefix  string
	serverIp   string
	env        string
	programKey string
	traceIdKey string
	debugKey   string
	timeFmt    string
}

const (
	DebugMode   Mode = "debug"
	ReleaseMode Mode = "release"
)

var (
	config  *Config
	setOnce sync.Once
)

func init() {
	if config != nil {
		return
	}
	serverIp, _ := xsys.GetHostIpv4()
	cfg := Config{
		mode:       DebugMode,
		logPrefix:  "[golib] ",
		serverIp:   serverIp,
		env:        "",
		programKey: "DEFAULT",
		traceIdKey: "X-Trace-Id",
		debugKey:   "X-Debug",
		timeFmt:    "2006-01-02 15:04:05",
	}
	config = &cfg

	// logger
	stdLogger = NewStdLogger(config.logPrefix)
}

type Option struct {
	Mode Mode

	LogPrefix  *string
	ServerIp   *string
	Env        *string
	ProgramKey *string
	TraceIdKey *string
	DebugKey   *string
	TimeFmt    *string
}

func SetConfig(opt *Option) {
	setOnce.Do(func() {
		config.mode = opt.Mode
		if opt.LogPrefix != nil {
			config.logPrefix = *opt.LogPrefix
			stdLogger.SetPrefix(config.logPrefix)
		}
		if opt.ServerIp != nil {
			config.serverIp = *opt.ServerIp
		}
		if opt.Env != nil {
			config.env = *opt.Env
		}
		if opt.ProgramKey != nil {
			config.programKey = *opt.ProgramKey
		}
		if opt.TraceIdKey != nil {
			config.traceIdKey = *opt.TraceIdKey
		}
		if opt.DebugKey != nil {
			config.debugKey = *opt.DebugKey
		}
		if opt.TimeFmt != nil {
			config.timeFmt = *opt.TimeFmt
		}
	})
}

func GetMode() Mode {
	return config.mode
}

func GetServerIp() string {
	return config.serverIp
}

func GetEnv() string {
	return config.env
}

func GetProgramKey() string {
	return config.programKey
}

func GetTraceIdKey() string {
	return config.traceIdKey
}

func GetDebugKey() string {
	return config.debugKey
}

func GetTimeFmt() string {
	return config.timeFmt
}
