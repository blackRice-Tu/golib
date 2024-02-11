package xlogger

import (
	"encoding/json"
	"fmt"
	"github.com/blackRice-Tu/golib/xkafka"

	"go.uber.org/zap/zapcore"
)

type KafkaHookConfig struct {
	Id    string `yaml:"id" json:"id"`
	Topic string `yaml:"topic" json:"topic"`
	Key   string `json:"key" json:"key"`
}

func NewKafkaHook(cfg *KafkaHookConfig) func(zapcore.Entry) error {
	return func(entry zapcore.Entry) error {
		if cfg == nil {
			return nil
		}
		producer, err := xkafka.LoadOrNewAsyncProducer(cfg.Id, nil, nil)
		if err != nil {
			return err
		}

		body := make(map[string]any)
		body["level"] = entry.Level
		body["time"] = entry.Time.Format("2006-01-02 15:04:05")
		body["msg"] = entry.Message
		body["caller"] = fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)

		msg, _ := json.Marshal(body)
		err = producer.Send(cfg.Topic, msg, &cfg.Key)
		if err != nil {
			return err
		}
		return nil
	}
}
