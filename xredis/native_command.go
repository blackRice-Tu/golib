package xredis

import (
	"context"
	"time"
)

// native command

func (t *Client) NativeDel(ctx context.Context, keys ...string) (int64, error) {
	return t.client.Del(ctx, keys...).Result()
}

func (t *Client) NativeExists(ctx context.Context, keys ...string) (int64, error) {
	return t.client.Exists(ctx, keys...).Result()
}

func (t *Client) NativeExpire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return t.client.Expire(ctx, key, expiration).Result()
}

func (t *Client) NativeExpireNX(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return t.client.ExpireNX(ctx, key, expiration).Result()
}

func (t *Client) NativeExpireXX(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return t.client.ExpireXX(ctx, key, expiration).Result()
}

func (t *Client) NativeExpireGT(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return t.client.ExpireGT(ctx, key, expiration).Result()
}

func (t *Client) NativeExpireLT(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return t.client.ExpireLT(ctx, key, expiration).Result()
}

func (t *Client) NativeExpireAt(ctx context.Context, key string, tm time.Time) (bool, error) {
	return t.client.ExpireAt(ctx, key, tm).Result()
}

func (t *Client) NativeKeys(ctx context.Context, pattern string) ([]string, error) {
	return t.client.Keys(ctx, pattern).Result()
}

func (t *Client) NativeTTL(ctx context.Context, key string) (time.Duration, error) {
	return t.client.TTL(ctx, key).Result()
}

func (t *Client) NativeType(ctx context.Context, key string) (string, error) {
	return t.client.Type(ctx, key).Result()
}

func (t *Client) NativeDecr(ctx context.Context, key string) (int64, error) {
	return t.client.Decr(ctx, key).Result()
}

func (t *Client) NativeDecrBy(ctx context.Context, key string, decrement int64) (int64, error) {
	return t.client.DecrBy(ctx, key, decrement).Result()
}

func (t *Client) NativeGet(ctx context.Context, key string) (string, error) {
	return t.client.Get(ctx, key).Result()
}

func (t *Client) NativeIncr(ctx context.Context, key string) (int64, error) {
	return t.client.Incr(ctx, key).Result()
}

func (t *Client) NativeIncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return t.client.IncrBy(ctx, key, value).Result()
}

func (t *Client) NativeIncrByFloat(ctx context.Context, key string, value float64) (float64, error) {
	return t.client.IncrByFloat(ctx, key, value).Result()
}

func (t *Client) NativeSet(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error) {
	return t.client.Set(ctx, key, value, expiration).Result()
}

func (t *Client) NativeSetEx(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error) {
	return t.client.SetEx(ctx, key, value, expiration).Result()
}

func (t *Client) NativeSetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return t.client.SetNX(ctx, key, value, expiration).Result()
}

func (t *Client) NativeSetXX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return t.client.SetXX(ctx, key, value, expiration).Result()
}

// ------

func (t *Client) NativeScan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return t.client.Scan(ctx, cursor, match, count).Result()
}

