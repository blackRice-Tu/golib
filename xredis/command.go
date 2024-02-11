package xredis

import (
	"context"
	"fmt"
	"strings"
	"time"
)

func (t *Client) getPrefix() string {
	prefix := ""
	if t.prefix != nil {
		prefix = *t.prefix
	}
	return prefix
}

func (t *Client) GetKey(key string) string {
	return fmt.Sprintf("%s%s", t.getPrefix(), key)
}

func (t *Client) GetKeys(keys []string) []string {
	keyList := make([]string, len(keys))
	for i, key := range keys {
		keyList[i] = t.GetKey(key)
	}
	return keyList
}

func (t *Client) Close() error {
	return t.client.Close()
}

// command

func (t *Client) Del(ctx context.Context, keys ...string) (int64, error) {
	return t.client.Del(ctx, t.GetKeys(keys)...).Result()
}

func (t *Client) Exists(ctx context.Context, keys ...string) (int64, error) {
	return t.client.Exists(ctx, t.GetKeys(keys)...).Result()
}

func (t *Client) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return t.client.Expire(ctx, t.GetKey(key), expiration).Result()
}

func (t *Client) ExpireNX(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return t.client.ExpireNX(ctx, t.GetKey(key), expiration).Result()
}

func (t *Client) ExpireXX(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return t.client.ExpireXX(ctx, t.GetKey(key), expiration).Result()
}

func (t *Client) ExpireGT(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return t.client.ExpireGT(ctx, t.GetKey(key), expiration).Result()
}

func (t *Client) ExpireLT(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return t.client.ExpireLT(ctx, t.GetKey(key), expiration).Result()
}

func (t *Client) ExpireAt(ctx context.Context, key string, tm time.Time) (bool, error) {
	return t.client.ExpireAt(ctx, t.GetKey(key), tm).Result()
}

func (t *Client) Keys(ctx context.Context, pattern string) ([]string, error) {
	keys, err := t.client.Keys(ctx, t.GetKey(pattern)).Result()
	if err != nil {
		return nil, err
	}
	keysWithoutPrefix := make([]string, len(keys))
	prefix := t.getPrefix()
	for i, _ := range keys {
		keysWithoutPrefix[i] = strings.TrimPrefix(keys[i], prefix)
	}
	return keysWithoutPrefix, nil
}

func (t *Client) TTL(ctx context.Context, key string) (time.Duration, error) {
	return t.client.TTL(ctx, t.GetKey(key)).Result()
}

func (t *Client) Type(ctx context.Context, key string) (string, error) {
	return t.client.Type(ctx, t.GetKey(key)).Result()
}

func (t *Client) Decr(ctx context.Context, key string) (int64, error) {
	return t.client.Decr(ctx, t.GetKey(key)).Result()
}

func (t *Client) DecrBy(ctx context.Context, key string, decrement int64) (int64, error) {
	return t.client.DecrBy(ctx, t.GetKey(key), decrement).Result()
}

func (t *Client) Get(ctx context.Context, key string) (string, error) {
	return t.client.Get(ctx, t.GetKey(key)).Result()
}

func (t *Client) Incr(ctx context.Context, key string) (int64, error) {
	return t.client.Incr(ctx, t.GetKey(key)).Result()
}

func (t *Client) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return t.client.IncrBy(ctx, t.GetKey(key), value).Result()
}

func (t *Client) IncrByFloat(ctx context.Context, key string, value float64) (float64, error) {
	return t.client.IncrByFloat(ctx, t.GetKey(key), value).Result()
}

func (t *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error) {
	return t.client.Set(ctx, t.GetKey(key), value, expiration).Result()
}

func (t *Client) SetEx(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error) {
	return t.client.SetEx(ctx, t.GetKey(key), value, expiration).Result()
}

func (t *Client) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return t.client.SetNX(ctx, t.GetKey(key), value, expiration).Result()
}

func (t *Client) SetXX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return t.client.SetXX(ctx, t.GetKey(key), value, expiration).Result()
}

// ------

func (t *Client) Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return t.client.Scan(ctx, cursor, t.GetKey(match), count).Result()
}

