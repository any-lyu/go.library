package redis

import (
	"github.com/garyburd/redigo/redis"
	"github.com/any-lyu/go.library/errors"
	"strconv"
	"time"
)

// Client redis client
type Client struct {
	Pool *redis.Pool
}

const redisNil = "redigo: nil returned" //redis正常返回

func redisOK(err error) (err2 error) {
	if err == nil || err.Error() == redisNil {
		return nil
	}
	return err
}

// TTL key
// 以秒为单位，返回给定 key 的剩余生存时间(TTL, time to live)。
func (pool *Client) TTL(key string) (ttl int, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	ttl, err = redis.Int(conn.Do("TTL", key))
	return ttl, err
}

// SADD 可以添加多个 返回成功数量
func (pool *Client) SADD(key string, value interface{}) (num int, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	num, err = redis.Int(conn.Do("SADD", key, value))
	return num, err
}

// Set 总是成功的
func (pool *Client) Set(key string, value interface{}) (err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	_, err = conn.Do("SET", key, value)
	return err
}

// SetNX 不存在则设置，存在则不设置
func (pool *Client) SetNX(key string, value interface{}) (err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	_, err = conn.Do("SET", key, value, "NX")
	return err
}

// SetNX2 不存在则设置，存在则不设置
func (pool *Client) SetNX2(key string, value interface{}) (num int, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	num, err = redis.Int(conn.Do("SETNX", key, value))
	return num, err
}

// Del 可以删除多个key 返回删除key的num和错误
func (pool *Client) Del(key ...interface{}) (num int, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	num, err = redis.Int(conn.Do("DEL", key...))
	return num, err
}

// Get redis get return string
func (pool *Client) Get(key string) (s string, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	s, err = redis.String(conn.Do("GET", key))
	return s, err
}

// GetInt redis get return int
func (pool *Client) GetInt(key string) (n int, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	n, err = redis.Int(conn.Do("GET", key))
	return n, err
}

// GetInt64 redis get return int64
func (pool *Client) GetInt64(key string) (n int64, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	n, err = redis.Int64(conn.Do("GET", key))
	return n, err
}

// EXISTS redis exist
func (pool *Client) EXISTS(key string) (ok bool, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	ok, err = redis.Bool(conn.Do("EXISTS", key))
	return ok, err
}

// KEYS redis range key
func (pool *Client) KEYS(pattern string) (keys []string, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	keys, err = redis.Strings(conn.Do("KEYS", pattern))
	return keys, err
}

// SCARD redis 返回集合中元素的数量
func (pool *Client) SCARD(key string) (num int, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	num, err = redis.Int(conn.Do("SCARD", key))
	return num, err
}

// SPOP 弹出被移除的元素, 当key不存在的时候返回 nil
func (pool *Client) SPOP(key string) (out string, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	out, err = redis.String(conn.Do("SPOP", key))
	return out, err
}

// SREM 移除集合 key 中的一个或多个 member 元素，不存在的 member 元素会被忽略
// 当 key 不是集合类型，返回一个错误。
func (pool *Client) SREM(key string, value interface{}) (num int, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	num, err = redis.Int(conn.Do("SREM", key, value))
	return num, err
}

// SISMEMBER 判断成员元素是否是集合的成员
func (pool *Client) SISMEMBER(key string, value interface{}) (ok bool, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	ok, err = redis.Bool(conn.Do("SISMEMBER", key, value))
	return ok, err
}

// SMEMBERS 返回集合 key 中的所有成员。
// 不存在的 key 被视为空集合。
func (pool *Client) SMEMBERS(key string) (reply []string, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	reply, err = redis.Strings(conn.Do("SMEMBERS", key))
	return reply, err
}

// LPOP 移除并返回列表 key 的头元素。
func (pool *Client) LPOP(key string) (out string, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	out, err = redis.String(conn.Do("LPOP", key))
	return out, err
}

// LPUSH 整型回复: 在 push 操作后的 list 长度。
func (pool *Client) LPUSH(key string, value ...interface{}) (num int, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	num, err = redis.Int(conn.Do("LPUSH", key, value))
	return num, err
}

