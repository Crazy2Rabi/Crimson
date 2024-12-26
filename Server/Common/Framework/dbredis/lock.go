package dbredis

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

func Lock(conn redis.Conn, key string, ttlSecs int) (err error) {
	args := redis.Args{}.Add(key, time.Now().UnixMilli(), "NX")

	if ttlSecs > 0 {
		args = args.Add("PX", ttlSecs*1000)
	}

	_, err = redis.String(conn.Do("SET", args...))

	return
}

func Unlock(conn redis.Conn, key string) (err error) {
	_, err = redis.String(conn.Do("DEL", key))
	return
}