func (t *Client) ScanAll(ctx context.Context, match string, count int64) ([]string, error) {
	var cursor uint64 = 0
	var keys []string
	var err error
	match = t.GetKey(match)
	for {
		_keys, _cursor, _err := t.client.Scan(ctx, cursor, match, count).Result()
		if _err != nil {
			err = _err
			break
		}
		if len(_keys) != 0 {
			keys = append(keys, _keys...)
		}
		if _cursor == 0 {
			break
		}
		cursor = _cursor
	}
	return keys, err
}

func (t *Client) ScanType(ctx context.Context, cursor uint64, match string, count int64, keyType string) ([]string, uint64, error) {
	return t.client.ScanType(ctx, cursor, t.GetKey(match), count, keyType).Result()
}

func (t *Client) SScan(ctx context.Context, key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return t.client.SScan(ctx, t.GetKey(key), cursor, match, count).Result()
}

func (t *Client) SScanAll(ctx context.Context, key string, match string, count int64) ([]string, error) {
	var cursor uint64 = 0
	var keys []string
	var err error
	key = t.GetKey(key)
	for {
		_keys, _cursor, _err := t.client.SScan(ctx, key, cursor, match, count).Result()
		if _err != nil {
			err = _err
			break
		}
		if len(_keys) != 0 {
			keys = append(keys, _keys...)
		}
		if _cursor == 0 {
			break
		}
		cursor = _cursor
	}
	return keys, err
}

func (t *Client) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return t.client.HScan(ctx, t.GetKey(key), cursor, match, count).Result()
}

func (t *Client) HScanAll(ctx context.Context, key string, match string, count int64) ([]string, error) {
	var cursor uint64 = 0
	var keys []string
	var err error
	key = t.GetKey(key)
	for {
		_keys, _cursor, _err := t.client.HScan(ctx, key, cursor, match, count).Result()
		if _err != nil {
			err = _err
			break
		}
		if len(_keys) != 0 {
			keys = append(keys, _keys...)
		}
		if _cursor == 0 {
			break
		}
		cursor = _cursor
	}
	return keys, err
}

func (t *Client) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return t.client.ZScan(ctx, t.GetKey(key), cursor, match, count).Result()
}

func (t *Client) ZScanAll(ctx context.Context, key string, match string, count int64) ([]string, error) {
	var cursor uint64 = 0
	var keys []string
	var err error
	key = t.GetKey(key)
	for {
		_keys, _cursor, _err := t.client.ZScan(ctx, key, cursor, match, count).Result()
		if _err != nil {
			err = _err
			break
		}
		if len(_keys) != 0 {
			keys = append(keys, _keys...)
		}
		if _cursor == 0 {
			break
		}
		cursor = _cursor
	}
	return keys, err
}

// ------

func (t *Client) HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	return t.client.HDel(ctx, t.GetKey(key), fields...).Result()
}

func (t *Client) HExists(ctx context.Context, key, field string) (bool, error) {
	return t.client.HExists(ctx, t.GetKey(key), field).Result()
}

func (t *Client) HGet(ctx context.Context, key, field string) (string, error) {
	return t.client.HGet(ctx, t.GetKey(key), field).Result()
}

func (t *Client) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return t.client.HGetAll(ctx, t.GetKey(key)).Result()
}

func (t *Client) HIncrBy(ctx context.Context, key, field string, incr int64) (int64, error) {
	return t.client.HIncrBy(ctx, t.GetKey(key), field, incr).Result()
}

func (t *Client) HIncrByFloat(ctx context.Context, key, field string, incr float64) (float64, error) {
	return t.client.HIncrByFloat(ctx, t.GetKey(key), field, incr).Result()
}

func (t *Client) HKeys(ctx context.Context, key string) ([]string, error) {
	return t.client.HKeys(ctx, t.GetKey(key)).Result()
}

func (t *Client) HLen(ctx context.Context, key string) (int64, error) {
	return t.client.HLen(ctx, t.GetKey(key)).Result()
}

func (t *Client) HSet(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return t.client.HSet(ctx, t.GetKey(key), values...).Result()
}