// LINDEX 当 key 位置的值不是一个列表的时候，会返回一个error
func (pool *Client) LINDEX(key string, index int) (out string, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	out, err = redis.String(conn.Do("LINDEX", key, index))
	return out, err
}

// HEXISTS 检查给定域 field 是否存在于哈希表 hash 当中。
func (pool *Client) HEXISTS(key, field string) (ok bool, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	ok, err = redis.Bool(conn.Do("HEXISTS", key, field))
	return ok, err
}

// HGET 该字段所关联的值。当字段不存在或者 key 不存在时返回nil。
func (pool *Client) HGET(key, field string) (out string, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	out, err = redis.String(conn.Do("HGET", key, field))
	return out, err
}

// HINCRBY 增值操作执行后的该字段的值。
func (pool *Client) HINCRBY(key, field string, in int) (num int, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	num, err = redis.Int(conn.Do("HINCRBY", key, field, in))
	return num, err
}

// HMGETSTRUCT 返回hash表中所有字段 并映射为结构体
func (pool *Client) HMGETSTRUCT(key, value interface{}) (err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	v, err := redis.Values(conn.Do("HGETALL", key))
	if len(v) == 0 {
		return redis.ErrNil
	}
	if err == nil {
		err = redis.ScanStruct(v, value)
	}

	return err
}

// HMGETMAP 返回hash表中所有字段 并映射为map[string]string
func (pool *Client) HMGETMAP(key string) (map[string]string, error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	m, err := redis.StringMap(conn.Do("HGETALL", key))
	return m, err
}

// HMGETINTMAP 返回hash表中所有字段 并映射为map[string]int
func (pool *Client) HMGETINTMAP(key string) (map[string]int, error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	m, err := redis.IntMap(conn.Do("HGETALL", key))
	return m, err
}

// HMGETINT64MAP 返回hash表中所有字段 并映射为map[string]int64
func (pool *Client) HMGETINT64MAP(key string) (map[string]int64, error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	m, err := redis.Int64Map(conn.Do("HGETALL", key))
	return m, err
}

// HGETALLMAP 返回hash表中所有字段
func (pool *Client) HGETALLMAP(key string) (interface{}, error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	data, err := conn.Do("HGETALL", key)
	return data, err
}

// HMSET 同时将多个 field-value (域-值)对设置到哈希表 key 中。
// 此命令会覆盖哈希表中已存在的域。
// 如果 key 不存在，一个空哈希表被创建并执行 HMSET 操作。
func (pool *Client) HMSET(key, value interface{}) (ok string, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	ok, err = redis.String(conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(value)...))
	return ok, err
}

// HMGET 返回哈希表 key 中，一个或多个给定域的值。
// 如果给定的域不存在于哈希表，那么返回一个 nil 值。
// 因为不存在的 key 被当作一个空哈希表来处理，所以对一个不存在的 key 进行 HMGET 操作将返回一个只带有 nil 值的表。
func (pool *Client) HMGET(key, feild string) (data []string, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	data, err = redis.Strings(conn.Do("HMGET", key, feild))
	return data, err
}

// HKEYS 返回哈希表 key 中的所有域
func (pool *Client) HKEYS(key string) (data []string, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	data, err = redis.Strings(conn.Do("HKEYS", key))
	return data, err
}

// HMGET2 返回哈希表 key 中，一个或多个给定域的值。
// 如果给定的域不存在于哈希表，那么返回一个 nil 值。
// 因为不存在的 key 被当作一个空哈希表来处理，所以对一个不存在的 key 进行 HMGET 操作将返回一个只带有 nil 值的表。
func (pool *Client) HMGET2(key string, feild ...string) (data []string, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	data, err = redis.Strings(conn.Do("HMGET", redis.Args{}.Add(key).AddFlat(feild)...))
	return data, err
}

// HSET 1如果field是一个新的字段  0如果field原来在map里面已经存在
func (pool *Client) HSET(key, field string, value interface{}) (ok bool, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	ok, err = redis.Bool(conn.Do("HSET", key, field, value))
	return ok, err
}

// HLEN 哈希集中字段的数量，当 key 指定的哈希集不存在时返回 0
func (pool *Client) HLEN(key string) (num int, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	num, err = redis.Int(conn.Do("HLEN", key))
	return num, err
}

