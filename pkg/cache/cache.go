package cache

import (
	"encoding/json"
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

type Response struct {
	Header []byte
	Body   []byte
}

func (m *Response) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}
func (m Response) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

type Memcache struct {
	Expires time.Time
	Store   Response
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
