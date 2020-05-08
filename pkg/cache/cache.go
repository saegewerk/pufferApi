package cache

import (
	"github.com/go-redis/redis/v7"
	"time"
)

type Cache struct {
	Expires time.Duration
	MaxSize int64
	Host    string
	client  *redis.Client
}
type Config struct {
	Expires time.Duration
	MaxSize int64
	Host    string
}

func Create(config Config) Cache {
	return Cache{
		Expires: config.Expires,
		MaxSize: config.MaxSize,
		Host:    config.Host,
	}
}
