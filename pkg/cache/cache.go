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
	HasApikey   bool
	Apikey      string
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
	Apikey   string
}

func Create(config Config) (cache Cache) {

	if config.Memcache {
		if config.Apikey == "" {
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
		}
		return Cache{
			Apikey:      config.Apikey,
			HasApikey:   true,
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
		if config.Apikey == "" {
			return Cache{
				Expires: config.Expires,
				MaxSize: config.MaxSize,
				Host:    config.Host,
			}
		}
		return Cache{
			Apikey:    config.Apikey,
			HasApikey: true,
			Expires:   config.Expires,
			MaxSize:   config.MaxSize,
			Host:      config.Host,
		}
	}
}
