package xfreecache

import (
	"encoding/json"
	"fmt"
	"time"
)

type Group struct {
	name        string
	expiredTime time.Duration
}

func NewGroup(name string, expiredTime time.Duration) *Group {
	group := &Group{
		name:        name,
		expiredTime: expiredTime,
	}
	return group
}

func (t *Group) getKey(key string) []byte {
	return []byte(fmt.Sprintf("<%s>%s", t.name, key))
}

func (t *Group) Set(key string, value any) error {
	k := t.getKey(key)
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return globalCache.Set(k, v, int(t.expiredTime.Seconds()))
}

func (t *Group) Get(key string, value any) (bool, error) {
	k := t.getKey(key)
	v, err := globalCache.Get(k)
	if err != nil {
		return false, err
	}
	return true, json.Unmarshal(v, &value)
}
