package godis

import (
	"github.com/garyburd/redigo/redis"
	"testing"
	"time"
)

func dial() *redis.Pool {
	return Dail("127.0.0.1:6379", "", 10, 10, 4, 300*time.Second)
}

func TestRedis(t *testing.T) {
	p := dial()

	t.Log(p)

	if err := Set("wangxianjin", 1); err != nil {
		t.Fail()
		t.Log(err.Error())
	}
	if _, err := GetInt("wangxianjin"); err != nil {
		t.Fail()
		t.Log(err.Error())
	}
}

func TestExpire(t *testing.T) {
	p := dial()

	t.Log(p)

	if err := Set("wangxianjin", 1); err != nil {
		t.Fail()
		t.Log(err.Error())
	}
	if err := Expire("wangxianjin", 60); err != nil {
		t.Fail()
		t.Log(err.Error())
	}
}

func TestHGet(t *testing.T) {
	p := dial()

	t.Log(p)

	if err := HSet("wangxianjin", "age", 26); err != nil {
		t.Fail()
		t.Log(err.Error())
		return
	}
	if err := HSet("wangxianjin", "name", "wangxianjin"); err != nil {
		t.Fail()
		t.Log(err.Error())
		return
	}
	var v int
	if err := HGet("wangxianjin", "age", &v); err != nil {
		t.Fail()
		t.Log(err.Error())
		return
	}
	t.Log(v)
	var m map[string]string
	if err := HGetAll("wangxianjin", &m); err != nil {
		t.Fail()
		t.Log(err.Error())
		return
	}
}