// ZREMRANGEBYRANK myzset 0 1  0 -200(保留200名)
func (pool *Client) ZREMRANGEBYRANK(key string, stop int) (num int, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	num, err = redis.Int(conn.Do("ZREMRANGEBYRANK", key, 0, stop))
	return num, err
}

// ZADD 将一个或多个 member 元素及其 score 值加入到有序集 key 当中。
// 如果某个 member 已经是有序集的成员，那么更新这个 member 的 score 值，并通过重新插入这个 member 元素，来保证该 member 在正确的位置上。
func (pool *Client) ZADD(key string, sorce int, member string) (num int, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	num, err = redis.Int(conn.Do("ZADD", key, sorce, member))
	return num, err
}

// ZFADD ZADD float64
func (pool *Client) ZFADD(key string, sorce float64, member string) (num int, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	num, err = redis.Int(conn.Do("ZADD", key, sorce, member))
	return num, err
}

// ZCARD cz
func (pool *Client) ZCARD(key string) (num int, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	num, err = redis.Int(conn.Do("ZCARD", key))
	return num, err
}

// ZRANGE cz
func (pool *Client) ZRANGE(key string, start, stop int) (list []string, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	list, err = redis.Strings(conn.Do("ZRANGE", key, start, stop))
	return list, err
}

// ZREVRANGE cz
func (pool *Client) ZREVRANGE(key string, start, stop int) (list []string, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	list, err = redis.Strings(conn.Do("ZREVRANGE", key, start, stop))
	return list, err
}

// ZSCORE cz
func (pool *Client) ZSCORE(key string, member string) (num int, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	num, err = redis.Int(conn.Do("ZSCORE", key, member))
	return num, err
}

// ZFSCORE ZSCORE cz
func (pool *Client) ZFSCORE(key string, member string) (num float64, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	num, err = redis.Float64(conn.Do("ZSCORE", key, member))
	return num, err
}

// ZREM cz
func (pool *Client) ZREM(key string, member string) (num int, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	num, err = redis.Int(conn.Do("ZREM", key, member))
	return num, err
}

// ZREVRANGEBYSCORE 逆序份数  获取的 前N个数据
func (pool *Client) ZREVRANGEBYSCORE(key string, limit int) (list map[string]string, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	list, err = redis.StringMap(conn.Do("ZREVRANGEBYSCORE", key, "+inf", "-inf", "WITHSCORES", "limit", 0, limit))
	return list, err
}

// ZREVRANGEBYSCORE2 ZREVRANGEBYSCORE 逆序份数  获取start len的数据
func (pool *Client) ZREVRANGEBYSCORE2(key string, start, len int) (list map[string]int, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	list, err = redis.IntMap(conn.Do("ZREVRANGEBYSCORE", key, "+inf", "-inf", "WITHSCORES", "limit", start, len))
	return list, err
}

// ZREVRANGEBYSCORE3 ZREVRANGEBYSCORE 逆序份数  获取start len的数据
func (pool *Client) ZREVRANGEBYSCORE3(key string, start, len int) (list map[string]float64, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	list, err = floatMap(conn.Do("ZREVRANGEBYSCORE", key, "+inf", "-inf", "WITHSCORES", "limit", start, len))
	return list, err
}

func floatMap(result interface{}, err error) (map[string]float64, error) {
	values, err := redis.Values(result, err)
	if err != nil {
		return nil, err
	}
	if len(values)%2 != 0 {
		return nil, errors.New("redigo: IntMap expects even number of values result")
	}
	m := make(map[string]float64, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].([]byte)
		if !ok {
			return nil, errors.New("redigo: IntMap key not a bulk string value")
		}
		value, err := redis.Float64(values[i+1], nil)
		if err != nil {
			return nil, err
		}
		m[string(key)] = value
	}
	return m, nil
}

// GetSearchKeys ZREVRANGEBYSCORE 逆序份数  获取的 前N个数据 不要scores
func (pool *Client) GetSearchKeys(key string, limit int) (list []string, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	list, err = redis.Strings(conn.Do("ZREVRANGEBYSCORE", key, "+inf", "-inf", "limit", 0, limit))
	return list, err
}

