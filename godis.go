package godis

import (
	`encoding/json`
	`errors`
	`github.com/garyburd/redigo/redis`
	`reflect`
	`time`
)

var (
	keyPrefix, keySubfix string
	pool                 *redis.Pool
)

func SetKeyPrefix(prefix string) {
	keyPrefix = prefix
}

func SetKeySubfix(subfix string) {
	keySubfix = subfix
}

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

func Get(key string, v interface{}) error {
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