func (t *Client) NativeScanAll(ctx context.Context, match string, count int64) ([]string, error) {
	var cursor uint64 = 0
	var keys []string
	var err error
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

func (t *Client) NativeScanType(ctx context.Context, cursor uint64, match string, count int64, keyType string) ([]string, uint64, error) {
	return t.client.ScanType(ctx, cursor, match, count, keyType).Result()
}

func (t *Client) NativeSScan(ctx context.Context, key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return t.client.SScan(ctx, key, cursor, match, count).Result()
}

func (t *Client) NativeSScanAll(ctx context.Context, key string, match string, count int64) ([]string, error) {
	var cursor uint64 = 0
	var keys []string
	var err error
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

func (t *Client) NativeHScan(ctx context.Context, key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return t.client.HScan(ctx, key, cursor, match, count).Result()
}

func (t *Client) NativeHScanAll(ctx context.Context, key string, match string, count int64) ([]string, error) {
	var cursor uint64 = 0
	var keys []string
	var err error
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

func (t *Client) NativeZScan(ctx context.Context, key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return t.client.ZScan(ctx, key, cursor, match, count).Result()
}

func (t *Client) NativeZScanAll(ctx context.Context, key string, match string, count int64) ([]string, error) {
	var cursor uint64 = 0
	var keys []string
	var err error
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

func (t *Client) NativeHDel(ctx context.Context, key string, fields ...string) (int64, error) {
	return t.client.HDel(ctx, key, fields...).Result()
}

func (t *Client) NativeHExists(ctx context.Context, key, field string) (bool, error) {
	return t.client.HExists(ctx, key, field).Result()
}

func (t *Client) NativeHGet(ctx context.Context, key, field string) (string, error) {
	return t.client.HGet(ctx, key, field).Result()
}

func (t *Client) NativeHGetAll(ctx context.Context, key string) (map[string]string, error) {
	return t.client.HGetAll(ctx, key).Result()
}

func (t *Client) NativeHIncrBy(ctx context.Context, key, field string, incr int64) (int64, error) {
	return t.client.HIncrBy(ctx, key, field, incr).Result()
}

func (t *Client) NativeHIncrByFloat(ctx context.Context, key, field string, incr float64) (float64, error) {
	return t.client.HIncrByFloat(ctx, key, field, incr).Result()
}

func (t *Client) NativeHKeys(ctx context.Context, key string) ([]string, error) {
	return t.client.HKeys(ctx, key).Result()
}

func (t *Client) NativeHLen(ctx context.Context, key string) (int64, error) {
	return t.client.HLen(ctx, key).Result()
}

func (t *Client) NativeHSet(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return t.client.HSet(ctx, key, values...).Result()
}

func (t *Client) NativeHVals(ctx context.Context, key string) ([]string, error) {
	return t.client.HVals(ctx, key).Result()
}

func (t *Client) NativeLLen(ctx context.Context, key string) (int64, error) {
	return t.client.LLen(ctx, key).Result()
}

func (t *Client) NativeLPop(ctx context.Context, key string) (string, error) {
	return t.client.LPop(ctx, key).Result()
}

func (t *Client) NativeLPopCount(ctx context.Context, key string, count int) ([]string, error) {
	return t.client.LPopCount(ctx, key, count).Result()
}

func (t *Client) NativeLPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return t.client.LPush(ctx, key, values).Result()
}

func (t *Client) NativeLRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return t.client.LRange(ctx, key, start, stop).Result()
}

func (t *Client) NativeLRem(ctx context.Context, key string, count int64, value interface{}) (int64, error) {
	return t.client.LRem(ctx, key, count, value).Result()
}

func (t *Client) NativeRPop(ctx context.Context, key string) (string, error) {
	return t.client.RPop(ctx, key).Result()
}

func (t *Client) NativeRPopCount(ctx context.Context, key string, count int) ([]string, error) {
	return t.client.RPopCount(ctx, key, count).Result()
}

func (t *Client) NativeRPopLPush(ctx context.Context, source, destination string) (string, error) {
	return t.client.RPopLPush(ctx, source, destination).Result()
}

func (t *Client) NativeRPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return t.client.RPush(ctx, key, values...).Result()
}

// ----

func (t *Client) NativeSAdd(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return t.client.SAdd(ctx, key, members...).Result()
}

func (t *Client) NativeSCard(ctx context.Context, key string) (int64, error) {
	return t.client.SCard(ctx, key).Result()
}

func (t *Client) NativeSDiff(ctx context.Context, keys ...string) ([]string, error) {
	return t.client.SDiff(ctx, keys...).Result()
}

func (t *Client) NativeSIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return t.client.SIsMember(ctx, key, member).Result()
}

func (t *Client) NativeSMembers(ctx context.Context, key string) ([]string, error) {
	return t.client.SMembers(ctx, key).Result()
}

func (t *Client) NativeSPop(ctx context.Context, key string) (string, error) {
	return t.client.SPop(ctx, key).Result()
}

func (t *Client) NativeSPopN(ctx context.Context, key string, count int64) ([]string, error) {
	return t.client.SPopN(ctx, key, count).Result()
}

func (t *Client) NativeSRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return t.client.SRem(ctx, key, members...).Result()
}

func (t *Client) NativeSUnion(ctx context.Context, keys ...string) ([]string, error) {
	return t.client.SUnion(ctx, keys...).Result()
}

func (t *Client) NativeSUnionStore(ctx context.Context, destination string, keys ...string) (int64, error) {
	return t.client.SUnionStore(ctx, destination, keys...).Result()
}

// -----

func (t *Client) NativeBZPopMax(ctx context.Context, timeout time.Duration, keys ...string) (*ZWithKey, error) {
	return t.client.BZPopMax(ctx, timeout, keys...).Result()
}

func (t *Client) NativeBZPopMin(ctx context.Context, timeout time.Duration, keys ...string) (*ZWithKey, error) {
	return t.client.BZPopMin(ctx, timeout, keys...).Result()
}

