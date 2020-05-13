package cache

import (
	"github.com/go-redis/redis/v7"
	"time"
)

type Pool struct {
	clients map[string]*redis.Client
}

func CreatePool() Pool {
	return Pool{
		clients: make(map[string]*redis.Client),
	}
}
func (pool *Pool) PrintPools() {
	for name, _ := range pool.clients {
		println(name)
	}
}
func (pool *Pool) SetCache(client, key, value string, expires time.Duration) (err error) {
	_, err = pool.clients[client].Set(key, value, expires).Result()
	return err
}

func (pool *Pool) GetCache(client, key string) (value string, err error) {
	println("-----")
	println(client)

	value, err = pool.clients[client].Get(key).Result()
	return value, err
}

func (pool *Pool) Add(host string) (client *redis.Client) {
	if c, ok := pool.clients[host]; !ok {
		pool.clients[host] = redis.NewClient(&redis.Options{
			Addr:     host,
			Password: "", // no password set
			DB:       0,  // use default DB
		})
		client = pool.clients[host]
	} else {
		client = c
	}
	return client
}