// GetSearchKeys2 ZREVRANGEBYSCORE 逆序份数  获取的 start,len 不要scores
func (pool *Client) GetSearchKeys2(key string, start, len int) (list []string, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	list, err = redis.Strings(conn.Do("ZREVRANGEBYSCORE", key, "+inf", "-inf", "limit", start, len))
	return list, err
}

// ZINCRBY +increment  如果没有key 插入
func (pool *Client) ZINCRBY(key string, increment int, member string) (num int, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	num, err = redis.Int(conn.Do("ZINCRBY", key, increment, member))
	return num, err
}

// ZRANK 判断一个member 在key中的索引 如果不在 返回nil ,在 返回索引
func (pool *Client) ZRANK(key string, member string) (num int, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	num, err = redis.Int(conn.Do("ZRANK", key, member))
	return num, err
}

// ZREVRANK cz
func (pool *Client) ZREVRANK(key string, member string) (num int, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	num, err = redis.Int(conn.Do("ZREVRANK", key, member))
	return num, err
}

// EXPIRE 设置一个key 的过期时间 返回值int 1 如果设置了过期时间 0 如果没有设置过期时间，或者不能设置过期时间
func (pool *Client) EXPIRE(key string, expireTime int) (num int, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	num, err = redis.Int(conn.Do("EXPIRE", key, expireTime))
	return num, err
}

// EXPIREAT 设置一个key 的在指定时间过期 返回值：如果生存时间设置成功，返回 1 ;当 key 不存在或没办法设置生存时间，返回 0 。
func (pool *Client) EXPIREAT(key string, expireAtTime int64) (num int, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	num, err = redis.Int(conn.Do("EXPIREAT", key, expireAtTime))
	return num, err
}

// SETEX key seconds value
func (pool *Client) SETEX(key string, seconds int, value interface{}) (err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	_, err = conn.Do("SETEX", key, seconds, value)
	return err
}

// INCR 为键 key 储存的数字值加上一。
// 如果键 key 不存在， 那么它的值会先被初始化为 0 ， 然后再执行 INCR 命令。
// 如果键 key 储存的值不能被解释为数字， 那么 INCR 命令将返回一个错误。
// 本操作的值限制在 64 位(bit)有符号数字表示之内。
func (pool *Client) INCR(key string) (err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	_, err = conn.Do("INCR", key)
	return err
}

// INCRRET INCR 为键 key 储存的数字值加上一。
// 如果键 key 不存在， 那么它的值会先被初始化为 0 ， 然后再执行 INCR 命令。
// 如果键 key 储存的值不能被解释为数字， 那么 INCR 命令将返回一个错误。
// 本操作的值限制在 64 位(bit)有符号数字表示之内。
func (pool *Client) INCRRET(key string) (num int, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	num, err = redis.Int(conn.Do("INCR", key))
	return num, err
}

// SETBIT 对 key 所储存的字符串值，设置或清除指定偏移量上的位(bit)。
// 位的设置或清除取决于 value 参数，可以是 0 也可以是 1 。
// 当 key 不存在时，自动生成一个新的字符串值。
// 字符串会进行伸展(grown)以确保它可以将 value 保存在指定的偏移量上。当字符串值进行伸展时，空白位置以 0 填充。
func (pool *Client) SETBIT(key string, bit, value int) (ret int, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	ret, err = redis.Int(conn.Do("SETBIT", key, bit, value))
	return ret, err
}

// GETBIT 获取指定偏移量上的位(bit)
func (pool *Client) GETBIT(key string, bit int) (ret int, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()

	ret, err = redis.Int(conn.Do("GETBIT", key, bit))
	return ret, err
}

// HMSETArgs HMSET args
func (pool *Client) HMSETArgs(key string, node interface{}) error {
	conn := pool.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(node)...)
	return err
}

// GetSet 将键 key 的值设为 value ， 并返回键 key 在被设置之前的旧值。
func (pool *Client) GetSet(key string, value interface{}) (s string, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()
	s, err = redis.String(conn.Do("GETSET", key, value))
	err = redisOK(err)
	return s, nil
}

