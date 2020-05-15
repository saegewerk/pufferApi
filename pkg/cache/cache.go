package cache

import (
	"github.com/go-redis/redis/v7"
	"time"
)

type Cache struct {
	Expires     time.Duration
	MaxSize     int64
	Client      *redis.Client
	Host        string
	HasMemcache bool
	Memcache    Memcache
}
type Memcache struct {
	Expires time.Time
	Store   string
}
type Config struct {
	Expires  time.Duration
	MaxSize  int64
	Host     string
	Memcache bool
}

func Create(config Config) Cache {
	println(config.Memcache)
	if config.Memcache {
		return Cache{
			Expires:     config.Expires,
			MaxSize:     config.MaxSize,
			Host:        config.Host,
			HasMemcache: true,
			Memcache: Memcache{
				Expires: time.Now(),
				Store:   "",
			},
		}
	} else {
		return Cache{
			Expires: config.Expires,
			MaxSize: config.MaxSize,
			Host:    config.Host,
		}
	}
}