func (t *Client) HVals(ctx context.Context, key string) ([]string, error) {
	return t.client.HVals(ctx, t.GetKey(key)).Result()
}

func (t *Client) LLen(ctx context.Context, key string) (int64, error) {
	return t.client.LLen(ctx, t.GetKey(key)).Result()
}

func (t *Client) LPop(ctx context.Context, key string) (string, error) {
	return t.client.LPop(ctx, t.GetKey(key)).Result()
}

func (t *Client) LPopCount(ctx context.Context, key string, count int) ([]string, error) {
	return t.client.LPopCount(ctx, t.GetKey(key), count).Result()
}

func (t *Client) LPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return t.client.LPush(ctx, t.GetKey(key), values).Result()
}

func (t *Client) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return t.client.LRange(ctx, t.GetKey(key), start, stop).Result()
}

func (t *Client) LRem(ctx context.Context, key string, count int64, value interface{}) (int64, error) {
	return t.client.LRem(ctx, t.GetKey(key), count, value).Result()
}

func (t *Client) RPop(ctx context.Context, key string) (string, error) {
	return t.client.RPop(ctx, t.GetKey(key)).Result()
}

func (t *Client) RPopCount(ctx context.Context, key string, count int) ([]string, error) {
	return t.client.RPopCount(ctx, t.GetKey(key), count).Result()
}

func (t *Client) RPopLPush(ctx context.Context, source, destination string) (string, error) {
	return t.client.RPopLPush(ctx, t.GetKey(source), t.GetKey(destination)).Result()
}

func (t *Client) RPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return t.client.RPush(ctx, t.GetKey(key), values...).Result()
}

// ----

func (t *Client) SAdd(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return t.client.SAdd(ctx, t.GetKey(key), members...).Result()
}

func (t *Client) SCard(ctx context.Context, key string) (int64, error) {
	return t.client.SCard(ctx, t.GetKey(key)).Result()
}

func (t *Client) SDiff(ctx context.Context, keys ...string) ([]string, error) {
	return t.client.SDiff(ctx, t.GetKeys(keys)...).Result()
}

func (t *Client) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return t.client.SIsMember(ctx, t.GetKey(key), member).Result()
}

func (t *Client) SMembers(ctx context.Context, key string) ([]string, error) {
	return t.client.SMembers(ctx, t.GetKey(key)).Result()
}

func (t *Client) SPop(ctx context.Context, key string) (string, error) {
	return t.client.SPop(ctx, t.GetKey(key)).Result()
}

func (t *Client) SPopN(ctx context.Context, key string, count int64) ([]string, error) {
	return t.client.SPopN(ctx, t.GetKey(key), count).Result()
}

func (t *Client) SRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return t.client.SRem(ctx, t.GetKey(key), members...).Result()
}

func (t *Client) SUnion(ctx context.Context, keys ...string) ([]string, error) {
	return t.client.SUnion(ctx, t.GetKeys(keys)...).Result()
}

func (t *Client) SUnionStore(ctx context.Context, destination string, keys ...string) (int64, error) {
	return t.client.SUnionStore(ctx, destination, t.GetKeys(keys)...).Result()
}

// -----

func (t *Client) BZPopMax(ctx context.Context, timeout time.Duration, keys ...string) (*ZWithKey, error) {
	return t.client.BZPopMax(ctx, timeout, t.GetKeys(keys)...).Result()
}

func (t *Client) BZPopMin(ctx context.Context, timeout time.Duration, keys ...string) (*ZWithKey, error) {
	return t.client.BZPopMin(ctx, timeout, t.GetKeys(keys)...).Result()
}

func (t *Client) ZAddArgs(ctx context.Context, key string, args ZAddArgs) (int64, error) {
	return t.client.ZAddArgs(ctx, t.GetKey(key), args).Result()
}

func (t *Client) ZAddArgsIncr(ctx context.Context, key string, args ZAddArgs) (float64, error) {
	return t.client.ZAddArgsIncr(ctx, t.GetKey(key), args).Result()
}

func (t *Client) ZAdd(ctx context.Context, key string, members ...Z) (int64, error) {
	return t.client.ZAdd(ctx, t.GetKey(key), members...).Result()
}

