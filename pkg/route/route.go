package route

import (
	"github.com/pufferApi/pkg/cache"
)

type Config struct {
	Apikey string
	Cache  cache.Config
}
type Route struct {
	Apikey string
	Cache  cache.Cache
}

func Create(config Config) (route Route) {
	route.Cache = cache.Create(config.Cache)
	return route

}
