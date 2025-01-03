package dbredis

import (
	"Common/Framework/config"
	"github.com/gomodule/redigo/redis"
	"log/slog"
	"time"
)

var pool *redis.Pool

func Init() error {
	pool = &redis.Pool{
		MaxIdle:     config.Instance().Redis.MaxIdle,
		MaxActive:   config.Instance().Redis.MaxActive,
		Wait:        true,
		IdleTimeout: 60 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				config.Instance().Redis.Addr,
				redis.DialDatabase(config.Instance().Redis.DB),
				redis.DialUsername(config.Instance().Redis.Username),
				redis.DialPassword(config.Instance().Redis.Password),
				redis.DialConnectTimeout(5*time.Second),
				redis.DialReadTimeout(5*time.Second),
				redis.DialWriteTimeout(5*time.Second),
			)
		},
	}

	var conn redis.Conn

	for {
		conn = pool.Get()
		if err := conn.Err(); err != nil {
			slog.Debug("redis connect ", slog.Any("error", err))
			_ = conn.Close()
			time.Sleep(2 * time.Second)
			continue
		}
	RETRY:
		if _, err := conn.Do("PING"); err != nil {
			slog.Debug("wait for redis up ", slog.Any("error", err))
			time.Sleep(2 * time.Second)
			goto RETRY
		} else {
			break
		}
	}

	slog.Info("redis init ok")
	return conn.Close()
}

func Conn() redis.Conn {
	return pool.Get()
}