func (t *Client) ZAddNX(ctx context.Context, key string, members ...Z) (int64, error) {
	return t.client.ZAddNX(ctx, t.GetKey(key), members...).Result()
}

func (t *Client) ZAddXX(ctx context.Context, key string, members ...Z) (int64, error) {
	return t.client.ZAddXX(ctx, t.GetKey(key), members...).Result()
}

func (t *Client) ZCard(ctx context.Context, key string) (int64, error) {
	return t.client.ZCard(ctx, t.GetKey(key)).Result()
}

func (t *Client) ZCount(ctx context.Context, key, min, max string) (int64, error) {
	return t.client.ZCount(ctx, t.GetKey(key), min, max).Result()
}

func (t *Client) ZLexCount(ctx context.Context, key, min, max string) (int64, error) {
	return t.client.ZLexCount(ctx, t.GetKey(key), min, max).Result()
}

func (t *Client) ZIncrBy(ctx context.Context, key string, increment float64, member string) (float64, error) {
	return t.client.ZIncrBy(ctx, t.GetKey(key), increment, member).Result()
}

func (t *Client) ZPopMax(ctx context.Context, key string, count ...int64) ([]Z, error) {
	return t.client.ZPopMax(ctx, t.GetKey(key), count...).Result()
}

func (t *Client) ZPopMin(ctx context.Context, key string, count ...int64) ([]Z, error) {
	return t.client.ZPopMin(ctx, t.GetKey(key), count...).Result()
}

func (t *Client) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return t.client.ZRange(ctx, t.GetKey(key), start, stop).Result()
}

func (t *Client) ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]Z, error) {
	return t.client.ZRangeWithScores(ctx, t.GetKey(key), start, stop).Result()
}

func (t *Client) ZRangeByScore(ctx context.Context, key string, opt *ZRangeBy) ([]string, error) {
	return t.client.ZRangeByScore(ctx, t.GetKey(key), opt).Result()
}

func (t *Client) ZRangeByLex(ctx context.Context, key string, opt *ZRangeBy) ([]string, error) {
	return t.client.ZRangeByLex(ctx, t.GetKey(key), opt).Result()
}

func (t *Client) ZRangeByScoreWithScores(ctx context.Context, key string, opt *ZRangeBy) ([]Z, error) {
	return t.client.ZRangeByScoreWithScores(ctx, t.GetKey(key), opt).Result()
}

func (t *Client) ZRank(ctx context.Context, key, member string) (int64, error) {
	return t.client.ZRank(ctx, t.GetKey(key), member).Result()
}

func (t *Client) ZRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return t.client.ZRem(ctx, t.GetKey(key), members...).Result()
}

func (t *Client) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) (int64, error) {
	return t.client.ZRemRangeByRank(ctx, t.GetKey(key), start, stop).Result()
}

func (t *Client) ZRemRangeByScore(ctx context.Context, key, min, max string) (int64, error) {
	return t.client.ZRemRangeByScore(ctx, t.GetKey(key), min, max).Result()
}

func (t *Client) ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return t.client.ZRevRange(ctx, t.GetKey(key), start, stop).Result()
}

func (t *Client) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) ([]Z, error) {
	return t.client.ZRevRangeWithScores(ctx, t.GetKey(key), start, stop).Result()
}

func (t *Client) ZRevRangeByScore(ctx context.Context, key string, opt *ZRangeBy) ([]string, error) {
	return t.client.ZRevRangeByScore(ctx, t.GetKey(key), opt).Result()
}

func (t *Client) ZRevRangeByScoreWithScores(ctx context.Context, key string, opt *ZRangeBy) ([]Z, error) {
	return t.client.ZRevRangeByScoreWithScores(ctx, t.GetKey(key), opt).Result()
}

func (t *Client) ZRevRank(ctx context.Context, key, member string) (int64, error) {
	return t.client.ZRevRank(ctx, t.GetKey(key), member).Result()
}

func (t *Client) ZScore(ctx context.Context, key, member string) (float64, error) {
	return t.client.ZScore(ctx, t.GetKey(key), member).Result()
}