// Unlink like delete redis 4.0 +
func (pool *Client) Unlink(key ...interface{}) (num int, err error) {
	conn := pool.Pool.Get()
	defer conn.Close()
	num, err = redis.Int(conn.Do("UNLINK", key...))
	return num, redisOK(err)
}

// 所有的过期时间都是毫秒（ms）单位
func getExpiredTime(timeout int) (int64, int64) {
	now := time.Now()
	return now.UnixNano(), now.Add(time.Millisecond * time.Duration(timeout)).UnixNano()
}

// Lock redis 实现的分布式锁 lock
func (pool *Client) Lock(key string, timeout int) (locked bool, expiredTime int64, err error) {
	// 根据过期时间毫秒数获取当前时间和过期时间
	now, expiredTime := getExpiredTime(timeout)

	// SetNX 设置过期时间，如果成功，加锁成功，如果失败，证明锁被占据
	if success, err := pool.SetNX2(key, expiredTime); err != nil {
		return false, 0, errors.Wrapf(err, "Lock:SetNX key:%v, value:%v", key, expiredTime)
	} else if success > 0 {
		return true, expiredTime, nil
	}

	// 如果锁被占据，获取锁的内容，判断是否过期等
	if value, err := pool.Get(key); err != nil {
		return false, 0, errors.Wrapf(err, "Lock:get key:%v", key)
	} else if value == "" {
		// 已经被删除的情况下，返回的是空
		if value, err := pool.GetSet(key, expiredTime); err != nil {
			return false, 0, errors.Wrapf(err, "Lock:GetSet key:%v, expireTime:%v", key, expiredTime)
		} else if value == "" {
			return true, expiredTime, nil
		}
	} else {
		// 如果里面有内容，查看是否过期即可
		if passTime, err := strconv.ParseInt(value, 10, 64); err != nil {
			// 解析时出现问题，这个时候问题较大，只能等待过期时间了（否则在这边删除也可以，不过正常情况系啊，不会走这边）
			return false, 0, errors.Wrapf(err, "Lock:ParseInt passTime:%v", value)
		} else if now > passTime {
			// 已经过期了，这个时候直接获取锁即可
			if valueNow, err := pool.GetSet(key, expiredTime); err != nil {
				return false, 0, errors.Wrapf(err, "GetSet key:%v, expireTime:%v", key, expiredTime)
			} else if valueNow == value {
				return true, expiredTime, nil
			}
		}
	}
	return false, 0, nil
}

// LockRetry redis 实现的分布式锁 lock retry
func (pool *Client) LockRetry(key string, timeout, retryTimes int) (locked bool, expiredTime int64, err error) {
	for i := 0; i < retryTimes; i++ {
		if locked, expiredTime, err = pool.Lock(key, timeout); err != nil || locked {
			return
		}
		time.Sleep(time.Millisecond * time.Duration(i+1))
	}
	return false, 0, nil
}

// LockMust redis 实现的分布式锁 lock must
func (pool *Client) LockMust(key string, timeout int) (locked bool, expiredTime int64, err error) {
	for i := 0; ; i++ {
		if locked, expiredTime, err = pool.Lock(key, timeout); err != nil || locked {
			return
		}
		time.Sleep(time.Millisecond * time.Duration(i+1))
	}
}

// UnLock redis 实现的分布式锁 unlock
func (pool *Client) UnLock(key string, safeDelTime int64) (bool, error) {
	if value, err := pool.Get(key); err != nil {
		// 获取KEY的时候报错，证明可能已经过期，或者别别人删除了
		return false, errors.Wrap(err, "UnLock:get")
	} else if expireTime, err := strconv.ParseInt(value, 10, 64); err != nil {
		// 过期时间解析错误
		return false, errors.Wrapf(err, "UnLock:ParseInt key:%v", value)
	} else if time.Now().UnixNano()+safeDelTime*1000000 > expireTime {
		// 就要到过期时间了，这个时候直接退出，等待过期时间即可
		return false, errors.Errorf("UnLock: Key:%v nearly to be expired.", key)
	} else if count, err := pool.Del(key); err != nil {
		return false, errors.Wrapf(err, "UnLock:Del key:%v ", key)
	} else if count == 0 {
		return false, errors.Errorf("UnLock: Del key:%v, Count:%v", key, count)
	}
	return true, nil
}
