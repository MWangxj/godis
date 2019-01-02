package godis

import (
	"encoding/json"
	"errors"
	"github.com/garyburd/redigo/redis"
	"reflect"
	"time"
)

var (
	keyPrefix, keySubfix string
	pool                 *redis.Pool
)

// SetKeyPrefix 设置key前缀
func SetKeyPrefix(prefix string) {
	keyPrefix = prefix
}

// SetKeySubfix 设置key后缀
func SetKeySubfix(subfix string) {
	keySubfix = subfix
}

// Dail 拨号
func Dail(addr, passwd string, maxidle, maxac, db int, to time.Duration) *redis.Pool {
	pool = &redis.Pool{
		MaxIdle:     maxidle,
		MaxActive:   maxac,
		IdleTimeout: to,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr, redis.DialPassword(passwd), redis.DialDatabase(db))
			if err != nil {
				return nil, err
			}
			if db > 0 && db < 16 {
				if _, err := c.Do("SELECT", db); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, nil
		},
	}
	return pool
}

// Set 调用redis set 命令
func Set(key string, v interface{}) error {
	if pool == nil {
		return errors.New("please dial redis server first.")
	}
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	conn := pool.Get()
	defer conn.Close()
	if _, err = conn.Do("SET", formatKey(key), data); err != nil {
		return err
	}
	return nil
}

func Expire(key string, ttl int) error {
	if pool == nil {
		return errors.New("please dial redis server first.")
	}
	conn := pool.Get()
	defer conn.Close()
	if _, err := conn.Do("EXPIRE", formatKey(key), ttl); err != nil {
		return err
	}
	return nil
}

// GetInt
func GetInt(key string) (int, error) {
	if pool == nil {
		return 0, errors.New("please dail redis server first.")
	}
	conn := pool.Get()
	defer conn.Close()
	return redis.Int(conn.Do("GET", formatKey(key)))
}

// GetString
func GetString(key string) (string, error) {
	if pool == nil {
		return "", errors.New("please dail redis server first.")
	}
	conn := pool.Get()
	defer conn.Close()
	return redis.String(conn.Do("GET", formatKey(key)))
}

// GetBool
func GetBool(key string) (bool, error) {
	if pool == nil {
		return false, errors.New("please dail redis server first.")
	}
	conn := pool.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("GET", formatKey(key)))
}

// GetObj 调用redis get 命令
func GetObj(key string, v interface{}) error {
	if pool == nil {
		return errors.New("please dail redis server first.")
	}
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		return errors.New("param2 is not a pointer")
	}
	conn := pool.Get()
	defer conn.Close()
	data, err := redis.Bytes(conn.Do("GET", formatKey(key)))
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// Del 调用redis del 命令
func Del(key string) error {
	if pool == nil {
		return errors.New("please dial redis server first.")
	}
	conn := pool.Get()
	defer conn.Close()
	if _, err := conn.Do("DEL", formatKey(key)); err != nil {
		return err
	}
	return nil
}

// HSet 调用redis hset 设置哈希数据结构
func HSet(key, region string, v interface{}) error {
	if pool == nil {
		return errors.New("please dial redis server first.")
	}
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	conn := pool.Get()
	defer conn.Close()
	if _, err = conn.Do("HSET", formatKey(key), region, data); err != nil {
		return err
	}
	return nil
}

// HDel 删除哈希数据结构中的一个region
func HDel(key, region string) error {
	if pool == nil {
		return errors.New("please dial redis server first.")
	}
	conn := pool.Get()
	defer conn.Close()
	if _, err := conn.Do("HDEL", formatKey(key), region); err != nil {
		return err
	}
	return nil
}

// HGet 获取哈希数据结构一个region的值
func HGet(key, region string, v interface{}) error {
	if pool == nil {
		return errors.New("please dial redis server first.")
	}
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		return errors.New("param2 is not a pointer")
	}
	conn := pool.Get()
	defer conn.Close()
	data, err := redis.Bytes(conn.Do("HGET", formatKey(key), region))
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func LPush(key string, v ...string) error {
	if pool == nil {
		return errors.New("please dial redis server first.")
	}
	conn := pool.Get()
	defer conn.Close()
	data, err := redis.Bytes(conn.Do("LPUSH", v))
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func LPop(key string) (string, error) {
	if pool == nil {
		return "", errors.New("please dial redis server first.")
	}
	conn := pool.Get()
	defer conn.Close()
	return redis.String(conn.Do("LPOP", key))
}

// HGetAll 获取该key下面所有哈希的集合，值以字符串表示
func HGetAll(key string, v *map[string]string) error {
	if pool == nil {
		return errors.New("please dial redis server first.")
	}
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		return errors.New("param2 is not a pointer")
	}
	conn := pool.Get()
	defer conn.Close()
	arr, err := redis.Strings(conn.Do("HGETALL", formatKey(key)))
	if err != nil {
		return err
	}
	if v, err = arr2map(arr); err != nil {
		return err
	}
	return nil
}

// SAdd 向该key下面添加set散列
func SAdd(key string, v interface{}) error {
	if pool == nil {
		return errors.New("please dial redis server first.")
	}
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	conn := pool.Get()
	defer conn.Close()
	if _, err = conn.Do("SADD", formatKey(key), data); err != nil {
		return err
	}
	return nil
}

// SMembers 获取该key下面 所有set集合
func SMembers(key string, v []string) error {
	if pool == nil {
		return errors.New("please dial redis server first.")
	}
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		return errors.New("param2 is not a pointer")
	}
	conn := pool.Get()
	defer conn.Close()
	var err error
	if v, err = redis.Strings(conn.Do("SMEMBERS", formatKey(key))); err != nil {
		return err
	}
	return nil
}

func arr2map(arr []string) (*map[string]string, error) {
	if len(arr)%2 != 0 {
		return nil, errors.New("array length is not right")
	}
	m := make(map[string]string, len(arr)/2)
	for i := 0; i < len(arr)-1; i += 2 {
		m[arr[i]] = arr[i+1]
	}
	return &m, nil
}

func formatKey(key string) string {
	if keyPrefix != "" {
		key = keyPrefix + ":" + key
	}
	if keySubfix != "" {
		key += ":" + keySubfix
	}
	return key
}
