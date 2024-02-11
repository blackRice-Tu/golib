package xkafka

import (
	"fmt"
	"os"
	"testing"
)

const (
	Key       = "test"
	Topic     = "ops_runtime_event_report_dev_1"
	BrokerKey = "KAFKA_BROKER"
)

func TestLoadOrNewSyncProducer(t *testing.T) {
	p, err := LoadOrNewSyncProducer(Key, nil, nil)
	if p != nil || err == nil {
		t.Fatalf("p is not nil or err is nil")
	}

	broker := os.Getenv(BrokerKey)
	config := &Config{
		Broker: broker,
	}
	p, err = LoadOrNewSyncProducer(Key, config, nil)
	if err != nil {
		t.Fatalf("err %+v", err)
	}
	fmt.Println(p)
}

func TestSendMsg(t *testing.T) {
	p, _ := LoadOrNewSyncProducer(Key, nil, nil)
	partition, offset, err := p.SendMsg(Topic, []byte("111"), nil)
	if err != nil {
		t.Fatalf("err %+v", err)
	}
	fmt.Println(partition, offset, err)
}