func (t *Client) NativeZAddArgs(ctx context.Context, key string, args ZAddArgs) (int64, error) {
	return t.client.ZAddArgs(ctx, key, args).Result()
}

func (t *Client) NativeZAddArgsIncr(ctx context.Context, key string, args ZAddArgs) (float64, error) {
	return t.client.ZAddArgsIncr(ctx, key, args).Result()
}

func (t *Client) NativeZAdd(ctx context.Context, key string, members ...Z) (int64, error) {
	return t.client.ZAdd(ctx, key, members...).Result()
}

func (t *Client) NativeZAddNX(ctx context.Context, key string, members ...Z) (int64, error) {
	return t.client.ZAddNX(ctx, key, members...).Result()
}

func (t *Client) NativeZAddXX(ctx context.Context, key string, members ...Z) (int64, error) {
	return t.client.ZAddXX(ctx, key, members...).Result()
}

func (t *Client) NativeZCard(ctx context.Context, key string) (int64, error) {
	return t.client.ZCard(ctx, key).Result()
}

func (t *Client) NativeZCount(ctx context.Context, key, min, max string) (int64, error) {
	return t.client.ZCount(ctx, key, min, max).Result()
}

func (t *Client) NativeZLexCount(ctx context.Context, key, min, max string) (int64, error) {
	return t.client.ZLexCount(ctx, key, min, max).Result()
}

func (t *Client) NativeZIncrBy(ctx context.Context, key string, increment float64, member string) (float64, error) {
	return t.client.ZIncrBy(ctx, key, increment, member).Result()
}

func (t *Client) NativeZPopMax(ctx context.Context, key string, count ...int64) ([]Z, error) {
	return t.client.ZPopMax(ctx, key, count...).Result()
}

func (t *Client) NativeZPopMin(ctx context.Context, key string, count ...int64) ([]Z, error) {
	return t.client.ZPopMin(ctx, key, count...).Result()
}

func (t *Client) NativeZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return t.client.ZRange(ctx, key, start, stop).Result()
}

func (t *Client) NativeZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]Z, error) {
	return t.client.ZRangeWithScores(ctx, key, start, stop).Result()
}

func (t *Client) NativeZRangeByScore(ctx context.Context, key string, opt *ZRangeBy) ([]string, error) {
	return t.client.ZRangeByScore(ctx, key, opt).Result()
}

func (t *Client) NativeZRangeByLex(ctx context.Context, key string, opt *ZRangeBy) ([]string, error) {
	return t.client.ZRangeByLex(ctx, key, opt).Result()
}

func (t *Client) NativeZRangeByScoreWithScores(ctx context.Context, key string, opt *ZRangeBy) ([]Z, error) {
	return t.client.ZRangeByScoreWithScores(ctx, key, opt).Result()
}

func (t *Client) NativeZRank(ctx context.Context, key, member string) (int64, error) {
	return t.client.ZRank(ctx, key, member).Result()
}

func (t *Client) NativeZRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return t.client.ZRem(ctx, key, members...).Result()
}

func (t *Client) NativeZRemRangeByRank(ctx context.Context, key string, start, stop int64) (int64, error) {
	return t.client.ZRemRangeByRank(ctx, key, start, stop).Result()
}

func (t *Client) NativeZRemRangeByScore(ctx context.Context, key, min, max string) (int64, error) {
	return t.client.ZRemRangeByScore(ctx, key, min, max).Result()
}

func (t *Client) NativeZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return t.client.ZRevRange(ctx, key, start, stop).Result()
}

func (t *Client) NativeZRevRangeWithScores(ctx context.Context, key string, start, stop int64) ([]Z, error) {
	return t.client.ZRevRangeWithScores(ctx, key, start, stop).Result()
}

func (t *Client) NativeZRevRangeByScore(ctx context.Context, key string, opt *ZRangeBy) ([]string, error) {
	return t.client.ZRevRangeByScore(ctx, key, opt).Result()
}

func (t *Client) NativeZRevRangeByScoreWithScores(ctx context.Context, key string, opt *ZRangeBy) ([]Z, error) {
	return t.client.ZRevRangeByScoreWithScores(ctx, key, opt).Result()
}

func (t *Client) NativeZRevRank(ctx context.Context, key, member string) (int64, error) {
	return t.client.ZRevRank(ctx, key, member).Result()
}

func (t *Client) NativeZScore(ctx context.Context, key, member string) (float64, error) {
	return t.client.ZScore(ctx, key, member).Result()
}
