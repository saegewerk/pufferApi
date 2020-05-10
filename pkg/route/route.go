package route

import (
	"github.com/pufferApi/pkg/cache"
)

type Config struct {
	Apikey string
	Cache  cache.Config
}

type Route struct {
	Apikey   string
	Cache    cache.Cache
	Wildcard bool
}

func Create(config Config, Wildcard bool) (route Route) {
	route.Cache = cache.Create(config.Cache)
	route.Apikey = config.Apikey
	route.Wildcard = Wildcard
	return route

}
