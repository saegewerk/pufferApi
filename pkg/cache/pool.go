package cache

import (
	"github.com/go-redis/redis/v7"
	"time"
)

type Pool struct {
	clients map[string]*redis.Client
}

func (pool *Pool) SetCache(client, key, value string, expires time.Duration) (err error) {
	_, err = pool.clients[client].Set(key, value, expires).Result()
	return err
}

func (pool *Pool) GetCache(client, key string) (value string, err error) {
	value, err = pool.clients[client].Get(key).Result()
	return value, err
}

func (pool *Pool) Add(config *Config) (client *redis.Client) {
	if c, ok := pool.clients[config.Host]; !ok {
		pool.clients[config.Host] = redis.NewClient(&redis.Options{
			Addr:     config.Host,
			Password: "", // no password set
			DB:       0,  // use default DB
		})
		client = pool.clients[config.Host]
	} else {
		client = c
	}
	return client
}
